package Handler

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
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

	//If country is specified search with isocode
	if isocode != "" {
		tmp, err := countrySearch(isocode)
		if err != nil {
			if err == ERRFILEREAD || err == ERRFILEPARSE {
				http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
				return
			} else {
				http.Error(w, "Error: "+err.Error(), http.StatusBadRequest)
				return
			}

		}
		countries = append(countries, tmp)
	} else { //If not specified all countries are added
		tmp, err := readFromCSV(CSVPATH)
		if err != nil {
			http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		countries = tmp
	}

	//Get neighbouring countries data if specified and country is specified
	if neighboursPar == "true" && isocode != "" {
		//Get request to country api, only need borders data
		response, err := http.Get(COUNTRIESAPIALPHA + isocode + "?fields=borders")
		if err != nil {
			fmt.Print(err.Error())
		}
		//Decode into struct:
		countryResp := Country{}
		decoder := json.NewDecoder(response.Body)
		err = decoder.Decode(&countryResp)
		//Error with decoding:
		if err != nil {
			http.Error(w, "Error: country: \""+isocode+"\" not found", http.StatusBadRequest)
			return
		}
		//Search file for all neighbour countries and append to slice
		//TODO: only one search with slice of country ISOs
		for _, c := range countryResp.Borders {
			tmp, err := countrySearch(c) //countrySearch() returns empty struct if not found
			//Only append if country struct is not empty
			if err != nil {
				http.Error(w, "Error: "+err.Error(), http.StatusBadRequest)
			} else {
				countries = append(countries, tmp)
			}
		}
	} else if neighboursPar == "true" && isocode == "" {
		http.Error(w, "Error: can't print neighbours if no country is specified", http.StatusBadRequest)
		return
	}

	//Remove all but the current data and adds to output slice (prettier print):
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
	searchInfoOutput := "Number of results: " + strconv.Itoa(len(outCountries)) + LINEBREAK
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
	// Assumption: No country comes twice in the "outCountries" for Currently Percentage Renewables
	// If so, it's a safe assumption to do as such:
	for f, u := range tempWebhooks {
		for _, y := range outCountries {
			if u.ISO == y.ISO {
				log.Println("DEBUG: Has detected allike ISOs (CURRENT), proceeds...")
				tempWebhooks[f].Invocations += 1
				if math.Mod(float64(u.Invocations), float64(u.Calls)) == 0 {
					log.Println("The result of the following SHOULD be 0: ")
					fmt.Println(math.Mod(float64(u.Invocations), float64(u.Calls)))
					log.Println("DEBUG: Has detected invocation call (CURRENT), proceeds...")
					invocationCall(w, u)
				}
				break
			}
		}
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
	//If country is specified search with isocode
	if isocode != "" {
		tmp, err := countrySearch(isocode)
		if err != nil {
			http.Error(w, "Error: "+err.Error(), http.StatusBadRequest)
			return
		}
		countries = append(countries, tmp)
	} else { //If not specified all countries are added
		tmp, err := readFromCSV(CSVPATH)
		if err != nil {
			http.Error(w, "Error: "+err.Error(), http.StatusBadRequest)
			return
		}
		countries = tmp
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
				if isocode != "" {

					tmpCountry.Name = c.Name
					tmpCountry.ISO = c.ISO
					tmpCountry.Percentage = y
					tmpCountry.Year = c.Year[i]
					outCountries = append(outCountries, tmpCountry)
				}
			}
		}

		if isocode == "" {
			tmpCountry.ISO = c.ISO
			tmpCountry.Name = c.Name
			if sum > 0 && num > 0 {
				tmpCountry.Percentage = sum / float64(num) //Divide sum by number of entries in that sum
			} else {
				tmpCountry.Percentage = 0
			}
			outCountries = append(outCountries, tmpCountry)
		}
	}

	//Formats in a pretty format
	output, err := json.MarshalIndent(outCountries, "", " ")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error during pretty printing", http.StatusInternalServerError)
		return
	}
	searchInfoOutput := "Number of results: " + strconv.Itoa(len(outCountries)) + LINEBREAK
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
	// In the situation of a specified isocode, all webhooks with selected webhook has invocation + 1
	if isocode != "" {
		for f, u := range tempWebhooks {
			if u.ISO == isocode {
				tempWebhooks[f].Invocations += 1
				if math.Mod(float64(u.Invocations), float64(u.Calls)) == 0 {
					log.Println("The result of the following SHOULD be 0: ")
					fmt.Println(math.Mod(float64(u.Invocations), float64(u.Calls)))
					log.Println("DEBUG: Has detected invocation call (HISTORY, SINGULAR ISO), proceeds...")
					invocationCall(w, u)
				}
			}
		}
	} else {
		// TODO: ASSUMPTION, this assumes that all webhooks have valid isocodes
		// A valid assumption which saves performance compared to before:
		for f, u := range tempWebhooks {
			tempWebhooks[f].Invocations += 1
			if math.Mod(float64(u.Invocations), float64(u.Calls)) == 0 {
				log.Println("The result of the following SHOULD be 0: ")
				fmt.Println(math.Mod(float64(u.Invocations), float64(u.Calls)))
				log.Println("DEBUG: Has detected invocation call (HISTORY, ALL), proceeds...")
				invocationCall(w, u)
			}
		}
	}
}
