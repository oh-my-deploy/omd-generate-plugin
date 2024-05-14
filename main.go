package main

import (
	"log"
	"os"

	"github.com/oh-my-deploy/omd-generate-plugin/cmd"
	"github.com/oh-my-deploy/omd-generate-plugin/utils"
)

func main() {
	rootCmd := cmd.CreateRootCmd()
	utils.RegisterSubCommands(rootCmd, cmd.InitVersionCmd, cmd.InitGenerateCmd)
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
