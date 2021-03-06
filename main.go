package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type Hits struct {
	author     string `json:"author"`
	created_at int32  `json:"created_at_i"`
	objectID   string `json:"objectID"`
	points     int    `json:"points"`
	title      string `json:"title"`
	url        string `json:"url"`
}

type TopPosts struct {
	Number int `json:"number"`
	Hits   []struct {
		Author    string `json:"author"`
		CreatedAt int32  `json:"created_at_i"`
		ObjectID  string `json:"objectID"`
		Points    int    `json:"points"`
		Title     string `json:"title"`
		Url       string `json:"url"`
	} `json:"hits"`
}

func SearchPosts(search string) error {
	url := fmt.Sprintf("https://hn.algolia.com/api/v1/search?query=%v&tags=story&hitsPerPage=100", search)
	return Get(url)
}

func LoadPosts(unixStart, unixEnd string) error {
	url := fmt.Sprintf("https://hn.algolia.com/api/v1/search_by_date?tags=story&numericFilters=created_at_i>%v,created_at_i<%v", unixStart, unixEnd)
	return Get(url)
}

func LoadTopPosts() error {
	url := "https://hn.algolia.com/api/v1/search?tags=front_page"
	return Get(url)
}

func Get(url string) error {
	dir, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	fmt.Println("Now retrieving latest HN posts, please wait...")
	// two variables (response and error) which stores the response from e GET request
	getRequest, err := http.Get(url)
	fmt.Println("The status code is", getRequest.StatusCode, http.StatusText(getRequest.StatusCode))

	if err != nil {
		return err
	}

	//close - this will be done at the end of the function
	// it's important to close the connection - we don't want the connection to leak
	defer getRequest.Body.Close()

	// read the body of the GET request
	rawData, err := ioutil.ReadAll(getRequest.Body)

	if err != nil {
		return err
	}

	top := TopPosts{}
	jsonErr := json.Unmarshal(rawData, &top)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	flag, err := ioutil.ReadDir(fmt.Sprintf("%v/dataset/hn/flag", dir))
	if err != nil {
		return err
	}

	neg, err := ioutil.ReadDir(fmt.Sprintf("%v/dataset/hn/neg", dir))
	if err != nil {
		return err
	}

	seen := make(map[string]bool)
	for _, file := range flag {
		if !file.IsDir() {
			seen[strings.TrimRight(file.Name(), ".txt")] = true
		}
	}

	for _, file := range neg {
		if !file.IsDir() {
			seen[strings.TrimRight(file.Name(), ".txt")] = true
		}
	}

	for _, v := range top.Hits {
		d := []byte(v.Title)
		if seen[v.ObjectID] == true {
			continue
		}
		path := fmt.Sprintf("%v/dataset/hn/unlabeled/%v.txt", dir, v.ObjectID)
		err := os.WriteFile(path, d, 0644)
		if err != nil {
			panic(err)
		}
	}
	return nil
}
