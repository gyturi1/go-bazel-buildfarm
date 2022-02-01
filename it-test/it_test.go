package ittest

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestEcho(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost:8080", nil)
	if err != nil {
		t.Fatal(err)
	}
	hk := "X-Whatever"
	hv := "good"
	req.Header.Add(hk, hv)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("Not 200 response %v", resp)
	}
	b := resp.Body
	defer b.Close()
	body, err := io.ReadAll(b)
	if err != nil {
		t.Error(err)
	}
	s := string(body)
	want := fmt.Sprintf("%s: %s", hk, hv)
	if !strings.Contains(s, want) {
		t.Errorf("Header value mismatch got: %s want containing: %s", s, want)
	}

}
