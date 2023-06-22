package config

import (
	"go-flow2sqlite/cmd/app/internal/services/updater"
	"os"
	"strings"

	"github.com/cristalhq/aconfig"
	"github.com/cristalhq/aconfig/aconfigyaml"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Loc                    string   `default:"UTC" usage:"Location for time"`
	SubNets                []string `default:"" usage:"List of subnets traffic between which will not be counted"`
	IgnorList              []string `default:"" usage:"List of lines that will be excluded from the final log"`
	BindAddr               string   `default:"0.0.0.0:30340" usage:"Listen address for HTTP-server"`
	GomtcAddr              string   `default:"http://127.0.0.1:3034" usage:"Address and port for connect to GOMTC API"`
	FlowAddr               string   `default:"0.0.0.0:2055" usage:"Address and port to listen NetFlow packets"`
	PathToBD               string   `default:"/var/go-flow2sqlite/access.sqlite" usage:"database file to which the stream is written"`
	Interval               string   `default:"1m" usage:"Interval to getting info from GOMTC"`
	LogLevel               string   `default:"info" usage:"Log level: panic, fatal, error, warn, info, debug, trace"`
	ReceiveBufferSizeBytes int      `default:"" usage:"Size of RxQueue, i.e. value for SO_RCVBUF in bytes"`
}

func NewConfig() *Config {
	// fix for https://github.com/cristalhq/aconfig/issues/82
	args := []string{}
	for _, a := range os.Args {
		if !strings.HasPrefix(a, "-test.") {
			args = append(args, a)
		}
	}
	// fix for https://github.com/cristalhq/aconfig/issues/82

	var cfg Config
	loader := aconfig.LoaderFor(&cfg, aconfig.Config{
		// feel free to skip some steps :)
		// SkipEnv:      true,
		SkipFiles:          false,
		AllowUnknownFields: true,
		SkipDefaults:       false,
		SkipFlags:          false,
		EnvPrefix:          "GO_FLOW2SQLITE",
		FlagPrefix:         "",
		Files: []string{
			"./config.yaml",
			"./config/config.yaml",
			"/etc/go-flow2sqlite/config.yaml",
			"/bin/local/bin/go-flow2sqlite/config.yaml",
			"/bin/local/bin/go-flow2sqlite/config/config.yaml",
			"/opt/go-flow2sqlite/config.yaml",
			"/opt/go-flow2sqlite/config/config.yaml",
			"./configs/app/default.yaml",
		},
		FileDecoders: map[string]aconfig.FileDecoder{
			// from `aconfigyaml` submodule
			// see submodules in repo for more formats
			".yaml": aconfigyaml.New(),
			// ".toml": aconfigtoml.New(),
		},
		Args: args[1:], // [1:] важно, см. доку к FlagSet.Parse
	})

	if err := loader.Load(); err != nil {
		panic(err)
	}

	lvl, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		logrus.Errorf("Error parse the level of logs (%v). Installed by default = Info", cfg.LogLevel)
		lvl, _ = logrus.ParseLevel("info")
	}
	logrus.SetLevel(lvl)

	logrus.Debugf("Config %#v:", cfg)

	return &cfg
}

func (cfg *Config) ToOptions() *updater.Options {
	return &updater.Options{
		Address: cfg.GomtcAddr,
	}
}
