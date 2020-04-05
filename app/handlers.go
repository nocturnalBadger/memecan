package app

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"encoding/base64"
	"bytes"
	"crypto/md5"
	"strings"
	"io"

	"github.com/go-chi/chi"

	"github.com/nocturnalBadger/memecan/connectors"
)

func Routes() *chi.Mux {
	router := chi.NewRouter()

	router.Route("/memes", func (r chi.Router) {
		r.Post("/", createMeme)
		r.Get("/{memeID}", getMeme)
	})

	router.Route("/images", func (r chi.Router) {
		r.Get("/{memeID}", getImage)
	})

	return router
}

func createMeme(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}


	var meme Meme
	err = json.Unmarshal(body, &meme)
	if err != nil {
		panic(err)
	}

	image := &meme.Image

	imageData, err := base64.StdEncoding.DecodeString(image.Base64)
	if err != nil {
		panic(err)
	}
	imageReader := bytes.NewReader(imageData)

	fileSum := md5.Sum(imageData)
	extension := strings.Split(image.FileName, ".")[1]

	image.FileName = fmt.Sprintf("%x.%s", fileSum, extension)

	_, err = connectors.SaveImage(image.FileName, imageReader)
	if err != nil {
		panic(err)
	}
	image.BucketName = connectors.BucketName

	meme.Text = connectors.GetImageText(image.Base64)

	// Remove base64 from response
	image.Base64 = ""

	// Marshal for storing in elasticsearch
	jsonData, err := json.Marshal(meme)
	if err != nil {
		panic(err)
	}

	// Store in elasticsearch
	memeID := GetULID()
	connectors.CreateDoc(memeID, bytes.NewReader(jsonData))

	// Need to marshal this again with the id (didn't want it insid the ES doc)
	meme.ID = memeID
	jsonData, err = json.MarshalIndent(meme, "", "    ")
	if err != nil {
		panic(err)
	}

	w.Write(jsonData)
}


func listMemes() {

}


func getMeme(w http.ResponseWriter, r *http.Request) {
	memeID := chi.URLParam(r, "memeID")

	var meme Meme
	result := connectors.GetDoc(memeID, &meme)

	if !result.Found {
		http.Error(w, http.StatusText(404), 404)
		return
	}

	jsonData, err := json.MarshalIndent(meme, "", "    ")
	if err != nil {
		panic(err)
	}

	w.Write(jsonData)
}

func getImage(w http.ResponseWriter, r *http.Request) {
	memeID := chi.URLParam(r, "memeID")

	var meme Meme
	result := connectors.GetDoc(memeID, &meme)

	if !result.Found {
		http.Error(w, http.StatusText(404), 404)
		return
	}

	object := connectors.GetObject(meme.Image.FileName)

	io.Copy(w, object)
}
