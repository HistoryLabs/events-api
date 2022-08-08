package data

type YearDto struct {
	TotalResults int         `json:"totalResults"`
	SourceUrl    string      `json:"sourceUrl"`
	Events       []YearEvent `json:"events"`
}

type YearEvent struct {
	Date    string `json:"date"`
	Content string `json:"content"`
}
