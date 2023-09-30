package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/alexditu/gogymapp/internal/server"
	"github.com/spf13/cobra"
)

type Arguments struct {
	settingsPath string
}

var args Arguments

func NewRootCmd(binaryName string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     binaryName,
		Version: "0.0.1",
		Short:   fmt.Sprintf("%s - HTTP REST API Server", binaryName),
		Long:    fmt.Sprintf("%s - HTTP REST API Server", binaryName),
		Run: func(cmd *cobra.Command, cmdlineArgs []string) {

			var setts server.Settings
			loadSettings(&setts, args.settingsPath)

			err := server.Run(binaryName, &setts)
			if err != nil {
				os.Exit(1)
			}
		},
	}

	rootCmd.Flags().StringVarP(&args.settingsPath, "settings", "s", "", "path to settings.jso file")

	return rootCmd
}

func loadSettings(setts *server.Settings, settsPath string) error {
	if settsPath == "" {
		setts.InitDefault()
		fmt.Printf("Using default server settings")
		return nil
	}

	// load from file

	rawData, err := os.ReadFile(settsPath)
	if err != nil {
		fmt.Printf("Error: failed to read settings file '%s': %s", settsPath, err)
		return err
	}

	err = json.Unmarshal(rawData, setts)
	if err != nil {
		fmt.Printf("Error: failed to parse json settings file '%s': %s", settsPath, err)
		return err
	}

	return nil
}
