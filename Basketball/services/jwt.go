package services

import (
	"github.com/gbrlsnchs/jwt/v3"
	"time"
)

// JWT default/fixed claims
var (
	exp = 24 * 30 * time.Hour // 30 days
	nbf = 30 * time.Second    // 30 minutes
)

type CustomPayload struct {
	jwt.Payload
	Foo string `json:"foo,omitempty"`
	Bar int    `json:"bar,omitempty"`
}

var hs = jwt.NewHS256([]byte("secret"))

// Validate JWT token
func ValidateToken(token []byte) (answer bool) {

	var (
		now = time.Now()
		aud = jwt.Audience{"https://golang.org"}

		// Validate claims "iat", "exp" and "aud".
		iatValidator = jwt.IssuedAtValidator(now)
		expValidator = jwt.ExpirationTimeValidator(now)
		audValidator = jwt.AudienceValidator(aud)

		// Use jwt.ValidatePayload to build a jwt.VerifyOption.
		// Validators are run in the order informed.
		pll             CustomPayload
		validatePayload = jwt.ValidatePayload(&pll.Payload, iatValidator, expValidator, audValidator)
	)

	_, e := jwt.Verify(token, hs, &pll, validatePayload)
	if e != nil {
		return false
	}
	return true
}

// Generate JWT token for user
func MakeToken() (token string, err error) {

	// Generate JWT Token for user
	now := time.Now()
	pl := CustomPayload{
		Payload: jwt.Payload{
			Issuer:         "gbrlsnchs",
			Subject:        "someone",
			Audience:       jwt.Audience{"https://golang.org", "https://jwt.io"},
			ExpirationTime: jwt.NumericDate(now.Add(exp * 12 * time.Hour)),
			NotBefore:      jwt.NumericDate(now.Add(nbf)),
			IssuedAt:       jwt.NumericDate(now),
			JWTID:          "foobar",
		},
		Foo: "foo",
		Bar: 1337,
	}

	t, e := jwt.Sign(pl, hs)
	if e != nil {
		return "", e
	}

	// Return token
	return string(t), nil
}
