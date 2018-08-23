package main

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"time"
	"zdrav/models"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func GenerateId() string {
	b := make([]byte, 16)
	rand.Read(b)
	fmt.Println("b: ", b)
	return fmt.Sprintf("%x", b)
}

func addToCol(data interface{}, col *mgo.Collection) {
	err := col.Insert(data)
	if err != nil {
		fmt.Println(err)
	}
}

func generateSession(w http.ResponseWriter, user models.User) {
	sessionId := GenerateId()
	cookie := &http.Cookie{
		Name:    "sessionId",
		Value:   sessionId,
		Expires: time.Now().Add(120 * time.Hour),
	}
	http.SetCookie(w, cookie)
	usersCol.Update(bson.M{"_id": (user.Id)}, bson.M{"$set": bson.M{"sessionId": sessionId}})
}

func checkSession(w http.ResponseWriter, r *http.Request) models.User {
	var user models.User
	cookieId, err := r.Cookie("sessionId")
	if cookieId == nil || err != nil {
		return user
	} else {
		usersCol.Find(bson.M{"sessionId": cookieId.Value}).One(&user)
	}
	return user
}
