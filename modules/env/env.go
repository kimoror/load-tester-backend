package env

import (
	"github.com/joho/godotenv"
	"log"
)

func Load() {
	err := godotenv.Load("local.env")
	if err != nil {
		log.Fatalf("Error occured while loading file: %s", err)
	}
}
