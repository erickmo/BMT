package notif

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const fcmEndpoint = "https://fcm.googleapis.com/fcm/send"

type FCMProvider struct {
	serverKey string
}

func NewFCMProvider(serverKey string) *FCMProvider {
	return &FCMProvider{serverKey: serverKey}
}

func (p *FCMProvider) Kirim(ctx context.Context, tujuan, subjek, pesan string) error {
	if p.serverKey == "" {
		return nil
	}

	body := map[string]any{
		"to": tujuan,
		"notification": map[string]string{
			"title": subjek,
			"body":  pesan,
		},
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("fcm marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fcmEndpoint, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("fcm buat request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "key="+p.serverKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("fcm kirim request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("fcm response status: %d", resp.StatusCode)
	}

	return nil
}
