package feed

import (
	"context"
	"net/http"
	"io"
	"encoding/xml"
	"html"
	// "fmt"
)

// type RSSLink struct {
// 	Text 			string		`xml:",chardata"`
// 	Href 			string		`xml:"href,attr"`
// 	Rel  			string		`xml:"rel,attr"`
// 	Type 			string		`xml:"type,attr"`
// }

type RSSItem struct {
	// Text        	string		`xml:",chardata"`
	Title       	string		`xml:"title"`
	Link        	string		`xml:"link"`
	PubDate     	string		`xml:"pubDate"`
	// Guid        	string		`xml:"guid"`
	Description 	string		`xml:"description"`
}

type RSSChannel struct{
	// Text  			string		`xml:",chardata"`
	Title 			string		`xml:"title"`
	Link  			string `xml:"link"`//RSSLink		`xml:"link"`
	Description   	string		`xml:"description"`
	// Generator     	string		`xml:"generator"`
	// Language      	string		`xml:"language"`
	// LastBuildDate 	string		`xml:"lastBuildDate"`
	Items		  	[]RSSItem 	`xml:"item"`
}

type RSSFeed struct {
	// XMLName 		xml.Name 	`xml:"rss"`
	// Text    		string   	`xml:",chardata"`
	// Version 		string   	`xml:"version,attr"`
	// Atom    		string   	`xml:"atom,attr"`
	Channel			RSSChannel	`xml:"channel"`
}

// fetch feed from url, return filled RSSFeed struct
func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {

	// request
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "gator")


	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var xmlData RSSFeed
	err = xml.Unmarshal(body, &xmlData)
	if err != nil {
		return nil, err
	}

	html.UnescapeString(xmlData.Channel.Title)
	html.UnescapeString(xmlData.Channel.Description)
	items := xmlData.Channel.Items
	for i, item := range items {
		items[i].Title = html.UnescapeString(item.Title)
		items[i].Description = html.UnescapeString(item.Description)
	}

	// fmt.Printf("xmlData\n%v\n", string(body))
	return &xmlData, nil
}