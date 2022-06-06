package p

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/xanzy/go-gitlab"
)

func PublishTagEvent(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	topicID := "gitlab-tag-event"

	var event gitlab.TagEvent

	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		fmt.Fprintf(w, "Error decoding tag event: %v\n", err)
		return
	}

	if event.ObjectKind == "tag_push" {
		encoding, err := json.Marshal(event)

		if err != nil {
			fmt.Fprintf(w, "Error encoding tag event: %v\n", err)
			return
		}

		id, err := PublishMessage(ctx, topicID, encoding)

		if err != nil {
			fmt.Fprintf(w, "Error publishing tag event message: %v\n", err)
			return
		}

		fmt.Fprintf(w, "Published gitlab tag event with message id: %v\n", id)
		return
	} else {
		fmt.Fprintf(w, "Ignoring tag event: %v\n", event.ObjectKind)
		return
	}

}
