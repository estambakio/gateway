package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	flags "github.com/jessevdk/go-flags"
	"github.com/spf13/viper"

	"github.com/estambakio/gateway/pkg/gateway"
)

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

	v := viper.New()

	v.SetConfigFile(opts.Config)
	v.SetConfigType("yaml")

	// watch config file for changes
	v.WatchConfig()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic(fmt.Errorf("Config file not found: %s", err))
		} else {
			panic(fmt.Errorf("Error reading config: %s", err))
		}
	}

	log.Fatal(http.ListenAndServe(":"+port(), gateway.Handler(v)))
}

func port() string {
	if v, ok := os.LookupEnv("PORT"); ok {
		return v
	}
	return "3000"
}
