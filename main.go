package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/stnokott/r6api"
	"github.com/stnokott/r6prom/config"
	"github.com/stnokott/r6prom/constants"
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
	logLevel, err := strconv.Atoi(constants.LOG_LEVEL)
	if err != nil {
		panic(err)
	}
	logger := zerolog.New(writer).Level(zerolog.Level(logLevel)).With().Timestamp().Str("name", "R6Prom").Logger()

	logger.Info().Str("version", constants.VERSION).Stringer("log_level", logger.GetLevel()).Msgf("Setting up %s", constants.NAME)

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
	logger.Info().Str("host", "localhost").Int("port", 2112).Msg("started Prometheus HTTP server")
	logger.Err(http.ListenAndServe(":2112", nil)).Msg("Prometheus HTTP server stopped")
}
