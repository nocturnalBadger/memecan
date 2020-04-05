package app

import (
	"time"
	"math/rand"

	"github.com/oklog/ulid"
)

type Image struct {
	BucketName string `json:"bucket_name"`
	FileName   string `json:"filename"`
	Base64     string `json:"base64,omitempty"`
}

type Meme struct {
	Tags  []string `json:"tags"`
	Text  string   `json:"text"`
	Image Image    `json:"image"`
	ID    string   `json:"id,omitempty"`
}

func GetULID() string {
	t := time.Now().UTC()
	entropy := rand.New(rand.NewSource(t.UnixNano()))

	return ulid.MustNew(ulid.Timestamp(t), entropy).String()
}
