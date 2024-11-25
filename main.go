package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

// Define the request structure
type SearchRequest struct {
	StartRow               int          `json:"StartRow"`
	PageSize               int          `json:"PageSize"`
	ShouldReturnHighlights bool         `json:"ShouldReturnHighlights"`
	ShouldCountFacets      bool         `json:"ShouldCountFacets"`
	SortOrder              string       `json:"SortOrder"`
	SearchTerms            []SearchTerm `json:"SearchTerms"`
	Contentsoorten         []string     `json:"Contentsoorten"`
	Rechtsgebieden         []string     `json:"Rechtsgebieden"`
	Instanties             []string     `json:"Instanties"`
	DatumPublicatie        []string     `json:"DatumPublicatie"`
	DatumUitspraak         []string     `json:"DatumUitspraak"`
	Advanced               Advanced     `json:"Advanced"`
	CorrelationId          string       `json:"CorrelationId"`
	Proceduresoorten       []string     `json:"Proceduresoorten"`
}

type SearchTerm struct {
	Term  string `json:"Term"`
	Field string `json:"Field"`
}

type Advanced struct {
	PublicatieStatus string `json:"PublicatieStatus"`
}

// Define the response structure
type SearchResponse struct {
	Results []Result `json:"Results"`
}

// Define RelatieVerwijzing struct
type RelatieVerwijzing struct {
	Id   string `json:"Id"`   // Adjust based on actual response structure
	Name string `json:"Name"` // Adjust based on actual response structure
}

type Result struct {
	Tekstfragment                     string              `json:"Tekstfragment"`
	Titel                             string              `json:"Titel"`
	TitelEmphasis                     string              `json:"TitelEmphasis"`
	InterneUrl                        string              `json:"InterneUrl"`
	DeeplinkUrl                       string              `json:"DeeplinkUrl"`
	Uitspraakdatum                    string              `json:"Uitspraakdatum"`
	UitspraakdatumType                string              `json:"UitspraakdatumType"`
	RelatieVerwijzingen               []RelatieVerwijzing `json:"RelatieVerwijzingen"` // Update to []RelatieVerwijzing
	Publicatiedatum                   string              `json:"Publicatiedatum"`
	GerechtelijkProductType           string              `json:"GerechtelijkProductType"`
	Publicatiestatus                  string              `json:"Publicatiestatus"`
	PublicatiedatumDate               string              `json:"PublicatiedatumDate"`
	Proceduresoorten                  []string            `json:"Proceduresoorten"`
	Vindplaatsen                      []Vindplaats        `json:"Vindplaatsen"`
	Rechtsgebieden                    []string            `json:"Rechtsgebieden"`
	InformatieNietGepubliceerdMessage string              `json:"InformatieNietGepubliceerdMessage"`
	IsInactief                        bool                `json:"IsInactief"`
}

type Vindplaats struct {
	Vindplaats          string `json:"Vindplaats"`
	VindplaatsAnnotator string `json:"VindplaatsAnnotator"`
	VindplaatsUrl       string `json:"VindplaatsUrl"`
}

func main() {
	url := "https://uitspraken.rechtspraak.nl/api/zoek"

	// Prompt user for a location keyword
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter a location keyword (e.g., 'urk'): ")
	location, _ := reader.ReadString('\n')
	location = strings.TrimSpace(location) // Remove any leading/trailing whitespace

	// Read words from the provided txt file
	words, err := readWordsFromFile("words.txt")
	if err != nil {
		fmt.Println("Error reading words from file:", err)
		return
	}

	// Initialize parameters for pagination
	startRow := 0
	pageSize := 10
	writtenToLogCount := 0

	// Create the request payload
	requestBody := SearchRequest{
		StartRow:               startRow,
		PageSize:               pageSize,
		ShouldReturnHighlights: true,
		ShouldCountFacets:      true,
		SortOrder:              "UitspraakDatumDesc",
		SearchTerms: []SearchTerm{
			{Term: location, Field: "AlleVelden"}, // Use the user-provided location keyword
		},
		Contentsoorten:  []string{},
		Rechtsgebieden:  []string{},
		Instanties:      []string{},
		DatumPublicatie: []string{},
		DatumUitspraak:  []string{},
		Advanced: Advanced{
			PublicatieStatus: "AlleenGepubliceerd",
		},
		CorrelationId:    "cc7a7c98d8bb4c6cb9de8babea65b402",
		Proceduresoorten: []string{},
	}

	// Create logs directory if it doesn't exist
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		os.Mkdir("logs", 0755)
	}

	// Generate log file with unique date
	logFileName := fmt.Sprintf("logs/results_%s.log", time.Now().Format("20060102_150405"))
	logFile, err := os.Create(logFileName)
	if err != nil {
		fmt.Println("Error creating log file:", err)
		return
	}
	defer logFile.Close()

	// Start timer
	startTime := time.Now()

	// Variable to track if any matches were found
	anyResultsFound := false

	for {
		// Update the StartRow for pagination
		requestBody.StartRow = startRow

		// Convert the request body to JSON
		jsonData, err := json.Marshal(requestBody)
		if err != nil {
			fmt.Println("Error marshalling JSON:", err)
			return
		}

		// Create the POST request
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Println("Error making POST request:", err)
			return
		}
		defer resp.Body.Close()

		// Check the response status
		if resp.StatusCode != http.StatusOK {
			fmt.Println("Received non-200 response status:", resp.Status)
			return
		}

		// Read the response body
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return
		}

		// Parse the response
		var response SearchResponse
		if err := json.Unmarshal(body, &response); err != nil {
			fmt.Println("Error unmarshalling response JSON:", err)
			return
		}

		// Handle the response data
		if len(response.Results) == 0 {
			break // Exit the loop if no results are returned
		}

		// Scan for keywords from the txt file and log results
		resultsFound := scanAndLogKeywords(response.Results, words, logFile)
		if resultsFound {
			anyResultsFound = true // Set flag if any results were found
			writtenToLogCount += 1
		}

		// Increment the start row for the next batch of results
		startRow += pageSize

		// Display progress
		elapsed := time.Since(startTime).Seconds()
		fmt.Printf("\rTime elapsed: %.2fs | Items written to log: %d", elapsed, writtenToLogCount)
	}

	// Print if no results were found at all
	if !anyResultsFound {
		fmt.Println("\n0 results found.")
	} else {
		fmt.Println("\nSearch completed.")
	}
}

func readWordsFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words = append(words, strings.TrimSpace(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return words, nil
}

func scanAndLogKeywords(results []Result, keywords []string, logFile *os.File) bool {
	found := false // Flag to track if any matches were found
	for _, result := range results {
		summary := strings.ToLower(result.Tekstfragment)
		for _, keyword := range keywords {
			if strings.Contains(summary, keyword) {
				logEntry := fmt.Sprintf("Found in Title: %s\nSummary: %s\nLink: %s\n-------------------------------\n",
					result.Titel, result.Tekstfragment, result.DeeplinkUrl)
				logFile.WriteString(logEntry)
				found = true
				break
			}
		}
	}
	return found
}
