package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stnokott/r6api/types/metadata"
	"github.com/stnokott/r6prom/store"
)

var labelsRanked = []string{"season", "username"}

var descRankedMMR = prometheus.NewDesc(
	"ranked_mmr",
	"Ranked MMR for a user in a season",
	labelsRanked,
	nil,
)

var descRankedRank = prometheus.NewDesc(
	"ranked_rank",
	"Ranked rank ID for a user in a season",
	labelsRanked,
	nil,
)

var descRankedConfidence = prometheus.NewDesc(
	"ranked_confidence",
	"Ranked confidence for a user in a season",
	labelsRanked,
	nil,
)

var allRankedDescs = []*prometheus.Desc{
	descRankedMMR,
	descRankedRank,
	descRankedConfidence,
}

type RankedMetricProvider struct{}

func (p RankedMetricProvider) GetDescs() []*prometheus.Desc {
	return allRankedDescs
}

func (p RankedMetricProvider) Collect(ch chan<- prometheus.Metric, s *store.StatsCollection, m *metadata.Metadata, username string) {
	rankedStats := s.RankedStats
	seasonSlug := m.SeasonSlugFromID(rankedStats.SeasonID)

	for _, v := range []struct {
		desc  *prometheus.Desc
		value float64
	}{
		{descRankedMMR, float64(rankedStats.MMR)},
		{descRankedMMR, float64(rankedStats.Rank)},
		{descRankedConfidence, float64(rankedStats.SkillStdev)},
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
