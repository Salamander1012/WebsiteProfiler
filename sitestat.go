package main

import (
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "sitestat",
		Usage: "A cli to test http requests to a website",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "url",
				Usage:    "(REQUIRED) The url for an http request",
				Required: true,
			},
			&cli.IntFlag{
				Name:  "profile",
				Value: 0,
				Usage: "The number of requests used to profile a url",
			},
		},
		Action: func(c *cli.Context) error {
			url, err := url.Parse(getURLString(c.String("url")))
			if err != nil {
				log.Fatal(err)
			}

			profileSize := c.Int("profile")
			if profileSize > 0 {
				handleProfile(*url, profileSize)
			} else {
				handlePrintResponse(*url)
			}

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func getURLString(url string) string {
	var result string = ""
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		result = url
	} else {
		result = "https://" + url
	}
	if !strings.HasSuffix(result, "/") {
		result = result + "/"
	}
	return result
}
