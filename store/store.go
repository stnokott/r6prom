package store

import (
	"time"

	"github.com/go-co-op/gocron"
	influxapi "github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/rs/zerolog"
	"github.com/stnokott/r6api"
	"github.com/stnokott/r6api/types/metadata"
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

	for i, username := range s.usernames {
		s.logger.Info().Str("username", username).Msgf("processing user %d/%d", i+1, len(s.usernames))
		s.sendUserStats(username, meta)
	}
}

type statSenderFunc func(profile *r6api.Profile, meta *metadata.Metadata, t time.Time) error

func (s *Store) sendUserStats(username string, meta *metadata.Metadata) {
	profile, err := s.api.ResolveUser(username)
	if err != nil {
		s.logger.Err(err).Msg("could not resolve profile")
		return
	}

	now := time.Now()
	statSenderFuncs := map[string]statSenderFunc{
		"maps":      s.sendMapStats,
		"matches":   s.sendMatchStats,
		"operators": s.sendOperatorStats,
		"ranked":    s.sendRankedStats,
	}
	for name, f := range statSenderFuncs {
		if err := f(profile, meta, now); err != nil {
			s.logger.Err(err).Str("stat_name", name).Msg("could not send metrics")
		}
	}
}
