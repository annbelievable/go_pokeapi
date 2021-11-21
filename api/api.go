package api

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type Pokemon struct {
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

type Encouters []Encouter

type Encouter struct {
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

type PokeCache struct {
	Id         int       `json:"id"`
	Name       string    `json:"name"`
	Types      []string  `json:"types"`
	Stats      []string  `json:"stats"`
	Encounters []string  `json:"encounters"`
	Date       time.Time `json:"date"`
}

type Error string

func (e Error) Error() string { return string(e) }

const (
	noResultErr   = Error("No result is found.")
	cacheFile     = "cache.txt"
	kantoRegion   = "kanto"
	pokemonApiUrl = "https://pokeapi.co/api/v2/pokemon/"
)

func SearchPokemon(nameOrId string) {
	cache, err := GetCachedResult(nameOrId)

	if err != nil {
		if err == noResultErr {
			newCache, err := CallPokemonApi(nameOrId)
			if err != nil {
				fmt.Println("An error occured:", err)
				return
			}
			go CacheResult(newCache)
			PrintResult(newCache)
			return
		} else {
			fmt.Println(err)
			return
		}
	}

	if MoreThanAWeek(cache.Date) {
		newCache, err := CallPokemonApi(nameOrId)
		if err != nil {
			fmt.Println("An error occured:", err)
			return
		}
		go UpdateCache(cache, newCache)
		PrintResult(newCache)
		return
	} else {
		PrintResult(cache)
		return
	}
}

func PrintResult(cache PokeCache) {
	fmt.Println("Id:", cache.Id)
	fmt.Println("Name:", cache.Name)

	fmt.Print("Types: ")
	for _, ctype := range cache.Types {
		fmt.Print(ctype)
	}
	fmt.Println()

	for _, stat := range cache.Stats {
		fmt.Println(stat)
	}

	fmt.Print("Encounter locations and methods:")
	if len(cache.Encounters) > 0 {
		fmt.Println()
		for _, encounter := range cache.Encounters {
			fmt.Println("  ", encounter)
		}
	} else {
		fmt.Printf(" -\n")
	}
}

func GetCachedResult(nameOrId string) (PokeCache, error) {
	var cache PokeCache
	f, err := os.Open(cacheFile)
	defer f.Close()

	if err != nil {
		return cache, err
	}

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		bytes := scanner.Bytes()
		json.Unmarshal(bytes, &cache)
		if nameOrId == fmt.Sprintf("%d", cache.Id) || nameOrId == cache.Name {
			return cache, nil
		}
	}

	return PokeCache{}, noResultErr
}

func CacheResult(cache PokeCache) {
	f, err := os.OpenFile(cacheFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}

	bytes, err := json.Marshal(cache)
	if err != nil {
		log.Println(err)
	}

	_, err = fmt.Fprintln(f, string(bytes))
	if err != nil {
		log.Println(err)
	}
}

func UpdateCache(oldCache, newCache PokeCache) {
	input, err := ioutil.ReadFile(cacheFile)
	if err != nil {
		log.Println(err)
	}

	lines := strings.Split(string(input), "\n")

	oldBytes, err := json.Marshal(oldCache)
	if err != nil {
		log.Println(err)
	}

	newBytes, err := json.Marshal(newCache)
	if err != nil {
		log.Println(err)
	}

	for i, line := range lines {
		if line == string(oldBytes) {
			lines[i] = string(newBytes)
		}
	}

	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(cacheFile, []byte(output), 0644)
	if err != nil {
		log.Println(err)
	}
}

func MoreThanAWeek(dt time.Time) bool {
	return int(time.Now().Sub(dt).Hours()/24) > 7
}

func CallPokemonApi(nameOrId string) (PokeCache, error) {
	var pokemon Pokemon
	apiUrl := pokemonApiUrl + nameOrId
	bodyBytes, err := CallApi(apiUrl)

	if err != nil {
		return PokeCache{}, err
	}

	json.Unmarshal(bodyBytes, &pokemon)

	cache := PokeCache{
		Id:         pokemon.Id,
		Name:       pokemon.Name,
		Types:      []string{},
		Stats:      []string{},
		Encounters: []string{},
		Date:       time.Now(),
	}

	for _, pokeType := range pokemon.Types {
		cache.Types = append(cache.Types, pokeType.Type.Name)
	}

	for _, pokeStat := range pokemon.Stats {
		formattedStats := fmt.Sprintf("%s: %d", pokeStat.StatStruct.StatName, pokeStat.StatNumber)
		cache.Stats = append(cache.Stats, formattedStats)
	}

	if len(pokemon.EncounterUrl) > 0 {
		cache.Encounters = GetEncounterLocation(pokemon.EncounterUrl)
	}

	return cache, nil
}

func GetEncounterLocation(encounterUrl string) []string {
	encounterBytes, err := CallApi(encounterUrl)

	if err != nil {
		return []string{""}
	}

	var encounters Encouters
	json.Unmarshal(encounterBytes, &encounters)

	var wg sync.WaitGroup
	ch := make(chan string)
	for _, encounter := range encounters {
		wg.Add(1)
		go GetKantoEncounterLocation(ch, encounter, &wg)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	result := []string{}
	for {
		res, ok := <-ch
		if ok == false {
			break
		}
		result = append(result, res)
	}

	return result
}

func GetKantoEncounterLocation(ch chan string, encounter Encouter, wg *sync.WaitGroup) {
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

	if location.Region.Name != kantoRegion {
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

	result := encounter.LocationArea.Name + " - "

	methodIx := 0
	methodsLen := len(methods)
	for key, _ := range methods {
		result += key
		methodIx++
		if methodIx != methodsLen {
			result += ", "
		}
	}

	ch <- result
}

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
		return nil, noResultErr
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
