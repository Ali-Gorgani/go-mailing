package auth

type Maker interface {
	CreateToken()
	VerifyToken()
}