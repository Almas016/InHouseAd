package app

import (
	"InHouseAd/internal/model"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

type WebsiteChecker struct {
	Websites map[string]time.Duration
	mu       sync.RWMutex
}

func NewWebsiteChecker() *WebsiteChecker {
	return &WebsiteChecker{
		Websites: make(map[string]time.Duration),
	}
}

func (wc *WebsiteChecker) LoadWebsitesFromFile(filename string) error {
	urls, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	lines := strings.Split(string(urls), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			if !strings.HasPrefix(line, "http://") && !strings.HasPrefix(line, "https://") {
				line = "http://" + line // Add "http://" if not present
			}

			parsedURL, err := url.Parse(line)
			if err != nil {
				log.Printf("Invalid URL: %s", line)
				continue
			}

			wc.mu.Lock()
			wc.Websites[parsedURL.String()] = -1
			wc.mu.Unlock()
			log.Printf("Loaded website: %s", parsedURL.String())
		}
	}

	return nil
}

func (wc *WebsiteChecker) CheckAvailability() {
	var wg sync.WaitGroup
	for {
		wc.mu.RLock()
		for url := range wc.Websites {
			wg.Add(1)
			go func(u string) {
				accessTime := wc.checkWebsiteAvailability(u)
				wc.mu.Lock()
				wc.Websites[u] = accessTime
				wc.mu.Unlock()
				wg.Done()
			}(url)
		}
		wc.mu.RUnlock()

		wg.Wait()
		time.Sleep(time.Minute)
	}
}

func (wc *WebsiteChecker) checkWebsiteAvailability(url string) time.Duration {
	startTime := time.Now()
	_, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to access %s: %s", url, err)
		return -1
	}

	accessTime := time.Since(startTime)
	log.Printf("Access time to %s: %s", url, accessTime)
	return accessTime
}

func (wc *WebsiteChecker) GetAccessTime(url string) (*model.Website, error) {
	wc.mu.RLock()
	defer wc.mu.RUnlock()
	t, ok := wc.Websites[url]
	if ok {
		return &model.Website{
			URL:        url,
			AccessTime: t,
		}, nil
	}
	return nil, fmt.Errorf("Website not found")
}

func (wc *WebsiteChecker) GetMinAccessTime() model.Website {
	wc.mu.RLock()
	defer wc.mu.RUnlock()
	var minWeb model.Website
	firstIteration := true
	for url, accessTime := range wc.Websites {
		if (firstIteration || accessTime < minWeb.AccessTime) && accessTime != -1 {
			minWeb = model.Website{
				URL:        url,
				AccessTime: accessTime,
			}
			firstIteration = false
		}
	}
	return minWeb
}

func (wc *WebsiteChecker) GetMaxAccessTime() model.Website {
	wc.mu.RLock()
	defer wc.mu.RUnlock()
	var maxWeb model.Website
	firstIteration := true
	for url, accessTime := range wc.Websites {
		if firstIteration || accessTime > maxWeb.AccessTime {
			maxWeb = model.Website{
				URL:        url,
				AccessTime: accessTime,
			}
			firstIteration = false
		}
	}
	return maxWeb
}
