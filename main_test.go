package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestSearchHamlet(t *testing.T) {
	searcher := Searcher{}
	err := searcher.Load("completeworks.txt")
	if err != nil {
		t.Fatal(err)
	}

	query := "Hamlet"
	req, err := http.NewRequest("GET", "/search?q="+query, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleSearch(searcher))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var results []string
	err = json.Unmarshal(rr.Body.Bytes(), &results)
	if err != nil {
		t.Fatal(err)
	}

	found := false
	for _, result := range results {
		if strings.Contains(strings.ToLower(result), strings.ToLower(query)) {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("expected result not found for query: %s", query)
	}
}

func TestSearchCaseSensitive(t *testing.T) {
	searcher := Searcher{}
	err := searcher.Load("completeworks.txt")
	if err != nil {
		t.Fatal(err)
	}

	query := "hAmLeT"
	req, err := http.NewRequest("GET", "/search?q="+query, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleSearch(searcher))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var results []string
	err = json.Unmarshal(rr.Body.Bytes(), &results)
	if err != nil {
		t.Fatal(err)
	}

	found := false
	for _, result := range results {
		if strings.Contains(strings.ToLower(result), strings.ToLower(query)) {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("expected result not found for query: %s", query)
	}
}

func TestSearchDrunk(t *testing.T) {
	searcher := Searcher{}
	err := searcher.Load("completeworks.txt")
	if err != nil {
		t.Fatal(err)
	}

	query := "drunk"
	req, err := http.NewRequest("GET", "/search?q="+query, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleSearch(searcher))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var results []string
	err = json.Unmarshal(rr.Body.Bytes(), &results)
	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 20 {
		t.Errorf("expected 20 results for query: %s, got %d", query, len(results))
	}
}

func TestSearchWithOffset(t *testing.T) {
	searcher := Searcher{}
	err := searcher.Load("completeworks.txt")
	if err != nil {
		t.Fatal(err)
	}

	query := "Florence" // 32 words in completeworks.txt

	var testCases = []struct {
		name                 string // name of test case
		queryOffset          int64  // query offset
		expectedNumOfResults int    // expected number of results
	}{
		{"query w/o offset should use default 0", -1, 20},
		{"query with negative should use default 0", -10, 20},
		{"query with offset=0", 0, 20},
		{"query with offset=20", 20, 12},
		{"query with offset=40", 40, 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var q strings.Builder
			q.WriteString("/search?q=")
			q.WriteString(query)
			if tc.queryOffset != -1 {
				q.WriteString("&offset=")
				q.WriteString(strconv.FormatInt(tc.queryOffset, 10))
			}

			req, err := http.NewRequest("GET", q.String(), nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(handleSearch(searcher))
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
			}

			var results []string
			err = json.Unmarshal(rr.Body.Bytes(), &results)
			if err != nil {
				t.Fatal(err)
			}

			if len(results) != tc.expectedNumOfResults {
				t.Errorf("expected %v results for query: %s, got %d", tc.expectedNumOfResults, query, len(results))
			}
		})
	}
}
