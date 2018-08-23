package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
	"zdrav/models"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var sickListCol *mgo.Collection
var usersCol *mgo.Collection

const (
	port = ":9087"
)

func sickListHandler(w http.ResponseWriter, r *http.Request) {
	user := checkSession(w, r)
	if user.Username == "" {
		http.Redirect(w, r, "/login", 302)
		return
	}
	fmt.Println("Method: ", r.Method)
	if r.Method == "GET" { //================GET===============
		user := checkSession(w, r)
		t, _ := template.ParseFiles("templates/addSickList.html", "templates/header.html", "templates/footer.html")
		data := []models.SickList{}
		if user.Lpu == "" {
			sickListCol.Find(bson.M{}).All(&data)
		} else {
			sickListCol.Find(bson.M{"lpu": user.Lpu}).All(&data)
		}
		fmt.Println("lpu: ", user.Lpu)
		t.ExecuteTemplate(w, "sickLists", data)
	} else { //==================POST===============
		r.ParseForm()
		if r.FormValue("sickList") != "" {
			var sickList models.SickList
			sickList.Id = bson.NewObjectId()
			sickList.SickList = r.FormValue("sickList")
			sickList.FirstName = r.FormValue("firstName")
			sickList.LastName = r.FormValue("lastName")
			sickList.MiddleName = r.FormValue("middleName")
			sickList.Lpu = user.Lpu
			sickList.Snils = r.FormValue("snils")
			sickList.Stazh = r.FormValue("stazh")
			sickList.Date = time.Now()
			addToCol(sickList, sickListCol)
		}
		http.Redirect(w, r, "/add", 302)
	}
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	user := checkSession(w, r)
	if user.Username == "" {
		http.Redirect(w, r, "/login", 302)
		return
	}
	if r.Method == "GET" { //================GET===============
		id := r.FormValue("id")
		fmt.Println(id)
		t, _ := template.ParseFiles("templates/editList.html", "templates/header.html", "templates/footer.html")
		data := []models.SickList{}
		if user.Lpu == "" {
			sickListCol.FindId(bson.ObjectIdHex(id)).All(&data)
		} else {
			sickListCol.Find(bson.M{"_id": bson.ObjectIdHex(id), "lpu": user.Lpu}).All(&data)
		}
		t.ExecuteTemplate(w, "editSickList", data)
	} else { //==================POST===============
		r.ParseForm()
		fmt.Println((r.FormValue("id")))
		sickListCol.Update(bson.M{"_id": bson.ObjectIdHex(r.FormValue("id"))},
			bson.M{"$set": bson.M{"firstName": r.FormValue("firstName"),
				"lastName": r.FormValue("lastName"), "middleName": r.FormValue("middleName"),
				"sickList": r.FormValue("sickList"), "stazh": r.FormValue("stazh"), "snils": r.FormValue("snils")}})
		http.Redirect(w, r, "/add", 302)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	user := checkSession(w, r)
	if user.Username != "" {
		http.Redirect(w, r, "/add", 302)
		return
	}
	if r.Method == "GET" {
		t, _ := template.ParseFiles("templates/login.html", "templates/header.html", "templates/footer.html")
		t.ExecuteTemplate(w, "login", nil)
	} else {
		r.ParseForm()
		username := r.FormValue("username")
		password := r.FormValue("password")
		var user models.User
		usersCol.Find(bson.M{"username": username, "password": password}).One(&user)
		if string(user.Id) == "" {
			fmt.Println("User: ", username, " not found")
		} else {
			generateSession(w, user)
			http.Redirect(w, r, "/add", 302)
		}
		http.Redirect(w, r, "/login", 302)
	}
}

func addUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Method add user: ", r.Method)
	user := checkSession(w, r)
	if user.Username != "admin" {
		http.Redirect(w, r, "/add", 302)
		return
	}
	if r.Method == "GET" { //================GET===============
		t, _ := template.ParseFiles("templates/adduser.html", "templates/header.html", "templates/footer.html")
		t.ExecuteTemplate(w, "adduser", nil)
	} else { //==================POST===============
		r.ParseForm()
		var user models.User
		user.Id = bson.NewObjectId()
		user.Username = r.FormValue("username")
		user.Password = r.FormValue("password")
		user.Lpu = r.FormValue("lpu")
		addToCol(user, usersCol)
		http.Redirect(w, r, "/adduser", 302)
		fmt.Println("  ")
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	user := checkSession(w, r)
	fmt.Println(user)
	usersCol.Update(bson.M{"_id": user.Id}, bson.M{"$set": bson.M{"sessionId": ""}})
	http.Redirect(w, r, "/login", 302)
}

//=============================MAIN=========================
func main() { //===================DB SETUP=================
	session, err := mgo.Dial("mongodb://127.0.0.1")
	if err != nil {
		fmt.Println(err)
	}
	defer session.Close()
	sickListCol = session.DB("localdb").C("sickLists")
	usersCol = session.DB("localdb").C("users")
	//==================Listener=====================
	http.Handle("/layout/", http.StripPrefix("/layout/", http.FileServer(http.Dir("templates/layout"))))
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/add", sickListHandler)
	//http.HandleFunc("/adminpart", adminPart)
	http.HandleFunc("/adduser", addUser)
	http.HandleFunc("/edit", editHandler)
	http.HandleFunc("/logout", logoutHandler)
	err1 := http.ListenAndServe(port, nil)
	if err1 != nil {
		log.Fatal("ListenAndServe: ", err1)
	}

}
