package views

// Alert represents the alert that shown inside temaplets
type Alert struct {
	Level   string
	Message string
}

// Alert levels
const (
	SuceessAlertLevel = "success"
	InfoAlertLeve     = "info"
	WarningAlertLevel = "warning"
	DangerAlertLevel  = "danger"
)

// Params represents the params that we can pass to the template
// when rendering it
type Params struct {
	Alert *Alert
	Data  interface{}
}