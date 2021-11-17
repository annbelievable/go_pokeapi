package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/cobra"
)

type PokemonData struct {
	Id           int           `json:"id"`
	Name         string        `json:"name"`
	EncounterUrl string        `json:"location_area_encounters"`
	Types        []PokemonType `json:"types"`
	Stats        []PokemonStat `json:"stats"`
}

type PokemonType struct {
	Type struct {
		Name string `json:"name"`
	} `json:"type"`
}

type PokemonStat struct {
	StatNumber int `json:"base_stat"`
	StatStruct struct {
		StatName string `json:"name"`
	} `json:"stat"`
}

// addCmd represents the add command
var pokemonCmd = &cobra.Command{
	Use:   "pokemon",
	Short: "Fetch a pokemon's data",
	Long:  `Fetch a pokemon's data by using name or id number`,
	Run: func(cmd *cobra.Command, args []string) {
		//if a command has any flags, get them here like the example below
		//the flag can be processed here or passed to the next function
		//you can run mutilple funtions
		// fstatus, _ := cmd.Flags().GetBool("float")

		//use this to get the full map/struct of the args
		// cmdflags := cmd.Flags()
		// fmt.Println(cmdflags)
		// cmdArgs := cmd.Flags().Args()
		// fmt.Println(cmdArgs)

		//using getstring
		nameOrId, _ := cmd.Flags().GetString("query")

		if len(nameOrId) == 0 {
			fmt.Println("Please enter a pokemon name or id number.")
		} else {
			searchPokemon(nameOrId)
		}
	},
}

func init() {
	rootCmd.AddCommand(pokemonCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	//make the name flag required
	pokemonCmd.Flags().StringP("query", "q", "", "Name or id number of the pokemon.")
}

func searchPokemon(nameOrId string) {
	pokemon := CallPokemonApi(nameOrId)

	fmt.Printf("Id: %d\n", pokemon.Id)
	fmt.Printf("Name: %s\n", pokemon.Name)

	fmt.Printf("Type: ")
	for i := 0; i < len(pokemon.Types); i++ {
		fmt.Printf("%s", pokemon.Types[i].Type.Name)
		if i < len(pokemon.Types)-1 {
			fmt.Print(", ")
		}
	}
	fmt.Println()

	for i := 0; i < len(pokemon.Stats); i++ {
		fmt.Printf("%s: %d\n", pokemon.Stats[i].StatStruct.StatName, pokemon.Stats[i].StatNumber)
	}
}

func CallPokemonApi(nameOrId string) PokemonData {
	client := &http.Client{}
	apiUrl := "https://pokeapi.co/api/v2/pokemon/" + nameOrId
	req, err := http.NewRequest("GET", apiUrl, nil)

	if err != nil {
		fmt.Print(err.Error())
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-type", "application/json")

	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		fmt.Print(err.Error())
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Print(err.Error())
	}

	var pokemon PokemonData
	json.Unmarshal(bodyBytes, &pokemon)

	return pokemon
}
