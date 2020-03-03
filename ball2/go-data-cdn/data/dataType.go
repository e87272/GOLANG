package data

type Exception struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type SysLog struct {
	UserUuid string `json:"userUuid"`
	Code     string `json:"code"`
	Message  string `json:"message"`
	Stamp    string `json:"stamp"`
}
type SysErrorLog struct {
	UserUuid string `json:"userUuid"`
	Code     string `json:"code"`
	Message  string `json:"message"`
	Error    string `json:"error"`
	Stamp    string `json:"stamp"`
}
