package handlers

import (
	"api/entities"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"
	"sync/atomic"

	"github.com/gorilla/mux"
)

const FilePath = "assets/pokemon.csv"

var readerCounter int32

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

func ReadConcurrently(w http.ResponseWriter, r *http.Request) {
	var numberOfWorkers int
	pokemonMap := make(map[string]string)

	queryParametersReader := r.URL.Query()
	readType := queryParametersReader.Get("type")
	items, _ := strconv.Atoi(queryParametersReader.Get("items"))
	items_per_workers, _ := strconv.ParseFloat(queryParametersReader.Get("items_per_worker"), 64)

	numberOfWorkers = int(math.Ceil(float64(items) / items_per_workers))

	if readType == "odd" {
		readerCounter = 1
	} else if readType == "even" {
		readerCounter = 0
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unsupported read type"))
		return
	}
	file, err := os.Open("berries.csv")
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

	jobs := make(chan int, items)
	results := make(chan []string, items)

	for w := 1; w <= numberOfWorkers; w++ {
		go worker(w, jobs, results, &records)
	}

	for j := 1; j <= items; j++ {
		jobs <- j
	}
	close(jobs)

	for a := 1; a <= items; a++ {
		line := <-results
		pokemonMap[line[0]] = line[1]
		fmt.Println(line[1])
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

func worker(id int, jobs <-chan int, results chan<- []string, data *[][]string) {
	for j := range jobs {
		fmt.Println("worker", id, "started  job", j)
		currentData := *data
		current := currentData[readerCounter]
		fmt.Println("worker", id, "finished job", j)
		atomic.AddInt32(&readerCounter, 2)
		results <- current
	}
}
