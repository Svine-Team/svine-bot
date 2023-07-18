package app

import "github.com/Svine-Team/svine-bot/internal/greeting"

func Init() string {
	return greeting.HelloWorld()
}
