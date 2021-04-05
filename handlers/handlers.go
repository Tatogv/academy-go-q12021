package handlers

import (
	"api/usecase"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

const FilePath = "assets/pokemon.csv"

func GetAll(w http.ResponseWriter, r *http.Request) {
	records, err := usecase.ReadCsv(FilePath)
	if err != nil {
		usecase.HandleError(w, err, http.StatusInternalServerError)
		return
	}
	response, err := usecase.ParseRecordsToJSON(records)
	if err != nil {
		usecase.HandleError(w, err, http.StatusInternalServerError)
		return
	}
	usecase.HandleSuccess(w, response)
}

func GetById(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	pokemonId, _ := strconv.Atoi(pathParams["pokemonId"])
	records, err := usecase.ReadCsv(FilePath)
	if err != nil {
		usecase.HandleError(w, err, http.StatusInternalServerError)
		return
	}
	pokemonMap, err := usecase.ParseRecordsToMap(records)

	if err != nil {
		usecase.HandleError(w, err, http.StatusInternalServerError)
		return
	}
	if _, ok := pokemonMap[pokemonId]; ok {
		result := make(map[string]string)
		result["id"] = strconv.Itoa(pokemonId)
		result["name"] = pokemonMap[pokemonId]
		response, err := json.Marshal(result)
		if err != nil {
			usecase.HandleError(w, err, http.StatusInternalServerError)
			return
		}
		usecase.HandleSuccess(w, response)
		return
	} else {
		err = errors.New("Resource Not Found")
		usecase.HandleError(w, err, http.StatusNotFound)
	}
}

func GetBerries(w http.ResponseWriter, r *http.Request) {
	berries, err := usecase.GetFromApi("https://pokeapi.co/api/v2/berry")
	if err != nil {
		usecase.HandleError(w, err, http.StatusInternalServerError)
		return
	}
	err = usecase.WriteCSVFromResponse("./berries.csv", berries)
	response, err := json.Marshal(berries.Results)
	if err != nil {
		usecase.HandleError(w, err, http.StatusInternalServerError)
		return
	}
	usecase.HandleSuccess(w, response)
}

func ReadConcurrently(w http.ResponseWriter, r *http.Request) {
	var err error
	berriesMap := make(map[string]string)

	queryParametersReader := r.URL.Query()
	readType := queryParametersReader.Get("type")
	items, _ := strconv.Atoi(queryParametersReader.Get("items"))
	itemsPerWorker, _ := strconv.Atoi(queryParametersReader.Get("items_per_worker"))

	initialCounterPosition, err := usecase.ValidateConcurrentReadParams(readType, items, itemsPerWorker)
	if err != nil {
		usecase.HandleError(w, err, http.StatusInternalServerError)
		return
	}
	records, err := usecase.ReadCsv("berries.csv")
	if err != nil {
		usecase.HandleError(w, err, http.StatusInternalServerError)
		return
	}
	berriesMap = usecase.ReadRecordsConcurrently(records, items, itemsPerWorker, initialCounterPosition)

	response, err := json.Marshal(berriesMap)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error parsing result"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
