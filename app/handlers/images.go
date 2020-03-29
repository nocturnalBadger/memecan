package handlers

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/nocturnalBadger/memecan/app/models"
)

type key int
const imageKey key = iota


func Routes() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", listImages)
	r.Post("/", uploadImage)

	r.Route("/{imageID}", func(r chi.Router) {
		r.Use(imageCtx)
		r.Get("/", getImage)
	})

	return r
}


func uploadImage(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20) // limit your max input length!

	// in your case file would be fileupload
	file, header, err := r.FormFile("file")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Parse extension from file name
	name := strings.Split(header.Filename, ".")
	extension := name[1]
	fmt.Printf("Extension:%s\n", name[1])

	// Take md5sum of uploaded file
	h := md5.New()
	io.Copy(h, file)
	fileSum := h.Sum(nil)
	// Rewind file Reader to start
	file.Seek(0, io.SeekStart)

	objectName := fmt.Sprintf("%x.%s", fileSum, extension)
	log.Printf("Creating object %s", objectName)

	// Store object in minio
	bucketName := "images"
	minioClient.PutObject(bucketName, objectName, file, -1, minio.PutObjectOptions{})

	image := models.Image{Bucket: bucketName, ObjectName: objectName}
	db.Create(&image)

	log.Printf("Created Image record %v\n", image)

	jsonImage, err := json.Marshal(image)
	if err != nil {
		panic(err)
	}

	w.Write(jsonImage)

	return
}

func imageCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		imageID := chi.URLParam(r, "imageID")
		fmt.Println("yo")
		log.Printf("Request for image %s", imageID)

		var image models.Image
		result := db.Where("id = ?", imageID).First(&image)
		log.Printf("Retrieved image %v", image)

		if result.RowsAffected == 0 {
			http.Error(w, http.StatusText(404), 404)
			return
		}

		ctx := context.WithValue(r.Context(), imageKey, &image)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getImage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	image, ok := ctx.Value(imageKey).(*Image)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	object, err := minioClient.GetObject(image.Bucket, image.ObjectName, minio.GetObjectOptions{})
	if err != nil {
		panic(err)
	}

	io.Copy(w, object)
}

func listImages(w http.ResponseWriter, r *http.Request) {
	var images []models.Image
	if result := db.Find(&images); result.Error != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	jsonValue, _ := json.Marshal(images)

	w.Write(jsonValue)
	w.Write([]byte("\n"))
}
