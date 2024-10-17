package models

import "go.mongodb.org/mongo-driver/v2/bson"

type LoginDetails struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type User struct {
	Username string        `json:"username,omitempty"`
	Password string        `json:"password,omitempty"`
	ID       bson.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
}

type NewURL struct {
	URL    string        `json:"url,omitempty"`
	Author bson.ObjectID `json:"author,omitempty"`
}
