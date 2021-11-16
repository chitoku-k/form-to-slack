form-to-slack
=============

[![][workflow-badge]][workflow-link]

Sends form submission to Slack.

## Requirements

- Go
- reCAPTCHA v3
- Slack Webhook URL

## Installation

```sh
$ docker buildx build .
```

```sh
# Port number (required)
export PORT=8080

# TLS certificate and private key (optional; if not specified, form-to-slack is served over HTTP)
export TLS_CERT=/path/to/tls/cert
export TLS_KEY=/path/to/tls/key

# reCAPTCHA secret key (required)
export RECAPTCHA_SECRET=

# Slack Webhook URL (required)
export SLACK_WEBHOOK_URL=

# Access-Control-Allow-Origin (optional; space-separated)
export ALLOWED_ORIGINS=
```

## Usage

```html
<form>
    <!-- https://developers.google.com/recaptcha/docs/v3 -->
    <input type="hidden" name="g-recaptcha-response" />
    <input type="text" name="name" placeholder="Your Name" required />
    <input type="text" name="subject" placeholder="Subject" required />
    <input type="text" name="email" placeholder="E-mail Address (optional)" />
    <textarea name="body" placeholder="Body" required></textarea>
</form>
```

[workflow-link]:    https://github.com/chitoku-k/form-to-slack/actions?query=branch:master
[workflow-badge]:   https://img.shields.io/github/workflow/status/chitoku-k/form-to-slack/CI%20Workflow/master.svg?style=flat-square
