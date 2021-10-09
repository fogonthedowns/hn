package main

import (
	//	"fmt"

	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/urfave/cli/v2" // imports as package "cli"
)

func main() {

	app := &cli.App{
		Name:  "greet",
		Usage: "fight the loneliness!",
		Action: func(c *cli.Context) error {
			err := processCSV()
			if err != nil {
				return err
			}

			fmt.Println("Hello friend!")
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

type Row struct {
	Title    string
	ObjectID string
	Class    int
}

func processCSV() error {
	file, err := os.Open("/home/jzollars/hack/data-sets/hn/file.tsv")
	if err != nil {
		return err
	}

	defer file.Close()
	reader := csv.NewReader(file)
	reader.Comma = '\t' // tsv
	reader.FieldsPerRecord = -1
	data, err := reader.ReadAll()

	if err != nil {
		return err
	}

	var rows []Row

	for _, v := range data {
		row := Row{Title: v[0], ObjectID: v[1]}
		fmt.Println(v[0])
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')

		text = strings.Replace(text, "\n", "", -1)
		row.Class, err = strconv.Atoi(text)
		if err != nil {
			return err
		}
		rows = append(rows, row)
	}
	fmt.Println(rows)

	return nil

}
