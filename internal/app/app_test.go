package app

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"strings"
	"testing"
)

func cli(t *testing.T, hostname string, args ...string) any {
	t.Helper()
	cmd := exec.Command("./bin/appctl", append([]string{"--hostname", hostname}, append(args, "-o", "json")...)...)
	cmd.Dir = "../.."
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("appctl %v: %v", args, err)
	}
	if len(output) == 0 {
		return nil
	}
	var value any
	if err := json.Unmarshal(output, &value); err != nil {
		t.Fatalf("decode appctl output: %v", err)
	}
	return value
}

func cliFailure(t *testing.T, hostname string, args ...string) int {
	t.Helper()
	cmd := exec.Command("./bin/appctl", append([]string{"--hostname", hostname}, append(args, "-o", "json")...)...)
	cmd.Dir = "../.."
	_, err := cmd.Output()
	if err == nil {
		t.Fatalf("appctl %v unexpectedly succeeded", args)
	}
	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("appctl %v: %v", args, err)
	}
	var response struct {
		Error struct {
			HTTP struct {
				Status int `json:"status"`
			} `json:"http"`
		} `json:"error"`
	}
	if err := json.Unmarshal(exitErr.Stderr, &response); err != nil {
		t.Fatalf("decode appctl error: %v\n%s", err, exitErr.Stderr)
	}
	return response.Error.HTTP.Status
}

func TestGeneratedCLIIsTheApplicationAcceptanceSurface(t *testing.T) {
	srv := httptest.NewServer(NewHandler())
	defer srv.Close()

	if got := cli(t, srv.URL, "health", "get"); !equalJSON(got, map[string]any{"status": "ok"}) {
		t.Fatalf("health = %#v", got)
	}
	created := cli(t, srv.URL, "tasks", "create", "--set", "title=Ship from the CLI").(map[string]any)
	if created["title"] != "Ship from the CLI" || created["completed"] != false {
		t.Fatalf("created = %#v", created)
	}
	id := created["id"].(string)
	if got := cli(t, srv.URL, "tasks", "list").([]any); len(got) != 1 {
		t.Fatalf("list = %#v", got)
	}
	if got := cli(t, srv.URL, "tasks", "get", "--id", id); !equalJSON(got, created) {
		t.Fatalf("get = %#v", got)
	}
	updated := cli(t, srv.URL, "tasks", "update", "--id", id, "--set", "completed=true").(map[string]any)
	if updated["completed"] != true {
		t.Fatalf("updated = %#v", updated)
	}
	cli(t, srv.URL, "tasks", "delete", "--id", id)
	if got := cli(t, srv.URL, "tasks", "list").([]any); len(got) != 0 {
		t.Fatalf("list after delete = %#v", got)
	}
}

func TestGeneratedCLISurfacesAPIErrors(t *testing.T) {
	srv := httptest.NewServer(NewHandler())
	defer srv.Close()

	if got := cliFailure(t, srv.URL, "tasks", "create", "--set-str", "title="); got != http.StatusBadRequest {
		t.Fatalf("create error status = %d", got)
	}
	if got := cliFailure(t, srv.URL, "tasks", "get", "--id", "missing"); got != http.StatusNotFound {
		t.Fatalf("get error status = %d", got)
	}
	created := cli(t, srv.URL, "tasks", "create", "--set", "title=Keep the contract honest").(map[string]any)
	if got := cliFailure(t, srv.URL, "tasks", "update", "--id", created["id"].(string), "--file", "test/empty.json"); got != http.StatusBadRequest {
		t.Fatalf("update error status = %d", got)
	}
}

func TestHTTPBoundaryRejectsUnsupportedInput(t *testing.T) {
	srv := httptest.NewServer(NewHandler())
	defer srv.Close()

	request, _ := http.NewRequest(http.MethodPut, srv.URL+"/tasks", nil)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	response.Body.Close()
	if response.StatusCode != http.StatusMethodNotAllowed || !strings.Contains(response.Header.Get("Allow"), "GET") || !strings.Contains(response.Header.Get("Allow"), "POST") {
		t.Fatalf("PUT /tasks = %d, allow=%q", response.StatusCode, response.Header.Get("Allow"))
	}

	assertHTTPError(t, srv.URL+"/tasks", "text/plain", `{"title":"wrong media type"}`, "content-type must be application/json")
	assertHTTPError(t, srv.URL+"/tasks", "application/json", `{"title":"`+strings.Repeat("x", 1_000_001)+`"}`, "request body exceeds 1 MB")
}

func assertHTTPError(t *testing.T, url, contentType, body, want string) {
	t.Helper()
	response, err := http.Post(url, contentType, bytes.NewBufferString(body))
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	data, _ := io.ReadAll(response.Body)
	if response.StatusCode != http.StatusBadRequest || !strings.Contains(string(data), want) {
		t.Fatalf("POST /tasks = %d %s", response.StatusCode, data)
	}
}

func equalJSON(a, b any) bool {
	left, _ := json.Marshal(a)
	right, _ := json.Marshal(b)
	return bytes.Equal(left, right)
}
