package main

import (
	"image"
	"image/color"
	"net/http"
	"sync"
	"unsafe"

	"github.com/puzpuzpuz/xsync/v2"
)

const MAX_LOCKS = 50

var (
	namedMutexPool = make([]sync.Mutex, MAX_LOCKS)
	imagesCache    = xsync.NewMapOf[image.Image]()
	neutralImage   = generateNeutralImage(color.RGBA{156, 62, 93, 255})
)

func generateNeutralImage(color color.Color) image.Image {
	const size = 1
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	for x := 0; x < size; x++ {
		for y := 0; y < size; y++ {
			img.Set(x, y, color)
		}
	}
	return img
}

// func addGroup(relayURL string, groupId string, relaysListWidget *widget.List, chatMessagesListWidget *widget.List) {
// 	chatRelay, ok := relays.Load(relayURL)
// 	if !ok {
// 		// TODO: Better handling
// 		fmt.Println("no relay to add group to:", relayURL)
// 		return
// 	}
//
// 	if g, ok := chatRelay.Groups.Load(groupId); ok {
// 		fmt.Println("group already there:", g)
// 		return
// 	}
//
// 	group := &ChatGroup{
// 		ID:           groupId,
// 		Name:         groupId,
// 		ChatMessages: make([]*nostr.Event, 0),
// 	}
// 	chatRelay.Groups.Store(groupId, group)
//
// 	ctx := context.Background()
// 	sub, err := chatRelay.Relay.Subscribe(ctx, []nostr.Filter{
// 		{
// 			Kinds: []int{9},
// 			Tags: nostr.TagMap{
// 				"g": {groupId},
// 			},
// 		},
// 		{
// 			Kinds: []int{39000, 39003},
// 			Tags: nostr.TagMap{
// 				"d": {groupId},
// 			},
// 		},
// 	}, nostr.WithLabel("chat"+groupId))
// 	if err != nil {
// 		fmt.Println("can't subscribe", chatRelay.Relay, groupId, err)
// 		return
// 	}
//
// 	chatRelay.Subscriptions.Store(groupId, sub)
// 	saveRelays()
// 	updateLeftMenuList(relaysListWidget)
//
// 	for idx, menuItem := range relayMenuData {
// 		if menuItem.GroupName == groupId {
// 			relaysListWidget.Select(idx)
// 			break
// 		}
// 	}
//
// 	if err := sub.Fire(); err != nil {
// 		// TODO: better handling
// 		panic(err)
// 	}
//
// 	go func() {
// 		for ev := range sub.Events {
// 			switch ev.Kind {
// 			case 39000:
// 				if tag := ev.Tags.GetFirst([]string{"name", ""}); tag != nil {
// 					group.Name = (*tag)[1]
// 				}
// 				if tag := ev.Tags.GetFirst([]string{"picture", ""}); tag != nil {
// 					group.Picture = (*tag)[1]
// 				}
// 				updateLeftMenuList(relaysListWidget)
// 			case 39003:
// 				for _, tag := range ev.Tags.GetAll([]string{"g", ""}) {
// 					group.Subgroups = append(group.Subgroups, tag[1])
// 				}
// 				updateLeftMenuList(relaysListWidget)
// 			case 9:
// 				group.ChatMessages = insertEventIntoAscendingList(group.ChatMessages, ev)
// 				chatMessagesListWidget.Refresh()
// 				chatMessagesListWidget.ScrollToBottom()
// 				updateLeftMenuList(relaysListWidget)
//
// 				go func(pubkey string) {
// 					metadata := <-ensurePersonMetadata(pubkey)
// 					if metadata == nil {
// 						// it will be nil if we didn't get any new metadata
// 						// so we don't have to update anything if
// 						return
// 					}
// 					chatMessagesListWidget.Refresh()
// 				}(ev.PubKey)
// 			}
// 		}
// 	}()
// }
//
// func addRelay(relayURL string) {
// 	if _, ok := relays.Load(relayURL); ok {
// 		return
// 	} else {
// 		fmt.Println("connecting to", relayURL)
// 		ctx := context.Background()
// 		relay, err := nostr.RelayConnect(ctx, relayURL)
// 		if err != nil {
// 			fmt.Println("Err connecting to: ", relayURL)
// 			return
// 		}
//
// 		go func() {
// 			// when we lose connectivity, connect again
// 			<-relay.Context().Done()
// 			relays.Delete(relayURL)
// 			addRelay(relayURL)
// 		}()
//
// 		chatRelay := &ChatRelay{
// 			Relay:         *relay,
// 			Subscriptions: xsync.NewMapOf[*nostr.Subscription](),
// 			Groups:        xsync.NewMapOf[*ChatGroup](),
// 		}
//
// 		relays.Store(relayURL, chatRelay)
// 	}
// }
//
// func updateLeftMenuList(relaysListWidget *widget.List) {
// 	relayMenuData = make([]LeftMenuItem, 0)
//
// 	relays.Range(func(_ string, chatRelay *ChatRelay) bool {
// 		relayMenuData = append(relayMenuData, LeftMenuItem{
// 			RelayURL: chatRelay.Relay.URL,
// 			IsRoot:   true,
// 			GroupID:  "/",
// 		})
//
// 		chatRelay.Groups.Range(func(_ string, group *ChatGroup) bool {
// 			relayMenuData = append(relayMenuData, LeftMenuItem{
// 				RelayURL:  chatRelay.Relay.URL,
// 				IsRoot:    false,
// 				GroupID:   group.ID,
// 				GroupName: group.Name,
// 				GroupIcon: group.Picture,
// 			})
// 			return true
// 		})
//
// 		return true
// 	})
//
// 	relaysListWidget.Refresh()
// }

func imageFromURL(u string) (res image.Image) {
	res, ok := imagesCache.Load(u)
	if ok {
		return res
	}

	// this is so that we only try to load the same url once
	unlock := namedLock(u)
	defer unlock()

	// store result on cache (even if it's nil)
	defer func() {
		imagesCache.Store(u, res)
	}()

	// load url
	response, err := http.Get(u)
	if err != nil {
		return nil
	}
	defer response.Body.Close()

	img, _, err := image.Decode(response.Body)
	if err != nil {
		return nil
	}

	return img
}

func namedLock(name string) (unlock func()) {
	sptr := unsafe.StringData(name)
	idx := uint64(memhash(unsafe.Pointer(sptr), 0, uintptr(len(name)))) % MAX_LOCKS
	namedMutexPool[idx].Lock()
	return namedMutexPool[idx].Unlock
}

//go:noescape
//go:linkname memhash runtime.memhash
func memhash(p unsafe.Pointer, h, s uintptr) uintptr
