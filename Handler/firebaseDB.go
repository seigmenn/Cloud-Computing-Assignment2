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

var ctx context.Context
var client *firestore.Client

/*
Reads string from the body and sends it to Firestore so it can be registered as a document
It's going to send 1 int and 3 strings.
*/
func addDocument(w http.ResponseWriter, r *http.Request, webhookInfo WebhookObject) {
	// Add element in embedded structure.
	// Note: this structure is defined by the client, not the server!; it exemplifies the use of a complex structure
	// and illustrates how you can use Firestore features such as Firestore timestamps.

	id, _, err := client.Collection(COLLECTION).Add(ctx,
		map[string]interface{}{
			"webhook_id":  webhookInfo.ID,
			"url":         webhookInfo.URL,
			"country":     webhookInfo.ISO,
			"invocations": webhookInfo.Invocations,
			"calls":       webhookInfo.Calls,
		})

	if err != nil {
		// Error handling
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
Returns all the documents as well as their info
*/
func returnWebhooks(w http.ResponseWriter, r *http.Request) []WebhookObject {
	// Collective retrieval of messages
	collection := client.Collection(COLLECTION)             // Loop through collection "webhooks"
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
		//Tries to add data to WebhookObject struct
		/*err = webhook.DataTo(&tempWebhook)
		//Error handling if it doesn't succeed
		if err != nil {
			fmt.Println("Error unmarshaling")
		}*/
		//Adds the webhook info to the tempInfo
		tempInfo = append(tempInfo, tempWebhook)
	}

	return tempInfo
}

/*
**
/*
Handler for all message-related operations

	func handleMessage(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			addDocument(w, r)
		case http.MethodGet:
			displayDocument(w, r)
		default:
			log.Println("Unsupported request method " + r.Method)
			http.Error(w, "Unsupported request method "+r.Method, http.StatusMethodNotAllowed)
			return
		}
	}

**
*/

/*
Function that returns firebase app client by taking in a context.Context
parameter and returns a tuple to the pointer to firebase.App and error object
*/

func GetFirebaseClient(ctx context.Context) (*firebase.App, error) {
	// Initialize a Firebase app using a service account file
	sa := option.WithCredentialsFile("group12-assignment2-sa.json")

	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
		return nil, err
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
