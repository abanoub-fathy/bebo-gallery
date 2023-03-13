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
	User  interface{}
}

// SetAlert is used to set alert in the the params
// that passed to templates
//
// it can return public error message or
// the generic error messge dependant on the
// kind of the error it receives.
func (p *Params) SetAlert(err error) {
	if pErr, ok := err.(PublicError); ok {
		// if the error is public error
		p.Alert = NewAlert(AlertLevelError, pErr.PublicErrMsg())
	} else {
		p.Alert = NewAlert(AlertLevelError, ErrMsgGeneric)
	}
}

// SetAlertWithErrMsg is used to set alert with
// specific err message.
//
// the error message will be shown to the end user
func (p *Params) SetAlertWithErrMsg(errMsg string) {
	p.Alert = NewAlert(AlertLevelError, errMsg)
}

// PublicError represnets error that can be
// used to be shown to end user
type PublicError interface {
	error

	// PublicErrMsg is the method we use to return
	// the public error message as string value
	PublicErrMsg() string
}
