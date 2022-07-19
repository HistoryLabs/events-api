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

func FetchDate(c *gin.Context) {
	month, monthValid := c.GetQuery("month")
	day, dayValid := c.GetQuery("day")

	if monthValid == false || dayValid == false {
		c.IndentedJSON(400, gin.H{"message": "You must provide a month and day"})
		return
	}

	minYearStr := c.Query("minYear")
	maxYearStr := c.Query("maxYear")

	var minYear = -500
	var maxYear = time.Now().Year()

	if len(strings.TrimSpace(minYearStr)) > 0 {
		minYearInt, err := strconv.Atoi(minYearStr)
		if err != nil {
			c.IndentedJSON(400, gin.H{"message": "'minYear' must be a valid integer"})
			return
		}

		minYear = minYearInt
	}

	if len(strings.TrimSpace(maxYearStr)) > 0 {
		maxYearInt, err := strconv.Atoi(maxYearStr)
		if err != nil {
			c.IndentedJSON(400, gin.H{"message": "'maxYear' must be a valid integer"})
			return
		}

		maxYear = maxYearInt
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

	matches := utils.EventsPattern.FindAllStringSubmatch(string(wikiData), -1)

	cleanMatches := make([]data.Event, 0)

	for _, match := range matches {
		cleanMatch := utils.RemoveHTMLPattern.ReplaceAllString(match[0], "")
		year := strings.TrimSpace(strings.Split(cleanMatch, "&#8211;")[0])
		event := strings.TrimSpace(strings.Split(cleanMatch, "&#8211;")[1])

		var yearInt int

		if strings.Contains(year, "BC") {
			cleanYear := strings.TrimSpace(strings.ReplaceAll(year, "BC", ""))
			yearInt, err = strconv.Atoi(cleanYear)
			yearInt = yearInt * -1
			if err != nil {
				return
			}
		} else {
			yearInt, err = strconv.Atoi(year)
			if err != nil {
				return
			}
		}

		if yearInt >= minYear && yearInt <= maxYear {
			cleanMatches = append(cleanMatches, data.Event{
				Year:    year,
				YearInt: yearInt,
				Event:   utils.CleanPattern.ReplaceAllString(event, ""),
			})
		}
	}

	c.IndentedJSON(200, cleanMatches)
}
