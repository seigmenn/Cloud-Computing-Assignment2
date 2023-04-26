package Testing

import (
	"assignment-2/Handler"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

const URL = Handler.BASEPATH + "renewables/current/nor?neighbours=true"
const URL2 = Handler.BASEPATH + "renewables/history/nor?begin=2015&end=2017"

/*
Tests students service, but requires manual start of service prior to invocation.
*/
func TestHttpGetCountryCurrentManual(t *testing.T) {

	// Create client instance
	client := http.Client{}

	// Retrieve content from server
	res, err := client.Get("http://localhost:8080" + URL)
	if err != nil {
		t.Fatal("Get request to URL failed. Check whether server has been started manually! Error:", err.Error())
	}

	// Decode array
	s := []Handler.CountryOut{}
	err2 := json.NewDecoder(res.Body).Decode(&s)
	if err2 != nil {
		t.Fatal("Error during decoding:", err2.Error())
	}

	// Perform content checks
	if len(s) != 4 {
		t.Fatal("Number of returned countries is wrong: " + strconv.Itoa(len(s)))
	}

	for _, country := range s {
		// Perform check of entries (randomly, since order of return may vary)
		switch country.ISO {
		case "NOR":
			// Specific students checks
			if country.Percentage != 71.55836486816406 || country.Name != "Norway" || country.Year != 2021 {
				t.Fatal("Country info about " + country.ISO + " is wrong")
			}
			break
		case "SWE":
			// Specific students checks
			if country.Percentage != 50.924007415771484 || country.Name != "Sweden" || country.Year != 2021 {
				t.Fatal("Country info about " + country.ISO + " is wrong")
			}
			break
		case "RUS":
			// Specific students checks
			if country.Percentage != 6.620289325714111 || country.Name != "Russia" || country.Year != 2021 {
				t.Fatal("Country info about " + country.ISO + " is wrong")
			}
			break
		case "FIN":
			// Specific students checks
			if country.Percentage != 34.611289978027344 || country.Name != "Finland" || country.Year != 2021 {
				t.Fatal("Country info about " + country.ISO + " is wrong")
			}
			break
		default:
			t.Fatal("Invalid country ISO:", country.ISO)
		}
	}

}

/*
Tests students service, and automated setup and tear down of infrastructure.
*/
func TestHttpGetCountryCurrentAutomated(t *testing.T) {

	// Set up infrastructure to be used for invocation - important: wrap handler function in http.HandlerFunc()
	server := httptest.NewServer(http.HandlerFunc(Handler.HandleRenewablesCurrent))

	// Ensure it is torn down properly at the end
	//defer server.Close()

	// Create client instance
	client := http.Client{}

	// Retrieve content from server
	res, err := client.Get(server.URL + URL)
	if err != nil {
		t.Fatal("Get request to URL failed:", err.Error())
	}

	// Decode array
	s := []Handler.CountryOut{}
	err2 := json.NewDecoder(res.Body).Decode(&s)
	if err2 != nil {
		t.Fatal("Error during decoding:", err2.Error())
	}

	// Perform content checks
	if len(s) != 4 {
		t.Fatal("Number of returned countries is wrong: " + strconv.Itoa(len(s)))
	}

	for _, country := range s {
		// Perform check of entries (randomly, since order of return may vary)
		switch country.ISO {
		case "NOR":
			// Specific country checks
			if country.Percentage != 71.55836486816406 || country.Name != "Norway" || country.Year != 2021 {
				t.Fatal("Country info about " + country.ISO + " is wrong")
			}
			break
		case "SWE":
			// Specific country checks
			if country.Percentage != 50.924007415771484 || country.Name != "Sweden" || country.Year != 2021 {
				t.Fatal("Country info about " + country.ISO + " is wrong")
			}
			break
		case "RUS":
			// Specific country checks
			if country.Percentage != 6.620289325714111 || country.Name != "Russia" || country.Year != 2021 {
				t.Fatal("Country info about " + country.ISO + " is wrong")
			}
			break
		case "FIN":
			// Specific country checks
			if country.Percentage != 34.611289978027344 || country.Name != "Finland" || country.Year != 2021 {
				t.Fatal("Country info about " + country.ISO + " is wrong")
			}
			break
		default:
			t.Fatal("Invalid country ISO:", country.ISO)
		}
	}
}

/*
Tests students service, and automated setup and tear down of infrastructure.
*/
func TestHttpGetCountryHistoryAutomated(t *testing.T) {

	// Set up infrastructure to be used for invocation
	server := httptest.NewServer(http.HandlerFunc(Handler.HandleRenewablesHistory))

	// Ensure it is torn down properly at the end
	//defer server.Close()

	// Create client instance
	client := http.Client{}

	// Retrieve content from server
	res, err := client.Get(server.URL + URL2)
	if err != nil {
		t.Fatal("Get request to URL failed:", err.Error())
	}

	// Decode array
	s := []Handler.CountryOut{}
	err2 := json.NewDecoder(res.Body).Decode(&s)
	if err2 != nil {
		t.Fatal("Error during decoding:", err2.Error())
	}

	// Perform content checks
	if len(s) != 3 {
		t.Fatal("Number of returned countries is wrong: " + strconv.Itoa(len(s)))
	}

	for _, country := range s {
		// Perform check of entries (randomly, since order of return may vary)
		switch country.Year {
		case 2015:
			// Specific country checks
			if country.Percentage != 68.87519073486328 || country.Name != "Norway" || country.ISO != "NOR" {
				t.Fatal("Country info about " + country.ISO + " is wrong")
			}
			break
		case 2016:
			// Specific country checks
			if country.Percentage != 69.86628723144531 || country.Name != "Norway" || country.ISO != "NOR" {
				t.Fatal("Country info about " + country.ISO + " is wrong")
			}
			break
		case 2017:
			// Specific country checks
			if country.Percentage != 69.26099395751953 || country.Name != "Norway" || country.ISO != "NOR" {
				t.Fatal("Country info about " + country.ISO + " is wrong")
			}
			break
		default:
			t.Fatal("Invalid country year:", country.Year)
		}
	}

}
