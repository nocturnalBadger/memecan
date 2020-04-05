package connectors

import (
	"bytes"
	"net/http"
	"encoding/json"
	"io/ioutil"
)

const getTextURL = "http://localhost:8080/base64"

type getTextRequest struct {
	Base64  string `json:"base64"`
	Trim    string `json:"trim"`
}

type getTextResponse struct {
	Result  string `json:"result"`
	Version string `json:"version"`
}

func GetImageText(base64Image string) string {
	requestPayload := getTextRequest{
		Base64: base64Image,
		Trim: "\n",
	}

	requestBytes, err := json.Marshal(&requestPayload)
	if err != nil {
		panic(err)
	}
	requestPayloadReader := bytes.NewReader(requestBytes)

	response, err := http.Post(getTextURL, "application/json", requestPayloadReader)
	if err != nil {
		panic(err)
	}

	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	var results getTextResponse
	err = json.Unmarshal(responseBytes, &results)
	if err != nil {
		panic(err)
	}

	return results.Result
}
