package notif

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	fonnteEndpoint  = "https://api.fonnte.com/send"
	wablasEndpoint  = "https://solo.wablas.com/api/send-message"
	providerFonnte  = "fonnte"
	providerWablas  = "wablas"
)

type WhatsAppProvider struct {
	token    string
	provider string // "fonnte" | "wablas"
}

func NewWhatsAppProvider(provider, token string) *WhatsAppProvider {
	return &WhatsAppProvider{provider: provider, token: token}
}

func (p *WhatsAppProvider) Kirim(ctx context.Context, tujuan, subjek, pesan string) error {
	if p.token == "" {
		return nil
	}

	switch p.provider {
	case providerFonnte:
		return p.kirimFonnte(ctx, tujuan, pesan)
	case providerWablas:
		return p.kirimWablas(ctx, tujuan, pesan)
	default:
		return p.kirimFonnte(ctx, tujuan, pesan)
	}
}

func (p *WhatsAppProvider) kirimFonnte(ctx context.Context, tujuan, pesan string) error {
	form := url.Values{}
	form.Set("target", tujuan)
	form.Set("message", pesan)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fonnteEndpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("whatsapp fonnte buat request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", p.token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("whatsapp fonnte kirim request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("whatsapp fonnte response status: %d", resp.StatusCode)
	}

	return nil
}

func (p *WhatsAppProvider) kirimWablas(ctx context.Context, tujuan, pesan string) error {
	body := map[string]string{
		"phone":   tujuan,
		"message": pesan,
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("whatsapp wablas marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, wablasEndpoint, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("whatsapp wablas buat request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", p.token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("whatsapp wablas kirim request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("whatsapp wablas response status: %d", resp.StatusCode)
	}

	return nil
}
