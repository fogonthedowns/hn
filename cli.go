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
		Commands: []*cli.Command{
			{
				Name:    "top",
				Aliases: []string{"t"},
				Usage:   "Load top HN Posts",
				Action: func(c *cli.Context) error {
					err := LoadTopPosts()
					if err != nil {
						return err
					}

					return nil
				},
			},
			{
				Name:    "timestamp",
				Aliases: []string{"d"},
				Usage:   "Load hn posts by unix start and stop ts",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "start",
						Aliases:  []string{"s"},
						Usage:    "Start Unix timestamp int",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "end",
						Aliases:  []string{"e"},
						Usage:    "End Unix timestamp int",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					err := LoadPosts(c.String("start"), c.String("end"))
					if err != nil {
						return err
					}

					return nil
				},
			},
			{
				Name:    "search",
				Aliases: []string{"d"},
				Usage:   "Search hn posts by search term",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "query",
						Aliases:  []string{"q"},
						Usage:    "Search term",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					err := SearchPosts(c.String("query"))
					if err != nil {
						return err
					}

					return nil
				},
			},

			{
				Name:    "flag",
				Aliases: []string{"c"},
				Usage:   "flag hn titles as political/policy/legal",
				Action: func(c *cli.Context) error {
					fmt.Println("Enter 1 to flag post. Do so if the post is related to politics, policy, regulations or legal:")
					err := process()
					if err != nil {
						return err
					}

					return nil
				},
			},
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
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	files, err := ioutil.ReadDir(fmt.Sprintf("%v/dataset/hn/unlabeled", dir))
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fmt.Println(file.Name(), file.IsDir())
		if file.IsDir() {
			break
		} else {
			content, err := ioutil.ReadFile(fmt.Sprintf("%v/dataset/hn/unlabeled/%v", dir, file.Name()))
			fmt.Println(string(content))
			reader := bufio.NewReader(os.Stdin)
			text, _ := reader.ReadString('\n')

			text = strings.Replace(text, "\n", "", -1)
			class, err := strconv.Atoi(text)

			if err != nil {
				return err
			}

			from := fmt.Sprintf("%v/dataset/hn/unlabeled/%v", dir, file.Name())
			switch class {
			case 0:
				err = os.Rename(from, fmt.Sprintf("%v/dataset/hn/neg/%v", dir, file.Name()))
				if err != nil {
					return err
				}
			case 1:
				err = os.Rename(from, fmt.Sprintf("%v/dataset/hn/flag/%v", dir, file.Name()))
				if err != nil {
					return err
				}
			}
		}
	}

	return nil

}
