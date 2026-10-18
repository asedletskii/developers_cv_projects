package domain

type TrafficType string

const (
	TrafficHuman   TrafficType = "human"
	TrafficBot     TrafficType = "bot"
	TrafficOwner   TrafficType = "owner"
	TrafficSuspect TrafficType = "suspect"
)

func (v TrafficType) IsValid() bool {
	switch v {
	case TrafficHuman, TrafficBot, TrafficOwner, TrafficSuspect:
		return true
	default:
		return false
	}
}

type DeviceType string

const (
	DeviceDesktop DeviceType = "desktop"
	DeviceMobile  DeviceType = "mobile"
	DeviceTablet  DeviceType = "tablet"
	DeviceUnknown DeviceType = "unknown"
)

func (v DeviceType) IsValid() bool {
	switch v {
	case DeviceDesktop, DeviceMobile, DeviceTablet, DeviceUnknown:
		return true
	default:
		return false
	}
}

type BrowserType string

const (
	BrowserChrome  BrowserType = "chrome"
	BrowserSafari  BrowserType = "safari"
	BrowserFirefox BrowserType = "firefox"
	BrowserEdge    BrowserType = "edge"
	BrowserYandex  BrowserType = "yandex"
	BrowserOther   BrowserType = "other"
)

func (v BrowserType) IsValid() bool {
	switch v {
	case BrowserChrome, BrowserSafari, BrowserFirefox, BrowserEdge, BrowserYandex, BrowserOther:
		return true
	default:
		return false
	}
}

type OSType string

const (
	OSWindows OSType = "windows"
	OSMacOS   OSType = "macos"
	OSLinux   OSType = "linux"
	OSIOS     OSType = "ios"
	OSAndroid OSType = "android"
	OSOther   OSType = "other"
)

func (v OSType) IsValid() bool {
	switch v {
	case OSWindows, OSMacOS, OSLinux, OSIOS, OSAndroid, OSOther:
		return true
	default:
		return false
	}
}
