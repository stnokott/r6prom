package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stnokott/r6api/types/metadata"
	"github.com/stnokott/r6prom/store"
)

var (
	labelsRanked  = []string{"season", "username"}
	descRankedMMR = prometheus.NewDesc(
		"ranked_mmr",
		"Ranked MMR by [user,season]",
		labelsRanked,
		nil,
	)
	descRankedRank = prometheus.NewDesc(
		"ranked_rank",
		"Ranked rank ID by [user,season]",
		labelsRanked,
		nil,
	)
	descRankedConfidence = prometheus.NewDesc(
		"ranked_confidence",
		"Ranked confidence by [user,season]",
		labelsRanked,
		nil,
	)
	descRankedGamesWon = prometheus.NewDesc(
		"ranked_games_won",
		"Ranked wins by [user,season]",
		labelsRanked,
		nil,
	)
	descRankedGamesLost = prometheus.NewDesc(
		"ranked_games_lost",
		"Ranked losses by [user,season]",
		labelsRanked,
		nil,
	)
)

var allRankedDescs = []*prometheus.Desc{
	descRankedMMR,
	descRankedRank,
	descRankedConfidence,
	descRankedGamesWon,
	descRankedGamesLost,
}

type RankedMetricProvider struct{}

func (p RankedMetricProvider) Collect(ch chan<- prometheus.Metric, s *store.StatsCollection, m *metadata.Metadata, username string) {
	rankedStats := s.RankedStats
	seasonSlug := m.SeasonSlugFromID(rankedStats.SeasonID)

	for _, v := range []struct {
		desc  *prometheus.Desc
		value float64
	}{
		{descRankedMMR, float64(rankedStats.MMR)},
		{descRankedRank, float64(rankedStats.Rank)},
		{descRankedConfidence, float64(rankedStats.SkillStdev)},
		{descRankedGamesWon, float64(rankedStats.Wins)},
		{descRankedGamesLost, float64(rankedStats.Losses)},
	} {
		ch <- prometheus.MustNewConstMetric(
			v.desc,
			prometheus.GaugeValue,
			v.value,
			seasonSlug,
			username,
		)
	}
}

func (p RankedMetricProvider) CollectErr(ch chan<- prometheus.Metric, err error) {
	for _, d := range allRankedDescs {
		ch <- prometheus.NewInvalidMetric(d, err)
	}
}
