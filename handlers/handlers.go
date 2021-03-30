package handlers

import (
	"api/entities"
	"encoding/csv"
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

const FilePath = "assets/pokemon.csv"

func GetAll(w http.ResponseWriter, r *http.Request) {
	pokemonMap := make(map[int]string)
	w.Header().Set("Content-Type", "application/json")

	file, err := os.Open(FilePath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("CSV is invalid"))
		return
	}
	for _, pokemon := range records {
		id, err := strconv.Atoi(pokemon[0])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error parsing ID"))
			return
		}
		pokemonMap[id] = pokemon[1]
	}
	response, err := json.Marshal(pokemonMap)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error parsing result"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func GetById(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	pokemonMap := make(map[string]string)
	pokemonId := pathParams["pokemonId"]
	w.Header().Set("Content-Type", "application/json")

	file, err := os.Open(FilePath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error reading CSV file"))
		return
	}
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("CSV is invalid"))
		return
	}
	for _, pokemon := range records {
		pokemonMap[pokemon[0]] = pokemon[1]
	}

	if _, ok := pokemonMap[pokemonId]; ok {
		result := make(map[string]string)
		result["id"] = pokemonId
		result["name"] = pokemonMap[pokemonId]
		response, err := json.Marshal(result)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error parsing result"))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Resource Not Found with id"))
	}
}

func GetBerries(w http.ResponseWriter, r *http.Request) {

	apiResponse, err := http.Get("https://pokeapi.co/api/v2/berry")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error parsing result"))
		return
	}

	berries := &entities.BerryResponse{}

	decoder := json.NewDecoder(apiResponse.Body)
	err = decoder.Decode(berries)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error decoding"))
		return
	}

	csvFile, err := os.Create("./berries.csv")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error creating CSV"))
		return
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	for index, berry := range berries.Results {
		var row []string
		row = append(row, strconv.Itoa(index))
		row = append(row, berry.Name)
		row = append(row, berry.Url)
		writer.Write(row)
	}
	writer.Flush()

	body, err := json.Marshal(berries.Results)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error reading result"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(body)
}
