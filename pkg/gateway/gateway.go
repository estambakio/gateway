// Package gateway holds decision-making logic for determining target URL
package gateway

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/spf13/viper"

	"github.com/estambakio/gateway/pkg/proxy"
)

// Config is a struct populated from configuration file
type Config struct {
	Rules []Rule `mapstructure:"rules"`
}

// Rule is a single proxy rule with source and target URIs
type Rule struct {
	From string `mapstructure:"from"`
	To   string `mapstructure:"to"`
}

// match finds proxy rule which matches request
func matchRule(c Config, r *http.Request) (Rule, error) {
	for _, rule := range c.Rules {
		if strings.HasPrefix(r.URL.Path, rule.From) {
			return rule, nil
		}
	}
	return Rule{}, fmt.Errorf("No gateway rule match %s", r.URL.Path)
}

// Handler returns http.Handler which dynamically reads config and proxies request to corresponding target
func Handler(v *viper.Viper) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var c Config

		err := v.Unmarshal(&c)
		if err != nil {
			http.Error(w, fmt.Sprintf("unable to unmarshal config: %v", err), http.StatusInternalServerError)
			return
		}

		proxyRule, err := matchRule(c, r)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		handler, err := proxy.New(proxyRule.From, proxyRule.To)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to create proxy: %v", err), http.StatusInternalServerError)
			return
		}

		handler.ServeHTTP(w, r)
	})
}
