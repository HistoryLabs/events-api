package utils

import "regexp"

var EventsPattern *regexp.Regexp
var RemoveHTMLPattern *regexp.Regexp
var CleanPattern *regexp.Regexp

func init() {
	EventsPattern = regexp.MustCompile(`<li>.*?<\/li>`)
	RemoveHTMLPattern = regexp.MustCompile(`<[^>]*>`)
	CleanPattern = regexp.MustCompile(`&#91;[1-9]*\d*\d&#93;`)
}
