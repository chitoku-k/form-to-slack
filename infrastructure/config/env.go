package config

import (
	"errors"
	"os"
	"strings"
)

const (
	defaultReCaptchaVerifyEndpoint = "https://www.google.com/recaptcha/api/siteverify"
)

type Environment struct {
	AllowedOrigins  string
	Port            string
	TLSCert         string
	TLSKey          string
	ReCaptchaURL    string
	ReCaptchaSecret string
	SlackWebhookURL string
}

func Get() (Environment, error) {
	var missing []string
	var env Environment

	for k, v := range map[string]*string{
		"ALLOWED_ORIGINS": &env.AllowedOrigins,
		"TLS_CERT":        &env.TLSCert,
		"TLS_KEY":         &env.TLSKey,
		"RECAPTCHA_URL":   &env.ReCaptchaURL,
	} {
		*v = os.Getenv(k)
	}

	for k, v := range map[string]*string{
		"PORT":              &env.Port,
		"RECAPTCHA_SECRET":  &env.ReCaptchaSecret,
		"SLACK_WEBHOOK_URL": &env.SlackWebhookURL,
	} {
		*v = os.Getenv(k)

		if *v == "" {
			missing = append(missing, k)
		}
	}

	if len(missing) > 0 {
		return env, errors.New("missing env(s): " + strings.Join(missing, ", "))
	}

	if env.ReCaptchaURL == "" {
		env.ReCaptchaURL = defaultReCaptchaVerifyEndpoint
	}

	return env, nil
}
