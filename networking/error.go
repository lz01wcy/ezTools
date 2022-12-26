package networking

type Error struct {
	IsSettingError  bool
	IsCheckingError bool
	IsLoadingError  bool
	errorDesc       string
}

func (e *Error) Error() string {
	return e.errorDesc
}
