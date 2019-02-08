/*
 * kearch specialist search engine gateway API
 *
 * kearch specialist search engine gateway API
 *
 * API version: 0.1.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const hostname = "localhost"

func MakeASummary() Summary {
	file, err := os.Open("en_default_dict.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	res := Summary{SpHost: hostname, EngineName: "kearch-sp-google", Dump: make(map[string]int32)}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		sep := strings.Fields(line)
		n, err := strconv.Atoi(sep[1])
		if err != nil {
			log.Fatal(err)
		}
		res.Dump[sep[0]] = int32(n)
	}
	return res
}

// AddAConnectionRequestPost - Add a connection request sent from meta server to specialist server.
func AddAConnectionRequestPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	decoder := json.NewDecoder(r.Body)

	var data ConnectionRequestOnSp
	err := decoder.Decode(&data)
	if err != nil {
		panic(err)
	}

	me_host := data.MeHost
	scheme := data.Scheme
	log.Printf("me_host = " + me_host + ", scheme = " + scheme)

	summary := MakeASummary()
	summaryStr, err := json.Marshal(&summary)
	_, err = http.Post("sample.com", "application/json", bytes.NewReader(summaryStr))
	if err != nil {
		panic(err)
	}

	res := InlineResponse200{MeHost: me_host}
	out, err := json.Marshal(&res)
	if err != nil {
		panic(err)
	}
	fmt.Fprint(w, string(out))
}

// DeleteAConnectionRequestDelete - Delete a connection request sent from meta server to specialist server.
func DeleteAConnectionRequestDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	me_host := r.URL.Query().Get("me_host")
	log.Printf("me_host = " + me_host)

	res := InlineResponse200{MeHost: me_host}
	out, err := json.Marshal(&res)
	if err != nil {
		panic(err)
	}
	fmt.Fprint(w, string(out))
}

// GetASummaryGet - Get summary of this specialist server.
func GetASummaryGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	summary := MakeASummary()
	out, err := json.Marshal(&summary)
	if err != nil {
		panic(err)
	}
	fmt.Fprint(w, string(out))
}

type Item struct {
	Title   string `json:"title"`
	Link    string `json:"link"`
	Snippet string `json:"snippet"`
}

type GoogleCustomSearchResult struct {
	Items []Item `json:"items`
}

// RetrieveGet - Retrieve search results.
func RetrieveGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	query := r.URL.Query().Get("queries")
	max_urls, err := strconv.Atoi(r.URL.Query().Get("max_urls"))
	if err != nil {
		panic(err)
	}

	log.Printf("query = " + query + ", max_urls = " + strconv.Itoa(max_urls))

	file, err := os.Open("GoogleCustomSearchAPIKey")
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	key := scanner.Text()
	log.Printf(key)
	file.Close()

	file, err = os.Open("GoogleCustomSearchEngineID")
	if err != nil {
		log.Fatal(err)
	}
	scanner = bufio.NewScanner(file)
	scanner.Scan()
	engine := scanner.Text()
	log.Printf(engine)
	file.Close()

	searchUrl := "https://www.googleapis.com/customsearch/v1?key=" + key + "&cx=" + engine + "&q=" + strings.Replace(query, " ", "+", -1)
	log.Printf(searchUrl)

	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(searchUrl)
	if err != nil {
		log.Fatal(err)
	}
	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	resp.Body.Close()

	// Comment in here to test without using Google API
	//
	// file, err = os.Open("responce-example.json")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer file.Close()
	//
	// raw, err := ioutil.ReadAll(file)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	var searchResult GoogleCustomSearchResult
	err = json.Unmarshal(raw, &searchResult)
	if err != nil {
		log.Fatal(err)
	}

	var responce []Document
	for i, r := range searchResult.Items {
		d := Document{Url: r.Link, Title: r.Title, Description: r.Snippet, Score: float32(len(searchResult.Items) - i)}
		responce = append(responce, d)
	}
	out, err := json.Marshal(&responce)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprint(w, string(out))
}
