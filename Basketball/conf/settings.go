package conf

import (
	"github.com/beego/beego/v2/core/logs"
	"github.com/joho/godotenv"
	"os"
)

// use godot package to load/read the .env file and
// return the value of the key

func GetEnvConst(key string) string {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		logs.Error("Error loading .env file ", err)
	}

	return os.Getenv(key)
}
