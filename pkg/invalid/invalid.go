package invalid

import (
	"chat-application/users/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type InvalidHelper struct {
	ctx *gin.Context
}

func NewInvalidHelper(ctx *gin.Context) *InvalidHelper {
	return &InvalidHelper{
		ctx: ctx,
	}
}

func (h *InvalidHelper) SendInvalidBody() {
	response := domain.NewResponse(
		http.StatusBadRequest,
		false,
		"Invalid request body",
		nil,
	)
	h.ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
}

func (h *InvalidHelper) SendErrorValidation(validate *validator.Validate, obj interface{}) {
	if err := validate.Struct(&obj); err != nil {
		errValidation := err.(validator.ValidationErrors)
		response := domain.NewResponse(
			http.StatusBadRequest,
			false,
			errValidation[0].Error(),
			nil,
		)
		h.ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
	}
}
