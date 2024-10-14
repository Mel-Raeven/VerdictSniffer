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

	// Initialize parameters for pagination
	startRow := 0
	pageSize := 10

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

		// Scan for the keywords "drug" or "drugs"
		resultsFound := scanForKeywords(response.Results)
		if resultsFound {
			anyResultsFound = true // Set flag if any results were found
		}

		// Increment the start row for the next batch of results
		startRow += pageSize
	}

	// Print if no results were found at all
	if !anyResultsFound {
		fmt.Println("0 results found.")
	}
}

// scanForKeywords scans case summaries for the keywords "drug" or "drugs".
func scanForKeywords(results []Result) bool {
	found := false // Flag to track if any matches were found
	for _, result := range results {
		summary := strings.ToLower(result.Tekstfragment)
		if strings.Contains(summary, "drug") || strings.Contains(summary, "drugs") {
			fmt.Println("Found in Title:", result.Titel)
			fmt.Println("Summary:", result.Tekstfragment)
			fmt.Println("Link:", result.DeeplinkUrl)
			fmt.Println("-------------------------------")
			found = true
		}
	}
	return found
}
