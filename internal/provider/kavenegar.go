package provider

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/Ali-Gorgani/go-mailing/internal/config"
	"github.com/Ali-Gorgani/go-mailing/internal/entity"
)

type KavenegarProvider struct {
	APIKey  string
	BaseURL string
}

func NewKavenegarProvider(provider entity.Provider) *KavenegarProvider {
	return &KavenegarProvider{
		APIKey:  provider.APIKey,
		BaseURL: provider.URL,
	}
}

func (k *KavenegarProvider) SendMail(receptor, message string) error {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	endpoint := fmt.Sprintf("%s/v1/%s/sms/send.json", k.BaseURL, k.APIKey)

	params := url.Values{}
	params.Set("receptor", receptor)
	params.Set("message", message)
	params.Set("sender", config.Kavenegar.Sender)

	resp, err := http.Get(endpoint + "?" + params.Encode())
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid status code: %d", resp.StatusCode)
	}

	return nil
}
