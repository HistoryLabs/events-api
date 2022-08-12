package routes

import (
	"github.com/HistoryLabs/events-api/data"
	"github.com/HistoryLabs/events-api/utils"
	"github.com/gin-gonic/gin"
	"html"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func FetchDate(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")

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
			c.IndentedJSON(400, gin.H{"message": "Parameter 'minYear' must be a valid integer"})
			return
		}

		minYear = minYearInt
	}

	if len(strings.TrimSpace(maxYearStr)) > 0 {
		maxYearInt, err := strconv.Atoi(maxYearStr)
		if err != nil {
			c.IndentedJSON(400, gin.H{"message": "Parameter 'maxYear' must be a valid integer"})
			return
		}

		maxYear = maxYearInt
	}

	monthInt, err := strconv.Atoi(month)
	if err != nil {
		c.AbortWithStatus(500)
		return
	}

	if monthInt > 12 || monthInt < 1 {
		c.IndentedJSON(400, gin.H{"message": "Parameter 'month' must be a valid integer from 1 to 12"})
		return
	}

	dayInt, err := strconv.Atoi(day)
	if err != nil {
		c.AbortWithStatus(500)
		return
	}

	if dayInt > 31 || dayInt < 0 {
		c.IndentedJSON(400, gin.H{"message": "Parameter 'day' must be a valid integer from 1 to 31"})
		return
	}

	dateStr := time.Month(monthInt).String() + "_" + day

	resp, err := http.Get("https://en.wikipedia.org/w/api.php?action=parse&format=json&section=1&page=" + dateStr)
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

	cleanMatches := make([]data.DateEvent, 0)

	for _, match := range matches {
		cleanMatch := utils.RemoveHTMLPattern.ReplaceAllString(match[0], "")
		splitMatch := strings.SplitN(cleanMatch, "&#8211;", 2)

		if len(splitMatch) == 2 {
			year := strings.TrimSpace(splitMatch[0])
			event := strings.TrimSpace(splitMatch[1])

			event, err = strconv.Unquote(`"` + event + `"`)
			if err != nil {
				continue
			}

			var yearInt int

			if strings.Contains(year, "BC") {
				cleanYear := strings.TrimSpace(strings.ReplaceAll(year, "BC", ""))
				yearInt, err = strconv.Atoi(cleanYear)
				yearInt = yearInt * -1
				if err != nil {
					continue
				}
			} else {
				yearInt, err = strconv.Atoi(year)
				if err != nil {
					continue
				}
			}

			if yearInt >= minYear && yearInt <= maxYear {
				cleanMatches = append(cleanMatches, data.DateEvent{
					Year:    year,
					YearInt: yearInt,
					Content: html.UnescapeString(utils.FormatPattern.ReplaceAllString(utils.CleanPattern.ReplaceAllString(event, ""), "–")),
				})
			}
		}
	}

	c.IndentedJSON(200, data.DateDto{
		TotalResults: len(cleanMatches),
		SourceUrl:    "https://en.wikipedia.org/wiki/" + dateStr,
		Events:       cleanMatches,
	})
}
