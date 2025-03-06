package hub

import (
	"fmt"
	"log/slog"
	"sync"
	"testing"
	"time"
)

func TestHub(t *testing.T) {

	logger := slog.Default()

	link := "test"

	hub := NewHub(nil, logger, nil, "")

	hub.AddTrack(link)
	if hub.linksCount[link] != 1 {
		t.Fatalf("count is not valid")
	}

	hub.AddTrack(link)
	if hub.linksCount[link] != 2 {
		t.Fatalf("count is not valid")
	}

	hub.RemoveTrack(link)
	if hub.linksCount[link] != 1 {
		t.Fatalf("expected 1 track after removal, got %d", hub.linksCount[link])
	}

	hub.RemoveTrack(link)
	if _, exists := hub.linksCount[link]; exists {
		t.Fatalf("expected link to be removed, but it's still there")
	}

	if _, exists := hub.linksCancel[link]; exists {
		t.Fatalf("expected cancel function to be removed")
	}
}

func TestHubConcurrency(t *testing.T) {

	logger := slog.Default()

	h := NewHub(nil, logger, nil, "")

	link := "http://example.com"

	for i := 0; i < 100; i++ {
		go h.AddTrack(link)
	}

	time.Sleep(100 * time.Millisecond)

	if h.linksCount[link] == 0 {
		t.Fatalf("expected at least one track, got %d", h.linksCount[link])
	}

	for i := 0; i < 100; i++ {
		go h.RemoveTrack(link)
	}

	time.Sleep(100 * time.Millisecond)

	if _, exists := h.linksCount[link]; exists {
		t.Fatalf("expected link to be fully removed after concurrent removals")
	}
}

func TestHubMaxLinks(t *testing.T) {

	logger := slog.Default()
	hub := NewHub(nil, logger, nil, "")

	for i := 0; i < 1000; i++ {
		link := fmt.Sprintf("http://example%d.com", i)
		hub.AddTrack(link)
		if hub.linksCount[link] != 1 {
			t.Fatalf("expected count to be 1 for link %s, got %d", link, hub.linksCount[link])
		}
	}

	for i := 0; i < 1000; i++ {
		link := fmt.Sprintf("http://example%d.com", i)
		if hub.linksCount[link] != 1 {
			t.Fatalf("link %s count is incorrect, expected 1, got %d", link, hub.linksCount[link])
		}
	}

	for i := 0; i < 1000; i++ {
		link := fmt.Sprintf("http://example%d.com", i)
		hub.RemoveTrack(link)
		if _, exists := hub.linksCount[link]; exists {
			t.Fatalf("expected link %s to be removed, but it still exists", link)
		}
	}
}


func TestHubEmptyLink(t *testing.T) {
	logger := slog.Default()
	hub := NewHub(nil, logger, nil, "")

	hub.AddTrack("")
	if _, exists := hub.linksCount[""]; exists {
		t.Fatalf("expected empty link not to be added")
	}

	hub.RemoveTrack("")
	if _, exists := hub.linksCount[""]; exists {
		t.Fatalf("expected empty link to be removed")
	}
}

func TestHubAddAndRemoveConcurrently(t *testing.T) {
	logger := slog.Default()
	hub := NewHub(nil, logger, nil, "")

	link := "http://example.com"

	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			hub.AddTrack(link)
		}()
	}

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			hub.RemoveTrack(link)
		}()
	}

	wg.Wait()

	if _, exists := hub.linksCount[link]; exists {
		t.Fatalf("expected link %s to be removed after concurrent add and remove", link)
	}
}

func TestHubRemoveNonExistentLink(t *testing.T) {
	logger := slog.Default()
	hub := NewHub(nil, logger, nil, "")

	link := "http://nonexistent.com"

	hub.RemoveTrack(link)

	if _, exists := hub.linksCount[link]; exists {
		t.Fatalf("expected link %s not to exist after remove", link)
	}
}


