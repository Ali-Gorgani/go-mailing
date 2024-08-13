package service

import (
	"errors"
	"strings"

	"github.com/Ali-Gorgani/go-mailing/internal/entity"
	"github.com/Ali-Gorgani/go-mailing/internal/provider"
	"github.com/Ali-Gorgani/go-mailing/internal/repository"
)

type MailService struct {
	repo *repository.ProviderRepository
}

func NewMailService() *MailService {
	return &MailService{
		repo: repository.GetProviderRepository(),
	}
}

func (s *MailService) GetAllProviders() []entity.Provider {
	return s.repo.GetAll()
}

func (s *MailService) SendMail(providerName, receptor, message string) error {
	providerName = strings.ToLower(providerName)
	if providerName == "" {
		providerName = "kavenegar"
	}

	providerEntity, found := s.repo.GetByName(providerName)
	if !found {
		return errors.New("provider not found")
	}

	mailProvider := provider.NewProvider(providerEntity)
	if mailProvider == nil {
		return errors.New("invalid provider")
	}

	return mailProvider.SendMail(receptor, message)
}
