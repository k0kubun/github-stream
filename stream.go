package stream

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type Stream struct {
	Events     chan interface{}
	tokens     []string
	tokenIndex int
	stopped    bool
	etag       string
	lastId     string
	done       chan bool
}

const (
	endpoint       = "https://api.github.com/events"
	bufferedEvents = 8000
)

// Start polling public events asynchronously.
func NewStream(tokens ...string) *Stream {
	if len(tokens) == 0 {
		panic("NewStream requires at least 1 token.")
	}

	s := &Stream{
		Events:     make(chan interface{}, bufferedEvents),
		tokens:     tokens,
		tokenIndex: 0,
		stopped:    false,
		done:       make(chan bool),
	}
	go s.pollEvents()
	return s
}

// Stop polling and returns pending events synchronously.
func (s *Stream) Stop() []interface{} {
	s.stopped = true
	<-s.done

	events := []interface{}{}
	for len(s.Events) > 0 {
		event := <-s.Events
		events = append(events, event)
	}
	return events
}

func (s *Stream) pollEvents() {
	for {
		if s.stopped {
			s.done <- true
			break
		}

		events := s.getEvents()
		for i := len(events) - 1; i >= 0; i-- {
			s.Events <- events[i].instantiate()
		}
	}
}

func (s *Stream) getEvents() []Event {
	includeLast := false
	events := s.getFirstPage()

	lastId := s.lastId
	if len(events) > 0 {
		lastId = events[0].Id
	}

	if len(s.lastId) > 0 {
		slice := []Event{}
		for _, event := range events {
			if event.Id == s.lastId {
				includeLast = true
				break
			}
			slice = append(slice, event)
		}

		events = slice
	}

	if !includeLast && len(events) > 0 {
		go s.loadPage(2, events[len(events)-1].Id, s.lastId)
	}

	s.lastId = lastId
	return events
}

func (s *Stream) getFirstPage() []Event {
	body := s.getRequest(endpoint)
	events := parseEvents(body)
	return events
}

func (s *Stream) loadPage(page int, maxId string, minId string) {
	if page > 10 {
		return
	}

	includeLast := false
	paginated := fmt.Sprintf("%s?page=%d", endpoint, page)

	body := s.getRequest(paginated)
	events := parseEvents(body)
	for _, event := range events {
		if event.Id >= maxId {
			continue
		}
		if event.Id <= minId {
			includeLast = true
			break
		}

		s.Events <- event.instantiate()
	}

	if !includeLast && len(events) > 0 {
		s.loadPage(page+1, events[len(events)-1].Id, minId)
	}
}

func (s *Stream) getRequest(url string) string {
	token := s.selectToken()
	req, _ := http.NewRequest("GET", endpoint, nil)
	req.Header.Set("Authorization", "token "+token)
	if len(s.etag) > 0 {
		req.Header.Set("If-None-Match", s.etag)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	s.updateETag(res)

	body, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	return string(body)
}

func (s *Stream) updateETag(res *http.Response) {
	var etag string
	if array, ok := res.Header["Etag"]; ok && len(array) > 0 {
		etag = array[0]
	}
	if len(etag) == 0 {
		return
	}
	s.etag = etag
}

func (s *Stream) selectToken() string {
	s.tokenIndex++
	if len(s.tokens) == s.tokenIndex {
		s.tokenIndex = 0
	}
	return s.tokens[s.tokenIndex]
}
