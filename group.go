package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip29"
	"golang.org/x/exp/slices"
)

type Group struct {
	Relay string

	nip29.Group

	messages      []*nostr.Event
	cancelContext context.CancelCauseFunc
}

const SAVED_GROUPS_KEY = "saved_groups"

func saveGroups() {
	j, _ := json.Marshal(state.groups)
	a.Preferences().SetString(SAVED_GROUPS_KEY, string(j))
}

func loadGroups() {
	jstr := a.Preferences().String(SAVED_GROUPS_KEY)
	var data []*Group
	json.Unmarshal([]byte(jstr), &data)
	state.groups = data
}

func joinGroup(relayURL string, groupId string) (*Group, error) {
	idx := slices.IndexFunc(state.groups, func(g *Group) bool { return g.ID == groupId })
	if idx >= 0 {
		return nil, fmt.Errorf("already in this group")
	}

	group := &Group{
		Relay: relayURL,
		Group: nip29.Group{
			ID:   groupId,
			Name: groupId,
		},
		messages: make([]*nostr.Event, 0, 300),
	}

	go func() {
		err := group.startListening()
		if err != nil {
			log.Println(err)
		}
	}()

	state.groups = append(state.groups, group)
	saveGroups()

	return group, nil
}

func (group *Group) startListening() error {
	relay, err := pool.EnsureRelay(group.Relay)
	if err != nil {
		return fmt.Errorf("error connecting to %s: %w", group.Relay, err)
	}

	ctx, cancel := context.WithCancelCause(context.Background())
	group.cancelContext = cancel

	sub, err := relay.Subscribe(ctx, []nostr.Filter{
		{
			Kinds: []int{9},
			Tags: nostr.TagMap{
				"h": {group.ID},
			},
			Limit: 100,
		},
		{
			Kinds: []int{39000, 39001, 39002},
			Tags: nostr.TagMap{
				"d": {group.ID},
			},
		},
	}, nostr.WithLabel("chat"+group.ID))
	if err != nil {
		return fmt.Errorf("failed to subscribe to '%s': %w", relay.URL, err)
	}

	messagesWidget := getMessagesWidget()

	eosed := false
	for {
		select {
		case <-sub.EndOfStoredEvents:
			eosed = true
			slices.SortFunc(group.messages, func(a, b *nostr.Event) int { return int(a.CreatedAt - b.CreatedAt) })
			messagesWidget.widget.Refresh()
			messagesWidget.widget.ScrollToBottom()
		case reason := <-sub.ClosedReason:
			return fmt.Errorf("subscription %s to %s closed: '%s'", sub.GetID(), sub.Relay.URL, reason)
		case <-ctx.Done():
			return fmt.Errorf("subscription %s to %s canceled: %w", sub.GetID(), sub.Relay.URL, ctx.Err())
		case evt := <-sub.Events:
			if evt == nil {
				// subscription closed
				return fmt.Errorf("subscription %s to %s closed abruptly", sub.GetID(), sub.Relay.URL)
			}

			switch evt.Kind {
			case 39000:
				group.MergeInMetadataEvent(evt)
				getGroupsWidget().widget.Refresh()
				saveGroups()
			case 39001:
			case 39002:
			case 9:
				if !eosed {
					// before eose we just add all messages very fast (they will be sorted when we get the eose)
					group.messages = append(group.messages, evt)
					continue
				}

				// now we assume most messages will be appended to the end
				if len(group.messages) == 0 || evt.CreatedAt > group.messages[len(group.messages)-1].CreatedAt {
					group.messages = append(group.messages, evt)
					messagesWidget.widget.Refresh()
					messagesWidget.widget.ScrollToBottom()
					continue
				}

				// otherwise insert it where it should be
				idx, ok := slices.BinarySearchFunc(group.messages, evt, func(e1, e2 *nostr.Event) int {
					return int(e1.CreatedAt - e2.CreatedAt)
				})
				if ok {
					// already have this
					continue
				}
				group.messages = append(group.messages, nil) // bogus, increase capacity
				copy(group.messages[idx+1:], group.messages[idx:])
				group.messages[idx] = evt

				messagesWidget.widget.Refresh()

				go func(pubkey string) {
					wasCached, isCachedNow := ensurePersonMetadataIsCached(pubkey, sub.Relay.URL)
					if wasCached {
						return
					}
					if !isCachedNow {
						return
					}
					messagesWidget.widget.Refresh()
				}(evt.PubKey)
			}
		}
	}
}
