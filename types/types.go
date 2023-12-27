package types

type Event struct {
	EventId     string `json:"eventId"`
	Start       int64  `json:"start"`
	End         int64  `json:"end"`
	Type        string `json:"type"`
	State       string `json:"state"`
	Problem     string `json:"problem"`
	Detail      string `json:"detail"`
	Severity    int    `json:"severity"`
	EntityName  string `json:"entityName"`
	EntityLabel string `json:"entityLabel"`
	EntityType  string `json:"entityType"`
}
