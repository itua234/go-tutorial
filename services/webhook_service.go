package services

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type IWebhookService interface {
	SendWebhook(webhookURL, event string, data any)
}

// WebhookService is responsible for sending webhooks.
type WebhookService struct{}

func NewWebhookService() *WebhookService {
	return &WebhookService{}
}

// SendWebhook sends a POST request with the given event and data to the specified URL.
// It is designed to be non-blocking when called with 'go'.
func (s *WebhookService) SendWebhook(webhookURL, event string, data any) {
	payload := map[string]any{
		"event": event,
		"data":  data,
	}

	// Marshal to JSON
	b, err := json.Marshal(payload)
	if err != nil {
		log.Println("Failed to marshal webhook payload:", err)
		return
	}

	// Send POST request
	response, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(b))
	if err != nil {
		log.Println("Failed to send webhook:", err)
		return
	}
	defer response.Body.Close()
	log.Println("Webhook sent, status:", response.Status)
}
