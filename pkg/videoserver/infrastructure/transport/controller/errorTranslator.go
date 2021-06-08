package controller

import (
	"errors"
	"github.com/bearname/videohost/pkg/common/infrarstructure/transport/controller"
	"github.com/bearname/videohost/pkg/videoserver/domain"
	"net/http"
)

func TranslateError(err error) controller.TransportError {
	if errors.Is(err, controller.ErrBadRequest) {
		return controller.TransportError{
			Status: http.StatusBadRequest,
			Response: controller.Response{
				Code:    101,
				Message: err.Error(),
			},
		}
	} else if errors.Is(err, domain.ErrFailedDeleteLike) {
		return controller.TransportError{
			Status: http.StatusInternalServerError,
			Response: controller.Response{
				Code:    102,
				Message: err.Error(),
			},
		}
	} else if errors.Is(err, domain.ErrFailedAddLike) {
		return controller.TransportError{
			Status: http.StatusBadRequest,
			Response: controller.Response{
				Code:    103,
				Message: err.Error(),
			},
		}
	} else if errors.Is(err, domain.ErrAlreadyLike) {
		return controller.TransportError{
			Status: http.StatusBadRequest,
			Response: controller.Response{
				Code:    104,
				Message: err.Error(),
			},
		}
	} else if errors.Is(err, domain.ErrAlreadyDisLike) {
		return controller.TransportError{
			Status: http.StatusBadRequest,
			Response: controller.Response{
				Code:    105,
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
