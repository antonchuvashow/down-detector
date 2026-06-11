package domain

type ErrorType int

const (
	ErrorTypeUnknown ErrorType = iota
	ErrorTypeServer
	ErrorTypeNetwork
	ErrorTypeWebAccess
	ErrorTypeMobileAccess
)

type SourceType int

const (
	SourceTypeUnknown SourceType = iota
	SourceTypeUser
	SourceTypeInspector
)
