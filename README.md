# gitlab-slack-bot
Slack Bot for Gitlab Webhook Events

## Architecture

Gitlab Event Webhook &rarr; API Gateway Endpoint &rarr; Event Publisher Function &rarr; Pub/Sub Event Topic &larr; Event Subscriber Function &rarr; Slack Channel

## Events
- Merge Requests
- Tags

## GCP Services
- API Gateway
- Cloud Functions
- Pub/Sub
- Secret Manager and Credentials
