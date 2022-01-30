package exception

import "github.com/solenovex/web-tutor/result"

type NewErrorException struct {
	res result.Status
}

func (e NewErrorException) Error() string {
	return e.res.Msg
}
func (e NewErrorException) Result() result.Status {
	return e.res
}
func CreateNewException(code int32, message string) NewErrorException {
	var exception = NewErrorException{}
	exception.res.Code = code
	exception.res.Msg = message
	return exception
}
