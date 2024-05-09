package userusecase

type userService interface {
	Registration(userName, email, password string) error
}

type UseCase struct {
	userService userService
}

func NewUseCase(userService userService) *UseCase {
	return &UseCase{
		userService: userService,
	}
}

func (uc *UseCase) Registration(userName, email, password string) error {
	return uc.userService.Registration(userName, email, password)
}
