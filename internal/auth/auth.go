//go:generate mockgen -source=auth.go -destination=mock/mock_auth.go -package=mock

package auth

var _ IAuthService = (*service)(nil)

type IAuthService interface {
	CreateUser(input CreateUserInput) (int64, error)
}
