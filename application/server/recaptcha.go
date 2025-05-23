package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type reCaptchaResponse struct {
	Success     bool      `json:"success"`
	ChallengeTS time.Time `json:"challenge_ts"`
	Hostname    string    `json:"hostname"`
	ErrorCodes  []string  `json:"error-codes"`
}

func (e *engine) verifyReCaptcha(ctx context.Context, response string) (bool, error) {
	body := url.Values{
		"secret":   {e.ReCaptchaSecret},
		"response": {response},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, e.ReCaptchaURL, strings.NewReader(body.Encode()))
	if err != nil {
		return false, fmt.Errorf("failed to construct a request: %w", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to verify reCAPTCHA token: %w", err)
	}
	defer func() {
		_ = res.Body.Close()
	}()

	var reCaptchaRes reCaptchaResponse
	err = json.NewDecoder(res.Body).Decode(&reCaptchaRes)
	if err != nil {
		return false, fmt.Errorf("failed to decode reCAPTCHA response: %w", err)
	}

	if !reCaptchaRes.Success {
		slog.Info("Error from reCAPTCHA", slog.Any("error-codes", reCaptchaRes.ErrorCodes))
	}
	return reCaptchaRes.Success, nil
}
