package Handler

import (
	"cloud.google.com/go/firestore" // Firestore-specific support
	"context"                       // State handling across API boundaries; part of native GoLang API
	"errors"
	"firebase.google.com/go"
	"fmt"
	"google.golang.org/api/option"
	"log"
	"net/http"
)

// Variables used throughout the program
var ctx context.Context
var client *firestore.Client

/*
Reads string from the body and sends it to Firestore so it can be registered as a document
*/
func addDocument(w http.ResponseWriter, webhookInfo WebhookObject) {
	// Add element in embedded structure. Adds the info by using the WebhookObject struct
	id, _, err := client.Collection(COLLECTION).Add(ctx,
		map[string]interface{}{
			"webhook_id":  webhookInfo.ID,
			"url":         webhookInfo.URL,
			"country":     webhookInfo.ISO,
			"invocations": webhookInfo.Invocations,
			"calls":       webhookInfo.Calls,
		})

	if err != nil {
		// Error handling prints to the terminal and postman console
		log.Println("Error when adding document " + webhookInfo.ID + ", Error: " + err.Error())
		http.Error(w, "Error when adding document "+webhookInfo.ID+", Error: "+err.Error(), http.StatusBadRequest)
		return
	} else {
		// Returns document ID in body
		log.Println("Document added to collection. Identifier of returned document: " + id.ID)
		http.Error(w, id.ID, http.StatusCreated)
		return
	}
}

/*
Returns all the documents as well as the information in the documents
*/
func returnWebhooks() []WebhookObject {
	ctx = context.Background()
	client, err := GetFirebaseClient(ctx)
	if err != nil {
		log.Fatal("Error getting Firebase client:", err)
	}
	if client == nil {
		log.Fatal("Firebase client is nil")
	}

	//Get the Firestore client so we can interact with the database
	fsClient, err := client.Firestore(ctx)
	if err != nil {
		log.Println("Error getting Firestore client: ", err.Error())
	}

	// Collective retrieval of messages
	collection := fsClient.Collection(COLLECTION)           // Loop through collection "webhooks"
	allDocuments, err := collection.Documents(ctx).GetAll() //Loops through all entries in collection
	if err != nil {
		fmt.Println("Error with collection")
	}

	var tempInfo []WebhookObject

	//For-loop that runs through all the entries in collection webhooks
	for _, webhook := range allDocuments {
		//Temp struct
		var tempWebhook WebhookObject
		data := webhook.Data()
		tempWebhook.ID = data["webhook_id"].(string)
		tempWebhook.URL = data["url"].(string)
		tempWebhook.ISO = data["country"].(string)
		tempWebhook.Calls = int(data["calls"].(int64))
		tempWebhook.Invocations = int(data["invocations"].(int64))

		//Adds the webhook info to the tempInfo
		tempInfo = append(tempInfo, tempWebhook)
	}

	return tempInfo

}

/*
Function that returns firebase app client by taking in a context.Context
parameter and returns a tuple to the pointer to firebase.App and error object
*/
func GetFirebaseClient(ctx context.Context) (*firebase.App, error) {
	// Initialize a Firebase app using a service account file
	sa := option.WithCredentialsFile("group12-assignment2-sa.json")

	//Creates a new Firebase app instance with the given information and the service account
	app, err := firebase.NewApp(ctx, nil, sa)
	//Error handling
	if err != nil {
		log.Println("Error initializing app: ", err.Error())
	}

	return app, nil
}

/*
Deletes a document and it's entries in the Firebase Database by
going through the Firestore collection and using the keyID
*/
func deleteDocument(w http.ResponseWriter, ctx context.Context, keyID string) error {
	//Get Firebase client using the service account provided from GetFirebaseClient
	client, err := GetFirebaseClient(ctx)
	if err != nil {
		log.Println("Error getting the Firebase client: ", err.Error())
	}

	//Get the Firestore client so we can interact with the database
	fsClient, err := client.Firestore(ctx)
	if err != nil {
		log.Println("Error getting Firestore client: ", err.Error())
	}
	//Finds webhook with key sent from parameter
	webhookKeyID := fsClient.Collection(COLLECTION).Where("webhook_id", "==", keyID)

	//Gets data from the query above
	data, err := webhookKeyID.Documents(ctx).GetAll()
	if err != nil {
		log.Println("Error getting the data", err.Error())
	}

	//Checking if there are any matches
	if len(data) == 0 {
		http.Error(w, "No webhooks found with this key,", http.StatusBadRequest)
		return errors.New("No webhooks found with this key" + keyID)
	}

	//Checking if there are more than one match
	if len(data) > 1 {
		http.Error(w, "Multiple webhooks found with this key", http.StatusBadRequest)
		return errors.New("Multiple webhooks found with this key" + keyID)
	}

	//Finds and deletes the document with the given keyID from the chosen Firestore collection
	_, err = fsClient.Collection(COLLECTION).Doc(data[0].Ref.ID).Delete(ctx)
	if err != nil {
		log.Println("Error deleting document", err.Error())
	}

	return nil
}

/*
Function that updates the number of invocations on the document in the Firebase-collection
*/
func invocationUpdate(w http.ResponseWriter, webhook WebhookObject) {
	//Get Firebase client using the service account provided from GetFirebaseClient
	client, err := GetFirebaseClient(ctx)
	if err != nil {
		log.Println("Error getting the Firebase client: ", err.Error())
	}

	//Get the Firestore client so we can interact with the database
	fsClient, err := client.Firestore(ctx)
	if err != nil {
		log.Println("Error getting Firestore client: ", err.Error())
	}
	//Finds webhook with key sent from parameter
	webhookKeyID := fsClient.Collection(COLLECTION).Where("webhook_id", "==", webhook.ID)

	//Gets data from the query above
	data, err := webhookKeyID.Documents(ctx).GetAll()
	if err != nil {
		log.Println("Error getting the data", err.Error())
	}
	//Saves the data reference from the data variable above in a new variable
	documentRef := data[0].Ref

	//Updates the invocations value in the document by using a Set-function
	_, err = documentRef.Set(ctx,
		map[string]interface{}{
			"invocations": webhook.Invocations + 1,
			//MergeAll-function that overwrites the info in the document with the updated one
		}, firestore.MergeAll)
	//Error handling
	if err != nil {
		log.Println("Error updating document: ", err.Error())
		http.Error(w, "Error updating document: ", http.StatusInternalServerError)
	}
}

/*
Firebase main connects the program to the firebase database
*/
func Firebasemain() {
	// Initialize a Firebase app using a service account file
	ctx = context.Background()
	sa := option.WithCredentialsFile("group12-assignment2-sa.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	// Check if app is nil before using it
	if app == nil {
		log.Fatalf("app is nil")
	}

	// Initialize a Firestore client
	client, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalf("error initializing firestore client: %v\n", err)
	}

	// Check if client is nil before using it
	if client == nil {
		log.Fatalf("client is nil")
	}
}
