package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"index/suffixarray"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

func main() {
	searcher := Searcher{}
	err := searcher.Load("completeworks.txt")
	if err != nil {
		log.Fatal(err)
	}

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	http.HandleFunc("/search", handleSearch(searcher))

	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	fmt.Printf("shakesearch available at http://localhost:%s...", port)
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}

const _MAX_RESULTS_PER_QUERY = 20

type Searcher struct {
	CompleteWorks string
	SuffixArray   *suffixarray.Index
}

func handleSearch(searcher Searcher) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// parse query parameter
		var query string
		if strQuery, ok := r.URL.Query()["q"]; ok && len(strQuery[0]) > 0 {
			query = strQuery[0]
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("missing search query in URL params"))
			return
		}

		// parse offset parameter, default 0
		var offset int
		if strOffset, ok := r.URL.Query()["offset"]; ok {
			o, err := strconv.ParseInt(strOffset[0], 10, 32)
			// make sure to pass only positive offset
			// otherwise in case of any error use default 0 offset
			if err == nil && o >= 0 {
				offset = int(o)
			}
		}

		// TODO: use `limit` instead of _MAX_RESULTS_PER_QUERY  as a query parameter

		// run search
		results, err := searcher.Search(query, offset, _MAX_RESULTS_PER_QUERY)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("bad query in URL params: %q", err)))
			return
		}

		// generate response
		var buf bytes.Buffer
		enc := json.NewEncoder(&buf)
		if err = enc.Encode(results); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("encoding failure"))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(buf.Bytes())
	}
}

func (s *Searcher) Load(filename string) error {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Load: %w", err)
	}
	s.CompleteWorks = string(dat)
	s.SuffixArray = suffixarray.New(dat)
	return nil
}

func sliceResults(results [][]int, offset, totalPerQuery int) [][]int {
	if offset >= len(results) {
		// Offset is beyond the length
		return nil
	}

	end := offset + totalPerQuery
	if end > len(results) {
		// Make sure the end doesn't exceed the length of results
		end = len(results)
	}
	return results[offset:end]
}

func (s *Searcher) Search(query string, offset, num int) ([]string, error) {
	rx, err := regexp.Compile(fmt.Sprintf("(?i)%v", regexp.QuoteMeta(query)))
	if err != nil {
		return nil, err
	}

	results := []string{}
	idxs := s.SuffixArray.FindAllIndex(rx, offset+num)
	for _, idx := range sliceResults(idxs, offset, num) {
		results = append(results, s.CompleteWorks[idx[0]-250:idx[1]+250])
	}
	return results, nil
}
