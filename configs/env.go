package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func EnvMongoUri() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loadin env file")
	}
	return os.Getenv("MONGO_URI")
}
