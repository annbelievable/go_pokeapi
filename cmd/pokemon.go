package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

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

type EncoutersData []EncouterData

type EncouterData struct {
	LocationArea struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"location_area"`
	VersionData []struct {
		EncounterDetails []EncounterDetail `json:"encounter_details"`
	} `json:"version_details"`
}

type EncounterDetail struct {
	Method struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"method"`
}

type LocationArea struct {
	Location struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"location"`
}

type Location struct {
	Region struct {
		Name string `json:"name"`
	} `json:"region"`
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

	if len(pokemon.EncounterUrl) > 0 {
		fmt.Println("Encounter locations and methods:")
		encounterLocation := GetEncounterLocation(pokemon.EncounterUrl)
		if encounterLocation != "" {
			fmt.Print(encounterLocation)
		} else {
			fmt.Print("-")
		}

	}
}

func CallPokemonApi(nameOrId string) (PokemonData, error) {
	var pokemon PokemonData
	apiUrl := "https://pokeapi.co/api/v2/pokemon/" + nameOrId

	bodyBytes, err := CallApi(apiUrl)

	if err != nil {
		return pokemon, err
	}

	json.Unmarshal(bodyBytes, &pokemon)

	return pokemon, nil
}

func GetEncounterLocation(encounterUrl string) string {
	encounterBytes, err := CallApi(encounterUrl)

	if err != nil {
		return ""
	}

	var encounters EncoutersData
	json.Unmarshal(encounterBytes, &encounters)
	result := ""

	var wg sync.WaitGroup
	ch := make(chan string)
	for _, encounter := range encounters {
		wg.Add(1)
		go CallLocationApi(ch, encounter, &wg)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for {
		res, ok := <-ch
		if ok == false {
			break
		}
		result += res
	}

	return result
}

func CallLocationApi(ch chan string, encounter EncouterData, wg *sync.WaitGroup) {
	defer wg.Done()
	locationAreaBytes, err := CallApi(encounter.LocationArea.Url)

	if err != nil {
		return
	}

	var locationArea LocationArea
	json.Unmarshal(locationAreaBytes, &locationArea)
	locationBytes, err := CallApi(locationArea.Location.Url)

	if err != nil {
		return
	}

	var location Location
	json.Unmarshal(locationBytes, &location)

	if location.Region.Name != "kanto" {
		return
	}

	var methods = make(map[string]bool)
	for _, version := range encounter.VersionData {
		for _, detail := range version.EncounterDetails {
			_, ok := methods[detail.Method.Name]
			if !ok {
				methods[detail.Method.Name] = true
			}
		}
	}

	result := "    " + encounter.LocationArea.Name + " - "
	methodIx := 0
	methodsLen := len(methods)
	for key, _ := range methods {
		result += key
		methodIx++
		if methodIx != methodsLen {
			result += ", "
		}
	}
	result += "\n"

	ch <- result
}

//a generic call api method that returns bytes
func CallApi(apiUrl string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", apiUrl, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-type", "application/json")

	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 404 {
		return nil, errors.New("No result is found.")
	}

	if resp.StatusCode != 200 {
		return nil, errors.New("Something wrong, please try again later.")
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return bodyBytes, nil
}
