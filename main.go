package main

import (
	"assignment-2/Handler"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
)

func readFromCSV(filePath string) {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Couldn't read file "+filePath, err)
		os.Exit(1)
	}
	csvReader := csv.NewReader(f)
	allData, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("CSV file could not be parsed "+filePath, err)
	}
	for _, c := range allData {
		fmt.Println(c)
	}
}

/*
Default handler displaying service information.
*/
func defaultHandler(w http.ResponseWriter, r *http.Request) {

	// Define content type, so browser renders links correctly
	w.Header().Add("content-type", "text/html")

	// Prepare output returned to client
	output := "This service offers the" +
		"<a href=\"" + Handler.RENEWABLESPATH + "/current\"> " + Handler.RENEWABLESPATH + "/current</a> endpoint, " + Handler.LINEBREAK +
		"<a href=\"" + Handler.RENEWABLESPATH + "/history\"> " + Handler.RENEWABLESPATH + "/history</a> endpoint, " + Handler.LINEBREAK +
		"<a href=\"" + Handler.NOTIFICATIONSPATH + "\"> " + Handler.NOTIFICATIONSPATH + "</a> endpoint, " + Handler.LINEBREAK +
		"<a href=\"" + Handler.STATUSPATH + "\"> " + Handler.STATUSPATH + "</a> endpoint, "

	// Write output to client
	_, err := fmt.Fprintf(w, "%v", output)

	// Deal with error if any
	if err != nil {
		http.Error(w, "Error when returning output", http.StatusInternalServerError)
	}
}

func main() {

	// Extract PORT variable from the environment variables
	port := os.Getenv("PORT")

	// Override port with default port if not provided (e.g. local deployment)
	if port == "" {
		log.Println("$PORT has not been set. Default: 8080")
		port = "8080"
	}
	readFromCSV("renewable-share-energy.csv")
	// Default handler for requests (just displays information and points to /diag)
	http.HandleFunc("/", defaultHandler)
	// Assign path for diagnostics handler (actual service feature)
	http.HandleFunc(Handler.RENEWABLESPATH, Handler.RenewablesHandler)
	http.HandleFunc(Handler.NOTIFICATIONSPATH, Handler.NotificationsHandler)
	http.HandleFunc(Handler.STATUSPATH, Handler.StatusHandler)

	// Start HTTP server
	log.Println("Starting server on port " + port + " ...")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
