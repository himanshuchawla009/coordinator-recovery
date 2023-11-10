package types

type AlertMessage string

const (
	AlertMessageJobFailed AlertMessage = "*Job Failed*\n\nRequest ID: %s\nRegion: %s\nError: %s"
)

func (m AlertMessage) String() string {
	return string(m)
}
