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
	Websites      map[string]time.Duration
	urls          []string
	MinAccessTime model.Website
	MaxAccessTime model.Website
	mu            sync.RWMutex
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
		if line == "" {
			continue
		}

		if !strings.HasPrefix(line, "http://") && !strings.HasPrefix(line, "https://") {
			line = "http://" + line // Add "http://" if not present
		}

		parsedURL, err := url.Parse(line)
		if err != nil {
			log.Printf("Invalid URL: %s", line)
			continue
		}

		wc.urls = append(wc.urls, parsedURL.String())
		log.Printf("Loaded website: %s", parsedURL.String())
	}

	return nil
}

func (wc *WebsiteChecker) CheckerWithTicker() {
	t := time.NewTicker(time.Minute)
	for range t.C {
		wc.CheckAvailability()
	}
}

func (wc *WebsiteChecker) CheckAvailability() {
	wc.MinAccessTime.AccessTime = time.Minute
	ch := make(chan string)
	go func() {
		for _, url := range wc.urls {
			ch <- url
		}
		close(ch)
	}()
	numWorker := 16
	var wg sync.WaitGroup
	for i := 0; i < numWorker; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for url := range ch {
				accessTime := wc.checkWebsiteAvailability(url)
				wc.mu.Lock()
				if wc.MinAccessTime.AccessTime > accessTime && accessTime > -1 {
					wc.MinAccessTime.AccessTime = accessTime
					wc.MinAccessTime.URL = url
				}
				if wc.MaxAccessTime.AccessTime < accessTime {
					wc.MaxAccessTime.AccessTime = accessTime
					wc.MaxAccessTime.URL = url
				}
				wc.Websites[url] = accessTime
				wc.mu.Unlock()
			}
		}()
	}
	wg.Wait()
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
	if !ok {
		return nil, fmt.Errorf("Website not found")
	}
	if t == -1 {
		return nil, fmt.Errorf("Website not found")
	}
	return &model.Website{
		URL:        url,
		AccessTime: t,
	}, nil
}

func (wc *WebsiteChecker) GetMinAccessTime() model.Website {
	wc.mu.RLock()
	defer wc.mu.RUnlock()
	return wc.MinAccessTime
}

func (wc *WebsiteChecker) GetMaxAccessTime() model.Website {
	wc.mu.RLock()
	defer wc.mu.RUnlock()
	return wc.MaxAccessTime
}
