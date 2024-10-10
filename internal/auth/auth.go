package auth

var _ IAuthService = (*service)(nil)

type IAuthService interface {
	CreateUser(input ServiceInput[CreateUser]) (int64, error)
}
