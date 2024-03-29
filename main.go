package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	// "os"
	"strings"
	link "github.com/jun-hf/link-extractor/link"
)

type SiteMap struct {
	href map[string]interface{}
	domainName string // https://example.com/
}

func getHtml(url string) (string, error){
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	html, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	return string(html), nil
}


// // takes in the url and return back a string slice which containe the
func (sm *SiteMap)buildSiteTree(url string) ([]link.Link, error) {
	requestUrl := url
	if !strings.HasPrefix(requestUrl, sm.domainName) {
		requestUrl = sm.domainName + requestUrl
	}
	html, err := getHtml(requestUrl)
	if err != nil {
		log.Print(err)
		return make([]link.Link, 0), err
	}
	links, _ := link.Parser(strings.NewReader(html))
	return links, nil
}

// // Takes in a url and build a site map return in xml format
func (sm *SiteMap)BuildSiteMap() { // (xml, error)
	link, _ := sm.buildSiteTree(sm.domainName)
	for _, l:= range link {
		fmt.Printf(l.Href+ "\n")
	}
}

func main() {
	// resp, err := http.Get("https://pkg.go.dev//net/http")
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
	// defer resp.Body.Close()
	// html, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	sm := SiteMap{domainName: "https://pkg.go.dev//net/http#Head"}
	sm.BuildSiteMap()

	// fmt.Println(string(html))
}