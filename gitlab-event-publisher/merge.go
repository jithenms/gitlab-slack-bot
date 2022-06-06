package p

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/xanzy/go-gitlab"
)

func PublishMergeEvent(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	topicID := "gitlab-merge-event"

	var event gitlab.MergeEvent

	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		fmt.Fprintf(w, "Error decoding merge event: %v\n", err)
		return
	}

	if event.ObjectAttributes.Action == "open" {
		encoding, err := json.Marshal(event)

		if err != nil {
			fmt.Fprintf(w, "Error encoding merge event: %v\n", err)
			return
		}

		id, err := PublishMessage(ctx, topicID, encoding)

		if err != nil {
			fmt.Fprintf(w, "Error publishing merge event message: %v\n", err)
			return
		}

		fmt.Fprintf(w, "Published gitlab merge event with message ID: %v\n", id)
		return
	} else {
		fmt.Fprintf(w, "Ignoring merge event: %v\n", event.ObjectAttributes.Action)
		return
	}
}
