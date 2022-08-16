package rule

import "fmt"

type Severity int

const (
	Error Severity = iota
	Warning
	Info
)

func (d Severity) String() string {
	return [...]string{"ERROR", "WARNING", "INFO"}[d]
}

func NewCoreViolation(message string) CoreViolation {
	return CoreViolation{
		message: message,
	}
}

type CoreViolation struct {
	message string
}

func (v CoreViolation) String() string {
	return v.message
}

func NewViolation(message string, severity Severity) Violation {
	return Violation{
		message:  message,
		severity: severity,
	}
}

type Violation struct {
	message  string
	severity Severity
}

func (v Violation) String() string {
	return fmt.Sprintf("[%s] %s", v.severity.String(), v.message)
}

func (v Violation) Message() string {
	return v.message
}

func (v Violation) Severity() string {
	return v.severity.String()
}
