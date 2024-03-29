package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	link "github.com/jun-hf/link-extractor/link"
)

type SiteMap struct {
	href map[string]string
	domainName string // https://example.com/
	depth int
	visitedDepth int
}

type url struct {
	Loc string `xml:"loc"`
}

type urlSet struct {
	Url []url `xml:"url"`
	Xmlns string `xml:"xmlns,attr"`
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
func (sm *SiteMap)buildSiteTree(url string) error {
	sm.visitedDepth++
	if sm.depth < sm.visitedDepth {
		return nil
	}
	html, err := getHtml(url)
	if err != nil {
		log.Print(err)
		return err
	}
	links, _ := link.Parser(strings.NewReader(html))
	

	for _, l := range links {
		if strings.HasPrefix(l.Href, "#") {
			continue
		}
		if (strings.HasPrefix(l.Href, "/") && l.Href != "/") || (strings.HasPrefix(l.Href, "?") && l.Href != "?"){
			l.Href = sm.domainName + l.Href
		}
		if !strings.HasPrefix(l.Href, sm.domainName) {
			continue
		}
		if _, ok := sm.href[l.Href]; ok {
			continue
		}
		sm.href[l.Href] = "member"
		if err := sm.buildSiteTree(l.Href); err != nil {
			fmt.Println(err)
		}
	}
	return nil
}
func (sm *SiteMap) createXML() ([]byte, error) {
	urlList := make([]url,0)
	for hr, _ := range sm.href {
		u := url{Loc:hr}
		urlList = append(urlList, u)
	}
	urlS := urlSet{urlList, "http://www.sitemaps.org/schemas/sitemap/0.9"}
	output, err := xml.MarshalIndent(urlS, " ", "  ")
	if err != nil {
		return nil, err
	}
	return []byte(xml.Header + string(output)), nil

}
// // Takes in a url and build a site map return in xml format
func (sm *SiteMap)BuildSiteMap() ([]byte, error){// (xml, error)
	sm.buildSiteTree(sm.domainName)
	return sm.createXML()
}

func main() {
	sm := SiteMap{domainName: "https://pkg.go.dev/net/http", href: make(map[string]string), depth: 3}
	output, err := sm.BuildSiteMap()
	if err != nil {
		log.Fatal(err)
	}
	os.Stdout.Write(output)

}