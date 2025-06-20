package fetch

import (
	"bufio"
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/samber/lo"
	"os"
	"strings"
	"time"
)

// Cookie represents a key-value pair for a typical HTTP cookie.
type Cookie struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// GetBrowserCookies retrieves HTTP cookies from a given URL by navigating to the page using a headless browser.
// It returns a slice of Cookie objects or an error if the process fails.
func GetBrowserCookies(url string, element string) ([]Cookie, error) {
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

func GetFileCookies(filePath string) ([]Cookie, error) {
	cookies := make([]Cookie, 0)

	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		// Skip comments & blank lines
		if len(line) == 0 || line[0] == '#' {
			continue
		}

		// Spec says 7 TAB-separated fields
		parts := strings.Split(line, "\t")
		if len(parts) != 7 {
			continue
		}

		name := parts[5]
		value := parts[6]

		cookies = append(cookies, Cookie{
			Name:  name,
			Value: value,
		})
	}

	return cookies, nil
}

// CookiesToHeader converts a slice of Cookie objects into a single HTTP header string in the "key=value" format.
func CookiesToHeader(cookies []Cookie) string {
	parts := lo.Map(cookies, func(cookie Cookie, index int) string {
		return fmt.Sprintf("%s=%s", cookie.Name, cookie.Value)
	})

	return strings.Join(parts, "; ")
}

// region - Private functions

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

// endregion
