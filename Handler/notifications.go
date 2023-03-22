package Handler

import (
	"net/http"
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
		http.Error(w, "ERROR: invalid request method. Endpoint has supported methods GET, POST and DELETE", http.StatusBadRequest)
		return
	}

}

func NotificationsDeleteHandler(w http.ResponseWriter, r *http.Request) {

}

func NotificationsPostHandler(w http.ResponseWriter, r *http.Request) {

}

func NotificationsGetHandler(w http.ResponseWriter, r *http.Request) {

}
