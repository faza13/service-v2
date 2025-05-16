package user

import (
	"service/pkg/otel"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUsecase IUserUsecase
}

func NewUserHandler(userUsecase IUserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
	}
}

// @Summary Get List
// @Description get list
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /users/list [get]
func (h *UserHandler) List(c *gin.Context) {
	ctx, span := otel.AddSpan(c.Request.Context(), "user.get")
	defer span.End()

	data := h.userUsecase.List(ctx)

	c.JSON(200, data)
}

func (h *UserHandler) Register(c *gin.Context) {
	ctx, span := otel.AddSpan(c.Request.Context(), "user.get")
	defer span.End()

	req := RegistrationRequest{}

	err := c.Bind(&req)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	data, err := h.userUsecase.Register(ctx, &req)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, data)
}
