package Handler

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func NotificationsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodDelete:
		NotificationsDeleteHandler(w, r)
		break

	case http.MethodPost:
		NotificationsPostHandler(w, r)
		break

	case http.MethodGet:
		NotificationsGetHandler(w, r)
		break

	default:
		http.Error(w, "ERROR: invalid request method. Endpoint has supported methods GET, POST and DELETE", http.StatusMethodNotAllowed)
		return
	}

}

/*
	func NotificationsDeleteHandler(...)

	Given a specific webhook through ID in URL, will attempt to delete the webhook from
	current local storage of webhooks if one of specified ID is found - else, send error.

	REQUEST:
		HTTP METHOD: DELETE
		Path: /energy/v1/notifications/(specified ID for deletion)
*/

func NotificationsDeleteHandler(w http.ResponseWriter, r *http.Request) {
	splitURL := strings.Split(r.URL.Path, "/")
	returnData := returnWebhooks()

	// If URL is wrong or if the supposed id is empty, checks for that
	if len(splitURL) != 5 || splitURL[4] == "" {
		http.Error(w, "Error; Incorrect usage of URL.", http.StatusBadRequest)
		log.Println("EXTERNAL ERROR: Attempted to delete webhook, but failed due to improper usage of URL.")
		return
	}
	// Also checks if the length of our registered webhooks is over 0...
	if len(returnData) <= 0 {
		http.Error(w, "Error; No webhooks registered at all.", http.StatusNotFound)
		log.Println("EXTERNAL ERROR: Attempted to find webhooks, but there are none registered at all.")
		return
	}

	//Uses deleteDocument to try to delete the webhook document via the
	ctx := context.Background()
	err := deleteDocument(w, ctx, splitURL[4])
	if err == nil {

		// As well as checks if the ID a user has given this handler is valid
		index := locateWebhookByID(splitURL[4], returnData)
		if index == -1 {
			http.Error(w, "Error, no webhook of such ID found.", http.StatusNotFound)
			log.Println("EXTERNAL ERROR: Attempted to delete webhook, but failed due to nonexisting id.")
			return
		}

		// If everything is valid, creates HTML text which explains the contents of
		// the webhook about to be deleted
		w.Header().Add("content-type", "text/html")
		output := "You are now deleting the webhook with the following information: <br>"
		output += "Identification: " + returnData[index].ID + "<br>"
		output += "URL: " + returnData[index].URL + "<br>"
		output += "In which it was focused to look at the country of" + tempWebhooks[index].ISO + ", and report a " +
			"notification every " + string(returnData[index].Calls) + " invocations."
		_, err := fmt.Fprintf(w, "%v", output)
		if err != nil {
			http.Error(w, "Error when returning output,", http.StatusInternalServerError)
		}
	}

	// Sends a "no content" status code to specify that there is no content
	// of specified ID left, thus a successful deletion
	w.WriteHeader(http.StatusNoContent)

}

/*
func NotificationsPostHandler(...)

Given a POST request with body of valid information, will create a webhook
with specifications of said information - else, will send errors.
*/
func NotificationsPostHandler(w http.ResponseWriter, r *http.Request) {
	// Converts temporal ID of numbers to a string, gets removed later with addition of proper ID
	// Creates a new webhook object with an already assigned ID automatically generated
	temporaryRetrieval := WebhookObject{}
	// Decodes the information from the body of request to object
	err := json.NewDecoder(r.Body).Decode(&temporaryRetrieval)
	if err != nil {
		// If problems, preforms error and stops
		http.Error(w, "Error in decoding POST request.", http.StatusBadRequest)
		log.Println("EXTERNAL ERROR: Cannot decode POST request.")
		return
	}
	if temporaryRetrieval.Calls <= 0 || temporaryRetrieval.URL == "" {
		http.Error(w, "Error; Invalid input.", http.StatusBadRequest)
		log.Println("EXTERNAL ERROR: User attempted to create webhook with unacceptable values, stopped registration.", http.StatusBadRequest)
		return
	}

	if temporaryRetrieval.ISO != "" {
		test, _, _ := countrySearch(temporaryRetrieval.ISO)
		if test.Name == "" {
			http.Error(w, "Error; Invalid isocode registered to no country. HINT: Have you written it properly? ", http.StatusBadRequest)
			log.Println("EXTERNAL ERROR: User attempted to create webhook with unacceptable ISOCODE, stopped registration.", http.StatusBadRequest)
			return
		}
	}

	hash := hmac.New(sha256.New, []byte{1, 3, 1, 1})
	_, err = hash.Write([]byte(temporaryRetrieval.URL))
	if err != nil {
		http.Error(w, "Error; Cannot generate HASH from URL.", http.StatusInternalServerError)
		log.Println("INTERNAL ERROR: Cannot generate hash. Look into problem.")
	}

	temporaryRetrieval.ID = hex.EncodeToString(hash.Sum(nil))

	// Adds to content-type and encoder to JSON
	w.Header().Add("content-type", "application/json")
	encoder := json.NewEncoder(w)
	// Creates a new webhook with automatically generated ID
	err = encoder.Encode(WebhookObject{ID: temporaryRetrieval.ID})
	if err != nil {
		// Checks related errors
		log.Println("Attempted to return JSON of ID, failed registration.", http.StatusBadRequest)
		fmt.Errorf("Error, has failed encoding of webhook ID, stopped registration.", err.Error())
		return
	}
	// If no errors, then append safely to firebase webhook
	addDocument(w, temporaryRetrieval)
	tempWebhooks = append(tempWebhooks, temporaryRetrieval)
	w.WriteHeader(http.StatusCreated)
	log.Println("Has successfully registered webhook to storage.")

}

func NotificationsGetHandler(w http.ResponseWriter, r *http.Request) {
	splitURL := strings.Split(r.URL.Path, "/")
	returnData := returnWebhooks()

	if len(splitURL) != 5 {
		http.Error(w, "Error; Incorrect usage of URL.", http.StatusBadRequest)
		log.Println("Attempted to get webhook, but failed due to improper usage of URL.")
		return
	}
	if len(returnData) <= 0 {
		http.Error(w, "Error; No webhooks registered at all.", http.StatusNotFound)
		log.Println("Attempted to find webhooks, but there are none registered at all.")
		return
	}

	var temp []WebhookObject
	for _, v := range returnData {
		obj := WebhookObject{ID: v.ID, URL: v.URL, Calls: v.Calls, ISO: v.ISO, Invocations: v.Invocations}
		temp = append(temp, obj)
	}

	if splitURL[4] == "" {
		w.Header().Add("content-type", "application/json")
		encoder := json.NewEncoder(w)
		err := encoder.Encode(returnData)
		if err != nil {
			log.Println("Attempted to return JSON of all registered webhooks, failed.")
			fmt.Errorf("Error, has failed encoding of all webhooks registered", err.Error())
		}
		log.Println("Has successfully returned all webhooks in GET request.")
		return
	}

	index := locateWebhookByID(splitURL[4], returnData)
	if index != -1 {
		selectedWebhook := returnData[index]
		w.Header().Add("content-type", "application/json")
		encoder := json.NewEncoder(w)
		err := encoder.Encode(selectedWebhook)
		if err != nil {
			log.Println("Attempted to return JSON of selected singular webhook, failed.")
			fmt.Errorf("Error, has failed encoding of singular webhook", err.Error())
			return
		}
		log.Println("Has successfully returned singular webhook in GET request.")
	} else {
		http.Error(w, "Error; No webhooks of such ID registered at all.", http.StatusNotFound)
		log.Println("Attempted to find webhooks with ID, but there are none registered at all with such.")
		return
	}
}

// Function which attempts to find a webhook object based on the ID;
// if it finds it, then it returns the index of the webhook in a local storage
// NB: This will probably change as we change it to Firebase, I just want the things
// sorted out locally for now
// If it doesn't find it, it returns a -1, impossible index
func locateWebhookByID(id string, data []WebhookObject) int {
	for i, v := range data {
		if v.ID == id {
			return i
		}
	}
	return -1
}

/*
Function to invoke webhooks
*/
func invocationCall(w http.ResponseWriter, webhook WebhookObject, countryName string) {
	// Generates an object with the information we want to send
	response := WebhookObject{ID: webhook.ID, Calls: webhook.Invocations, ISO: countryName}
	// Marshals it in order to format it to a sendable format (in []byte later)
	stringified, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error with marshalling content.", http.StatusBadRequest)
		log.Println("Invocation Call Problem, Marhshalling.")
		return
	}
	// Attempts to create a new HTTP POST to the webhook's URL, with JSON we want to send
	res, err := http.NewRequest("POST", webhook.URL, bytes.NewReader([]byte(stringified)))
	if err != nil {
		http.Error(w, "Error in creating new HTTP Request.", http.StatusBadRequest)
		log.Println("Error in creating new HTTP Request")
		return
	}
	// Initializes client, and preforms the request
	client := http.Client{}
	_, err = client.Do(res)
	if err != nil {
		http.Error(w, "Error; Invocation Call Problem, HTTP Request.", http.StatusBadRequest)
		log.Println("Invocation Call Problem")
		return
	}
	log.Println("Success: Invocation Notification to", webhook.URL, "has been achieved.")
}
