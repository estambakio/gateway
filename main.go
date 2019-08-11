package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/estambakio/gateway/pkg/proxy"
	"github.com/spf13/viper"
)

type config struct {
	Rules []rule `mapstructure:"rules"`
}

type rule struct {
	From string `mapstructure:"from"`
	To   string `mapstructure:"to"`
}

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic(fmt.Errorf("Config file not found: %s", err))
		} else {
			panic(fmt.Errorf("Error reading config: %s", err))
		}
	}

	var C config

	err := viper.Unmarshal(&C)
	if err != nil {
		panic(fmt.Errorf("unable to decode into struct, %v", err))
	}

	for _, cfg := range C.Rules {
		handler, err := proxy.NewProxy(cfg.From, cfg.To)
		if err != nil {
			panic(err)
		}
		// register proxy handler for every path
		// trying to set the same path (cfg.From) as already existing will cause panic
		http.Handle(cfg.From, handler)
	}

	log.Fatal(http.ListenAndServe(":"+port(), nil))
}

func port() string {
	if v, ok := os.LookupEnv("PORT"); ok {
		return v
	}
	return "3000"
}
