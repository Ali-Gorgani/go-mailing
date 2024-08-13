package handler

import (
    "encoding/json"
    "net/http"

    "github.com/Ali-Gorgani/go-mailing/internal/service"
)

func GetProviders(mailService *service.MailService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        providers := mailService.GetAllProviders()

        providerNames := make([]string, len(providers))
        for i, provider := range providers {
            providerNames[i] = provider.Name
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(providerNames)
    }
}
