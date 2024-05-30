package ssErrors

type DataCrcError struct {
	err string
}

func NewDataCrcError(msg string) *DataCrcError {
	return &DataCrcError{msg}
}

func (err *DataCrcError) Error() string {
	return err.err
}
