package user

import (
	"github.com/gin-gonic/gin"
	"service/pkg/otel"
)

type UserHandler struct {
	userUsecase IUserUsecase
}

func NewUserHandler(userUsecase IUserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
	}
}

func (h *UserHandler) List(c *gin.Context) {
	ctx, span := otel.AddSpan(c.Request.Context(), "user.get")
	defer span.End()

	data := h.userUsecase.List(ctx)

	c.JSON(200, data)
}
