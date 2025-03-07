package hub

import (
	"fmt"
	"log/slog"
	"sync"
	"testing"
)

func TestHub(t *testing.T) {
	logger := slog.Default()
	link := "test"
	userID := uint(1)

	hub := NewHub(nil, logger, nil, "")

	hub.AddTrack(link, userID)
	if len(hub.linksUsers[link]) != 1 {
		t.Fatalf("expected 1 user tracking link, got %d", len(hub.linksUsers[link]))
	}

	hub.AddTrack(link, 2)
	if len(hub.linksUsers[link]) != 2 {
		t.Fatalf("expected 2 users tracking link, got %d", len(hub.linksUsers[link]))
	}

	hub.RemoveTrack(link, userID)
	if len(hub.linksUsers[link]) != 1 {
		t.Fatalf("expected 1 user after removal, got %d", len(hub.linksUsers[link]))
	}

	hub.RemoveTrack(link, 2)
	if _, exists := hub.linksUsers[link]; exists {
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

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(uid uint) {
			defer wg.Done()
			h.AddTrack(link, uid)
		}(uint(i))
	}

	wg.Wait()

	if len(h.linksUsers[link]) == 0 {
		t.Fatalf("expected at least one user tracking, got %d", len(h.linksUsers[link]))
	}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(uid uint) {
			defer wg.Done()
			h.RemoveTrack(link, uid)
		}(uint(i))
	}

	wg.Wait()

	if _, exists := h.linksUsers[link]; exists {
		t.Fatalf("expected link to be fully removed after concurrent removals")
	}
}

func TestHubMaxLinks(t *testing.T) {
	logger := slog.Default()
	hub := NewHub(nil, logger, nil, "")

	for i := 0; i < 1000; i++ {
		link := fmt.Sprintf("http://example%d.com", i)
		hub.AddTrack(link, uint(i))
		if len(hub.linksUsers[link]) != 1 {
			t.Fatalf("expected 1 user tracking %s, got %d", link, len(hub.linksUsers[link]))
		}
	}

	for i := 0; i < 1000; i++ {
		link := fmt.Sprintf("http://example%d.com", i)
		if len(hub.linksUsers[link]) != 1 {
			t.Fatalf("link %s count is incorrect, expected 1, got %d", link, len(hub.linksUsers[link]))
		}
	}

	for i := 0; i < 1000; i++ {
		link := fmt.Sprintf("http://example%d.com", i)
		hub.RemoveTrack(link, uint(i))
		if _, exists := hub.linksUsers[link]; exists {
			t.Fatalf("expected link %s to be removed, but it still exists", link)
		}
	}
}

func TestHubEmptyLink(t *testing.T) {
	logger := slog.Default()
	hub := NewHub(nil, logger, nil, "")

	hub.AddTrack("", 1)
	if _, exists := hub.linksUsers[""]; exists {
		t.Fatalf("expected empty link not to be added")
	}

	hub.RemoveTrack("", 1)
	if _, exists := hub.linksUsers[""]; exists {
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
		go func(uid uint) {
			defer wg.Done()
			hub.AddTrack(link, uid)
		}(uint(i))
	}

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(uid uint) {
			defer wg.Done()
			hub.RemoveTrack(link, uid)
		}(uint(i))
	}

	wg.Wait()

	if _, exists := hub.linksUsers[link]; exists {
		t.Fatalf("expected link %s to be removed after concurrent add and remove", link)
	}
}

func TestHubRemoveNonExistentLink(t *testing.T) {
	logger := slog.Default()
	hub := NewHub(nil, logger, nil, "")
	link := "http://nonexistent.com"

	hub.RemoveTrack(link, 1)

	if _, exists := hub.linksUsers[link]; exists {
		t.Fatalf("expected link %s not to exist after remove", link)
	}
}
