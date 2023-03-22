package Handler

import (
	"net/http"
	"strings"
)

func RenewablesHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	request := parts[len(parts)]
	switch request {
	case "current":
		HandleRenewablesCurrent(w, r)
		break
	case "history":
		HandleRenewablesHistory(w, r)
		break
	default:
		http.Error(w, "Endpoint does not exist. Renewables has endpoints /current and /history", http.StatusBadRequest)
		return
	}
}

func HandleRenewablesCurrent(w http.ResponseWriter, r *http.Request) {

	/*
		//Formats in a pretty format
		output, err := json.MarshalIndent(uniRespNeighbour, "", " ")
		if err != nil {
			http.Error(w, "Error during pretty printing", http.StatusInternalServerError)
			return
		}
		searchInfoOutput := "Number of matches: " + strconv.Itoa(len(uniRespNeighbour)) + LINEBREAK
		searchInfoOutput += "Found in " + Uptime(handlingTime).Round(10000000).String() + LINEBREAK
		//writes to responseWriter
		_, err = fmt.Fprintf(w, "%v", searchInfoOutput)
		if err != nil {
			http.Error(w, "Error when returning InfoOutput", http.StatusInternalServerError)
		}
		_, err = fmt.Fprintf(w, "%v", string(output))
		if err != nil {
			http.Error(w, "Error when returning UniOutput", http.StatusInternalServerError)
		}
	*/
}

func HandleRenewablesHistory(w http.ResponseWriter, r *http.Request) {

	/*
		//Formats in a pretty format
		output, err := json.MarshalIndent(uniRespNeighbour, "", " ")
		if err != nil {
			http.Error(w, "Error during pretty printing", http.StatusInternalServerError)
			return
		}
		searchInfoOutput := "Number of matches: " + strconv.Itoa(len(uniRespNeighbour)) + LINEBREAK
		searchInfoOutput += "Found in " + Uptime(handlingTime).Round(10000000).String() + LINEBREAK
		//writes to responseWriter
		_, err = fmt.Fprintf(w, "%v", searchInfoOutput)
		if err != nil {
			http.Error(w, "Error when returning InfoOutput", http.StatusInternalServerError)
		}
		_, err = fmt.Fprintf(w, "%v", string(output))
		if err != nil {
			http.Error(w, "Error when returning UniOutput", http.StatusInternalServerError)
		}
	*/
}
