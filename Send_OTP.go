package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func random(min int, max int) int {
	return rand.Intn(max-min) + min
}

func main() {
	// Set account keys & information
	accountSid := "AC380ea09dce91c893ff0890bf74a32451"
	authToken := "64495615ce8f849021307b7731dac4c7"
	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Messages.json"

	// Create possible message bodies
	// quotes := [7]string{"I urge you to please notice when you are happy, and exclaim or murmur or think at some point, 'If this isn't nice, I don't know what is.'",
	// 	"Peculiar travel suggestions are dancing lessons from God.",
	// 	"There's only one rule that I know of, babies—God damn it, you've got to be kind.",
	// 	"Many people need desperately to receive this message: 'I feel and think much as you do, care about many of the things you care about, although most people do not care about them. You are not alone.'",
	// 	"That is my principal objection to life, I think: It's too easy, when alive, to make perfectly horrible mistakes.",
	// 	"So it goes.",
	// 	"We must be careful about what we pretend to be."}

	//quotes := rand.Intn(9999)
	rand.Seed(time.Now().UnixNano())
	randomNum := strconv.Itoa(random(1000, 2000))

	// Set up rand

	// Pack up the data for our message
	msgData := url.Values{}
	msgData.Set("To", "+919582712685")
	msgData.Set("From", "+15623569987")
	msgData.Set("Body", randomNum+"is your otp that you have to varyfied") //[rand.Intn(len(quotes))]
	msgDataReader := *strings.NewReader(msgData.Encode())
	fmt.Print(msgDataReader)
	// Create HTTP request client
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
