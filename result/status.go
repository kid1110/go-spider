package result

type Status struct {
	Code int32
	Msg  string
	Data interface{}
}

const (
	RequestErrorException   = -400
	ParameterErrorException = -401
	ParaseUrlError          = -502
	NoTitleData             = -404
	Success                 = 1
)

var statusText = map[int32]string{
	RequestErrorException:   "请求方法不支持",
	ParameterErrorException: "参数错误",
	ParaseUrlError:          "爬取地址失败",
	NoTitleData:             "根据此title暂无数据",
	Success:                 "成功",
}

func GetStatusText(code int32) Status {
	var r = Status{}
	r.Code = code
	r.Msg = statusText[code]
	return r
}

func GetStautsString(code int32) string {
	a := statusText[code]
	return a
}
func GetStatusByMsg(code int32, message string) Status {
	var r = Status{}
	r.Code = code
	r.Msg = message
	return r
}
