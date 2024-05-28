package config

import (
	"encoding/json"
	"flag"
	"log/slog"
	"os"

	"github.com/go-playground/validator/v10"
)

var confFlag = flag.String("c", "../config.json", "Configaration file path")

func LoadConfig() error {

	flag.Parse()
	var data []byte
	var err error

	exit := func(err error) {
		slog.Error(err.Error())
		os.Exit(1)
	}

	if data, err = os.ReadFile(*confFlag); err != nil {
		exit(err)
	}

	if err = json.Unmarshal(data, &config); err != nil {
		exit(err)
	}

	v := validator.New()
	if err = v.Struct(config); err != nil {
		exit(err)
	}

	return nil
}
