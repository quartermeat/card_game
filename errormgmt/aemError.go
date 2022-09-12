// Package 'errormgmt' Error management, probably the beginnings of what is really a debug log
// TODO: refactor to debug log
package errormgmt

// IAemError interface for the custom error, may change to 'log entry'
// may allow a return of 'status' or 'error' based on type
type IAemError interface {
	Error() string
}

// Errors is instantiated in main to be global list of errors/debug log statments made
type Errors []IAemError

// AemError is the entry structure
type AemError struct {
	Message string
}

// Error returns the message in the entry struct
func (err AemError) Error() string {
	return err.Message
}
