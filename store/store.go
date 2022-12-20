package store

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"github.com/stnokott/r6api"
	"github.com/stnokott/r6api/types/metadata"
	"github.com/stnokott/r6api/types/ranked"
	"github.com/stnokott/r6api/types/stats"
)

type StatsCollection struct {
	RankedStats     ranked.SeasonStats
	SummarizedStats *stats.SummarizedStats
}

type MetricProvider interface {
	Collect(ch chan<- prometheus.Metric, s *StatsCollection, m *metadata.Metadata, username string)
	CollectErr(chan<- prometheus.Metric, error)
}

type Store struct {
	api             *r6api.R6API
	metaMutex       sync.Mutex
	meta            *metadata.Metadata
	metricProviders []MetricProvider
	cache           *cache
	opts            Opts
	logger          *zerolog.Logger
}

type Opts struct {
	// ObservedUsernames specifies the Uplay usernames to track metrics for
	ObservedUsernames []string
	// ChanErrors will receive errors from asynchronous operations, e.g. metadata refresh
	ChanErrors chan<- error
	// MetadataTimeout defines refresh interval for metadata
	MetadataTimeout time.Duration
}

func New(api *r6api.R6API, logger *zerolog.Logger, opts Opts, ctx context.Context) (*Store, error) {
	store := &Store{
		api:       api,
		metaMutex: sync.Mutex{},
		cache:     newCache(api, logger, ctx),
		opts:      opts,
		logger:    logger,
	}

	if opts.MetadataTimeout.Seconds() < 10 {
		return nil, errors.New("please set metadata timeout to >=10s, since it is an expensive operation and should not be called frequently")
	}

	if err := store.RefreshMetadata(); err != nil {
		return nil, fmt.Errorf("could not get first metadata: %w", err)
	}
	go func() {
		for {
			time.Sleep(opts.MetadataTimeout)
			if err := store.RefreshMetadata(); err != nil {
				opts.ChanErrors <- err
			}
		}
	}()
	return store, nil
}

func (s *Store) Register(provider MetricProvider) {
	s.metricProviders = append(s.metricProviders, provider)
}

func (s *Store) RefreshMetadata() error {
	s.logger.Debug().Msg("refreshing metadata")
	s.metaMutex.Lock()
	meta, err := s.api.GetMetadata()
	if err != nil {
		return err
	}
	if len(meta.Seasons) == 0 {
		return fmt.Errorf("no seasons found")
	}
	s.meta = meta
	s.metaMutex.Unlock()
	return nil
}

func (s *Store) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(s, ch)
}

func (s *Store) Collect(ch chan<- prometheus.Metric) {
	for _, username := range s.opts.ObservedUsernames {
		s.collectUser(ch, username)
	}
}

func (s *Store) collectUser(ch chan<- prometheus.Metric, username string) {
	s.logger.Debug().Str("username", username).Msg("collecting")

	var err error
	var profile *r6api.Profile
	profile, err = s.cache.GetProfile(username)
	if err != nil || profile == nil {
		if err == nil && profile == nil {
			err = fmt.Errorf("could not resolve profile for %s", username)
		}
		for _, m := range s.metricProviders {
			m.CollectErr(ch, err)
		}
		return
	}

	var sc *StatsCollection
	if sc, err = s.getStats(profile); err != nil {
		for _, m := range s.metricProviders {
			m.CollectErr(ch, err)
		}
	}

	for _, provider := range s.metricProviders {
		provider.Collect(ch, sc, s.meta, username)
	}
}

func (s *Store) getStats(profile *r6api.Profile) (sc *StatsCollection, err error) {
	// length of metadata already checked in RefreshMetadata, no need to check here
	season := s.meta.Seasons[len(s.meta.Seasons)-1]
	stats := new(stats.SummarizedStats)
	if err = s.api.GetStats(profile, season.Slug, stats); err != nil {
		return
	}

	var skillHistory ranked.SkillHistory
	if skillHistory, err = s.api.GetRankedHistory(profile, 1); err != nil {
		return
	}

	sc = &StatsCollection{
		RankedStats:     skillHistory[0],
		SummarizedStats: stats,
	}
	return
}
