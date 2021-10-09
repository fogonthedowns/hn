package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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

func main() {

	url := "https://hn.algolia.com/api/v1/search?tags=front_page"

	fmt.Println("Now retrieving Underground line status, please wait...")
	// two variables (response and error) which stores the response from e GET request
	getRequest, err := http.Get(url)
	fmt.Println("The status code is", getRequest.StatusCode, http.StatusText(getRequest.StatusCode))

	if err != nil {
		fmt.Println("Error!")
		fmt.Println(err)
	}

	//close - this will be done at the end of the function
	// it's important to close the connection - we don't want the connection to leak
	defer getRequest.Body.Close()

	// read the body of the GET request
	rawData, err := ioutil.ReadAll(getRequest.Body)

	if err != nil {
		fmt.Println("Error!")
		fmt.Println(err)
	}

	top := TopPosts{}
	jsonErr := json.Unmarshal(rawData, &top)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	file, err := os.Create("/home/jzollars/hack/data-sets/hn/index/file.tsv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	writer.WriteString("Title	ObjectID	Class\n")
	for _, v := range top.Hits {
		d := []byte(v.Title)
		path := fmt.Sprintf("/home/jzollars/hack/data-sets/hn/unlabeled/%v.txt", v.ObjectID)
		err := os.WriteFile(path, d, 0644)
		if err != nil {
			panic(err)
		}
		row := fmt.Sprintf("%v	%v	0\n", v.Title, v.ObjectID)
		writer.WriteString(row)
		if err != nil {
			panic(err)
		}

	}
	writer.Flush()
	file.Close()
}
