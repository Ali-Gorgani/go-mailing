package provider

import "github.com/Ali-Gorgani/go-mailing/internal/entity"

type MailProvider interface {
	SendMail(receptor, message string) error
}

func NewProvider(entity entity.Provider) MailProvider {
	switch entity.Name {
	case "kavenegar":
		return NewKavenegarProvider(entity)
	default:
		return nil
	}
}
