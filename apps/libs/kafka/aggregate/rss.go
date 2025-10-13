package aggregate

type RSSAggregate struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Link        string `json:"link"`
	Active      bool   `json:"active"`
}
