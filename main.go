package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/stnokott/r6api"
	"github.com/stnokott/r6prom/config"
	"github.com/stnokott/r6prom/store"
)

func main() {
	// setup
	writer := zerolog.ConsoleWriter{
		Out:           os.Stdout,
		TimeFormat:    time.RFC3339,
		PartsOrder:    []string{"time", "level", "name", "message"},
		FieldsExclude: []string{"name"},
	}
	logger := zerolog.New(writer).Level(zerolog.DebugLevel).With().Timestamp().Str("name", "R6Prom").Logger()

	conf, err := config.Load()
	if err != nil {
		logger.Fatal().Err(err).Msg("error setting up")
	}

	// create API instance
	r6Logger := logger.With().Str("name", "R6API").Logger()
	a := r6api.NewR6API(conf.Email, conf.Password, r6Logger)

	chanStoreErrs := make(chan error)
	go func() {
		for {
			err := <-chanStoreErrs
			logger.Err(err).Msg("error in store operation")
		}
	}()
	storeOpts := store.Opts{
		ObservedUsernames: conf.ObservedUsernames,
		ChanErrors:        chanStoreErrs,
		MetadataTimeout:   12 * time.Hour,
	}
	store, err := store.New(a, &logger, storeOpts)
	if err != nil {
		logger.Fatal().Err(err).Msg("error creating store")
	}

	registry := prometheus.NewRegistry()
	registry.MustRegister(store)

	http.Handle(
		"/metrics",
		promhttp.HandlerFor(
			registry,
			promhttp.HandlerOpts{
				ErrorLog: log.New(logger.Level(zerolog.ErrorLevel), "", 0),
			},
		),
	)
	logger.Info().Msg("started Prometheus HTTP server on localhost:2112")
	logger.Err(http.ListenAndServe(":2112", nil)).Msg("Prometheus HTTP server stopped")
}
