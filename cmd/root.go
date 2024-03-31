/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "data-generator",
	Short: "you can use commands for generate/read data",
	Long:  `Data-Generator is a tool that gives you necessary data for our team`,
}

// Execute adds all child commands to the root command and sets flags appropriately.

func init() {
}
