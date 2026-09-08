package domain

// HostingMode is how a Listello instance is hosted.
type HostingMode string

const (
	HostingModeLocal         HostingMode = "local"
	HostingModeStandaloneWeb HostingMode = "standalone-web"
)

// IsValid reports whether the hosting mode is a supported value.
func (m HostingMode) IsValid() bool {
	switch m {
	case HostingModeLocal, HostingModeStandaloneWeb:
		return true
	default:
		return false
	}
}
