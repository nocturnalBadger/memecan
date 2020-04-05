package connectors

import (
	"log"
	"io"

	"github.com/minio/minio-go/v6"
)

var MinioClient *minio.Client
const BucketName = "images"

func InitMinio() {
	var err error
	MinioClient, err = minio.New("localhost:9000", "minioadmin", "minioadmin", false)
	if err != nil {
		panic(err)
	}

	location := "us-east-1"

	err = MinioClient.MakeBucket(BucketName, location)
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := MinioClient.BucketExists(BucketName)
		if errBucketExists == nil && exists {
			log.Printf("Bucket %s already exists.\n", BucketName)
		} else {
			log.Fatal(err)
		}
	} else {
		log.Printf("Successfully created bucket %s\n", BucketName)
	}
}

func SaveImage(name string, data io.Reader) (int64, error) {
	return MinioClient.PutObject(BucketName, name, data, -1, minio.PutObjectOptions{})
}

func GetObject(name string) io.Reader {
	object, err := MinioClient.GetObject(BucketName, name, minio.GetObjectOptions{})
	if err != nil {
		panic(err)
	}

	return object
}
