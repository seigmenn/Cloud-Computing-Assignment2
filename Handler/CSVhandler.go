package Handler

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func readFromCSV(filePath string) ([]Country, error) {
	nameIndex := 0
	ISOIndex := 1
	yearIndex := 2
	percentageIndex := 3

	f, err := os.Open(filePath)
	//If file couldn't be opened return empty slice
	if err != nil {
		log.Fatal("Couldn't read file "+filePath, err)
		return []Country{}, ERRFILEREAD
	}
	csvReader := csv.NewReader(f)
	allData, err := csvReader.ReadAll()
	//If file couldn't be parsed return empty slice
	if err != nil {
		log.Fatal("CSV file could not be parsed "+filePath, err)
		return []Country{}, ERRFILEPARSE
	}
	oldName := ""

	var countries []Country
	var tmpCountry Country
	for _, c := range allData {
		newName := c[nameIndex]
		//If new country:
		if newName != oldName {
			//Append last read Country struct if it has a valid name
			if tmpCountry.Name != "" {
				countries = append(countries, tmpCountry)
			}
			//Set name and ISOcode for new country
			tmpCountry = Country{}
			tmpCountry.Name = c[nameIndex]
			tmpCountry.ISO = c[ISOIndex]
		}
		//Reading year and percentage and appending to slices
		year, _ := strconv.Atoi(c[yearIndex])
		tmpCountry.Year = append(tmpCountry.Year, year)
		//Trying to parse from string to float
		if percentage, err := strconv.ParseFloat(c[percentageIndex], 32); err == nil {
			tmpCountry.Percentage = append(tmpCountry.Percentage, percentage)
		}
		oldName = newName
	}
	//Return slice of all read countries
	return countries, nil
}

func countrySearch(ISOcode string) (Country, string, error) {
	countries, err := readFromCSV(CSVPATH)
	if err != nil {
		//No match found: return empty struct and error
		return Country{}, "", err
	}
	strings.ToUpper(ISOcode)
	for _, c := range countries {
		//If ISO codes OR name match: return struct
		if strings.ToUpper(c.ISO) == ISOcode || strings.ToUpper(c.Name) == ISOcode {
			return c, c.ISO, nil
		}
	}
	return Country{}, "", ERRCOUNTRYNOTFOUND
}

func countrySearchSlice(ISOcode []string) ([]Country, error) {
	countries, err := readFromCSV(CSVPATH)
	returnCountries := []Country{}
	if err != nil {
		//No CSV data found: return empty struct and error
		return []Country{}, err
	}

	for _, c := range countries {
		//If ISO codes OR name match: return struct
		for _, iso := range ISOcode {
			strings.ToUpper(iso)
			if strings.ToUpper(c.ISO) == iso || strings.ToUpper(c.Name) == iso {
				returnCountries = append(returnCountries, c)
			}
		}
	}
	if len(returnCountries) == 0 {
		//No countries matched, return empty slice and error
		return []Country{}, ERRCOUNTRYNOTFOUND
	} else {
		//Found at least one country, return slice and nil
		return returnCountries, nil
	}
}

func printCountries() {

	tmp, err := readFromCSV(CSVPATH)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	for _, c := range tmp {
		fmt.Println(c.Name)
		fmt.Println(c.ISO)
		for y := 0; y < len(c.Year); y++ {
			year := strconv.Itoa(c.Year[y])
			percent := strconv.FormatFloat(c.Percentage[y], 'f', -1, 32)
			fmt.Println("Year: " + year + "\tPercent: " + percent)
		}
		fmt.Println("")
	}
}
