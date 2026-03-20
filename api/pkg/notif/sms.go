package notif

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

const (
	zenziveEndpoint = "https://gsm.zenziva.net/apps/smsapi.php"
	providerZenziva = "zenziva"
)

type SMSProvider struct {
	provider string // "zenziva" | "twilio"
	apiKey   string
	apiUser  string // untuk zenziva
}

func NewSMSProvider(provider, apiUser, apiKey string) *SMSProvider {
	return &SMSProvider{provider: provider, apiUser: apiUser, apiKey: apiKey}
}

func (p *SMSProvider) Kirim(ctx context.Context, tujuan, subjek, pesan string) error {
	if p.apiKey == "" {
		return nil
	}

	switch p.provider {
	case providerZenziva:
		return p.kirimZenziva(ctx, tujuan, pesan)
	default:
		return p.kirimZenziva(ctx, tujuan, pesan)
	}
}

func (p *SMSProvider) kirimZenziva(ctx context.Context, tujuan, pesan string) error {
	params := url.Values{}
	params.Set("userkey", p.apiUser)
	params.Set("passkey", p.apiKey)
	params.Set("nohp", tujuan)
	params.Set("pesan", pesan)

	endpoint := zenziveEndpoint + "?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return fmt.Errorf("sms zenziva buat request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("sms zenziva kirim request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("sms zenziva response status: %d", resp.StatusCode)
	}

	return nil
}
