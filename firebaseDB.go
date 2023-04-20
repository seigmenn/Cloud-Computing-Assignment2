package main

import (
	"cloud.google.com/go/firestore"   // Firestore-specific support
	"context"                         // State handling across API boundaries; part of native GoLang API
	firebase "firebase.google.com/go" // Generic firebase support
	"google.golang.org/api/option"
	"log"
)

var ctx context.Context
var client *firestore.Client

func firebasemain() {
	// Use a service account
	sa := option.WithCredentialsFile("./group12-service-account.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()
}

/***

func firebase() {
	// Firebase initialisation
	ctx = context.Background()

	sa := option.WithCredentialsFile("./group12-service-account.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	//Instantiate client
	client, err = app.Firestore(ctx)

	//Error checking for this
	if err != nil {
		log.Fatalln(err)
	}

	//When function ends, close down the client
	defer func() {
		err := client.Close()
		if err != nil {
			log.Fatal("Closing of firebase client failed. Error: ", err)
		}

		//Heroku-compaible
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}

		addr := ":" + port

		log.Printf("Firestore REST service listening on %s ...\n", addr)
		if err := http.ListenAndServe(addr, nil); err != nil {
			panic(err)
		}
	}

}
***/
