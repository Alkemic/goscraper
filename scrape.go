package goscraper

import (
	"io/ioutil"
	"net/http"

	"github.com/moovweb/gokogiri"
	"github.com/moovweb/gokogiri/html"
	"github.com/moovweb/gokogiri/xml"
	"github.com/moovweb/gokogiri/xpath"
)

// Type Field defines how data should be selected and process
type Field struct {
	// XPath selector
	Selector string
	// Callback is used to post process selected field
	Callback func(*xml.Node) interface{}
}

// ProcessField method fetches data from passed document
func (f *Field) ProcessField(d *html.HtmlDocument) interface{} {
	var value interface{}
	var node xml.Node
	selector := xpath.Compile(f.Selector)
	result, _ := d.Root().Search(selector)

	if len(result) > 0 {
		node = result[0]
	} else {
		return ""
	}

	if f.Callback != nil {
		value = f.Callback(&node)
	} else {
		value = node.Content()
	}

	return value
}

// Scrape type defines what to scrape and headers used in HTTP request, and
// allows to process data
type Scrape struct {
	// definitions of fields that will be scraped
	Fields map[string]*Field
	// HTTP headers used in query
	Headers map[string]string
}

// ProcessURL method is used to fetch website contents and process it's data
func (s *Scrape) ProcessURL(url string) map[string]interface{} {
	client := &http.Client{}

	req, _ := http.NewRequest("GET", url, nil)
	for k, v := range s.Headers {
		req.Header.Add(k, v)
	}

	resp, _ := client.Do(req)
	page, _ := ioutil.ReadAll(resp.Body)
	doc, _ := gokogiri.ParseHtml(page)
	defer doc.Free()

	return s.ProcessDocument(doc)
}

// ProcessDocument is used to process document fetched from website
func (s *Scrape) ProcessDocument(d *html.HtmlDocument) map[string]interface{} {
	data := make(map[string]interface{})
	for key := range s.Fields {
		data[key] = s.Fields[key].ProcessField(d)
	}

	return data
}
