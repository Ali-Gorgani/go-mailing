package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Ali-Gorgani/go-mailing/internal/service"
)

type SendMailRequest struct {
	Provider string `json:"provider"`
	Receptor string `json:"receptor"`
	Message  string `json:"message"`
}

func SendMail(mailService *service.MailService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SendMailRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		err := mailService.SendMail(req.Provider, req.Receptor, req.Message)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Email sent successfully"))
	}
}
