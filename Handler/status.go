package Handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

func Uptime(t time.Time) time.Duration {
	return time.Since(t)
}

/*
Diagnostic handler to showcases access to request content (headers, body, method, parameters, etc.)
*/
func StatusHandler(w http.ResponseWriter, r *http.Request) {

	countriesResp, err := http.Get(COUNTRIESAPIURL)
	if err != nil {
		log.Fatal(err)
	}

	output := "{" + LINEBREAK
	// Prepares return info with API status codes, version number and uptime
	output += "\tcountries_api: " + strconv.Itoa(countriesResp.StatusCode) + " - " + http.StatusText(countriesResp.StatusCode) + LINEBREAK
	output += "\tnotification_db: " + LINEBREAK
	output += "\twebhooks: " + strconv.Itoa(len(returnWebhooks())) + LINEBREAK
	output += "\tversion: v1" + LINEBREAK
	output += "\tuptime: " + Uptime(StartTime).Round(100000000).String() + LINEBREAK + "}" //Converting time.Duration to string

	// For all options for Printf see https://yourbasic.org/golang/fmt-printf-reference-cheat-sheet/
	_, err = fmt.Fprintf(w, "%v", output)
	if err != nil {
		http.Error(w, "Error when returning output", http.StatusInternalServerError)
	}

}
