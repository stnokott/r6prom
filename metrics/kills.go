package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stnokott/r6api/types/metadata"
	"github.com/stnokott/r6api/types/stats"
)

var DescKills = prometheus.NewDesc(
	"kills",
	"Kills in a season for a user in a gamemode",
	[]string{"season", "username", "gamemode"},
	nil,
)

func CollectKills(ch chan<- prometheus.Metric, s *stats.SummarizedStats, meta *metadata.Metadata, username string) {
	// summarized
	if s == nil {
		return
	}
	gameModes := map[string]*stats.SummarizedGameModeStats{
		"casual":   s.Casual,
		"unranked": s.Unranked,
		"ranked":   s.Ranked,
	}
	for gameModeName, gameMode := range gameModes {
		if gameMode == nil {
			continue
		}
		seasonSlug := "n/a"
		kills := 0
		if gameMode.Attack != nil {
			kills += gameMode.Attack.Kills
			seasonSlug = s.SeasonSlug
		}
		if gameMode.Defence != nil {
			kills += gameMode.Defence.Kills
			seasonSlug = s.SeasonSlug
		}
		ch <- prometheus.MustNewConstMetric(
			DescKills,
			prometheus.GaugeValue,
			float64(kills),
			seasonSlug,
			username,
			gameModeName,
		)
	}
}
