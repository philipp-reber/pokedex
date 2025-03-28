package pokeapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestListLocations_Success(t *testing.T) {
	mockResponse := RespShallowLocations{
		Count: 1,
		Results: []struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		}{
			{Name: "pallet-town", URL: "https://pokeapi.co/api/v2/location-area/1"},
		},
	}
	respBody, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(respBody)
	}))
	defer server.Close()

	client := NewClient(5*time.Second, 10*time.Minute)
	locations, err := client.ListLocations(&server.URL)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if locations.Count != mockResponse.Count {
		t.Errorf("Expected Count %d, got: %d", mockResponse.Count, locations.Count)
	}
	if locations.Results[0].Name != mockResponse.Results[0].Name {
		t.Errorf("Expected Name %s, got: %s", mockResponse.Results[0].Name, locations.Results[0].Name)
	}
}
