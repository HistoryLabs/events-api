package routes

import (
	"github.com/HistoryLabs/events-api/data"
	"github.com/HistoryLabs/events-api/utils"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func FetchEvents(c *gin.Context) {
	month, monthValid := c.GetQuery("month")
	day, dayValid := c.GetQuery("day")

	if monthValid == false || dayValid == false {
		c.IndentedJSON(400, gin.H{"message": "You must provide a month and day"})
	}

	monthInt, err := strconv.Atoi(month)

	if err != nil {
		return
	}

	dateStr := time.Month(monthInt).String() + "_" + day

	resp, err := http.Get("https://en.wikipedia.org/w/api.php?action=parse&format=json&section=1&page=" + dateStr)

	if err != nil {
		return
	}

	wikiData, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return
	}

	dataStr := string(wikiData)

	matches := utils.EventsPattern.FindAllStringSubmatch(dataStr, -1)

	cleanMatches := make([]data.Event, 0)

	for _, match := range matches {
		cleanMatch := utils.RemoveHTMLPattern.ReplaceAllString(match[0], "")
		year := strings.Split(cleanMatch, "&#8211;")[0]
		event := strings.Split(cleanMatch, "&#8211;")[1]

		if err != nil {
			return
		}

		cleanMatches = append(cleanMatches, data.Event{
			Year:  strings.TrimSpace(year),
			Event: strings.TrimSpace(event),
		})
	}

	c.IndentedJSON(200, cleanMatches)
}
