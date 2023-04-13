package Handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func RenewablesHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	endpoint := ""
	country := ""
	//If country is not specified
	if len(parts) == 5 {
		endpoint = parts[len(parts)-1]
	}
	//If country is specified
	if len(parts) == 6 {
		endpoint = parts[len(parts)-2]
		country = strings.ToUpper(parts[len(parts)-1])
	}

	switch endpoint {
	case "current":
		HandleRenewablesCurrent(w, r, country)
		break
	case "history":
		HandleRenewablesHistory(w, r, country)
		break
	default:
		http.Error(w, "Endpoint does not exist. Renewables has endpoints \"/current\" and \"/history\"", http.StatusBadRequest)
		return
	}
}

func HandleRenewablesCurrent(w http.ResponseWriter, r *http.Request, isocode string) {
	handlingTime := time.Now()
	neighboursPar := r.URL.Query().Get("neighbours")
	var countries []Country
	//If country is specified, add only the last
	if isocode != "" {
		countries = append(countries, countrySearch(isocode))
	} else {
		countries = readFromCSV(CSVPATH)
	}

	//Get neighbouring countries data if specified
	if neighboursPar == "true" {
		response, err := http.Get(COUNTRIESAPIALPHA + countries[0].ISO + "?fields=borders")
		if err != nil {
			fmt.Print(err.Error())
		}
		countryResp := Country{}
		decoder := json.NewDecoder(response.Body)
		err = decoder.Decode(&countryResp)
		if err != nil {
			http.Error(w, "Error: country: \""+countries[0].Name+"\" not found", http.StatusBadRequest)
			return
		}
		for _, c := range countryResp.Borders {
			countries = append(countries, countrySearch(c))
		}
	}

	//Remove all but the current data:
	for i, c := range countries {
		countries[i].Year = c.Year[len(c.Year)-1:]
		countries[i].Percentage = c.Percentage[len(c.Percentage)-1:]
	}

	//Formats in a pretty format
	output, err := json.MarshalIndent(countries, "", " ")
	if err != nil {
		http.Error(w, "Error during pretty printing", http.StatusInternalServerError)
		return
	}
	searchInfoOutput := "Number of matches: " + strconv.Itoa(len(countries)) + LINEBREAK
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

}

func HandleRenewablesHistory(w http.ResponseWriter, r *http.Request, isocode string) {
	handlingTime := time.Now()
	http.Header.Add(w.Header(), "content-type", "application/json")
	var countries []Country

	if isocode != "" {
		countries = append(countries, countrySearch(isocode))
	} else {
		countries = readFromCSV(CSVPATH)
	}

	//Formats in a pretty format
	output, err := json.MarshalIndent(countries, "", " ")
	if err != nil {
		http.Error(w, "Error during pretty printing", http.StatusInternalServerError)
		return
	}
	searchInfoOutput := "Number of matches: " + strconv.Itoa(len(countries)) + LINEBREAK
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
}
