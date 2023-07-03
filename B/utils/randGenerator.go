package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func RandInt(min int, max int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	num := r.Intn((max - min)) + min
	if num < min || num > max {
		RandInt(min, max)
	}
	return fmt.Sprintf("%d", num)
}
