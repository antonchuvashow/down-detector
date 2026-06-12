package report

type ErrorType string

const (
	ErrorTypeUnknown      ErrorType = "unknown"
	ErrorTypeServer       ErrorType = "server"
	ErrorTypeNetwork      ErrorType = "network"
	ErrorTypeWebAccess    ErrorType = "web_access"
	ErrorTypeMobileAccess ErrorType = "mobile_access"
)

type SourceType string

const (
	SourceTypeUnknown   SourceType = "unknown"
	SourceTypeUser      SourceType = "user"
	SourceTypeInspector SourceType = "inspector"
)
