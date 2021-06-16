package controller

import (
	"errors"
	"github.com/bearname/videohost/internal/common/infrarstructure/transport/controller"
	"github.com/bearname/videohost/internal/user/domain"
	"net/http"
)

func TranslateError(err error) controller.TransportError {
	if errors.Is(err, controller.ErrRouteNotFound) {
		return controller.TransportError{
			Status: http.StatusNotFound,
			Response: controller.Response{
				Code:    100,
				Message: err.Error(),
			},
		}
	} else if errors.Is(err, controller.ErrBadRequest) {
		return controller.TransportError{
			Status: http.StatusBadRequest,
			Response: controller.Response{
				Code:    101,
				Message: err.Error(),
			},
		}
	} else if errors.Is(err, domain.ErrUserNotExist) {
		return controller.TransportError{
			Status: http.StatusNotFound,
			Response: controller.Response{
				Code:    102,
				Message: err.Error(),
			},
		}
	} else if errors.Is(err, domain.ErrInvalidAuthorizationHeader) {
		return controller.TransportError{
			Status: http.StatusBadRequest,
			Response: controller.Response{
				Code:    103,
				Message: err.Error(),
			},
		}
	} else if errors.Is(err, domain.ErrInvalidAccessToken) {
		return controller.TransportError{
			Status: http.StatusUnauthorized,
			Response: controller.Response{
				Code:    104,
				Message: err.Error(),
			},
		}
	} else if errors.Is(err, domain.ErrInvalidRefreshToken) {
		return controller.TransportError{
			Status: http.StatusUnauthorized,
			Response: controller.Response{
				Code:    105,
				Message: err.Error(),
			},
		}
	} else if errors.Is(err, domain.ErrFailedCreateAccessToken) {
		return controller.TransportError{
			Status: http.StatusInternalServerError,
			Response: controller.Response{
				Code:    106,
				Message: err.Error(),
			},
		}
	} else if errors.Is(err, domain.ErrFailedUpdateAccessToken) {
		return controller.TransportError{
			Status: http.StatusInternalServerError,
			Response: controller.Response{
				Code:    107,
				Message: err.Error(),
			},
		}
	} else if errors.Is(err, domain.ErrFailedSaveUser) {
		return controller.TransportError{
			Status: http.StatusInternalServerError,
			Response: controller.Response{
				Code:    108,
				Message: err.Error(),
			},
		}
	} else if errors.Is(err, domain.ErrFailedCreateUserID) {
		return controller.TransportError{
			Status: http.StatusInternalServerError,
			Response: controller.Response{
				Code:    109,
				Message: err.Error(),
			},
		}
	} else if errors.Is(err, domain.ErrDuplicateUser) {
		return controller.TransportError{
			Status: http.StatusConflict,
			Response: controller.Response{
				Code:    110,
				Message: err.Error(),
			},
		}
	} else if errors.Is(err, domain.ErrWrongPassword) {
		return controller.TransportError{
			Status: http.StatusUnauthorized,
			Response: controller.Response{
				Code:    111,
				Message: err.Error(),
			},
		}
	}
	return controller.TransportError{
		Status: http.StatusInternalServerError,
		Response: controller.Response{
			Code:    100,
			Message: "unexpected error",
		},
	}
}
