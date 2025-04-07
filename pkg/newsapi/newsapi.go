package newsapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

func NewsAPIEverything(api_key, query_req NewsAPIEverythingReq) (NewsAPIResponse, error) {

	var newsAPIResponse NewsAPIResponse

	return newsAPIResponse, nil
}

func NewsAPITopHeadlines(api_key string, query_req NewsAPITopHeadlinesReq) (NewsAPIResponse, error) {

	if api_key == "" {
		return NewsAPIResponse{}, nil
	}

	// * make a request to the NewsAPI with the given parameters
	httpClient := &http.Client{}
	req, err := http.NewRequest("GET", NewsAPIURLTopHeadlines, nil)
	if err != nil {
		return NewsAPIResponse{}, err
	}

	req.Header.Set("Authorization", "Bearer "+api_key)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// query req checker
	pageSize := query_req.PageSize
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	page := query_req.Page
	if page < 1 {
		page = 1
	}

	// Get the current query values
	q := req.URL.Query()

	// Set the query parameters
	q.Set("pageSize", strconv.Itoa(pageSize))
	q.Set("page", strconv.Itoa(page))

	if query_req.Category != "" {
		if query_req.Category != "business" && query_req.Category != "entertainment" && query_req.Category != "general" && query_req.Category != "health" && query_req.Category != "science" && query_req.Category != "sports" && query_req.Category != "technology" {
			return NewsAPIResponse{}, fmt.Errorf("invalid category: %s. valid categories are business, entertainment, general, health, science, sports, technology", query_req.Category)
		}

		q.Set("category", query_req.Category)
	}

	if query_req.Country != "" {
		if query_req.Country != "us" {
			return NewsAPIResponse{}, fmt.Errorf("invalid country code: %s, valid country code is us", query_req.Country)
		}

		q.Set("country", query_req.Country)
	}

	if query_req.Q != "" {
		q.Set("q", query_req.Q)
	}

	// Assign the modified query values back to the request URL
	req.URL.RawQuery = q.Encode()

	resp, err := httpClient.Do(req)
	if err != nil {
		return NewsAPIResponse{}, err
	}

	// * read the response body and unmarshal it into the NewsAPIResponse struct
	// * if the response is not 200 OK, read the body and return an error
	defer func() {
		if resp.StatusCode != http.StatusOK {
			io.ReadAll(resp.Body)
		}
		resp.Body.Close()
	}()

	// decode response from json to struct
	var newsAPIResponse NewsAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&newsAPIResponse); err != nil {
		return NewsAPIResponse{}, err
	}

	return newsAPIResponse, nil
}
