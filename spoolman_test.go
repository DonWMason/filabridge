package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type recordedRequest struct {
	Method string
	Path   string
	Auth   string
	Body   map[string]interface{}
}

// newFakeFilaMan serves FilaMan's native consumption endpoint (returning consumptionStatus)
// and the Spoolman-compatible /use endpoint, recording every request it receives.
func newFakeFilaMan(t *testing.T, consumptionStatus int) (*httptest.Server, *[]recordedRequest) {
	t.Helper()
	var requests []recordedRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]interface{}
		_ = json.NewDecoder(r.Body).Decode(&body)
		requests = append(requests, recordedRequest{r.Method, r.URL.Path, r.Header.Get("Authorization"), body})

		switch r.URL.Path {
		case "/api/v1/spools/7/consumptions":
			w.WriteHeader(consumptionStatus)
			_, _ = w.Write([]byte(`{}`))
		case "/spoolman/api/v1/spool/7/use":
			_, _ = w.Write([]byte(`{}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	t.Cleanup(server.Close)
	return server, &requests
}

func TestUpdateSpoolUsageUsesFilaManConsumptionsWithNote(t *testing.T) {
	server, requests := newFakeFilaMan(t, http.StatusOK)
	client := NewSpoolmanClient(server.URL+"/spoolman/", 5, "", "", "uak.1.secret")

	if err := client.UpdateSpoolUsage(7, 12.5, "Prusa XL: benchy.bgcode"); err != nil {
		t.Fatalf("UpdateSpoolUsage returned error: %v", err)
	}

	if len(*requests) != 1 {
		t.Fatalf("expected 1 request, got %d: %+v", len(*requests), *requests)
	}
	req := (*requests)[0]
	if req.Method != "POST" || req.Path != "/api/v1/spools/7/consumptions" {
		t.Errorf("unexpected request %s %s", req.Method, req.Path)
	}
	if req.Auth != "ApiKey uak.1.secret" {
		t.Errorf("unexpected Authorization header %q", req.Auth)
	}
	if req.Body["delta_weight_g"] != 12.5 || req.Body["note"] != "Prusa XL: benchy.bgcode" {
		t.Errorf("unexpected body %+v", req.Body)
	}
}

func TestUpdateSpoolUsageFallsBackToUseWhenFilaManRejects(t *testing.T) {
	server, requests := newFakeFilaMan(t, http.StatusUnauthorized)
	client := NewSpoolmanClient(server.URL+"/spoolman", 5, "", "", "bad-key")

	if err := client.UpdateSpoolUsage(7, 3, "Prusa XL: benchy.bgcode"); err != nil {
		t.Fatalf("UpdateSpoolUsage returned error: %v", err)
	}

	if len(*requests) != 2 || (*requests)[1].Path != "/spoolman/api/v1/spool/7/use" {
		t.Fatalf("expected consumptions then /use fallback, got %+v", *requests)
	}
	if (*requests)[1].Body["use_weight"] != 3.0 {
		t.Errorf("unexpected /use body %+v", (*requests)[1].Body)
	}
}

func TestUpdateSpoolUsageDoesNotFallBackOnServerError(t *testing.T) {
	server, requests := newFakeFilaMan(t, http.StatusInternalServerError)
	client := NewSpoolmanClient(server.URL+"/spoolman", 5, "", "", "uak.1.secret")

	if err := client.UpdateSpoolUsage(7, 3, "note"); err == nil {
		t.Fatal("expected error on FilaMan 500")
	}
	if len(*requests) != 1 {
		t.Fatalf("expected no /use fallback after a 500 (could double-count), got %+v", *requests)
	}
}

func TestUpdateSpoolUsageWithoutAPIKeyUsesSpoolmanUse(t *testing.T) {
	server, requests := newFakeFilaMan(t, http.StatusOK)
	client := NewSpoolmanClient(server.URL+"/spoolman", 5, "", "", "")

	if err := client.UpdateSpoolUsage(7, 4, "ignored"); err != nil {
		t.Fatalf("UpdateSpoolUsage returned error: %v", err)
	}
	if len(*requests) != 1 || (*requests)[0].Path != "/spoolman/api/v1/spool/7/use" {
		t.Fatalf("expected only /use, got %+v", *requests)
	}
}
