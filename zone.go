package udnssdk

import "fmt"

// ZoneKey is the key for an UltraDNS zone
type ZoneKey string

// URI generates the URI for a task
func (z ZoneKey) URI() string {
	return fmt.Sprintf("zones/%s", z)
}
