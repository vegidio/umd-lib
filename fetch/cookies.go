package fetch

import (
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/samber/lo"
	"strings"
	"time"
)

// Cookie represents a key-value pair for a typical HTTP cookie.
type Cookie struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// GetCookies retrieves HTTP cookies from a given URL by navigating to the page using a headless browser.
// It returns a slice of Cookie objects or an error if the process fails.
func GetCookies(url string, element string) ([]Cookie, error) {
	browser := rod.New().MustConnect()
	defer browser.Close()

	page, err := browser.Page(proto.TargetCreateTarget{})
	if err != nil {
		return nil, fmt.Errorf("could not create page: %w", err)
	}

	defer page.Close()

	err = page.Timeout(10 * time.Second).Navigate(url)
	if err != nil {
		return nil, fmt.Errorf("could not navigate: %w", err)
	}

	el, err := page.Timeout(10 * time.Second).Element("header.user-header")
	if err == nil {
		el.MustWaitVisible()
		return extractCookies(page)
	}

	waitNav := page.Timeout(10 * time.Second).WaitNavigation(proto.PageLifecycleEventNameLoad)
	waitNav()

	el, err = page.Timeout(10 * time.Second).Element("header.user-header")
	if err != nil {
		return nil, fmt.Errorf("could not find element: %w", err)
	}

	return extractCookies(page)
}

func extractCookies(page *rod.Page) ([]Cookie, error) {
	rodCookies, err := page.Cookies(nil)
	if err != nil {
		return nil, fmt.Errorf("could not get cookies: %w", err)
	}

	return lo.Map(rodCookies, func(cookie *proto.NetworkCookie, index int) Cookie {
		return Cookie{
			Name:  cookie.Name,
			Value: cookie.Value,
		}
	}), nil
}

// CookiesToHeader converts a slice of Cookie objects into a single HTTP header string in the "key=value" format.
func CookiesToHeader(cookies []Cookie) string {
	parts := lo.Map(cookies, func(cookie Cookie, index int) string {
		return fmt.Sprintf("%s=%s", cookie.Name, cookie.Value)
	})

	return strings.Join(parts, "; ")
}
