package ui

import (
	"time"
)

type timeOption struct {
	Label    string
	Duration time.Duration
}

var timeOptions = []timeOption{
	{"Infinite", 0},
	{"10 seconds", 10 * time.Second},
	{"15 seconds", 15 * time.Second},
	{"30 seconds", 30 * time.Second},
	{"1 minute", time.Minute},
	{"5 minutes", 5 * time.Minute},
}

type messages struct {
	messages []message
	ttl      time.Duration
}

type message struct {
	contents   string
	expiration *time.Time
}

func newMessages() *messages {
	return &messages{
		messages: make([]message, 0),
		ttl:      0,
	}
}

func (s *messages) AppendMessage(contents string) {
	var expiration *time.Time
	if s.ttl == 0 {
		expiration = nil
	} else {
		t := time.Now().Add(s.ttl)
		expiration = &t
	}

	s.messages = append(s.messages, message{
		contents:   contents,
		expiration: expiration,
	})
}

func (s *messages) SetTTL(option time.Duration) {
	// Validate if present in avaliable options
	for _, available := range timeOptions {
		if option == available.Duration {
			s.ttl = available.Duration
		}
	}
	// Ignore otherwise
}

func (s *messages) String() string {
	result := ""
	var cleanedMessages []message
	for _, m := range s.messages {
		if m.expiration == nil || m.expiration.After(time.Now()) {
			result += m.contents
			cleanedMessages = append(cleanedMessages, m)
			continue
		}
	}

	s.messages = cleanedMessages
	return result
}
