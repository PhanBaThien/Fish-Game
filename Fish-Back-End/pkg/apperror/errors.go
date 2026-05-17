package apperror

import "net/http"

type AppError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	HTTPStatus int    `json:"-"`
}

func (e *AppError) Error() string {
	return e.Message
}

var (
	ErrBadRequest         = &AppError{Code: "BAD_REQUEST", Message: "dữ liệu yêu cầu không hợp lệ", HTTPStatus: http.StatusBadRequest}
	ErrInvalidCredentials = &AppError{Code: "INVALID_CREDENTIALS", Message: "tài khoản hoặc mật khẩu không đúng", HTTPStatus: http.StatusUnauthorized}
	ErrUsernameExisted    = &AppError{Code: "USERNAME_EXISTED", Message: "tài khoản đã tồn tại", HTTPStatus: http.StatusConflict}
	ErrUserNotFound       = &AppError{Code: "USER_NOT_FOUND", Message: "không tìm thấy người dùng", HTTPStatus: http.StatusNotFound}
	ErrInvalidToken       = &AppError{Code: "INVALID_TOKEN", Message: "token không hợp lệ", HTTPStatus: http.StatusUnauthorized}
	ErrExpiredToken       = &AppError{Code: "EXPIRED_TOKEN", Message: "token đã hết hạn", HTTPStatus: http.StatusUnauthorized}
	ErrForbidden          = &AppError{Code: "FORBIDDEN", Message: "bạn không có quyền thực hiện thao tác này", HTTPStatus: http.StatusForbidden}
	ErrInternalServer     = &AppError{Code: "INTERNAL_SERVER_ERROR", Message: "lỗi máy chủ nội bộ", HTTPStatus: http.StatusInternalServerError}
)
