package debuglog

import "sync"

var instance *Entries
var once sync.Once

// GetDebugLog returns the singleton instance of the debug log
func GetDebugLog() *Entries {
	once.Do(func() {
		instance = &Entries{}
	})
	return instance
}
