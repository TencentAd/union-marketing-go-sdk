package httpx

import (
	"encoding/json"
	"net/http"
)

// AuthResponse 授权输出
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// ServeErrorResponse http处理错误时填充错误信息
func ServeErrorResponse(w http.ResponseWriter, err error) {
	resp := &Response{
		Code:    -1,
		Message: err.Error(),
	}

	ServerResponse(w, resp)
}

// ServerResponse 填充http response body
func ServerResponse(w http.ResponseWriter, resp *Response) {
	data, _ := json.Marshal(resp)
	_, _ = w.Write(data)
}
