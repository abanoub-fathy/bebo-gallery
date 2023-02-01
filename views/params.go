package views

// Alert represents the alert that shown inside temaplets
type Alert struct {
	Level   string
	Message string
}

// NewAlert is used to create a new alert
// with level and message
func NewAlert(level string, message string) *Alert {
	return &Alert{
		Level:   level,
		Message: message,
	}
}

const (
	AlertLevelSuccess = "success"
	AlertLevelInfo    = "info"
	AlertLevelWarning = "warning"
	AlertLevelError   = "danger"

	// ErrMsgGeneric is usually used to represents
	// unknown or unfiltered error message
	ErrMsgGeneric = "something went wrong. please try again later or contact support"
)

// Params represents the params that we can pass to the template
// when rendering it
type Params struct {
	Alert *Alert
	Data  interface{}
}
