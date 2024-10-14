package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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

type Result struct {
	Tekstfragment                     string       `json:"Tekstfragment"`
	Titel                             string       `json:"Titel"`
	TitelEmphasis                     string       `json:"TitelEmphasis"`
	InterneUrl                        string       `json:"InterneUrl"`
	DeeplinkUrl                       string       `json:"DeeplinkUrl"`
	Uitspraakdatum                    string       `json:"Uitspraakdatum"`
	UitspraakdatumType                string       `json:"UitspraakdatumType"`
	RelatieVerwijzingen               []string     `json:"RelatieVerwijzingen"`
	Publicatiedatum                   string       `json:"Publicatiedatum"`
	GerechtelijkProductType           string       `json:"GerechtelijkProductType"`
	Publicatiestatus                  string       `json:"Publicatiestatus"`
	PublicatiedatumDate               string       `json:"PublicatiedatumDate"`
	Proceduresoorten                  []string     `json:"Proceduresoorten"`
	Vindplaatsen                      []Vindplaats `json:"Vindplaatsen"`
	Rechtsgebieden                    []string     `json:"Rechtsgebieden"`
	InformatieNietGepubliceerdMessage string       `json:"InformatieNietGepubliceerdMessage"`
	IsInactief                        bool         `json:"IsInactief"`
}

type Vindplaats struct {
	Vindplaats          string `json:"Vindplaats"`
	VindplaatsAnnotator string `json:"VindplaatsAnnotator"`
	VindplaatsUrl       string `json:"VindplaatsUrl"`
}

func main() {
	url := "https://uitspraken.rechtspraak.nl/api/zoek"

	// Create the request payload
	requestBody := SearchRequest{
		StartRow:               0,
		PageSize:               10,
		ShouldReturnHighlights: true,
		ShouldCountFacets:      true,
		SortOrder:              "UitspraakDatumDesc",
		SearchTerms: []SearchTerm{
			{Term: "urk", Field: "Samenvatting"},
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
	for _, result := range response.Results {
		fmt.Println("Title:", result.Titel)
		fmt.Println("Summary:", result.Tekstfragment)
		fmt.Println("Publication Date:", result.Publicatiedatum)
		fmt.Println("Link:", result.DeeplinkUrl)
		fmt.Println("-------------------------------")
	}
}
