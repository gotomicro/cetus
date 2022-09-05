package kerror

import "fmt"

// 1	00	01
// 服务级别	模块	具体错误
var (
	OK                  = &KError{Code: 0, Message: "OK"}
	InternalServerError = &KError{Code: 10001, Message: "Internal server error."}
	ErrBind             = &KError{Code: 10002, Message: "Error occurred while binding the request body to the struct."}
)

// KError 定义错误码
type KError struct {
	Code    int
	Message string
}

func (err KError) Error() string {
	return err.Message
}

// Err 定义错误
type Err struct {
	Code     int    // 错误码
	Message  string // 展示给用户看的
	ErrInner error  // 保存内部错误信息
}

func (err *Err) Error() string {
	return fmt.Sprintf("Err - code: %d, message: %s, error: %s", err.Code, err.Message, err.ErrInner)
}

// New 使用 错误码 和 error 创建新的 错误
func New(errno *KError, err error) *Err {
	return &Err{
		Code:     errno.Code,
		Message:  errno.Message,
		ErrInner: err,
	}
}

// DecodeErr 解码错误, 获取 Code 和 Message
func DecodeErr(err error) (int, string) {
	if err == nil {
		return OK.Code, OK.Message
	}
	switch typed := err.(type) {
	case *Err:
		if typed.Code == ErrBind.Code {
			typed.Message = typed.Message + " 具体是 " + typed.ErrInner.Error()
		}
		return typed.Code, typed.Message
	case *KError:
		return typed.Code, typed.Message
	default:
	}
	return InternalServerError.Code, err.Error()
}
