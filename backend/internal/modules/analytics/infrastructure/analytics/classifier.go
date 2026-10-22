package analytics

import (
	"container/list"
	"context"
	"strings"
	"sync"
	"time"

	analyticsports "developers-cv/backend/internal/modules/analytics/application/ports"
	analyticsdomain "developers-cv/backend/internal/modules/analytics/domain"
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

func (c *TrafficClassifier) Classify(_ context.Context, input analyticsports.TrafficClassificationInput) (analyticsports.TrafficClassificationResult, error) {
	if !input.OwnerID.IsZero() {
		return classifyResult(analyticsdomain.TrafficOwner, "owner", input.UserAgent), nil
	}
	if isBotUA(input.UserAgent) {
		return classifyResult(analyticsdomain.TrafficBot, "ua:"+botReason(input.UserAgent), input.UserAgent), nil
	}
	if c != nil && c.recordAndCheckBurst(input.IPHash, input.LinkID, input.Now) {
		return classifyResult(analyticsdomain.TrafficSuspect, "burst", input.UserAgent), nil
	}
	return classifyResult(analyticsdomain.TrafficHuman, "", input.UserAgent), nil
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

func classifyResult(traffic analyticsdomain.TrafficType, reason, ua string) analyticsports.TrafficClassificationResult {
	return analyticsports.TrafficClassificationResult{Traffic: traffic, Reason: reason, Device: classifyDevice(ua), Browser: classifyBrowser(ua), OS: classifyOS(ua)}
}
func classifyDevice(ua string) analyticsdomain.DeviceType {
	ua = strings.ToLower(ua)
	if strings.TrimSpace(ua) == "" {
		return analyticsdomain.DeviceUnknown
	}
	if strings.Contains(ua, "ipad") || strings.Contains(ua, "tablet") || strings.Contains(ua, "kindle") || strings.Contains(ua, "silk/") {
		return analyticsdomain.DeviceTablet
	}
	if strings.Contains(ua, "macintosh") && strings.Contains(ua, "mobile/") {
		return analyticsdomain.DeviceTablet
	}
	if strings.Contains(ua, "android") && !strings.Contains(ua, "mobile") {
		return analyticsdomain.DeviceTablet
	}
	if strings.Contains(ua, "mobile") || strings.Contains(ua, "iphone") || strings.Contains(ua, "android") {
		return analyticsdomain.DeviceMobile
	}
	return analyticsdomain.DeviceDesktop
}
func classifyBrowser(ua string) analyticsdomain.BrowserType {
	ua = strings.ToLower(ua)
	switch {
	case strings.Contains(ua, "yabrowser"):
		return analyticsdomain.BrowserYandex
	case strings.Contains(ua, "edg/") || strings.Contains(ua, "edge/"):
		return analyticsdomain.BrowserEdge
	case strings.Contains(ua, "firefox/"):
		return analyticsdomain.BrowserFirefox
	case strings.Contains(ua, "chrome/") || strings.Contains(ua, "chromium/"):
		return analyticsdomain.BrowserChrome
	case strings.Contains(ua, "safari/") && !strings.Contains(ua, "chrome/") && !strings.Contains(ua, "chromium/"):
		return analyticsdomain.BrowserSafari
	default:
		return analyticsdomain.BrowserOther
	}
}
func classifyOS(ua string) analyticsdomain.OSType {
	ua = strings.ToLower(ua)
	switch {
	case strings.Contains(ua, "iphone") || strings.Contains(ua, "ipad") || strings.Contains(ua, "ipod"):
		return analyticsdomain.OSIOS
	case strings.Contains(ua, "macintosh") && strings.Contains(ua, "mobile/"):
		return analyticsdomain.OSIOS
	case strings.Contains(ua, "android"):
		return analyticsdomain.OSAndroid
	case strings.Contains(ua, "windows"):
		return analyticsdomain.OSWindows
	case strings.Contains(ua, "mac os x") || strings.Contains(ua, "macintosh"):
		return analyticsdomain.OSMacOS
	case strings.Contains(ua, "linux"):
		return analyticsdomain.OSLinux
	default:
		return analyticsdomain.OSOther
	}
}
