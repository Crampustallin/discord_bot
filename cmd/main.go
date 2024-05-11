package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Crampustallin/discord_bot/internal/aws"
	"github.com/Crampustallin/discord_bot/internal/bot"
	"github.com/joho/godotenv"
)

const (
	AWS_REGION_VAR_NAME      = "AWS_REGION"
	AWS_BUCKET_NAME_VAR_NAME = "AWS_BUCKET_NAME"

	DISCORD_BOT_TOKEN_VAR_NAME = "DISCORD_BOT_TOKEN"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	awsReg := checkEnv(AWS_REGION_VAR_NAME)
	bucketName := checkEnv(AWS_BUCKET_NAME_VAR_NAME)
	aws := aws.NewAws(awsReg, bucketName)

	token := checkEnv(DISCORD_BOT_TOKEN_VAR_NAME)
	bot := bot.NewBot(token)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	go func() { bot.Start() }()

	for {
		select {
		case file := <-bot.FileNameSend:
			err := aws.Upload(file)
			if err != nil {
				fmt.Println(err)
			}
		case <-signalChan:
			fmt.Println("exiting")
			bot.Close()
			return
		}
	}

}

func checkEnv(envName string) string {
	str := os.Getenv(envName)
	if str == "" {
		panic(errors.New("No variable " + envName + " found\n"))
	}
	return str
}
