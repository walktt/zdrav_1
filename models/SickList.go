package models

import (
	"gopkg.in/mgo.v2/bson"
)

type SLForEdit struct {
	SLs     []SickList
	UserLpu string
}

type SickList struct {
	Id         bson.ObjectId `bson:"_id"`
	SickList   string        `bson:"sickList"`
	FirstName  string        `bson:"firstName"`
	LastName   string        `bson:"lastName"`
	MiddleName string        `bson:"middleName"`
	Lpu        string        `bson:"lpu"`
	Date       string        `bson:"time"`
	Snils      string        `bson:"snils"`
	Stazh      string        `bson:"stazh"`
	Pass       bool          `bson:"pass"`
}
