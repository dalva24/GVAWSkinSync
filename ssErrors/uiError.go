package ssErrors

type UiError struct {
	Title string
	Err   string
}

func NewUiError(title string, msg string) *UiError {
	return &UiError{
		Title: title,
		Err:   msg,
	}
}

func (err *UiError) Error() string {
	return err.Title + " - " + err.Err
}
