package connectors

import (
	"fmt"
	"log"
	"net/http"
	"io/ioutil"
	"io"
	"strings"
	"encoding/json"
)

const mapping = `
{
  "mappings": {
    "properties": {
      "text": {
        "type": "text"
      },
      "tags": {
        "type": "keyword"
      },
      "image": {
        "properties": {
          "bucket_name": {
            "type": "keyword"
          },
          "filename": {
            "type": "keyword"
          }
        }
      }
    }
  }
}
`

type DocResponse struct {
	Found  bool        `json:"found"`
	Source interface{} `json:"_source"`
}

const esBaseUrl = "http://localhost:9200"
const index = "memecan"

func InitES() {
	url := fmt.Sprintf("%s/%s", esBaseUrl, index)

	req, err := http.NewRequest("PUT", url, strings.NewReader(mapping))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	respText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	log.Printf("%s\n", respText)
}

func CreateDoc(docID string, jsonData io.Reader) []byte {
	url := fmt.Sprintf("%s/%s/_doc/%s", esBaseUrl, index, docID)

	resp, err := http.Post(url, "application/json", jsonData)
	if err != nil {
		panic(err)
	}

	responseText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	log.Printf("%s\n", responseText)
	return responseText
}


func GetDoc(docID string, target interface{}) DocResponse {
	url := fmt.Sprintf("%s/%s/_doc/%s", esBaseUrl, index, docID)

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	result := DocResponse{
		Found: false,
		Source: target,
	}
	err = json.Unmarshal(respBytes, &result)
	if err != nil {
		panic(err)
	}

	return result
}
