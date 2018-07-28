package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/Alkemic/goscraper"
	"github.com/moovweb/gokogiri/xml"
)

func main() {
	processPrice := func(n *xml.Node) (interface{}, error) {
		f, err := strconv.ParseFloat((*n).Content()[2:], 64)
		if err != nil {
			return nil, err
		}
		return f, nil
	}
	processImage := func(n *xml.Node) (interface{}, error) {
		var f map[string]interface{}
		// this part quite naive, and there should be added some
		// validation, but hey! it's only an example ;-)
		attr := (*n).Attribute("data-a-dynamic-image").String()
		err := json.Unmarshal([]byte(attr), &f)
		if err != nil {
			return nil, err
		}

		for k := range f {
			return k, nil
		}

		return nil, nil
	}
	processTitle := func(n *xml.Node) (interface{}, error) {
		return strings.TrimSpace((*n).Content()), nil
	}

	scrape := &goscraper.Scrape{
		Fields: map[string]goscraper.Field{
			"title": {Selector: "//*[@id=\"productTitle\"]", Callback: processTitle},
			"price": {Selector: "//*[@id=\"unqualifiedBuyBox\"]/div/div[1]/span", Callback: processPrice},
			"image": {Selector: "//*[@id=\"landingImage\"]", Callback: processImage},
		},
		Headers: map[string]string{
			"accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
			"accept-language":           "pl-PL,pl;q=0.9,en-US;q=0.8,en;q=0.7",
			"cache-control":             "no-cache",
			"cookie":                    "session-id=259-7520693-7366814; session-id-time=2082787201l; ubid-acbuk=259-2899942-1087152",
			"dnt":                       "1",
			"pragma":                    "no-cache",
			"upgrade-insecure-requests": "1",
			"user-agent":                "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.87 Safari/537.36",
		},
	}

	url := "http://www.amazon.co.uk/Lenovo-ThinkPad-T430s-i5-3320M-N1RLQGE/dp/B00AQPBA8U/ref=sr_1_1?ie=UTF8&qid=1455007797&sr=8-1&keywords=t430s"
	scraped, err := scrape.ProcessURL(url)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(scraped)
}
