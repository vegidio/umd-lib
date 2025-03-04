package fetch

import (
	"github.com/cavaliergopher/grab/v3"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
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

func TestFetch_DownloadFile_Error(t *testing.T) {
	fetch := New(nil, 0)
	request, _ := grab.NewRequest("testfile.txt", "http://invalid-url")
	resp := fetch.DownloadFile(request)

	assert.Error(t, resp.Err())
}
