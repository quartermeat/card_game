// Package 'errormgmt' Error management, probably the beginnings of what is really a debug log
// TODO: refactor to debug log
package debuglog

// IEntry interface for the custom error, may change to 'log entry'
// may allow a return of 'status' or 'error' based on type
type IEntry interface {
	GetMessage() string
}

// Entries is instantiated in main to be global list of errors/debug log statments made
type Entries []IEntry

// Entry is the entry structure
type Entry struct {
	Message string
}

// GetMessage returns the message in the entry struct
func (entry Entry) GetMessage() string {
	return entry.Message
}
