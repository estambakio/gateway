package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/estambakio/gateway/pkg/proxy"
	flags "github.com/jessevdk/go-flags"
	"github.com/spf13/viper"
)

type config struct {
	Rules []rule `mapstructure:"rules"`
}

type rule struct {
	From string `mapstructure:"from"`
	To   string `mapstructure:"to"`
}

func dynamicProxyHandler(rules []rule) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, rule := range rules {
			if strings.HasPrefix(r.URL.Path, rule.From) {
				handler, err := proxy.New(rule.From, rule.To)
				if err != nil {
					http.Error(w, "failed to proxy request", http.StatusInternalServerError)
					return
				}
				handler.ServeHTTP(w, r)
				return
			}
		}
		// if request doesn't match any rule.From then serve 404 error
		http.NotFound(w, r)
	})
}

// command-line options
var opts struct {
	Config string `short:"c" long:"config" description:"Path to config.yaml" required:"true"`
}

func main() {
	parser := flags.NewParser(&opts, flags.Default)
	_, err := parser.Parse()
	if err != nil {
		flagsErr, ok := err.(*flags.Error)
		if ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else if ok && flagsErr.Type == flags.ErrCommandRequired {
			parser.WriteHelp(os.Stdout)
			os.Exit(1)
		} else {
			os.Exit(1)
		}
	}

	viper.SetConfigFile(opts.Config)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic(fmt.Errorf("Config file not found: %s", err))
		} else {
			panic(fmt.Errorf("Error reading config: %s", err))
		}
	}

	var C config

	err = viper.Unmarshal(&C)
	if err != nil {
		panic(fmt.Errorf("unable to decode into struct, %v", err))
	}

	log.Fatal(http.ListenAndServe(":"+port(), dynamicProxyHandler(C.Rules)))
}

func port() string {
	if v, ok := os.LookupEnv("PORT"); ok {
		return v
	}
	return "3000"
}
