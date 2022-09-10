package errormgmt

type IAemError interface {
	Error() string
}

type Errors []IAemError

type AemError struct {
	Message string
}

func (err AemError) Error() string {
	return err.Message
}
