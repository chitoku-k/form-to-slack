form-to-slack
=============

Sends form submission to Slack.

## Requirements

- PHP >= 7.3 (+ composer)
- reCAPTCHA v3
- Slack Webhook URL

## Installation

```sh
$ composer install
```

```sh
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
