package response

import (
	"errors"
	"fmt"
	"net/http"
)

type (
	iCode interface {
		// Code returns the integer number of current error code.
		Code() int
		// Message returns the brief message for current error code.
		Message() string
		// Detail returns the detailed information of current error code,
		// which is mainly designed as an extension field for error code.
		Detail() interface{}
		Error() string
		Status() int
	}

	// IUnwrap is the interface for Unwrap feature.
	iUnwrap interface {
		Error() string
		Unwrap() error
	}

	Code struct {
		code    int         // Error code, usually an integer.
		message string      // Brief message for this error code.
		detail  interface{} // As type of interface, it is mainly designed as an extension field for error code.
		status  int         // Http status
	}
)

func DefaultSuccess() *Code {
	return CodeDefaultSuccess
}

func DefaultFailure() *Code {
	return CodeDefaultFailure
}

func NewCode(bizCode int, message string, httpStatus int, detail ...interface{}) *Code {
	var dtl interface{}
	if len(detail) > 0 {
		dtl = detail[0]
	}
	return &Code{
		code:    bizCode,
		message: message,
		detail:  dtl,
		status:  httpStatus,
	}
}

func CodeFromError(err error) *Code {
	if err == nil {
		return nil
	}
	var e iCode
	if errors.As(err, &e) {
		return &Code{
			code:    e.Code(),
			message: e.Message(),
			detail:  e.Detail(),
			status:  e.Status(),
		}
	}
	var unwrap iUnwrap
	if errors.As(err, &unwrap) {
		return CodeFromError(unwrap.Unwrap())
	}
	return DefaultFailure().WithDetail(err.Error())
}

func (c *Code) Code() int {
	return c.code
}

func (c *Code) Message() string {
	return c.message
}

func (c *Code) Detail() interface{} {
	return c.detail
}

func (c *Code) Error() string {
	msg := c.message
	if msg == "" {
		msg = http.StatusText(c.code)
	}

	if c.detail != nil {
		if msg == "" {
			msg = fmt.Sprintf("%v", c.detail)
		} else {
			msg = fmt.Sprintf("%s: %v", msg, c.detail)
		}
	}
	return msg
}

func (c *Code) Status() int {
	return c.status
}

func (c *Code) Clone() *Code {
	return &Code{
		code:    c.code,
		message: c.message,
		detail:  c.detail,
		status:  c.status,
	}
}

func (c *Code) WithHttpStatus(status int) *Code {
	cc := c.Clone()
	cc.status = status
	return cc
}

func (c *Code) WithCode(code int) *Code {
	cc := c.Clone()
	cc.code = code
	return cc
}

func (c *Code) WithMessage(message string) *Code {
	cc := c.Clone()
	cc.message = message
	return cc
}

func (c *Code) WithDetail(detail interface{}) *Code {
	cc := c.Clone()
	cc.detail = detail
	return cc
}
