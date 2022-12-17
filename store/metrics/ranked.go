package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stnokott/r6api/types/metadata"
	"github.com/stnokott/r6api/types/ranked"
)

var DescRankedMMR = prometheus.NewDesc(
	"ranked_mmr",
	"Ranked MMR for a user in a season",
	[]string{"season", "username"},
	nil,
)

var DescRankedRank = prometheus.NewDesc(
	"ranked_rank",
	"Ranked rank ID for a user in a season",
	[]string{"season", "username"},
	nil,
)

var DescRankedConfidence = prometheus.NewDesc(
	"ranked_confidence",
	"Ranked confidence for a user in a season",
	[]string{"season", "username"},
	nil,
)

func CollectRank(ch chan<- prometheus.Metric, s ranked.SeasonStats, meta *metadata.Metadata, username string) {
	seasonSlug := meta.SeasonSlugFromID(s.SeasonID)
	ch <- prometheus.MustNewConstMetric(
		DescRankedMMR,
		prometheus.GaugeValue,
		float64(s.MMR),
		seasonSlug,
		username,
	)
	ch <- prometheus.MustNewConstMetric(
		DescRankedRank,
		prometheus.GaugeValue,
		float64(s.Rank),
		seasonSlug,
		username,
	)
	ch <- prometheus.MustNewConstMetric(
		DescRankedConfidence,
		prometheus.GaugeValue,
		float64(s.SkillStdev),
		seasonSlug,
		username,
	)
}
