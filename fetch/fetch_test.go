package fetch

import (
	"bytes"
	"github.com/cavaliergopher/grab/v3"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFetch_GetText(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	}))

	defer server.Close()

	fetch := New(nil, 0)
	body, err := fetch.GetText(server.URL)

	assert.NoError(t, err)
	assert.Equal(t, "Hello, World!", body)
}

func TestFetch_GetText_TooManyRequests(t *testing.T) {
	// Create a buffer and redirect global log output to it
	var buf bytes.Buffer
	log.SetOutput(&buf)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte("Hello, World!"))
	}))

	defer server.Close()

	fetch := New(nil, 3)
	body, err := fetch.GetText(server.URL)

	// Check the log output
	output := buf.String()
	assert.Equal(t, strings.Count(output, "failed to get data; retrying in"), 3)

	assert.Errorf(t, err, "429 Too Many Requests")
	assert.Equal(t, "", body)
}

func TestFetch_GetText_UserAgent(t *testing.T) {
	fetch := New(nil, 0)
	body, err := fetch.GetText("https://httpbin.org/get")

	assert.NoError(t, err)
	assert.Contains(t, body, "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
}

func TestFetch_GetText_Error(t *testing.T) {
	fetch := New(nil, 0)
	_, err := fetch.GetText("http://invalid-url")

	assert.Error(t, err)
}

func TestFetch_DownloadFile(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("file content"))
	}))

	defer server.Close()

	fetch := New(nil, 0)
	request, _ := grab.NewRequest("testfile.txt", server.URL)
	resp := fetch.DownloadFile(request)

	assert.NoError(t, resp.Err())
	assert.Equal(t, int64(len("file content")), resp.Size())
}

func TestFetch_DownloadFile_TooManyRequests(t *testing.T) {
	// Create a buffer and redirect global log output to it
	var buf bytes.Buffer
	log.SetOutput(&buf)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte("file content"))
	}))

	defer server.Close()

	fetch := New(nil, 3)
	request, _ := grab.NewRequest("testfile.txt", server.URL)
	resp := fetch.DownloadFile(request)

	// Check the log output
	output := buf.String()
	assert.Equal(t, strings.Count(output, "failed to download file; retrying in"), 3)

	assert.Errorf(t, resp.Err(), "429 Too Many Requests")
	assert.Equal(t, int64(0), resp.Size())
}

func TestFetch_DownloadFile_Error(t *testing.T) {
	fetch := New(nil, 0)
	request, _ := grab.NewRequest("testfile.txt", "http://invalid-url")
	resp := fetch.DownloadFile(request)

	assert.Error(t, resp.Err())
}
