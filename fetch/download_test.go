package fetch

import (
	"fmt"
	"github.com/cavaliergopher/grab/v3"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

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

func TestFetch_DownloadFiles(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("file content"))
	}))

	defer server.Close()

	requests := lo.Map([]int{1, 2, 3}, func(i int, _ int) *grab.Request {
		r, _ := grab.NewRequest(fmt.Sprintf("testfile%d.txt", i), server.URL)
		return r
	})

	fetch := New(nil, 0)
	result := fetch.DownloadFiles(requests, 1)

	for resp := range result {
		assert.NoError(t, resp.Err())
		assert.Equal(t, int64(len("file content")), resp.Size())
	}
}
