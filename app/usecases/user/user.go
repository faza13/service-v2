package user

type UserUsecase struct {
	userRepo IUserRepo
}

func NewUserUsecase(userRepo IUserRepo) *UserUsecase {
	return &UserUsecase{
		userRepo: userRepo,
	}
}
