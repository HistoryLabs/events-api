package routes

import (
	"github.com/HistoryLabs/events-api/data"
	"github.com/HistoryLabs/events-api/utils"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func FetchYear(c *gin.Context) {
	yearStr, yearValid := c.GetQuery("year")
	onlyDated := c.Query("onlyDated")

	if yearValid == false {
		c.IndentedJSON(400, gin.H{"message": "You must provide a year to find events for"})
		return
	}

	yearInt, err := strconv.Atoi(yearStr)
	if err != nil {
		c.IndentedJSON(400, gin.H{"message": "You must provide a valid year (as an integer; negative for BC/BCE)"})
		return
	}

	if yearInt < -500 || yearInt > time.Now().Year() {
		c.IndentedJSON(400, gin.H{"message": "'year' must be greater than -500 and less than the current year"})
		return
	}

	var wikiYear string

	if yearInt < 0 {
		wikiYear = strconv.Itoa(int(math.Abs(float64(yearInt)))) + "_BC"
	} else {
		wikiYear = "AD_" + yearStr
	}

	resp, err := http.Get("https://en.wikipedia.org/w/api.php?action=parse&format=json&section=1&redirects=true&page=" + wikiYear)
	if err != nil {
		c.AbortWithStatus(500)
		return
	}

	wikiData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.AbortWithStatus(500)
		return
	}

	matches := utils.EventsPattern.FindAllStringSubmatch(string(wikiData), -1)
	cleanMatches := make([]data.YearEvent, 0)

	for _, match := range matches {
		cleanMatch := utils.RemoveHTMLPattern.ReplaceAllString(match[0], "")
		splitMatch := strings.SplitN(cleanMatch, "&#8211;", 2)
		if onlyDated != "true" || len(splitMatch) >= 2 {
			var date string
			var event string

			if len(splitMatch) >= 2 {
				date = strings.TrimSpace(splitMatch[0])
				event = strings.TrimSpace(splitMatch[1])
			} else {
				event = cleanMatch
			}

			cleanMatches = append(cleanMatches, data.YearEvent{
				Date:  date,
				Event: utils.FormatPattern.ReplaceAllString(utils.CleanPattern.ReplaceAllString(event, ""), "–"),
			})
		}
	}

	c.IndentedJSON(200, data.YearDto{
		TotalResults: len(cleanMatches),
		SourceUrl:    "https://en.wikipedia.org/wiki/" + wikiYear,
		Events:       cleanMatches,
	})
}
