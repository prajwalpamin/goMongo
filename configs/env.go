package configs

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func EnvMongoUri() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loadin env file")
	}
	fmt.Println("env file loaded successfully %v", os.Getenv("MONGO_URI"))
	return os.Getenv("MONGO_URI")
}
