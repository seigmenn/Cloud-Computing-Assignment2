package Handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Note; I'll clean it up when the functionality is complete.

// Temporal id; will be changed later with a proper ID creation,
// but for now, this is a temporary ID which gets incremented
// for every new webhook to support multiple webhooks with unique IDs
// to demonstrate requests with specified ID
var tempId = 0

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

func NotificationsDeleteHandler(w http.ResponseWriter, r *http.Request) {
	splitURL := strings.Split(r.URL.Path, "/")
	// If URL is wrong or if the supposed id is empty, checks for that
	if len(splitURL) != 5 || splitURL[4] == "" {
		http.Error(w, "Error; Incorrect usage of URL.", http.StatusBadRequest)
		log.Println("Attempted to delete webhook, but failed due to improper usage of URL.")
		return
	}
	// Also checks if the length of our registered webhooks is over 0...
	if len(tempWebhooks) <= 0 {
		http.Error(w, "Error; No webhooks registered at all.", http.StatusNotFound)
		log.Println("Attempted to find webhooks, but there are none registered at all.")
		return
	}
	// As well as checks if the ID an user has given this handler is valid
	index := locateWebhookByID(splitURL[4])
	if index == -1 {
		http.Error(w, "Error, no webhook of such ID found.", http.StatusNotFound)
		log.Println("Attempted to delete webhook, but failed due to nonexisting id.")
		return
	}

	// If everything is valid, creates HTML text which explains the contents of
	// the webhook about to be deleted
	w.Header().Add("content-type", "text/html")
	output := "You are now deleting the webhook with the following information: <br>"
	output += "Identification: " + tempWebhooks[index].ID + "<br>"
	output += "URL: " + tempWebhooks[index].URL + "<br>"
	output += "In which it was focused to look at the country of" + tempWebhooks[index].ISO + ", and report a " +
		"notification every " + string(tempWebhooks[index].Calls) + " invocations."
	_, err := fmt.Fprintf(w, "%v", output)
	if err != nil {
		http.Error(w, "Error when returning output,", http.StatusInternalServerError)
	}

	// Proceeds to update local storage of webhooks to not contain this deleted webhook
	// replace with more efficient delete afterwards through append and slicing [:index], [index:] later
	var temp []WebhookObject
	for _, v := range tempWebhooks {
		if v.ID != tempWebhooks[index].ID {
			temp = append(temp, v)
		}
	}
	tempWebhooks = temp
	w.WriteHeader(http.StatusNoContent)

}

func NotificationsPostHandler(w http.ResponseWriter, r *http.Request) {
	// Converts temporal ID of numbers to a string, gets removed later with addition of proper ID
	temporal := strconv.Itoa(tempId)
	// Creates a new webhook object with an already assigned ID automatically generated
	newWebhook := WebhookObject{ID: temporal}
	// Decodes the information from the body of request to object
	err := json.NewDecoder(r.Body).Decode(&newWebhook)
	if err != nil {
		// If problems, preforms error and stops
		http.Error(w, "Error in decoding POST request", http.StatusBadRequest)
		return
	}
	test, _, _ := countrySearch(newWebhook.ISO)
	if test.Name == "" {
		http.Error(w, "Error; Invalid isocode registered to no country. HINT: Have you written it properly? ", http.StatusBadRequest)
		log.Println("User attempted to create webhook with unacceptable ISOCODE, stopped registration.", http.StatusBadRequest)
		return
	}

	// Adds to content-type and encoder to JSON
	w.Header().Add("content-type", "application/json")
	encoder := json.NewEncoder(w)
	// Creates a new webhook with automatically generated ID
	err = encoder.Encode(WebhookObject{ID: newWebhook.ID})
	if err != nil {
		// Checks related errors
		log.Println("Attempted to return JSON of ID, failed registration.", http.StatusBadRequest)
		fmt.Errorf("Error, has failed encoding of webhook ID, stopped registration.", err.Error())
		return
	}
	// If no errors, then append safely to local storage of webhooks
	tempWebhooks = append(tempWebhooks, newWebhook)
	w.WriteHeader(http.StatusCreated)
	log.Println("Has successfully registered webhook to storage.")
	tempId += 1

}

func NotificationsGetHandler(w http.ResponseWriter, r *http.Request) {
	splitURL := strings.Split(r.URL.Path, "/")
	if len(splitURL) != 5 {
		http.Error(w, "Error; Incorrect usage of URL.", http.StatusBadRequest)
		log.Println("Attempted to get webhook, but failed due to improper usage of URL.")
		return
	}
	if len(tempWebhooks) <= 0 {
		http.Error(w, "Error; No webhooks registered at all.", http.StatusNotFound)
		log.Println("Attempted to find webhooks, but there are none registered at all.")
		return
	}

	var temp []WebhookObject
	for _, v := range tempWebhooks {
		obj := WebhookObject{URL: v.URL, ID: v.ID, Calls: v.Calls}
		temp = append(temp, obj)
	}

	if splitURL[4] == "" {
		w.Header().Add("content-type", "application/json")
		encoder := json.NewEncoder(w)
		err := encoder.Encode(temp)
		if err != nil {
			log.Println("Attempted to return JSON of all registered webhooks, failed.")
			fmt.Errorf("Error, has failed encoding of all webhooks registered", err.Error())
		}
		log.Println("Has successfully returned all webhooks in GET request.")
		return
	}

	index := locateWebhookByID(splitURL[4])
	if index != -1 {
		selectedWebhook := temp[index]
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
func locateWebhookByID(id string) int {
	for i, v := range tempWebhooks {
		if v.ID == id {
			return i
		}
	}
	return -1
}

// Remember, rough steps need more refining, for later
// "He who has not tasted grapes says sour"
func invocationCall(w http.ResponseWriter, webhook WebhookObject, countryName string) {
	// Generates an object with the information we want to send
	response := WebhookObject{ID: webhook.ID, Calls: webhook.Calls, ISO: countryName}
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
