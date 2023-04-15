package Handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Note; I'll clean it up when the functionality is complete.

// Until Firebase gets sorted, this is the main storage for webhooks;
// will be changed to work throughout the entire project to work for
// invocation within renewables.go
var tempWebhooks []webhookObject

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
		http.Error(w, "ERROR: invalid request method. Endpoint has supported methods GET, POST and DELETE", http.StatusBadRequest)
		return
	}

}

func NotificationsDeleteHandler(w http.ResponseWriter, r *http.Request) {
	splitURL := strings.Split(r.URL.Path, "/")
	fmt.Print(splitURL)
	fmt.Print(splitURL[0])
	if len(splitURL) != 5 || splitURL[4] == "" {
		http.Error(w, "Error; Incorrect usage of URL.", http.StatusBadRequest)
		log.Println("Attempted to delete webhook, but failed due to improper usage of URL.")
		return
	}
	if len(tempWebhooks) <= 0 {
		http.Error(w, "Error; No webhooks registered at all.", http.StatusBadRequest)
		log.Println("Attempted to find webhooks, but there are none registered at all.")
		return
	}
	index := locateWebhookByID(splitURL[4])
	if index == -1 {
		http.Error(w, "Error, no webhook of such ID found.", http.StatusBadRequest)
		log.Println("Attempted to delete webhook, but failed due to nonexisting id.")
		return
	}

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

	// replace with more efficient delete afterwards through append and slicing [:index], [index:] later
	var temp []webhookObject
	for _, v := range tempWebhooks {
		if v.ID != tempWebhooks[index].ID {
			temp = append(temp, v)
		}
	}
	tempWebhooks = temp

}

func NotificationsPostHandler(w http.ResponseWriter, r *http.Request) {
	temporal := strconv.Itoa(tempId)
	newWebhook := webhookObject{ID: temporal}
	err := json.NewDecoder(r.Body).Decode(&newWebhook)
	if err != nil {
		http.Error(w, "Error in decoding POST request", http.StatusBadRequest)
	}
	w.Header().Add("content-type", "application/json")
	encoder := json.NewEncoder(w)
	err = encoder.Encode(webhookObject{ID: newWebhook.ID})
	if err != nil {
		log.Println("Attempted to return JSON of ID, failed registration.", http.StatusBadRequest)
		fmt.Errorf("Error, has failed encoding of webhook ID, stopped registration.", err.Error())
	}
	tempWebhooks = append(tempWebhooks, newWebhook)
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
		http.Error(w, "Error; No webhooks registered at all.", http.StatusBadRequest)
		log.Println("Attempted to find webhooks, but there are none registered at all.")
		return
	}

	if splitURL[4] == "" {
		w.Header().Add("content-type", "application/json")
		encoder := json.NewEncoder(w)
		err := encoder.Encode(tempWebhooks)
		if err != nil {
			log.Println("Attempted to return JSON of all registered webhooks, failed.")
			fmt.Errorf("Error, has failed encoding of all webhooks registered", err.Error())
		}
		log.Println("Has successfully returned all webhooks in GET request.")
		return
	}

	index := locateWebhookByID(splitURL[4])
	if index != -1 {
		selectedWebhook := tempWebhooks[index]
		w.Header().Add("content-type", "application/json")
		encoder := json.NewEncoder(w)
		err := encoder.Encode(selectedWebhook)
		if err != nil {
			log.Println("Attempted to return JSON of selected singular webhook, failed.")
			fmt.Errorf("Error, has failed encoding of singular webhook", err.Error())
		}
		log.Println("Has successfully returned singular webhook in GET request.")
	} else {
		http.Error(w, "Error; No webhooks of such ID registered at all.", http.StatusBadRequest)
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
