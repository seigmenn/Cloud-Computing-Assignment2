package Handler

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
)

func readFromCSV(filePath string) []Country {
	nameIndex := 0
	ISOIndex := 1
	yearIndex := 2
	percentageIndex := 3

	f, err := os.Open(filePath)
	//If file couldn't be opened return empty slice
	if err != nil {
		log.Fatal("Couldn't read file "+filePath, err)
		return []Country{}
	}
	csvReader := csv.NewReader(f)
	allData, err := csvReader.ReadAll()
	//If file couldn't be parsed return empty slice
	if err != nil {
		log.Fatal("CSV file could not be parsed "+filePath, err)
		return []Country{}
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
	return countries
}

func countrySearch(ISOcode string) Country {
	countries := readFromCSV(CSVPATH)
	for _, c := range countries {
		//If ISO codes match: return struct
		if c.ISO == ISOcode {
			return c
		}
	}
	//No match found: return empty struct
	return Country{}
}

func printCountries() {

	tmp := readFromCSV(CSVPATH)

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
