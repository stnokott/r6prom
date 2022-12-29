package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stnokott/r6api/types/metadata"
	"github.com/stnokott/r6api/types/ranked"
)

var (
	labelsRanked    = []string{"season", "username"}
	metricRankedMMR = metricDetails{
		desc: prometheus.NewDesc(
			"ranked_mmr",
			"Ranked MMR by [user,season]",
			labelsRanked,
			nil,
		),
		metricType: prometheus.GaugeValue,
	}
	metricRankedRank = metricDetails{
		desc: prometheus.NewDesc(
			"ranked_rank",
			"Ranked rank ID by [user,season]",
			labelsRanked,
			nil,
		),
		metricType: prometheus.GaugeValue,
	}
	metricRankedConfidence = metricDetails{
		desc: prometheus.NewDesc(
			"ranked_confidence",
			"Ranked confidence by [user,season]",
			labelsRanked,
			nil,
		),
		metricType: prometheus.GaugeValue,
	}
	metricRankedMatchesWon = metricDetails{
		desc: prometheus.NewDesc(
			"ranked_matches_won",
			"Ranked wins by [user,season]",
			labelsRanked,
			nil,
		),
		metricType: prometheus.CounterValue,
	}
	metricRankedMatchesLost = metricDetails{
		desc: prometheus.NewDesc(
			"ranked_matches_lost",
			"Ranked losses by [user,season]",
			labelsRanked,
			nil,
		),
		metricType: prometheus.CounterValue,
	}
)

var allRankedDescs = []metricDetails{
	metricRankedMMR,
	metricRankedRank,
	metricRankedConfidence,
	metricRankedMatchesWon,
	metricRankedMatchesLost,
}

type RankedMetricProvider struct {
	Stats    *ranked.SeasonStats
	Meta     *metadata.Metadata
	Username string
}

func (p RankedMetricProvider) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(p, ch)
}

func (p RankedMetricProvider) Collect(ch chan<- prometheus.Metric) {
	rankedStats := p.Stats
	seasonSlug := p.Meta.SeasonSlugFromID(rankedStats.SeasonID)

	for _, v := range []metricInstance{
		{metricRankedMMR, float64(rankedStats.MMR)},
		{metricRankedRank, float64(rankedStats.Rank)},
		{metricRankedConfidence, float64(rankedStats.SkillStdev)},
		{metricRankedMatchesWon, float64(rankedStats.Wins)},
		{metricRankedMatchesLost, float64(rankedStats.Losses)},
	} {
		ch <- prometheus.MustNewConstMetric(
			v.details.desc,
			v.details.metricType,
			v.value,
			seasonSlug,
			p.Username,
		)
	}
}

func RankedErr(ch chan<- prometheus.Metric, err error) {
	for _, d := range allRankedDescs {
		ch <- prometheus.NewInvalidMetric(d.desc, err)
	}
}
