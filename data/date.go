package data

type DateDto struct {
	TotalResults int         `json:"totalResults"`
	SourceUrl    string      `json:"sourceUrl"`
	Events       []DateEvent `json:"events"'`
}

type DateEvent struct {
	Year    string `json:"year"`
	YearInt int    `json:"yearInt"`
	Event   string `json:"event"`
}
