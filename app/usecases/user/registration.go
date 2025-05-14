package user

import (
	"context"
	"service/app/controllers/restapi/user"
	"service/pkg/otel"
)

func (u *UserUsecase) Registration(ctx context.Context, request *user.RegistrationRequest) (error, map[string]interface{}) {
	ctx, span := otel.AddSpan(ctx, "user_usecase.Registration")
	defer span.End()

	err := u.transactor.WithTx(ctx, func(ctx context.Context) error {
		// proses data

		return nil
	})

	return err, map[string]interface{}{"test": "test"}
}
