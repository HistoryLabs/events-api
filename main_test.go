package main

import (
	"encoding/json"
	"github.com/HistoryLabs/events-api/data"
	"github.com/HistoryLabs/events-api/server/routes"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http/httptest"
	"strconv"
	"testing"
)

func initRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode) // to disable logging
	router := gin.New()
	return router
}

func TestDate(t *testing.T) {
	router := initRouter()
	router.GET("/date", routes.FetchDate)

	minYear := 400
	maxYear := 1900

	requestUrl := "/date?month=3&day=15&minYear=" + strconv.Itoa(minYear) + "&maxYear=" + strconv.Itoa(maxYear)

	req := httptest.NewRequest("GET", requestUrl, nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Result().StatusCode != 200 {
		t.Errorf("Expected status code 200 but got: %s\nRequest: GET %s", rr.Result().Status, requestUrl)
		return
	}

	body, _ := ioutil.ReadAll(rr.Body)
	var dateData data.DateDto
	json.Unmarshal(body, &dateData)

	if dateData.TotalResults <= 0 {
		t.Errorf("Expected at least one Event but got 0\nRequest: GET %s", requestUrl)
		return
	}

	for _, event := range dateData.Events {
		if event.YearInt > maxYear {
			t.Errorf("Expected maximum Year to be %s but found Event with Year: %s", strconv.Itoa(maxYear), event.Year)
		}

		if event.YearInt < minYear {
			t.Errorf("Expected minimum Year to be %s but found Event with Year: %s", strconv.Itoa(minYear), event.Year)
		}
	}
}

func TestYear(t *testing.T) {
	router := initRouter()
	router.GET("/year", routes.FetchYear)

	req := httptest.NewRequest("GET", "/year?year=1776&onlyDated=true", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Result().StatusCode != 200 {
		t.Errorf("Expected status code 200 but got: %s", rr.Result().Status)
		return
	}

	body, _ := ioutil.ReadAll(rr.Body)
	var yearData data.YearDto
	json.Unmarshal(body, &yearData)

	if yearData.TotalResults <= 0 {
		t.Errorf("Expected at least one Event but got 0")
		return
	}

	for _, event := range yearData.Events {
		if event.Date == "" {
			prettyEvent, _ := json.MarshalIndent(event, "", "\t")
			t.Errorf("Expected all Events to have non-empty Date, but found Event without Date:\n%s", prettyEvent)
		}
	}
}
