package proxy_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"go.uber.org/zap"

	"{{{ .Package }}}/app/lib/proxy"
	"{{{ .Package }}}/app/util"
)

func TestNewService(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		urlPrefix      string
		initialProxies map[string]string
		expectedList   []string
	}{
		{
			name:           "empty proxies",
			urlPrefix:      "/proxy",
			initialProxies: map[string]string{},
			expectedList:   []string{},
		},
		{
			name:           "with initial proxies",
			urlPrefix:      "/api/proxy",
			initialProxies: map[string]string{"svc1": "http://localhost:8080", "svc2": "http://localhost:9090"},
			expectedList:   []string{"svc1", "svc2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := proxy.NewService(tt.urlPrefix, tt.initialProxies)
			if svc == nil {
				t.Fatal("NewService() returned nil")
			}
			list := svc.List()
			if !reflect.DeepEqual(list, tt.expectedList) {
				t.Errorf("NewService().List() = %v, want %v", list, tt.expectedList)
			}
		})
	}
}

func TestService_SetURL(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		initialProxies map[string]string
		setService     string
		setURL         string
		expectedList   []string
	}{
		{
			name:           "add new service",
			initialProxies: map[string]string{},
			setService:     "newsvc",
			setURL:         "http://localhost:3000",
			expectedList:   []string{"newsvc"},
		},
		{
			name:           "update existing service",
			initialProxies: map[string]string{"existing": "http://old-url"},
			setService:     "existing",
			setURL:         "http://new-url",
			expectedList:   []string{"existing"},
		},
		{
			name:           "add to existing proxies",
			initialProxies: map[string]string{"svc1": "http://localhost:8080"},
			setService:     "svc2",
			setURL:         "http://localhost:9090",
			expectedList:   []string{"svc1", "svc2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := proxy.NewService("/proxy", copyMap(tt.initialProxies))
			svc.SetURL(tt.setService, tt.setURL)
			list := svc.List()
			if !reflect.DeepEqual(list, tt.expectedList) {
				t.Errorf("After SetURL(), List() = %v, want %v", list, tt.expectedList)
			}
		})
	}
}

func TestService_Remove(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		initialProxies map[string]string
		removeService  string
		expectedList   []string
	}{
		{
			name:           "remove existing service",
			initialProxies: map[string]string{"svc1": "http://localhost:8080", "svc2": "http://localhost:9090"},
			removeService:  "svc1",
			expectedList:   []string{"svc2"},
		},
		{
			name:           "remove non-existing service",
			initialProxies: map[string]string{"svc1": "http://localhost:8080"},
			removeService:  "nonexistent",
			expectedList:   []string{"svc1"},
		},
		{
			name:           "remove from empty",
			initialProxies: map[string]string{},
			removeService:  "any",
			expectedList:   []string{},
		},
		{
			name:           "remove last service",
			initialProxies: map[string]string{"only": "http://localhost:8080"},
			removeService:  "only",
			expectedList:   []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := proxy.NewService("/proxy", copyMap(tt.initialProxies))
			svc.Remove(tt.removeService)
			list := svc.List()
			if !reflect.DeepEqual(list, tt.expectedList) {
				t.Errorf("After Remove(), List() = %v, want %v", list, tt.expectedList)
			}
		})
	}
}

func TestService_List(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		initialProxies map[string]string
		expectedList   []string
	}{
		{
			name:           "empty list",
			initialProxies: map[string]string{},
			expectedList:   []string{},
		},
		{
			name:           "single item",
			initialProxies: map[string]string{"alpha": "http://localhost:8080"},
			expectedList:   []string{"alpha"},
		},
		{
			name:           "multiple items sorted",
			initialProxies: map[string]string{"charlie": "http://c", "alpha": "http://a", "bravo": "http://b"},
			expectedList:   []string{"alpha", "bravo", "charlie"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := proxy.NewService("/proxy", copyMap(tt.initialProxies))
			list := svc.List()
			if !reflect.DeepEqual(list, tt.expectedList) {
				t.Errorf("List() = %v, want %v", list, tt.expectedList)
			}
		})
	}
}

func TestService_Handle_SimpleResponse(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Hello, World!"))
	}))
	defer server.Close()

	svc := proxy.NewService("/proxy", map[string]string{"testsvc": server.URL})
	req := httptest.NewRequest(http.MethodGet, "/proxy/testsvc/test", http.NoBody)
	w := httptest.NewRecorder()

	err := svc.Handle(context.Background(), "testsvc", w, req, "/test", nopLogger())
	if err != nil {
		t.Fatalf("Handle() unexpected error: %v", err)
	}

	resp := w.Result()
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if string(body) != "Hello, World!" {
		t.Errorf("Handle() body = %q, want %q", string(body), "Hello, World!")
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Handle() status = %d, want %d", resp.StatusCode, http.StatusOK)
	}
}

func TestService_Handle_URLRewriting(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		response     string
		expectedBody string
	}{
		{
			name:         "href rewriting",
			response:     `<a href="/page">Link</a>`,
			expectedBody: `<a href="/proxy/testsvc/page">Link</a>`,
		},
		{
			name:         "src rewriting",
			response:     `<img src="/image.png">`,
			expectedBody: `<img src="/proxy/testsvc/image.png">`,
		},
		{
			name:         "multiple rewrites",
			response:     `<a href="/link1"></a><img src="/img1"><a href="/link2"></a>`,
			expectedBody: `<a href="/proxy/testsvc/link1"></a><img src="/proxy/testsvc/img1"><a href="/proxy/testsvc/link2"></a>`,
		},
		{
			name:         "relative paths preserved",
			response:     `<a href="relative">Link</a>`,
			expectedBody: `<a href="relative">Link</a>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "text/html")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(tt.response))
			}))
			defer server.Close()

			svc := proxy.NewService("/proxy", map[string]string{"testsvc": server.URL})
			req := httptest.NewRequest(http.MethodGet, "/proxy/testsvc/", http.NoBody)
			w := httptest.NewRecorder()

			err := svc.Handle(context.Background(), "testsvc", w, req, "/", nopLogger())
			if err != nil {
				t.Fatalf("Handle() unexpected error: %v", err)
			}

			resp := w.Result()
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)

			if string(body) != tt.expectedBody {
				t.Errorf("Handle() body = %q, want %q", string(body), tt.expectedBody)
			}
		})
	}
}

func TestService_Handle_ErrorStatus(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("Not Found"))
	}))
	defer server.Close()

	svc := proxy.NewService("/proxy", map[string]string{"testsvc": server.URL})
	req := httptest.NewRequest(http.MethodGet, "/proxy/testsvc/missing", http.NoBody)
	w := httptest.NewRecorder()

	err := svc.Handle(context.Background(), "testsvc", w, req, "/missing", nopLogger())
	if err != nil {
		t.Fatalf("Handle() unexpected error: %v", err)
	}

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Handle() status = %d, want %d", resp.StatusCode, http.StatusNotFound)
	}
}

func TestService_Handle_UnregisteredService(t *testing.T) {
	t.Parallel()

	svc := proxy.NewService("/proxy", map[string]string{})
	req := httptest.NewRequest(http.MethodGet, "/proxy/unknown/path", http.NoBody)
	w := httptest.NewRecorder()

	err := svc.Handle(context.Background(), "unknown", w, req, "/path", nopLogger())

	if err == nil {
		t.Error("Handle() expected error for unregistered service but got nil")
	}
	if !strings.Contains(err.Error(), "unknown") {
		t.Errorf("Handle() error = %q, want error containing 'unknown'", err.Error())
	}
}

func TestService_Handle_HeaderPropagation(t *testing.T) {
	t.Parallel()

	var receivedHeaders http.Header
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedHeaders = r.Header
		w.Header().Set("X-Response-Header", "response-value")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	}))
	defer server.Close()

	svc := proxy.NewService("/proxy", map[string]string{"testsvc": server.URL})

	req := httptest.NewRequest(http.MethodGet, "/proxy/testsvc/test", http.NoBody)
	req.Header.Set("X-Custom-Header", "custom-value")
	req.Header.Set("Accept", "application/json")
	w := httptest.NewRecorder()

	err := svc.Handle(context.Background(), "testsvc", w, req, "/test", nopLogger())
	if err != nil {
		t.Fatalf("Handle() unexpected error: %v", err)
	}

	if receivedHeaders.Get("X-Custom-Header") != "custom-value" {
		t.Errorf("Request header X-Custom-Header = %q, want %q", receivedHeaders.Get("X-Custom-Header"), "custom-value")
	}
	if receivedHeaders.Get("Proxied") == "" {
		t.Error("Request header Proxied was not set")
	}

	resp := w.Result()
	defer resp.Body.Close()

	if resp.Header.Get("X-Response-Header") != "response-value" {
		t.Errorf("Response header X-Response-Header = %q, want %q", resp.Header.Get("X-Response-Header"), "response-value")
	}
}

func TestService_Handle_Methods(t *testing.T) {
	t.Parallel()

	methods := []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			t.Parallel()

			var receivedMethod string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				receivedMethod = r.Method
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("OK"))
			}))
			defer server.Close()

			svc := proxy.NewService("/proxy", map[string]string{"testsvc": server.URL})

			req := httptest.NewRequest(method, "/proxy/testsvc/test", http.NoBody)
			w := httptest.NewRecorder()

			err := svc.Handle(context.Background(), "testsvc", w, req, "/test", nopLogger())
			if err != nil {
				t.Fatalf("Handle() unexpected error: %v", err)
			}

			if receivedMethod != method {
				t.Errorf("Handle() method = %q, want %q", receivedMethod, method)
			}
		})
	}
}

func TestService_Handle_RequestBody(t *testing.T) {
	t.Parallel()

	var receivedBody string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		receivedBody = string(body)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	}))
	defer server.Close()

	svc := proxy.NewService("/proxy", map[string]string{"testsvc": server.URL})

	requestBody := `{"key": "value"}`
	req := httptest.NewRequest(http.MethodPost, "/proxy/testsvc/test", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	err := svc.Handle(context.Background(), "testsvc", w, req, "/test", nopLogger())
	if err != nil {
		t.Fatalf("Handle() unexpected error: %v", err)
	}

	if receivedBody != requestBody {
		t.Errorf("Handle() request body = %q, want %q", receivedBody, requestBody)
	}
}

func TestService_Handle_PathNormalization(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		inputPath    string
		expectedPath string
	}{
		{name: "path with leading slash", inputPath: "/api/endpoint", expectedPath: "/api/endpoint"},
		{name: "path without leading slash", inputPath: "api/endpoint", expectedPath: "/api/endpoint"},
		{name: "root path", inputPath: "/", expectedPath: "/"},
		{name: "empty path", inputPath: "", expectedPath: "/"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var receivedPath string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				receivedPath = r.URL.Path
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("OK"))
			}))
			defer server.Close()

			svc := proxy.NewService("/proxy", map[string]string{"testsvc": server.URL})

			req := httptest.NewRequest(http.MethodGet, "/proxy/testsvc"+tt.expectedPath, http.NoBody)
			w := httptest.NewRecorder()

			err := svc.Handle(context.Background(), "testsvc", w, req, tt.inputPath, nopLogger())
			if err != nil {
				t.Fatalf("Handle() unexpected error: %v", err)
			}

			if receivedPath != tt.expectedPath {
				t.Errorf("Handle() received path = %q, want %q", receivedPath, tt.expectedPath)
			}
		})
	}
}

// nopLogger returns a no-op logger for testing.
func nopLogger() util.Logger {
	return zap.NewNop().Sugar()
}

// copyMap creates a copy of a string map to avoid test interference.
func copyMap(m map[string]string) map[string]string {
	result := make(map[string]string, len(m))
	for k, v := range m {
		result[k] = v
	}
	return result
}
