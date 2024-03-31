/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	clients "datageneratorbookapi/clients"
	"fmt"
	"os"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

var environment string
var bookApiUrl string

// getBooksCmd represents the book command
var getBooksCmd = &cobra.Command{
	Use:   "getbooks",
	Short: "Get Book List",
	Long:  `Getbooks command that gives you book list`,
	Run: func(cmd *cobra.Command, args []string) {

		checkEnvironment()
		getUrl()

		requestInfo := clients.RequestInfo{
			Environment: environment,
		}

		books := clients.GetBooksRequest(requestInfo)

		fmt.Println(books)

		os.Exit(1)
	},
}

// postBookCmd represents the book command
var postBookCmd = &cobra.Command{
	Use:   "postbook",
	Short: "Post book",
	Long:  `Postbook command that create you random book`,
	Run: func(cmd *cobra.Command, args []string) {

		checkEnvironment()
		getUrl()

		requestInfo := clients.RequestInfo{
			Environment: environment,
		}

		bookRequest := clients.CreateBook() // Call book api (stage,dev,local whatever)

		bookBarcode := clients.PostBookRequest(requestInfo, bookRequest)

		fmt.Println(bookBarcode)
		clients.WriteAll(bookBarcode)

		os.Exit(1)
	},
}

func Execute() {
	err := getBooksCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(getBooksCmd)
	rootCmd.AddCommand(postBookCmd)
	cobra.OnInitialize(initConfig)
	postBookCmd.Flags().StringVarP(&environment, "url", "u", "", "Url for your api")
	postBookCmd.Flags().StringVarP(&environment, "environment", "e", "", "Environment for the authentication. Default value is stage")
	getBooksCmd.Flags().StringVarP(&environment, "environment", "e", "", "Environment for the authentication. Default value is stage")
	getBooksCmd.Flags().StringVarP(&environment, "url", "u", "", "Url for your api")

}

func checkEnvironment() {
	if len(environment) == 0 {
		environment = viper.GetString("environment")
	}

	if len(environment) == 0 {
		fmt.Println("Setting default environment")
		environment = "local"
	}

}

func getUrl() {
	if len(bookApiUrl) == 0 {
		if err := viper.ReadInConfig(); err == nil {

			if viper.IsSet("url") && len(bookApiUrl) == 0 {
				bookApiUrl = viper.GetString("url")
			}

		} else {
			fmt.Println("If you want to use the command without parameters, you need to create the datagenerator.yaml file under the &HOME path.")
			os.Exit(1)
		}
	}
}

func initConfig() {
	viper.SetDefault("environment", "stage")
	viper.AddConfigPath("$HOME")
	viper.SetConfigFile("datagenerator.yaml")
	viper.AutomaticEnv()
	viper.ReadInConfig()
}
