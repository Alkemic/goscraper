package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/Alkemic/goscraper"
	"github.com/moovweb/gokogiri/xml"
)

func main() {
	scrape := &goscraper.Scrape{
		Fields: map[string]*goscraper.Field{
			"title": &goscraper.Field{Selector: "//*[@id=\"productTitle\"]"},
			"price": &goscraper.Field{
				Selector: "//*[@id=\"unqualifiedBuyBox\"]/div/div[1]/span",
				Callback: func(n *xml.Node) interface{} {
					f, err := strconv.ParseFloat((*n).Content()[2:], 64)
					if err != nil {
						log.Println(err)
						return nil
					}
					return f
				},
			},
			"image": &goscraper.Field{
				Selector: "//*[@id=\"landingImage\"]",
				Callback: func(n *xml.Node) interface{} {
					var f map[string]interface{}
					// this part quite naive, and there should be added some
					// validation, but hey! it's only an example ;-)
					attr := (*n).Attribute("data-a-dynamic-image").String()
					err := json.Unmarshal([]byte(attr), &f)
					if err != nil {
						log.Println(err)
						return nil
					}

					for k := range f {
						return k
					}

					return nil
				},
			},
		},
		Headers: map[string]string{
			"accept-encoding": "gzip, deflate, sdch",
			"cache-control":   "no-cache",
			"accept-language": "pl-PL,pl;q=0.8,en-US;q=0.6,en;q=0.4",
		},
	}

	fmt.Println(scrape.ProcessURL(
		"http://www.amazon.co.uk/Lenovo-ThinkPad-T430s-i5-3320M-N1RLQGE/" +
			"dp/B00AQPBA8U/ref=sr_1_1?ie=UTF8&qid=1455007797&" +
			"sr=8-1&keywords=t430s",
	))
}
