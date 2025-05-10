package user

import (
	"context"
	"service/pkg/otel"
)

func (u *UserUsecase) List(ctx context.Context) map[string]interface{} {
	ctx, span := otel.AddSpan(ctx, "user_usecase.list")
	defer span.End()
	return map[string]interface{}{"test": "test"}
}
