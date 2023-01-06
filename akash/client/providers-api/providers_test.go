package providers_api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestGetAllProviders(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/provider/", func(writer http.ResponseWriter, request *http.Request) {
		file, err := os.Open("../../../testdata/providers.json")
		if err != nil {
			t.Fatalf("Could not open file: %s", err)
		}
		if _, err = io.Copy(writer, file); err != nil {
			t.Fatalf("Could not copy file content: %s", err)
		}
	})

	t.Run("should return all providers", func(t *testing.T) {
		server := httptest.NewServer(mux)
		expectedProviders := 51
		mockClient := New(server.URL)
		defer server.Close()
		providers, err := mockClient.GetAllProviders()
		if err != nil {
			t.Fatal(err)
		}

		if len(providers) != expectedProviders {
			t.Fatalf("Expected %d providers, got %d", expectedProviders, len(providers))
		}
	})
}

func TestGetActiveProviders(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/provider/", func(writer http.ResponseWriter, request *http.Request) {
		file, err := os.Open("../../../testdata/providers.json")
		if err != nil {
			t.Fatalf("Could not open file: %s", err)
		}
		if _, err = io.Copy(writer, file); err != nil {
			t.Fatalf("Could not copy file content: %s", err)
		}
	})

	t.Run("should return active providers", func(t *testing.T) {
		server := httptest.NewServer(mux)
		expectedProviders := 41
		mockClient := New(server.URL)
		defer server.Close()
		providers, _ := mockClient.GetActiveProviders()

		if len(providers) != expectedProviders {
			t.Fatalf("Expected %d providers, got %d", expectedProviders, len(providers))
		}
	})

	t.Run("should break if its not 200 OK", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/provider/", func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(500)
		})
		server := httptest.NewServer(mux)
		mockClient := New(server.URL)
		defer server.Close()

		if _, err := mockClient.GetActiveProviders(); err == nil {
			t.Fatalf("Expected an error")
		}
	})
}
