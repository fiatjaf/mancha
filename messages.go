package main

import (
	"fmt"
	"time"

	"github.com/nbd-wtf/go-nostr"
	"golang.org/x/net/context"
)

func publishChat(group *Group, message string) error {
	evt := nostr.Event{
		CreatedAt: nostr.Now(),
		Kind:      9,
		Tags: nostr.Tags{
			nostr.Tag{"h", group.ID},
			// TODO: "previous"
		},
		Content: message,
	}
	if err := k.Sign(&evt); err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	relay, err := pool.EnsureRelay(group.Relay)
	if err != nil {
		return fmt.Errorf("failed to ensure relay '%s' in order to publish: %w", group.Relay, err)
	}

	fmt.Println("publishing", evt)
	if err := relay.Publish(ctx, evt); err != nil {
		return fmt.Errorf("failed to publish: %w", err)
	}

	return nil
}
