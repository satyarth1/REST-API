package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"

	. "newapp/config"
	. "newapp/dao"
	. "newapp/models"

	"github.com/dgrijalva/jwt-go"
)

var config = Config{}
var dao = UserDAO{}
var Num1 string
var num2 int

// POST a new user
func random(min int, max int) int {
	return rand.Intn(max-min) + min
}

func SendSms(PhoneNum int64) {
	userno := strconv.FormatInt(PhoneNum, 10)
	accountSid := "AC380ea09dce91c893ff0890bf74a32451"
	authToken := "64495615ce8f849021307b7731dac4c7"
	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Messages.json"

	rand.Seed(time.Now().UnixNano())
	randomNum := strconv.Itoa(random(1000, 2000))
	Num1 = randomNum
	msgData := url.Values{}
	msgData.Set("To", "+91"+userno) //9582712685
	msgData.Set("From", "+15623569987")
	msgData.Set("Body", randomNum+"  is your otp that you have to varyfied") //[rand.Intn(len(quotes))]
	msgDataReader := *strings.NewReader(msgData.Encode())
	fmt.Print(msgDataReader)
	client := &http.Client{}
	req, _ := http.NewRequest("POST", urlStr, &msgDataReader)
	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Make HTTP POST request and return message SID
	resp, _ := client.Do(req)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&data)
		if err == nil {
			fmt.Println(data["sid"])
		}
	} else {
		fmt.Println(resp.Status)
	}

}

func signup(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	pass := user.Password
	pwd := []byte(pass)
	hash, err := bcrypt.GenerateFromPassword(pwd, 10)
	if err != nil {
		log.Println(err)
	}
	user.Password = string(hash)
	fmt.Println(pass)
	user.ID = bson.NewObjectId()
	phonenum := user.PhoneNumber

	SendSms(phonenum)
	result, err := dao.Find(user)
	fmt.Println(result.PhoneNumber)
	fmt.Println(user.PhoneNumber)
	if err != nil {
		//log.Fatal(err)
		fmt.Println(err)
	}
	//num2 = 0
	// if user.PhoneNumber == result.PhoneNumber {
	// 	num2 = 1
	// 	fmt.Fprintf(w, "USER PHONE NUMBER ALLREADY EXIST")
	// }
	// if num2 == 0 {
	// 	if err := dao.Insert(user); err != nil {
	// 		respondWithError(w, http.StatusInternalServerError, err.Error())
	// 		return

	// 	}
	// 	respondWithJson(w, http.StatusCreated, user)

	// }
	if _, err := dao.Find(user); err == nil {
		respondWithError(w, http.StatusInternalServerError, "phone no aready exist")
	} else {
		if err := dao.Insert(user); err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return

		}
		respondWithJson(w, http.StatusCreated, user)
	}
}

func verifyOtp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	randotp := Num1
	var otp Otp
	if err := json.NewDecoder(r.Body).Decode(&otp); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	newotp := otp.UserOtp
	if randotp != newotp {
		respondWithError(w, http.StatusBadRequest, "INVALID OTP")
	} else {
		respondWithJson(w, http.StatusOK, "OTP VERFYED")
	}

}

func login(w http.ResponseWriter, r *http.Request) {
	//fmt.Println(r.URL.Query())
	defer r.Body.Close()
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	//fmt.Println(user)
	pass := user.Password
	result, err := dao.Find(user)
	//fmt.Println(user)
	if err != nil {
		//log.Fatal(err)
		fmt.Println(err)
	}
	//err := vr.Password
	//fmt.Println(result, "pass")
	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(pass))
	if err != nil {
		//log.Fatal(err)
		fmt.Fprintf(w, "invalid phonenumber or password")
	} else {

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"phone_number": user.PhoneNumber,
			"password":     user.Password,
		})
		fmt.Println(token)
		tokenString, error := token.SignedString([]byte("rajat"))
		if error != nil {
			fmt.Println(error)
		}
		//json.NewEncoder(w).Encode(JwtToken{Token: tokenString})
		fmt.Println("Token:", tokenString)
		fmt.Fprintf(w, "message:login successful")
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

// Parse the configuration file 'config.toml', and establish a connection to DB
func init() {
	config.Read()

	dao.Server = config.Server
	dao.Database = config.Database
	dao.Connect()
}

// Define HTTP request routes
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/login", login).Methods("POST")
	r.HandleFunc("/signup", signup).Methods("POST")
	r.HandleFunc("/verifyOtp", verifyOtp).Methods("POST")
	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal(err)
	}
}
