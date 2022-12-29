package providers_api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestGetAllProviders(t *testing.T) {
	handler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		file, err := os.Open("../../../testdata/providers.json")
		if err != nil {
			t.Fatalf("Could not open file: %s", err)
		}
		if _, err = io.Copy(writer, file); err != nil {
			t.Fatalf("Could not copy file content: %s", err)
		}
	})

	t.Run("should return all providers", func(t *testing.T) {
		server := httptest.NewServer(handler)
		expectedProviders := 51
		mockClient := New(server.URL)
		providers, _ := mockClient.GetAllProviders()

		if len(providers) != expectedProviders {
			t.Fatalf("Expected %d providers, got %d", expectedProviders, len(providers))
		}
	})
}

func TestGetActiveProviders(t *testing.T) {
	handler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		file, err := os.Open("../../../testdata/providers.json")
		if err != nil {
			t.Fatalf("Could not open file: %s", err)
		}
		if _, err = io.Copy(writer, file); err != nil {
			t.Fatalf("Could not copy file content: %s", err)
		}
	})

	t.Run("should return active providers", func(t *testing.T) {
		server := httptest.NewServer(handler)
		expectedProviders := 41
		mockClient := New(server.URL)
		providers, _ := mockClient.GetActiveProviders()

		if len(providers) != expectedProviders {
			t.Fatalf("Expected %d providers, got %d", expectedProviders, len(providers))
		}
	})

	t.Run("should break if its not 200 OK", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(500)
		}))

		mockClient := New(server.URL)

		if _, err := mockClient.GetActiveProviders(); err == nil {
			t.Fatalf("Expected an error")
		}
	})
}
