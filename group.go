package main

import (
	"encoding/json"

	"github.com/nbd-wtf/go-nostr"
)

type Group struct {
	Relay string

	ID      string
	Name    string
	Picture string
	About   string
	Members map[string]*Role
	Private bool
	Closed  bool

	messages []*nostr.Event
}

type Role struct {
	Name        string
	Permissions map[Permission]struct{}
}

type Permission = string

const (
	PermAddUser          Permission = "add-user"
	PermEditMetadata     Permission = "edit-metadata"
	PermDeleteEvent      Permission = "delete-event"
	PermRemoveUser       Permission = "remove-user"
	PermAddPermission    Permission = "add-permission"
	PermRemovePermission Permission = "remove-permission"
	PermEditGroupStatus  Permission = "edit-group-status"
)

const SAVED_GROUPS_KEY = "saved_groups"

func saveGroups(groups []*Group) {
	j, _ := json.Marshal(groups)
	a.Preferences().SetString(SAVED_GROUPS_KEY, string(j))
}

func loadGroups() []Group {
	jstr := a.Preferences().String(SAVED_GROUPS_KEY)
	var data []Group
	json.Unmarshal([]byte(jstr), &data)
	return data
}
