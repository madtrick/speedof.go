package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"syscall"

	"github.com/dghubble/oauth1"
	"github.com/dghubble/oauth1/twitter"
)

type Configuration struct {
	Token       string
	TokenSecret string
}

func main() {
	_, err := os.Open("./config.json")

	if err != nil {
		if e, ok := err.(*os.PathError); ok && e.Err == syscall.ENOENT {
			var oauthConfig oauth1.Config

			CONSUMER_KEY := os.Getenv("CONSUMER_KEY")
			CONSUMER_SECRET := os.Getenv("CONSUMER_SECRET")

			if CONSUMER_KEY == "" {
				log.Fatal("Empty CONSUMER_KEY env variable")
			}

			if CONSUMER_SECRET == "" {
				log.Fatal("Empty CONSUMER_SECRET env variable")
			}

			oauthConfig = oauth1.Config{
				ConsumerKey:    CONSUMER_KEY,
				ConsumerSecret: CONSUMER_SECRET,
				CallbackURL:    "oob",
				Endpoint:       twitter.AuthorizeEndpoint,
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

			data, _ := json.Marshal(configuration)
			ioutil.WriteFile("config.json", data, 0600)

			fmt.Println("vim-go")
		} else {
			log.Fatal(err)
		}
	}
}
