package fetch

import (
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
	filePath := "testfile.txt"
	size, err := fetch.DownloadFile(server.URL, filePath)

	assert.NoError(t, err)
	assert.Equal(t, int64(len("file content")), size)
}

func TestFetch_DownloadFile_Error(t *testing.T) {
	fetch := New(nil, 0)
	_, err := fetch.DownloadFile("http://invalid-url", "testfile.txt")

	assert.Error(t, err)
}
