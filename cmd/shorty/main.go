package main

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/kellegous/go/backend"
	"github.com/kellegous/go/backend/leveldb"
	"github.com/kellegous/go/web"
	log "github.com/sirupsen/logrus"
)

func main() {
	pflag.String("addr", "127.0.0.1:8067", "default bind address")
	pflag.Bool("admin", false, "allow admin-level requests, e.g. /.hidden_adminz/dumps")
	pflag.String("backend", "leveldb", "backing store to use - 'leveldb' is supported")
	pflag.String("data", "data", "the location of the leveldb data directory")

	pflag.Parse()

	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		log.Panic(err)
	}

	// allow env vars to set pflags
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	var backend backend.Backend

	switch viper.GetString("backend") {
	case "leveldb":
		var err error
		log.Infof("Using leveldb database at: %+v", viper.GetString("data"))
		backend, err = leveldb.New(viper.GetString("data"))
		if err != nil {
			log.Panic(err)
		}
	default:
		log.Panic(fmt.Sprintf("unknown backend %s", viper.GetString("backend")))
	}

	defer backend.Close()
	log.Info("Starting up...")
	log.Panic(web.ListenAndServe(backend))
}
