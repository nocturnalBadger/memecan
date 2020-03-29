package main

import (
	"crypto/md5"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/minio/minio-go/v6"
	"io"
	"log"
	"net/http"
	//	"bytes"
	"strings"
	//	"os"
	"context"
	"encoding/json"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/oklog/ulid"
	"math/rand"
	"time"
	"io/ioutil"
	"errors"
)

// Image an image
type Image struct {
	ID         string `json:"id" gorm:"primary_key"`
	Bucket     string `json:"bucket"`
	ObjectName string `json:"object_name"`
}

// Tag tag
type Tag struct {
	Text string
}

// Meme a meme
type Meme struct {
	ID    string `json:"id" gorm:"primary_key"`
	Image Image  `gorm:"foreignkey:ID" json:"-"`
	Text  string `json:"text,omitempty"`
	Tags  []Tag  `gorm:"foreignkey:string" json:"-"`
}

type MemeAlias Meme
type MemeForm struct {
	MemeAlias
	Image  string   `json:"image_id"`
	Tags   []string `json:"tags"`
}


func (m Meme) MarshalJSON() ([]byte, error) {
	// Alias used to bring all of Meme's existing fields over
	tagNames := []string{}
	for _, tag := range m.Tags {
		tagNames = append(tagNames, tag.Text)
	}

	return json.Marshal(MemeForm{
		MemeAlias: MemeAlias(m),
		Image: m.Image.ID,
		Tags: tagNames,
	})
}


func (m *Meme) UnmarshalJSON(data []byte) error {

	var form MemeForm
	if err := json.Unmarshal(data, &form); err != nil {
		return err
	}

	result := db.Where("id = ?", form.Image).First(&m.Image)
	if result.Error != nil {
		return errors.New("Image not found")
	}


	for _, tagStr := range form.Tags {
		var tag Tag
		db.FirstOrCreate(&tag, Tag{Text: tagStr})
		m.Tags = append(m.Tags, tag)
	}

	return nil
}


type memeRequest struct {
	ImageID string   `json:"image_id"`
	Text    string   `json:"text"`
	Tags    []string `json:"tags"`
}

type key int
const imageKey key = iota


// BeforeCreate Run before creating Image
func (base *Image) BeforeCreate(scope *gorm.Scope) error {
	ulid := getULID()
	return scope.SetColumn("ID", ulid)
}

// BeforeCreate Run before creating Image
func (base *Meme) BeforeCreate(scope *gorm.Scope) error {
	ulid := getULID()
	return scope.SetColumn("ID", ulid)
}

var db *gorm.DB
var minioClient *minio.Client

func main() {
	log.Println("Starting memecan server")
	r := chi.NewRouter()
	minioClient = initMinio()

	var err error
	db, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=memecan sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db.AutoMigrate(&Image{})
	db.AutoMigrate(&Tag{})
	db.AutoMigrate(&Meme{})
	log.Println("Database ready.")

	r.Route("/images", func(r chi.Router) {
		r.Get("/", listImages)
		r.Post("/", uploadImage)

		r.Route("/{imageID}", func(r chi.Router) {
			r.Use(imageCtx)
			r.Get("/", getImage)
		})
	})

	r.Post("/memes", createMeme)

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
	t := time.Now().UTC()
	entropy := rand.New(rand.NewSource(t.UnixNano()))

	return ulid.MustNew(ulid.Timestamp(t), entropy).String()
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

	image := Image{Bucket: bucketName, ObjectName: objectName}
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

		var image Image
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
	var images []Image
	if result := db.Find(&images); result.Error != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	jsonValue, _ := json.Marshal(images)

	w.Write(jsonValue)
	w.Write([]byte("\n"))
}

func createMeme(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	var meme Meme
	err = json.Unmarshal(body, &meme)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), 400)
		return
	}
	log.Printf("Unmarshalled meme: %v", meme)

	db.Create(&meme)
	jsonValue, _ := json.Marshal(meme)

	w.Write(jsonValue)
	w.Write([]byte("\n"))
}
