package main

import (
	_ "embed"
	"log"
	"os"
	"panteks/cmd"
	"github.com/joho/godotenv"
)
//go:embed .env
var envFile string

func main() {

	envMap, err := godotenv.Unmarshal(envFile)
	if err != nil {
		log.Fatal("Error parsing embedded .env:", err)
	}
	for k, v := range envMap {
		os.Setenv(k, v)
	}

	cmd.ApiKey = os.Getenv("API_KEY")
	cmd.ApiURL = os.Getenv("API_URL")
	cmd.CommandString = os.Getenv("COMMAND_STRING")

	cmd.Execute()
}
