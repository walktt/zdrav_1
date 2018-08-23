package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type SickList struct {
	Id         bson.ObjectId `bson:"_id"`
	SickList   string        `bson:"sickList"`
	FirstName  string        `bson:"firstName"`
	LastName   string        `bson:"lastName"`
	MiddleName string        `bson:"middleName"`
	Lpu        string        `bson:"lpu"`
	Date       time.Time     `bson:"time"`
	Snils      string        `bson:"snils"`
	Stazh      string        `bson:"stazh"`
}
