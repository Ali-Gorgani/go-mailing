package repository

import (
	"log"
	"sync"

	"github.com/Ali-Gorgani/go-mailing/internal/config"
	"github.com/Ali-Gorgani/go-mailing/internal/entity"
)

type ProviderRepository struct {
	providers []entity.Provider
}

var instance *ProviderRepository
var once sync.Once

func GetProviderRepository() *ProviderRepository {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}
	once.Do(func() {
		instance = &ProviderRepository{
			providers: []entity.Provider{
				{
					Name:   "kavenegar",
					URL:    config.Kavenegar.URL,
					APIKey: config.Kavenegar.APIKey,
				},
				// Add more providers here
			},
		}
	})
	return instance
}

func (r *ProviderRepository) GetAll() []entity.Provider {
	return r.providers
}

func (r *ProviderRepository) GetByName(name string) (entity.Provider, bool) {
	for _, provider := range r.providers {
		if provider.Name == name {
			return provider, true
		}
	}
	return entity.Provider{}, false
}
