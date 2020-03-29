package models

// Image an image
type Image struct {
	BaseObject
	Bucket     string `json:"bucket"`
	ObjectName string `json:"object_name"`
}
