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
	var tmpCountry CountryOut
	var outCountries []CountryOut
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

	//Remove all but the current data and adds to output slice:
	for _, c := range countries {
		tmpCountry = CountryOut{}
		tmpCountry.ISO = c.ISO
		tmpCountry.Name = c.Name
		tmpCountry.Year = c.Year[len(c.Year)-1]
		tmpCountry.Percentage = c.Percentage[len(c.Percentage)-1]
		outCountries = append(outCountries, tmpCountry)
	}

	//Formats in a pretty format
	output, err := json.MarshalIndent(outCountries, "", " ")
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
		http.Error(w, "Error when returning DataOutput", http.StatusInternalServerError)
	}

}

func HandleRenewablesHistory(w http.ResponseWriter, r *http.Request, isocode string) {
	handlingTime := time.Now() //Start clock
	http.Header.Add(w.Header(), "content-type", "application/json")
	var countries []Country
	var outCountries []CountryOut
	var tmpCountry CountryOut

	//Get parameters from request
	beginTime, _ := strconv.Atoi(r.URL.Query().Get("begin"))
	endTime, _ := strconv.Atoi(r.URL.Query().Get("end"))
	if beginTime == 0 {
		beginTime = 1950 //Default value 1950
	}
	if endTime == 0 {
		endTime = 2023 //Default value 2023
	}
	//If country is specified, add only it to the slice
	if isocode != "" {
		countries = append(countries, countrySearch(isocode))
	} else {
		//Add all countries to slice
		countries = readFromCSV(CSVPATH)
	}

	//Get average within given time slice and move to output struct
	for _, c := range countries {
		tmpCountry = CountryOut{}
		sum := 0.0
		num := 0
		for i, y := range c.Percentage {
			//If country has an entry within time slice
			if c.Year[i] >= beginTime && c.Year[i] <= endTime {
				sum += y //Add value to sum
				num++    //Add number of entries that matched criteria
			}
		}
		tmpCountry.ISO = c.ISO
		tmpCountry.Name = c.Name
		if sum > 0 && num > 0 {
			tmpCountry.Percentage = sum / float64(num) //Divide sum by number of entries in that sum
		} else {
			tmpCountry.Percentage = 0
		}
		outCountries = append(outCountries, tmpCountry)
	}

	//Formats in a pretty format
	output, err := json.MarshalIndent(outCountries, "", " ")
	if err != nil {
		fmt.Println(err)
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
		http.Error(w, "Error when returning DataOutput", http.StatusInternalServerError)
	}
}
