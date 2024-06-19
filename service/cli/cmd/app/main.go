package main

import (
	"log"

	"github.com/javascriptizer1/grpc-cli-chat.backend/service/cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatalf("failed to execute root command: %s", err.Error())
	}
}
