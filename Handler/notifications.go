package Handler

import (
	"net/http"
	"strings"
)

func NotificationsHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")

	switch r.Method {
	case "DELETE":
		if len(parts) == 5 {
			NotificationsDeleteHandler(w, r)
		} else {
			http.Error(w, "ERROR: invalid URL for method. Correct path : "+BASEPATH+NOTIFICATIONSPATH+"/{id}", http.StatusBadRequest)
			return
		}
		break

	case "POST":
		NotificationsPostHandler(w, r)
		break

	case "GET":
		NotificationsGetHandler(w, r)
		break
	default:
		http.Error(w, "ERROR: invalid request method. Endpoint has supported methods GET, POST and DELETE", http.StatusBadRequest)
		return
	}

	http.Error(w, "OK", http.StatusOK)
}

func NotificationsDeleteHandler(w http.ResponseWriter, r *http.Request) {

}

func NotificationsPostHandler(w http.ResponseWriter, r *http.Request) {

}

func NotificationsGetHandler(w http.ResponseWriter, r *http.Request) {

}
