package main

import (
	"context"
	"fmt"
	"time"

	"github.com/nbd-wtf/go-nostr"
	sdk "github.com/nbd-wtf/nostr-sdk"
	"github.com/puzpuzpuz/xsync/v2"
	"golang.org/x/sync/singleflight"
)

var (
	people              = xsync.NewMapOf[*sdk.ProfileMetadata]()
	profileMetadataPool = new(singleflight.Group)
)

func ensurePersonMetadata(pubkey string) chan *sdk.ProfileMetadata {
	ch := make(chan *sdk.ProfileMetadata)

	go func() {
		person, loaded := people.LoadOrCompute(pubkey, func() *sdk.ProfileMetadata {
			v, err, _ := profileMetadataPool.Do(pubkey, func() (any, error) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
				defer cancel()

				events := pool.SubManyEose(ctx, []string{
					"wss://purplepag.es",
					"wss://relay.noswhere.com",
					"wss://relay.nostr.band",
				}, nostr.Filters{
					{
						Kinds:   []int{0, 10002},
						Authors: []string{pubkey},
					},
				})
				if events == nil {
					return nil, fmt.Errorf("subscriptions couldn't be created")
				}

				var kind10002 *nostr.Event
				for ie := range events {
					switch ie.Event.Kind {
					case 0:
						// we got the metadata directly, so just use that
						if metadata, err := sdk.ParseMetadata(ie.Event); err == nil {
							return metadata, nil
						}
					case 10002:
						// we got a relay list, we may use this if we don't get any metadata
						if kind10002 == nil || kind10002.CreatedAt < ie.Event.CreatedAt {
							kind10002 = ie.Event
						}
					}
				}

				if kind10002 == nil {
					return nil, fmt.Errorf("couldn't find metadata for %s anywhere", pubkey)
				}

				// if we reach this point we only have a relay list, so use that
				relays := make([]string, 0, len(kind10002.Tags))
				for _, tag := range kind10002.Tags {
					if len(tag) >= 2 {
						relays = append(relays, tag[1])
					}
				}

				fmt.Println("searching metadata for", pubkey, "on", relays)
				events = pool.SubManyEose(ctx, relays, nostr.Filters{
					{
						Kinds:   []int{0},
						Authors: []string{pubkey},
					},
				})
				if events == nil {
					return nil, fmt.Errorf("subscriptions (second) couldn't be created")
				}

				for ie := range events {
					if metadata, err := sdk.ParseMetadata(ie.Event); err == nil {
						return metadata, nil
					}
				}

				return nil, fmt.Errorf("couldn't find metadata for %s anywhere (2)", pubkey)
			})

			if err != nil {
				fmt.Println("failed to load metadata for", pubkey)
				return &sdk.ProfileMetadata{} // an empty thing so we don't try to load again
			}

			return v.(*sdk.ProfileMetadata)
		})
		if person != nil && !loaded {
			// this means we got something new
			ch <- person
		}
		close(ch)
	}()

	return ch
}
