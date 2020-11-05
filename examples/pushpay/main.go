package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/mobilemoney/mpesa"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	fmt.Println(mpesa.Version())
}