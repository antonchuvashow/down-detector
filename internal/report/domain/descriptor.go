package report

import "net"

type PlatformType string

const (
	PlatformTypeUnknown PlatformType = "unknown"
	PlatformTypeLinux   PlatformType = "linux"
	PlatformTypeWindows PlatformType = "windows"
	PlatformTypeAndroid PlatformType = "android"
	PlatformTypeIOS     PlatformType = "ios"
)

type Descriptor struct {
	Source    SourceType
	Latitude  float64
	Longitude float64
	IP        net.IP
	Platform  PlatformType
}
