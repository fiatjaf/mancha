package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/nbd-wtf/go-nostr"
	sdk "github.com/nbd-wtf/nostr-sdk"
	"golang.org/x/sync/singleflight"
)

var people = People{cache: make(map[string]*Metadata)}

type People struct {
	sync.Mutex
	singleflight.Group
	cache          map[string]*Metadata
	storeScheduled atomic.Bool
}

type Metadata struct {
	sdk.ProfileMetadata
	When nostr.Timestamp
}

const SAVED_PEOPLE_KEY = "saved_people"

func savePeople() {
	j, _ := json.Marshal(people.cache)
	a.Preferences().SetString(SAVED_PEOPLE_KEY, string(j))
}

func loadPeople() {
	jstr := a.Preferences().String(SAVED_PEOPLE_KEY)
	json.Unmarshal([]byte(jstr), &people.cache)
}

func (people *People) Load(pubkey string) (*Metadata, bool) {
	people.Lock()
	defer people.Unlock()
	m, ok := people.cache[pubkey]
	return m, ok
}

func (people *People) ensurePersonMetadata(pubkey string, currentRelay string) (wasCached bool, isCachedNow bool) {
	metadata, ok := people.Load(pubkey)
	if ok {
		if metadata.When > nostr.Now()-60*60*3 /* 3 hours */ {
			return true, true
		}
	}
	isCachedNow = people.fetchAndCacheMetadata(pubkey, currentRelay)
	return
}

func (people *People) fetchAndCacheMetadata(pubkey string, currentRelay string) bool {
	people.Lock()
	defer people.Unlock()

	v, err, _ := people.Do(pubkey, func() (any, error) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		events := pool.SubManyEose(ctx, metadataRelays, nostr.Filters{
			{
				Kinds:   []int{10002},
				Authors: []string{pubkey},
			},
		})
		if events == nil {
			return nil, fmt.Errorf("subscriptions couldn't be created")
		}

		var kind10002 *nostr.Event
		for ie := range events {
			// we got a relay list, we may use this if we don't get any metadata
			if kind10002 == nil || kind10002.CreatedAt < ie.Event.CreatedAt {
				kind10002 = ie.Event
			}
		}

		if kind10002 == nil {
			kind10002 = &nostr.Event{Tags: nil}
		}

		// if we reach this point we only have a relay list, so use that
		relays := make([]string, 0, len(kind10002.Tags)+1+len(metadataRelays))
		for _, tag := range kind10002.Tags {
			if len(tag) >= 2 {
				relays = append(relays, tag[1])
			}
		}
		relays = append(relays, currentRelay)
		relays = append(relays, metadataRelays...)

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

		var latest nostr.Timestamp
		var metadata sdk.ProfileMetadata
		for ie := range events {
			parsed, err := sdk.ParseMetadata(ie.Event)
			if err != nil {
				continue
			}

			if ie.Relay.URL == nostr.NormalizeURL(currentRelay) {
				// prioritize metadata stored in this group's relay
				metadata = parsed
				return metadata, nil
			}

			if ie.Event.CreatedAt > latest {
				latest = ie.Event.CreatedAt
				metadata = parsed
			}
		}

		return metadata, nil
	})

	// we will return an err if there was an error while trying to fetch this
	// if we tried to fetch but still got nothing we will cache the nil thing
	if err != nil {
		return false
	}
	people.cache[pubkey] = &Metadata{
		ProfileMetadata: v.(sdk.ProfileMetadata),
		When:            nostr.Now(),
	}

	if people.storeScheduled.CompareAndSwap(false, true) {
		go func() {
			time.Sleep(time.Second * 5)
			savePeople()
			people.storeScheduled.Store(false)
		}()
	}

	return true
}
