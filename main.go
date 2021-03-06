package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"syscall"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	twitterOauth1 "github.com/dghubble/oauth1/twitter"
)

func initStream(tokens *Configuration, consumerKey, consumerSecret, track string) {
	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(tokens.Token, tokens.TokenSecret)
	httpClient := config.Client(context.TODO(), token)
	client := twitter.NewClient(httpClient)

	params := &twitter.StreamFilterParams{
		Track:         []string{track},
		StallWarnings: twitter.Bool(true),
	}

	count := 0

	go meter(track, &count)

	stream, err := client.Streams.Filter(params)

	if err != nil {
		log.Fatal(err)
	}

	for range stream.Messages {
		count = count + 1
	}
}

type Configuration struct {
	Token       string
	TokenSecret string
}

func meter(track string, count *int) {
	for {
		fmt.Printf("\r %d %ss/s", *count, track)
		*count = 0
		time.Sleep(1 * time.Second)
	}
}

func getToken(consumerKey, consumerSecret string) *Configuration {
	var oauthConfig oauth1.Config

	oauthConfig = oauth1.Config{
		ConsumerKey:    consumerKey,
		ConsumerSecret: consumerSecret,
		CallbackURL:    "oob",
		Endpoint:       twitterOauth1.AuthorizeEndpoint,
	}

	var requestToken string
	// var requestSecret string
	// var err error

	requestToken, _, _ = oauthConfig.RequestToken()

	var authorizationUrl *url.URL

	authorizationUrl, _ = oauthConfig.AuthorizationURL(requestToken)

	fmt.Println(authorizationUrl)

	fmt.Printf("Paste your PIN here: ")
	var verifier string
	_, err := fmt.Scanf("%s", &verifier)
	if err != nil {
		log.Fatal(err)
	}

	accessToken, accessSecret, err := oauthConfig.AccessToken(requestToken, "LOL", verifier)
	token := oauth1.NewToken(accessToken, accessSecret)

	configuration := &Configuration{Token: token.Token, TokenSecret: token.TokenSecret}

	return configuration
}

func main() {

	CONSUMER_KEY := os.Getenv("CONSUMER_KEY")
	CONSUMER_SECRET := os.Getenv("CONSUMER_SECRET")
	TRACK := os.Getenv("TRACK")

	if CONSUMER_KEY == "" {
		log.Fatal("Empty CONSUMER_KEY env variable")
	}

	if CONSUMER_SECRET == "" {
		log.Fatal("Empty CONSUMER_SECRET env variable")
	}

	if TRACK == "" {
		log.Fatal("Empty TRACK env variable")
	}

	_, err := os.Open("./config.json")

	if err != nil {
		if e, ok := err.(*os.PathError); ok && e.Err == syscall.ENOENT {

			configuration := getToken(CONSUMER_KEY, CONSUMER_SECRET)

			data, _ := json.Marshal(configuration)
			ioutil.WriteFile("config.json", data, 0600)

			initStream(configuration, CONSUMER_KEY, CONSUMER_SECRET, TRACK)
		} else {
			log.Fatal(err)
		}
	} else {
		data, _ := ioutil.ReadFile("config.json")
		tokens := &Configuration{}
		json.Unmarshal(data, tokens)

		initStream(tokens, CONSUMER_KEY, CONSUMER_SECRET, TRACK)
	}
}
