package errors

import (
	"fmt"
	"net/http"
)

// AppError 应用错误
type AppError struct {
	Code       int    // 业务错误码
	Message    string // 错误信息
	HTTPStatus int    // HTTP状态码
	Err        error  // 原始错误
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// 业务错误码定义
const (
	CodeSuccess           = 0
	CodeInvalidParameter  = 10001
	CodeAIServiceError    = 10002
	CodeRateLimitError    = 10003
	CodeTimeoutError      = 10004
	CodeInternalError     = 10005
	CodeProviderNotFound  = 10006
	CodeConfigError       = 10007

	// 认证相关错误码 20xxx
	CodeUserExists       = 20001 // 用户已存在
	CodeInvalidPassword  = 20002 // 密码错误
	CodeUserNotFound     = 20003 // 用户不存在
	CodeTokenInvalid     = 20004 // Token无效
	CodeTokenExpired     = 20005 // Token过期
	CodePasswordWeak     = 20006 // 密码强度不够
	CodeAccountBanned    = 20007 // 账号已被封禁
	CodeUnauthorized     = 20008 // 未授权
	CodeForbidden        = 20009 // 禁止访问
)

// NewInvalidParameterError 参数错误
func NewInvalidParameterError(message string) *AppError {
	return &AppError{
		Code:       CodeInvalidParameter,
		Message:    message,
		HTTPStatus: http.StatusBadRequest,
	}
}

// NewAIServiceError AI服务错误
func NewAIServiceError(message string, err error) *AppError {
	return &AppError{
		Code:       CodeAIServiceError,
		Message:    message,
		HTTPStatus: http.StatusInternalServerError,
		Err:        err,
	}
}

// NewRateLimitError 限流错误
func NewRateLimitError(message string) *AppError {
	return &AppError{
		Code:       CodeRateLimitError,
		Message:    message,
		HTTPStatus: http.StatusTooManyRequests,
	}
}

// NewTimeoutError 超时错误
func NewTimeoutError(message string, err error) *AppError {
	return &AppError{
		Code:       CodeTimeoutError,
		Message:    message,
		HTTPStatus: http.StatusGatewayTimeout,
		Err:        err,
	}
}

// NewInternalError 内部错误
func NewInternalError(message string, err ...error) *AppError {
	var e error
	if len(err) > 0 {
		e = err[0]
	}
	return &AppError{
		Code:       CodeInternalError,
		Message:    message,
		HTTPStatus: http.StatusInternalServerError,
		Err:        e,
	}
}

// NewProviderNotFoundError 提供商不存在错误
func NewProviderNotFoundError(provider string) *AppError {
	return &AppError{
		Code:       CodeProviderNotFound,
		Message:    fmt.Sprintf("AI provider '%s' not found", provider),
		HTTPStatus: http.StatusBadRequest,
	}
}

// NewConfigError 配置错误
func NewConfigError(message string, err error) *AppError {
	return &AppError{
		Code:       CodeConfigError,
		Message:    message,
		HTTPStatus: http.StatusInternalServerError,
		Err:        err,
	}
}

// NewBadRequestError 请求参数错误
func NewBadRequestError(message string) *AppError {
	return &AppError{
		Code:       CodeInvalidParameter,
		Message:    message,
		HTTPStatus: http.StatusBadRequest,
	}
}

// NewUnauthorizedError 未授权错误
func NewUnauthorizedError(message string) *AppError {
	return &AppError{
		Code:       CodeUnauthorized,
		Message:    message,
		HTTPStatus: http.StatusUnauthorized,
	}
}

// NewForbiddenError 禁止访问错误
func NewForbiddenError(message string) *AppError {
	return &AppError{
		Code:       CodeForbidden,
		Message:    message,
		HTTPStatus: http.StatusForbidden,
	}
}

// NewNotFoundError 资源不存在错误
func NewNotFoundError(message string) *AppError {
	return &AppError{
		Code:       CodeUserNotFound,
		Message:    message,
		HTTPStatus: http.StatusNotFound,
	}
}
