package gateway

import (
	"net/http"
	"testing"
)

func Test_MatchRule(t *testing.T) {
	from, to := "/one/", "http://domain.com/two/"
	c := Config{Rules: []Rule{{From: from, To: to}}}

	r, err := http.NewRequest("GET", from, nil)
	if err != nil {
		t.Errorf("failed to create a test request: %v", err)
	}
	rule, err := matchRule(c, r)
	if err != nil {
		t.Errorf("failed to find rule: %v", err)
	}
	if rule.To != to {
		t.Errorf("wrong target: want %s, got: %v", to, rule.To)
	}

	badRequest, err := http.NewRequest("GET", "/something/not/in/rules/", nil)
	if err != nil {
		t.Errorf("failed to create a test request: %v", err)
	}
	_, err = matchRule(c, badRequest)
	if err == nil {
		t.Error("should've failed")
	}
}
