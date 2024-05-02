package data

import (
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/text/language"
)

type langData struct {
	Language string `json:"language"`
}

type Language struct {
	language.Tag
}

// MarshalBSON ...
func (l *Language) MarshalBSON() ([]byte, error) {
	var data langData
	data.Language = l.String()

	return bson.Marshal(data)
}

// UnmarshalBSON ...
func (l *Language) UnmarshalBSON(b []byte) error {
	var data langData

	if err := bson.Unmarshal(b, &data); err != nil {
		return err
	}

	l.Tag = language.Make(data.Language)
	return nil
}
