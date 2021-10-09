package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/urfave/cli/v2" // imports as package "cli"
)

func main() {

	app := &cli.App{
		Name:  "organize",
		Usage: "categorize hn titles as political/policy or not",
		Action: func(c *cli.Context) error {
			err := process()
			if err != nil {
				return err
			}

			fmt.Println("Enter 1 if post is related to politics, policy, regulations or legal:")
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

func process() error {
	files, err := ioutil.ReadDir("/home/jzollars/hack/data-sets/hn/unlabeled")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fmt.Println(file.Name(), file.IsDir())
		if file.IsDir() {
			break
		} else {
			content, err := ioutil.ReadFile(fmt.Sprintf("/home/jzollars/hack/data-sets/hn/unlabeled/%v", file.Name()))
			fmt.Println(string(content))
			reader := bufio.NewReader(os.Stdin)
			text, _ := reader.ReadString('\n')

			text = strings.Replace(text, "\n", "", -1)
			class, err := strconv.Atoi(text)

			if err != nil {
				return err
			}

			from := fmt.Sprintf("/home/jzollars/hack/data-sets/hn/unlabeled/%v", file.Name())
			switch class {
			case 0:
				fmt.Println(fmt.Sprintf("/home/jzollars/hack/data-sets/hn/neg/%v", file.Name()))
				err = os.Rename(from, fmt.Sprintf("/home/jzollars/hack/data-sets/hn/neg/%v", file.Name()))
				if err != nil {
					return err
				}
			case 1:
				fmt.Println(fmt.Sprintf("/home/jzollars/hack/data-sets/hn/neg/%v", file.Name()))
				err = os.Rename(from, fmt.Sprintf("/home/jzollars/hack/data-sets/hn/pos/%v", file.Name()))
				if err != nil {
					return err
				}
			}
		}
	}

	return nil

}
