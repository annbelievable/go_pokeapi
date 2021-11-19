package cmd

import (
	"fmt"
	"pokeapi/api"

	"github.com/spf13/cobra"
)

var pokemonCmd = &cobra.Command{
	Use:   "pokemon",
	Short: "Fetch a pokemon's data",
	Long:  `Fetch a pokemon's data by using name or id number`,
	Run: func(cmd *cobra.Command, args []string) {

		nameOrId, _ := cmd.Flags().GetString("query")

		if len(nameOrId) == 0 {
			fmt.Println("Please enter a pokemon name or id number.")
		} else {
			api.SearchPokemon(nameOrId)
		}
	},
}

func init() {
	rootCmd.AddCommand(pokemonCmd)
	pokemonCmd.Flags().StringP("query", "q", "", "Name or id number of the pokemon.")
}
