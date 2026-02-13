package image

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestDownloadImage(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("fake image data"))
	}))
	defer ts.Close()

	// 使用主函数（默认配置）
	msg := downloadImage(ts.URL)

	switch msg := msg.(type) {
	case imageDownloadedMsg:
		if msg.path == "" {
			t.Error("Downloaded image path should not be empty")
		}
	case imageErrMsg:
		t.Errorf("Download failed: %v", msg.err)
	}
}

func TestDownloadImageTimeout(t *testing.T) {
	// 使用选项模式设置 1 秒超时
	msg := downloadImage("http://10.255.255.1/test.jpg", WithTimeout(1*time.Second))

	if _, ok := msg.(imageErrMsg); !ok {
		t.Errorf("Expected error for timeout, got %T", msg)
	}
}

func TestDownloadImageTooLarge(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "99999999999")
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	// 使用主函数（默认配置）
	msg := downloadImage(ts.URL)

	if _, ok := msg.(imageErrMsg); !ok {
		t.Error("Expected error for oversized image")
	}
}

func TestDownloadImageWithCustomConfig(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("data"))
	}))
	defer ts.Close()

	// 使用选项模式配置多个参数
	msg := downloadImage(ts.URL,
		WithTimeout(5*time.Second),
		WithMaxSize(1024*1024), // 1MB
	)

	switch msg := msg.(type) {
	case imageDownloadedMsg:
		if msg.path == "" {
			t.Error("Downloaded image path should not be empty")
		}
	case imageErrMsg:
		t.Errorf("Download failed: %v", msg.err)
	}
}

func TestImagePageCreation(t *testing.T) {
	page := NewImagePage()

	if page == nil {
		t.Fatal("NewImagePage should not return nil")
	}
}

func TestCleanup(t *testing.T) {
	page := NewImagePage()

	err := page.cleanup()
	if err != nil {
		t.Errorf("Cleanup should not return error: %v", err)
	}
}
