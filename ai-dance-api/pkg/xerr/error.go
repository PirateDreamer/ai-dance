package xerr

type BizError struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

func (e *BizError) Error() string {
	return e.Msg
}

func NewBizErr(code, msg string) *BizError {
	return &BizError{Code: code, Msg: msg}
}

func NewCommBizErr(msg string) *BizError {
	return &BizError{Code: "1", Msg: msg}
}
