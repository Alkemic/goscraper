package goscraper

import (
	"io/ioutil"
	"net/http"

	"github.com/moovweb/gokogiri"
	"github.com/moovweb/gokogiri/html"
	"github.com/moovweb/gokogiri/xml"
	"github.com/moovweb/gokogiri/xpath"
	"github.com/pkg/errors"
)

// Field defines how data should be selected and process
type Field struct {
	// XPath selector
	Selector string
	// Callback is used to post process selected field
	Callback func(n *xml.Node) (interface{}, error)
}

// processField method fetches data from passed document
func (f Field) processField(d *html.HtmlDocument) (interface{}, error) {
	selector := xpath.Compile(f.Selector)
	result, err := d.Root().Search(selector)
	if err != nil {
		return nil, err
	}

	var node xml.Node
	if len(result) > 0 {
		node = result[0]
	} else {
		return "", nil
	}

	var value interface{}
	if f.Callback != nil {
		value, err = f.Callback(&node)
		if err != nil {
			return nil, errors.Wrap(err, "error processing callback")
		}
	} else {
		value = node.Content()
	}

	return value, nil
}

// Scrape type defines what to scrape and headers used in HTTP request, and
// allows to process data
type Scrape struct {
	// definitions of fields that will be scraped
	Fields map[string]Field
	// HTTP headers used in query
	Headers map[string]string
}

// ProcessURL method is used to fetch website contents and process it's data
func (s *Scrape) ProcessURL(url string) (map[string]interface{}, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "can't create new request")
	}
	for k, v := range s.Headers {
		req.Header.Add(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "can't create new request")
	}

	page, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "can't read response body")
	}

	doc, err := gokogiri.ParseHtml(page)
	if err != nil {
		return nil, errors.Wrap(err, "can't parse html")
	}
	defer doc.Free()

	return s.processDocument(doc)
}

// processDocument is used to process document fetched from website
func (s *Scrape) processDocument(d *html.HtmlDocument) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	var err error
	for key := range s.Fields {
		data[key], err = s.Fields[key].processField(d)
		if err != nil {
			return nil, errors.Wrapf(err, "can't process field %s", key)
		}
	}

	return data, nil
}
