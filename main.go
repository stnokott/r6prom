package main

import (
	"context"
	"os"
	"strconv"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/domain"
	influxlog "github.com/influxdata/influxdb-client-go/v2/log"
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

	logger.Info().Str("version", constants.VERSION).Stringer("log_level", logger.GetLevel()).Msgf("setting up %s", constants.NAME)

	conf, err := config.Load()
	if err != nil {
		logger.Fatal().Err(err).Msg("error setting up")
	}

	// create API instance
	r6Logger := logger.With().Str("name", "R6API").Logger()
	a := r6api.NewR6API(conf.Email, conf.Password, r6Logger)

	influxClient := influxdb2.NewClientWithOptions(
		conf.InfluxURL,
		conf.InfluxAuthToken,
		influxdb2.DefaultOptions().
			SetBatchSize(10000). // high batch size to allow for manual flushing
			SetApplicationName(constants.NAME).
			SetLogLevel(influxlog.ErrorLevel).
			SetPrecision(time.Second),
	)
	health, err := influxClient.Health(context.Background())
	if err != nil {
		logger.Fatal().Err(err).Msg("could not get InfluxDB health")
	}
	if health.Status != domain.HealthCheckStatusPass {
		logger.Fatal().Msg("InfluxDB server unhealthy, aborting")
	}
	logger.Info().Str("version", *health.Version).Str("msg", *health.Message).Str("db_name", health.Name).Msg("connected to InfluxDB")
	writeAPI := influxClient.WriteAPI(conf.InfluxOrg, conf.InfluxBucket)
	influxErrs := writeAPI.Errors()
	defer influxClient.Close()

	// create store
	storeOpts := store.Opts{
		ObservedUsernames: conf.ObservedUsernames,
		InfluxWriteAPI:    writeAPI,
		RefreshCron:       conf.RefreshCron,
	}
	store, err := store.New(a, &logger, storeOpts)
	if err != nil {
		logger.Fatal().Err(err).Msg("error creating store")
	}
	store.RunAsync()
	for err := range influxErrs {
		logger.Err(err).Msg("encountered Influx write error")
	}
}
