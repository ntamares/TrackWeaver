package parser

import (
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/html"
)

// TODO remove?
// func FindStringInHTML(r io.Reader, search string) bool {
// 	doc, err := html.Parse(r)
// 	if err != nil {
// 		return false
// 	}
// 	return searchNode(doc, search)
// }

func searchNode(n *html.Node, search string) bool {
	if n.Type == html.TextNode && strings.Contains(n.Data, search) {
		return true
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if searchNode(c, search) {
			return true
		}
	}
	return false
}

func GetDateFromHTML(r io.Reader) (string, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return "", err
	}

	var date string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "span" {
			for _, attr := range n.Attr {
				if attr.Key == "class" && attr.Val == "minidim" {
					if n.FirstChild != nil && n.FirstChild.Type == html.TextNode {
						date = n.FirstChild.Data
					}
					return
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	if date == "" {
		return "", fmt.Errorf("date not found")
	}
	return date, nil
}

func FindStartingPoint(r io.Reader) (bool, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return false, err
	}

	found := searchNode(doc, "Big Sonic Chill for")

	return found, nil
}

// TODO
// func ProcessSongs(r io.Reader) ([]Song, []Song, error){
// 	// TODO
// 	// find starting point
// 	// loop through divs after
// 	// ignore empty divs, time slots
// 	// if song, add song to raw data song list (CSV?)
// 	// call CheckAmbiguity
// }

// func CheckAmbiguity(str string) bool {
// 	// check edge cases of whether artist may be ambiguous (multiple hyphens)
// 	// if all good, parse at last hyphen delimiter and add to master CSV for song list
// 	// if not, add to secondary CSV to manually review after the fact
// }
