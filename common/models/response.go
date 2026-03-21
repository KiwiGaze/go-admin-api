package models

type Response struct {
	Code      int         `json:"code" example:"200"`
	Data      interface{} `json:"data"`
	Msg       string      `json:"msg"`
	RequestId string      `json:"requestId"`
}

type Page struct {
	List      interface{} `json:"list"`
	PageIndex int         `json:"pageIndex"`
	PageSize  int         `json:"pageSize"`
	Count     int         `json:"count"`
}

// Return OK
func (r *Response) ReturnOK() *Response {
	r.Code = 200
	return r
}

// Return Error
func (r *Response) ReturnError(code int) *Response {
	r.Code = code
	return r
}