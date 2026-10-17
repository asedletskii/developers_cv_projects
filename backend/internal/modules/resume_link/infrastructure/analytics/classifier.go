package analytics

import (
	"container/list"
	"context"
	"strings"
	"sync"
	"time"

	"developers-cv/backend/internal/modules/resume_link/application/ports"
	"developers-cv/backend/internal/modules/resume_link/domain"
	sharedvalueobjects "developers-cv/backend/internal/shared/domain/valueobjects"
)

const defaultSuspectCacheCapacity = 10_000

type burstEntry struct {
	ipHash []byte
	links  map[sharedvalueobjects.ID]time.Time
}

type TrafficClassifier struct {
	burst    int
	window   time.Duration
	capacity int
	mu       sync.Mutex
	items    map[string]*list.Element
	lru      *list.List
}

func NewTrafficClassifier(burst int, window time.Duration) *TrafficClassifier {
	if burst < 1 {
		burst = 5
	}
	if window <= 0 {
		window = time.Minute
	}
	return &TrafficClassifier{
		burst: burst, window: window, capacity: defaultSuspectCacheCapacity,
		items: make(map[string]*list.Element), lru: list.New(),
	}
}

func (c *TrafficClassifier) Classify(_ context.Context, input ports.TrafficClassificationInput) (ports.TrafficClassificationResult, error) {
	if !input.OwnerID.IsZero() {
		return classifyResult(domain.TrafficOwner, "owner", input.UserAgent), nil
	}
	if isBotUA(input.UserAgent) {
		return classifyResult(domain.TrafficBot, "ua:"+botReason(input.UserAgent), input.UserAgent), nil
	}
	if c != nil && c.recordAndCheckBurst(input.IPHash, input.LinkID, input.Now) {
		return classifyResult(domain.TrafficSuspect, "burst", input.UserAgent), nil
	}
	return classifyResult(domain.TrafficHuman, "", input.UserAgent), nil
}

func (c *TrafficClassifier) recordAndCheckBurst(ipHash []byte, linkID sharedvalueobjects.ID, now time.Time) bool {
	if c == nil || len(ipHash) != 32 || linkID.IsZero() || now.IsZero() {
		return false
	}
	key := string(ipHash)
	cutoff := now.UTC().Add(-c.window)

	c.mu.Lock()
	defer c.mu.Unlock()

	var entry *burstEntry
	if element, ok := c.items[key]; ok {
		entry, ok = element.Value.(*burstEntry)
		if !ok || entry == nil {
			delete(c.items, key)
			c.lru.Remove(element)
			entry = nil
		} else {
			c.lru.MoveToFront(element)
		}
	}
	if entry == nil {
		entry = &burstEntry{ipHash: append([]byte(nil), ipHash...), links: make(map[sharedvalueobjects.ID]time.Time)}
		element := c.lru.PushFront(entry)
		c.items[key] = element
		if c.lru.Len() > c.capacity {
			old := c.lru.Back()
			if old != nil {
				if oldEntry, ok := old.Value.(*burstEntry); ok && oldEntry != nil {
					delete(c.items, string(oldEntry.ipHash))
				}
				c.lru.Remove(old)
			}
		}
	}

	for id, seenAt := range entry.links {
		if seenAt.Before(cutoff) {
			delete(entry.links, id)
		}
	}
	entry.links[linkID] = now.UTC()
	return len(entry.links) > c.burst
}

func isBotUA(ua string) bool {
	ua = strings.ToLower(strings.TrimSpace(ua))
	if ua == "" {
		return true
	}

	markers := []string{
		"googlebot", "bingbot", "yandexbot", "duckduckbot", "baiduspider",
		"facebookexternalhit", "facebookcatalog", "linkedinbot", "twitterbot",
		"slackbot", "telegrambot", "discordbot", "whatsapp", "bingpreview",
		"crawler", "spider", "slurp", "headless", "phantomjs", "curl/", "wget/",
		"python-requests", "httpclient", "go-http-client",
	}
	for _, marker := range markers {
		if strings.Contains(ua, marker) {
			return true
		}
	}
	return hasUserAgentToken(ua, "bot")
}

func hasUserAgentToken(ua, token string) bool {
	for start := 0; start < len(ua); {
		relative := strings.Index(ua[start:], token)
		if relative < 0 {
			return false
		}
		index := start + relative
		end := index + len(token)
		leftBoundary := index == 0 || !isUserAgentWordByte(ua[index-1])
		rightBoundary := end == len(ua) || !isUserAgentWordByte(ua[end])
		if leftBoundary && rightBoundary {
			return true
		}
		start = index + len(token)
	}
	return false
}

func isUserAgentWordByte(value byte) bool {
	return value >= 'a' && value <= 'z' || value >= '0' && value <= '9'
}

func botReason(ua string) string {
	ua = strings.ToLower(ua)
	for _, marker := range []string{
		"googlebot", "bingbot", "yandexbot", "duckduckbot", "baiduspider",
		"facebookexternalhit", "facebookcatalog", "linkedinbot", "twitterbot",
		"slackbot", "telegrambot", "discordbot", "whatsapp", "bingpreview",
		"crawler", "spider", "slurp", "headless", "phantomjs", "curl/", "wget/",
		"python-requests", "httpclient", "go-http-client",
	} {
		if strings.Contains(ua, marker) {
			return marker
		}
	}
	if hasUserAgentToken(ua, "bot") {
		return "bot"
	}
	return "automated"
}

func classifyResult(traffic domain.TrafficType, reason, ua string) ports.TrafficClassificationResult {
	return ports.TrafficClassificationResult{Traffic: traffic, Reason: reason, Device: classifyDevice(ua), Browser: classifyBrowser(ua), OS: classifyOS(ua)}
}
func classifyDevice(ua string) domain.DeviceType {
	ua = strings.ToLower(ua)
	if strings.TrimSpace(ua) == "" {
		return domain.DeviceUnknown
	}
	if strings.Contains(ua, "ipad") || strings.Contains(ua, "tablet") || strings.Contains(ua, "kindle") || strings.Contains(ua, "silk/") {
		return domain.DeviceTablet
	}
	if strings.Contains(ua, "macintosh") && strings.Contains(ua, "mobile/") {
		return domain.DeviceTablet
	}
	if strings.Contains(ua, "android") && !strings.Contains(ua, "mobile") {
		return domain.DeviceTablet
	}
	if strings.Contains(ua, "mobile") || strings.Contains(ua, "iphone") || strings.Contains(ua, "android") {
		return domain.DeviceMobile
	}
	return domain.DeviceDesktop
}
func classifyBrowser(ua string) domain.BrowserType {
	ua = strings.ToLower(ua)
	switch {
	case strings.Contains(ua, "yabrowser"):
		return domain.BrowserYandex
	case strings.Contains(ua, "edg/") || strings.Contains(ua, "edge/"):
		return domain.BrowserEdge
	case strings.Contains(ua, "firefox/"):
		return domain.BrowserFirefox
	case strings.Contains(ua, "chrome/") || strings.Contains(ua, "chromium/"):
		return domain.BrowserChrome
	case strings.Contains(ua, "safari/") && !strings.Contains(ua, "chrome/") && !strings.Contains(ua, "chromium/"):
		return domain.BrowserSafari
	default:
		return domain.BrowserOther
	}
}
func classifyOS(ua string) domain.OSType {
	ua = strings.ToLower(ua)
	switch {
	case strings.Contains(ua, "iphone") || strings.Contains(ua, "ipad") || strings.Contains(ua, "ipod"):
		return domain.OSIOS
	case strings.Contains(ua, "macintosh") && strings.Contains(ua, "mobile/"):
		return domain.OSIOS
	case strings.Contains(ua, "android"):
		return domain.OSAndroid
	case strings.Contains(ua, "windows"):
		return domain.OSWindows
	case strings.Contains(ua, "mac os x") || strings.Contains(ua, "macintosh"):
		return domain.OSMacOS
	case strings.Contains(ua, "linux"):
		return domain.OSLinux
	default:
		return domain.OSOther
	}
}
