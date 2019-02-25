<?php

use Dotenv\Dotenv;
use ReCaptcha\ReCaptcha;

require_once __DIR__ . '/vendor/autoload.php';

try {
    $dotenv = Dotenv::create(__DIR__);
    $dotenv->load();
    $dotenv->required([
        'RECAPTCHA_SECRET',
        'SLACK_WEBHOOK_URL',
    ]);
} catch (\RuntimeException $e) {
    http_response_code(500);
    echo json_encode([ 'code' => 500 ]);
    exit;
}

if (isset($_ENV['ALLOWED_ORIGINS'], $_SERVER['HTTP_ORIGIN']) &&
    in_array($_SERVER['HTTP_ORIGIN'], explode(' ', $_ENV['ALLOWED_ORIGINS']), true)
) {
    header('Access-Control-Allow-Origin: ' . $_SERVER['HTTP_ORIGIN']);
}

if (!isset(
    $_POST['g-recaptcha-response'],
    $_POST['name'],
    $_POST['subject'],
    $_POST['body'],
)) {
    http_response_code(400);
    echo json_encode([ 'code' => 400 ]);
    exit;
}

if (!(new ReCaptcha($_ENV['RECAPTCHA_SECRET']))
    ->verify($_POST['g-recaptcha-response'])
) {
    http_response_code(401);
    echo json_encode([ 'code' => 401 ]);
    exit;
}

if (!(new \SlackMessage(new \Slack($_ENV['SLACK_WEBHOOK_URL'])))
    ->setText('New message has arrived:')
    ->addAttachment(
        (new SlackAttachment($_POST['body']))
            ->setTitle($_POST['subject'])
            ->setText($_POST['body'])
            ->setTimestamp($_SERVER['REQUEST_TIME'])
            ->setAuthor($_POST['name'] . (!empty($_POST['email']) ? "（{$_POST['email']}）" : '')))
    ->send()
) {
    http_response_code(503);
    echo json_encode([ 'code' => 503 ]);
    exit;
}

echo json_encode([ 'code' => 200 ]);
