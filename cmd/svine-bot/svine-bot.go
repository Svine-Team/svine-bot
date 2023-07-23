package main

import (
	"fmt"
	"os"
    "log"

	"github.com/Svine-Team/svine-bot/pkg/app"
    "github.com/joho/godotenv"
)

// Variables used for command line parameters
var (
	Token string
)

func init() {

	// flag.StringVar(&Token, "t", "", "Bot Token")
	// flag.Parse()
    // godotenv.Load()
}

// Use godot package to load/read the .env file and
//   return the value of the key.
func getEnvVariable(key string) string {

  // load .env file
  err := godotenv.Load()

  if err != nil {
    log.Fatalf("Error loading .env file")
  }

  return os.Getenv(key)
}

func main() {
  // godotenv package
  envVariable := getEnvVariable("BOT_TOKEN")

  fmt.Printf("godotenv : %s = %s \n", "BOT_TOKEN", envVariable)
}
