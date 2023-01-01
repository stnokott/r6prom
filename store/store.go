package store

import (
	"sync"
	"time"

	"github.com/go-co-op/gocron"
	influxapi "github.com/influxdata/influxdb-client-go/v2/api"

	"github.com/rs/zerolog"
	"github.com/stnokott/r6api"
	"github.com/stnokott/r6api/types/metadata"
	"github.com/stnokott/r6prom/metrics"
)

type Store struct {
	usernames []string
	api       *r6api.R6API
	influxAPI influxapi.WriteAPI
	scheduler *gocron.Scheduler
	logger    *zerolog.Logger
}

type Opts struct {
	// ObservedUsernames specifies the Uplay usernames to track metrics for
	ObservedUsernames []string
	// InfluxClient handles the connection with the InfluxDB v2
	InfluxWriteAPI influxapi.WriteAPI
	// RefreshCron defines the interval at which the application checks for new stats
	RefreshCron string
}

func New(api *r6api.R6API, logger *zerolog.Logger, opts Opts) (*Store, error) {
	sched := gocron.NewScheduler(time.Local)

	store := &Store{
		usernames: opts.ObservedUsernames,
		api:       api,
		influxAPI: opts.InfluxWriteAPI,
		scheduler: sched,
		logger:    logger,
	}

	if _, err := sched.Cron(opts.RefreshCron).Do(store.sendAll); err != nil {
		return nil, err
	}
	store.scheduler = store.scheduler.SingletonMode().StartImmediately()

	logger.Info().Str("cron", opts.RefreshCron).Int("numUsernames", len(opts.ObservedUsernames)).Msg("initialized store")

	return store, nil
}

// Run starts the scheduler as a blocking call
func (s *Store) Run() {
	s.scheduler.StartBlocking()
}

// RunAsync starts the scheduler as a non-blocking call
func (s *Store) RunAsync() {
	s.scheduler.StartAsync()
}

func (s *Store) sendAll() {
	s.logger.Info().Msg("sending all metrics")
	defer func() {
		s.influxAPI.Flush()
		_, nextRun := s.scheduler.NextRun()
		s.logger.Info().Msgf("done, next run at %v", nextRun)
	}()
	meta, err := s.api.GetMetadata()
	if err != nil {
		s.logger.Err(err).Msg("could not get metadata")
		return
	}

	now := time.Now()
	var wg sync.WaitGroup

	for _, username := range s.usernames {
		wg.Add(1)
		go func(username string) {
			s.logger.Info().Str("username", username).Msgf("processing user %s", username)
			s.sendUserStats(username, meta, now)
			wg.Done()
		}(username)
	}
	wg.Wait()
}

func (s *Store) sendUserStats(username string, meta *metadata.Metadata, t time.Time) {
	profile, err := s.api.ResolveUser(username)
	if err != nil {
		s.logger.Err(err).Msg("could not resolve profile")
		return
	}

	statSenderFuncs := []metrics.StatSenderFunc{
		metrics.SendMapStats,
		metrics.SendMatchStats,
		metrics.SendOperatorStats,
		metrics.SendRankedStats,
	}

	running := len(statSenderFuncs)
	chData := make(chan metrics.StatResponse, 3)

	for _, f := range statSenderFuncs {
		go f(s.api, profile, meta, t, chData)
	}

	for running > 0 {
		data := <-chData
		if data.Done {
			running -= 1
		} else if data.Err != nil {
			s.logger.Err(data.Err).Msg("error sending statistics")
			running -= 1
		} else if data.P != nil {
			s.influxAPI.WritePoint(data.P)
		} else {
			s.logger.Warn().Msg("got invalid data from data channel")
		}
	}
	close(chData)
}
