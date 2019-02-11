package main

import (
	"encoding/json"
	"log"
	"net/http"

	//"io/ioutil"
	"fmt"

	"gopkg.in/mgo.v2/bson"

	// "go/scanner"
	"github.com/gorilla/mux"
	. "github.com/mlabouardy/movies-restapi/models"
	"golang.org/x/crypto/bcrypt"
	mgo "gopkg.in/mgo.v2"
)

var db *mgo.Database

const (
	COLLECTION = "test"
)

type Movie struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	Name        string        `bson:"name" json:"name"`
	Password    string        `bson:"password" json:"password"`
	Description string        `bson:"description" json:"description"`
}
type forgot struct {
	Name string `bson:"name" json:"name"`
}

func signup(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var movie Movie
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	//	nm := movie.Name
	if err := db.C(COLLECTION).Find(bson.M{"name": movie.Name}).One(&movie); err == nil {
		fmt.Print(err)

		respondWithError(w, http.StatusInternalServerError, " Exists")
		return
	}
	pass := movie.Password
	pwd := []byte(pass)
	hash, err := bcrypt.GenerateFromPassword(pwd, 10)
	if err != nil {
		log.Println(err)
	}
	movie.Password = string(hash)
	fmt.Println(pass)
	// if err := db.C(COLLECTION).Find(nm); err != nil {
	// 	fmt.Print(err)

	// 	respondWithError(w, http.StatusInternalServerError, "Already Exists")
	// 	return
	// }
	movie.ID = bson.NewObjectId()

	if err := db.C(COLLECTION).Insert(&movie); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, movie)
}
func login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var movie Movie
	redirectTarget := "/"
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	pass := movie.Password
	err := db.C(COLLECTION).Find(bson.M{"name": movie.Name}).One(&movie)
	if err != nil {
		log.Fatal(err)
	}
	//err := vr.Password
	fmt.Println(movie.Password, "pass")
	err = bcrypt.CompareHashAndPassword([]byte(movie.Password), []byte(pass))
	if err != nil {
		log.Fatal(err)
	} else {
		redirectTarget = "https://dzone.com/articles/build-restful-api-in-go-and-mongodb"
	}
	// fmt.Printf("hash: %v\n", hash)

	//movie.ID = bson.NewObjectId()
	//	if err := db.C(COLLECTION).Insert(&movie); err != nil {
	//	respondWithError(w, http.StatusInternalServerError, err.Error())
	//	return
	//}
	//respondWithJson(w, http.StatusCreated, movie)
	http.Redirect(w, r, redirectTarget, 302)
}
func Forgot(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var username forgot
	//redirectTarget := "/"
	if err := json.NewDecoder(r.Body).Decode(&username); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	//nam := username.Name
	err := db.C(COLLECTION).Find(bson.M{"name": username.Name}).One(&username)
	if err != nil {
		log.Fatal(err)
	}
	//err := vr.Password
	fmt.Println(username.Name, "nam")

	// err = bcrypt.CompareHashAndPassword([]byte(movie.Password), []byte(pass))
	// if err != nil {
	// 	log.Fatal(err)
	// } else {
	// 	redirectTarget = "https://dzone.com/articles/build-restful-api-in-go-and-mongodb"
	// }
	// // fmt.Printf("hash: %v\n", hash)

	// //movie.ID = bson.NewObjectId()
	// //	if err := db.C(COLLECTION).Insert(&movie); err != nil {
	// //	respondWithError(w, http.StatusInternalServerError, err.Error())
	// //	return
	// //}
	////respondWithJson(w, http.StatusCreated, username)
	// http.Redirect(w, r, redirectTarget, 302)
}
func init() {
	//config.Read()

	Server := "localhost"
	Database := "movies_db"
	session, err := mgo.Dial(Server)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(Database)
}

// Define HTTP request routes
func main() {

	r := mux.NewRouter()
	r.HandleFunc("/movies", login).Methods("GET")
	r.HandleFunc("/movies", signup).Methods("POST")
	r.HandleFunc("/movies", Forgot).Methods("VIEW")
	//r.HandleFunc("/movies", DeleteMovieEndPoint).Methods("DELETE")
	//r.HandleFunc("/movies/{name}", FindMovieEndpoint).Methods("GET")
	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal(err)  
	}
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

//fmt.Printf("%s\n", string(pwd))
