package p

import (
	"context"
	"log"

	"cloud.google.com/go/pubsub"
)

func PublishMessage(ctx context.Context, topicID string, data []byte) (string, error) {
	client, err := pubsub.NewClient(ctx, projectID)

	if err != nil {
		log.Fatalf("pubsub: NewClient: %v\n", err)
		return "", err
	}

	if err != nil {
		log.Fatalf("pubsub: Marshal: %v\n", err)
		return "", err
	}

	defer client.Close()

	if err != nil {
		log.Fatalf("pubsub: ReadAll: %v\n", err)
		return "", err
	}

	t := client.Topic(topicID)
	result := t.Publish(ctx, &pubsub.Message{
		Data: []byte(data),
	})

	id, err := result.Get(ctx)

	if err != nil {
		log.Fatalf("pubsub: result.Get: %v\n", err)
		return "", err
	}

	return id, nil
}
