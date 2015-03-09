package stream

import (
	"encoding/json"
	"strings"
)

type Event struct {
	Id        string
	Type      string
	Actor     User
	Repo      Repo
	Payload   Payload
	Public    bool
	CreatedAt string `json:"created_at"`
}

type WatchEvent struct {
	Id        string
	Actor     User
	Repo      Repo
	Payload   Payload
	Public    bool
	CreatedAt string `json:"created_at"`
}

type User struct {
	Id        int64
	Login     string
	Url       string
	AvatarUrl string `json:"avatar_url"`
}

type Repo struct {
	Id   int64
	Name string
	Url  string
}

type Payload struct {
	Action string
}

func parseEvents(body string) []Event {
	events := []Event{}

	decoder := json.NewDecoder(strings.NewReader(body))
	decoder.Decode(&events)
	return events
}

// not implemented yet
func (e *Event) instantiate() interface{} {
	switch e.Type {
	case "WatchEvent":
		return &WatchEvent{
			Id:        e.Id,
			Actor:     e.Actor,
			Repo:      e.Repo,
			Payload:   e.Payload,
			Public:    e.Public,
			CreatedAt: e.CreatedAt,
		}
	default:
		return e
	}
}
