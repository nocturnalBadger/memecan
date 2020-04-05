package connectors

import (
	"fmt"
	"strconv"
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


type SearchResult struct {
	Hits  struct{
		Hits  []SearchHit `json:"hits"`
	} `json:"hits"`
}

type SearchHit struct {
	Index  string      `json:"_index"`
	ID     string      `json:"_id"`
	Score  float32     `json:"_score"`
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

func Search(query string, limit int) []SearchHit {
	url := fmt.Sprintf("%s/%s/_search", esBaseUrl, index)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	reqQuery := req.URL.Query()
	reqQuery.Add("size", strconv.Itoa(limit))
	if query != "" {
		reqQuery.Add("q", query)
	}
	req.URL.RawQuery = reqQuery.Encode()
	log.Printf("%s\n", req.URL.RawQuery)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	log.Printf("%s\n", respBody)

	var result SearchResult
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		panic(err)
	}

	return result.Hits.Hits
}
