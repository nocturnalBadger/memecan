package main

import (
	"github.com/go-chi/chi"
	"github.com/minio/minio-go/v6"
	"log"
	"net/http"
	"io"
	"fmt"
	"crypto/md5"
//	"bytes"
	"strings"
//	"os"
	"github.com/oklog/ulid"
	"time"
	"math/rand"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Image an image
type Image struct {
	ID         string  `sql: "type:ulid;primary_key;default:getULID()"`
	Bucket     string
	ObjectName string
}

// BeforeCreate Run before creating Image
func (base *Image) BeforeCreate(scope *gorm.Scope) error {
	ulid := getULID()
	return scope.SetColumn("ID", ulid)
}


func main() {
	log.Println("Starting memecan server")
	r := chi.NewRouter()
	minioClient := initMinio()

	db, err := gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=memecan sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db.AutoMigrate(&Image{})
	log.Println("Database ready.")

	r.Post("/images", func(w http.ResponseWriter, r *http.Request) {
		uploadImage(w, r, minioClient, db)
	})

	http.ListenAndServe(":3000", r)
}

func initMinio() *minio.Client {
	minioClient, err := minio.New("localhost:9000", "minioadmin", "minioadmin", false)
	if err != nil {
		log.Fatal(err)
	}

	bucketName := "images"
	location := "us-east-1"

	err = minioClient.MakeBucket(bucketName, location)
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(bucketName)
		if errBucketExists == nil && exists {
			log.Printf("Bucket %s already exists.\n", bucketName)
		} else {
			log.Fatal(err)
		}
	} else {
		log.Printf("Successfully created bucket %s\n", bucketName)
	}

	return minioClient
}

func getULID() string {
	t := time.Unix(1000000, 0)
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)

	return ulid.MustNew(ulid.Timestamp(t), entropy).String()
}

func uploadImage(w http.ResponseWriter, r *http.Request, minioClient *minio.Client, db *gorm.DB) {
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

	db.Create(&Image{Bucket: bucketName, ObjectName: objectName})

	fmt.Fprintf(w, "ok")
	// etc write header
	return
}
