package app

import (
	"github.com/go-chi/chi"
	"github.com/jinzhu/gorm"
	"github.com/minio/minio-go/v6"

	"github.com/nocturnalBadger/memecan/app/handlers"
)

// App does stuff
type App struct {
	Router      *chi.Router
	DB          *gorm.DB
	MinioClient *minio.Client
}

func Routes() *chi.Mux {
	r := chi.NewRouter()

	r.Mount("/images", handlers.images.Routes())

	return r
}

func Initialize *App {
	var app App

	app.Router = Routes()
	app.DB = DBConnect()
	app.MinioClient = initMinio()
}

func (a *App) Run(host string) {
	log.Fatal(http.ListenAndServe(host, a.Router)
}

func AppContext(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx := context.WithValue(r.Context(), imageKey, &image)
		next.ServeHTTP(w, r.WithContext(ctx))
	})

}

func DBConnect() *gorm.DB, error {
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=memecan sslmode=disable")
	if err != nil {
		return db, err
	}
	defer db.Close()

	models.DBMigrate(db)

	log.Println("Database ready.")

	return db, nil
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
