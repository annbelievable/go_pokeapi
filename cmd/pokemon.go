package cmd

import (
	"encoding/json"
	"errors"
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

var pokemonCmd = &cobra.Command{
	Use:   "pokemon",
	Short: "Fetch a pokemon's data",
	Long:  `Fetch a pokemon's data by using name or id number`,
	Run: func(cmd *cobra.Command, args []string) {

		nameOrId, _ := cmd.Flags().GetString("query")

		if len(nameOrId) == 0 {
			fmt.Println("Please enter a pokemon name or id number.")
		} else {
			SearchPokemon(nameOrId)
		}
	},
}

func init() {
	rootCmd.AddCommand(pokemonCmd)
	pokemonCmd.Flags().StringP("query", "q", "", "Name or id number of the pokemon.")
}

func SearchPokemon(nameOrId string) {
	pokemon, err := CallPokemonApi(nameOrId)

	if err != nil {
		fmt.Println("An error occured:", err)
		return
	}

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

func CallPokemonApi(nameOrId string) (PokemonData, error) {
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

	var pokemon PokemonData

	if resp.StatusCode != 200 {
		respErr := errors.New("Something wrong, please try again later.")
		if resp.StatusCode == 404 {
			respErr = errors.New("No pokemon is found.")
		}
		return pokemon, respErr
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Print(err.Error())
	}

	json.Unmarshal(bodyBytes, &pokemon)

	return pokemon, nil
}
