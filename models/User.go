package models

import "gopkg.in/mgo.v2/bson"

type User struct {
	Id        bson.ObjectId `bson:"_id"`
	Username  string        `bson:"username"`
	Password  string        `bson:"password"`
	Lpu       string        `bson:"lpu"`
	SessionId string        `bson:"sessionId"`
}
