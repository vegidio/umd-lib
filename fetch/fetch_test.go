package fetch

import (
	"bytes"
	"fmt"
	"github.com/cavaliergopher/grab/v3"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

type Test struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

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

	const RetryCount = 3
	fetch := New(nil, RetryCount)
	body, err := fetch.GetText(server.URL)

	// Check the log output
	output := buf.String()
	assert.Equal(t, RetryCount, strings.Count(output, "failed to get data; retrying in"))

	assert.Errorf(t, err, "429 Too Many Requests")
	assert.Equal(t, "", body)
}

func TestFetch_GetText_UserAgent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(r.UserAgent()))
	}))

	defer server.Close()

	fetch := New(nil, 0)
	body, err := fetch.GetText(server.URL)

	assert.NoError(t, err)
	assert.Contains(t, body, "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko)")
}

func TestFetch_GetText_Error(t *testing.T) {
	fetch := New(nil, 0)
	_, err := fetch.GetText("http://invalid-url")

	assert.Error(t, err)
}

func TestFetch_GetResult(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"name":"Vinicius","age":44}`))
	}))

	defer server.Close()

	var test Test
	fetch := New(nil, 0)
	_, err := fetch.GetResult(server.URL, nil, &test)

	assert.NoError(t, err)
	assert.Equal(t, Test{"Vinicius", 44}, test)
}

func TestFetch_GetResult_SetHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"name":"%s"}`, r.Header.Get("X-Name"))))
	}))

	defer server.Close()

	var test Test
	fetch := New(nil, 0)
	_, err := fetch.GetResult(server.URL, map[string]string{"X-Name": "Egidio"}, &test)

	assert.NoError(t, err)
	assert.Equal(t, Test{"Egidio", 0}, test)
}

func TestFetch_DownloadFile(t *testing.T) {
	// Delete any previous file before continuing
	const FilePath = "testfile.txt"
	_ = os.Remove(FilePath)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("file content"))
	}))

	defer server.Close()

	fetch := New(nil, 0)
	request, _ := grab.NewRequest(FilePath, server.URL)
	resp := fetch.DownloadFile(request)

	assert.NoError(t, resp.Err())
	assert.Equal(t, int64(len("file content")), resp.Size())
}

func TestFetch_DownloadFile_UserAgent(t *testing.T) {
	// Delete any previous file before continuing
	const FilePath = "testfile.txt"
	_ = os.Remove(FilePath)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(r.UserAgent()))
	}))

	defer server.Close()

	fetch := New(nil, 0)
	request, _ := grab.NewRequest(FilePath, server.URL)
	resp := fetch.DownloadFile(request)

	assert.NoError(t, resp.Err())

	byteArray, _ := resp.Bytes()
	assert.Contains(t, string(byteArray), "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko)")
}

func TestFetch_DownloadFile_Error(t *testing.T) {
	// Delete any previous file before continuing
	const FilePath = "testfile.txt"
	_ = os.Remove(FilePath)

	fetch := New(nil, 0)
	request, _ := grab.NewRequest(FilePath, "http://invalid-url")
	resp := fetch.DownloadFile(request)

	assert.Error(t, resp.Err())
}
