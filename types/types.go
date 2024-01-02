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

type ApplicationPerspective struct {
	Id    string `json:"id"`
	Label string `json:"label"`
}

type ApplicationPerspectiveResponse struct {
	Items []ApplicationPerspective `json:"items"`
}

type CreateMaintenanceWindowRequest struct {
	Id         string   `json:"id"`
	Name       string   `json:"name"`
	Query      string   `json:"query"`
	Scheduling Schedule `json:"scheduling"`
}

type Schedule struct {
	Duration Duration `json:"duration"`
	Start    int64    `json:"start"`
	Type     string   `json:"type"`
}
type Duration struct {
	Amount int64  `json:"amount"`
	Unit   string `json:"unit"`
}
