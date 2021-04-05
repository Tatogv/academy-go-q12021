package usecase

import (
	"api/entities"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"sync/atomic"
)

var readerCounter int32

func ReadCsv(filename string) ([][]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	return records, nil
}

func ParseRecordsToJSON(records [][]string) ([]byte, error) {
	recordsMap, err := ParseRecordsToMap(records)
	if err != nil {
		return nil, err
	}
	response, err := json.Marshal(recordsMap)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func ParseRecordsToMap(records [][]string) (map[int]string, error) {
	recordsMap := make(map[int]string)
	for _, line := range records {
		id, err := strconv.Atoi(line[0])
		if err != nil {
			return nil, err
		}
		recordsMap[id] = line[1]
	}
	return recordsMap, nil
}

func WriteCSVFromResponse(filename string, data *entities.Response) error {
	csvFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	for index, berry := range data.Results {
		var row []string
		row = append(row, strconv.Itoa(index))
		row = append(row, berry.Name)
		row = append(row, berry.Url)
		writer.Write(row)
	}
	writer.Flush()
	return nil
}

func ValidateConcurrentReadParams(readType string, items int, itemsPerWorker int) (int32, error) {
	var initialCounterPosition int32
	var err error
	if readType == "odd" {
		initialCounterPosition = 1
	} else if readType == "even" {
		initialCounterPosition = 0
	} else {
		err = errors.New("Invalid type")
		return 0, err
	}

	if items == 0 || itemsPerWorker == 0 {
		err = errors.New("Invalid items")
		return 0, err
	}
	return initialCounterPosition, nil
}

func ReadRecordsConcurrently(records [][]string, items int, itemsPerWorker int, initialCounterPosition int32) map[string]string {
	var numberOfWorkers int
	readableItems := len(records) / 2
	recordsMap := make(map[string]string)

	readerCounter = initialCounterPosition

	if items > readableItems {
		items = readableItems
	}

	if itemsPerWorker > readableItems {
		itemsPerWorker = readableItems
	}

	numberOfWorkers = int(math.Ceil(float64(items) / float64(itemsPerWorker)))

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
		recordsMap[line[0]] = line[1]
		fmt.Println(line[1])
	}

	close(results)

	return recordsMap
}

func worker(id int, jobs <-chan int, results chan<- []string, records *[][]string) {
	for range jobs {
		csvData := *records
		csvLine := csvData[readerCounter]
		atomic.AddInt32(&readerCounter, 2)
		results <- csvLine
	}
}
