package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "pokeapi",
	Short: "A terminal application for pokeapi",
	Long:  `A terminal application that searches for pokemon data using API provided by PokeAPI.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.pokeapi.yaml)")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	_, err := os.Stat("cache.txt")

	if err != nil {
		if os.IsNotExist(err) {
			file, err := os.Create("cache.txt")
			if err != nil {
				fmt.Println(err.Error())
			}
			file.Close()
		} else {
			fmt.Println(err)
		}
	}
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".pokeapi")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
