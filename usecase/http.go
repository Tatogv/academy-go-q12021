package usecase

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"api/entities"
)

func HandleError(w http.ResponseWriter, err error, status int) {
	w.WriteHeader(status)
	w.Write([]byte(err.Error()))
	return
}

func HandleSuccess(w http.ResponseWriter, response []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
	return
}

func GetFromApi(url string) (*entities.Response, error) {
	response := &entities.Response{}

	apiResponse, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if apiResponse.StatusCode != http.StatusOK {
		errorMessage := "Request to external API failed with status " + strconv.Itoa(apiResponse.StatusCode)
		err = errors.New(errorMessage)
		return nil, err
	}
	decoder := json.NewDecoder(apiResponse.Body)
	err = decoder.Decode(response)
	return response, nil
}
