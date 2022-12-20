package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stnokott/r6api/types/metadata"
	"github.com/stnokott/r6api/types/stats"
	"github.com/stnokott/r6prom/store"
)

var labelsActions = []string{"season", "username", "gamemode", "role"}

var descKills = prometheus.NewDesc(
	"kills",
	"Kills in a season for a user in a gamemode for a role",
	labelsActions,
	nil,
)

var descDeaths = prometheus.NewDesc(
	"deaths",
	"Deaths in a season for a user in a gamemode for a role",
	labelsActions,
	nil,
)

var allKillDescs = []*prometheus.Desc{
	descKills,
	descDeaths,
}

type ActionsMetricProvider struct{}

func (p ActionsMetricProvider) GetDescs() []*prometheus.Desc {
	return allKillDescs
}

func (p ActionsMetricProvider) Collect(ch chan<- prometheus.Metric, s *store.StatsCollection, m *metadata.Metadata, username string) {
	if s == nil {
		return
	}
	summarizedStats := s.SummarizedStats

	gameModes := map[string]*stats.SummarizedGameModeStats{
		"casual":   summarizedStats.Casual,
		"unranked": summarizedStats.Unranked,
		"ranked":   summarizedStats.Ranked,
	}
	for gameModeName, gameMode := range gameModes {
		if gameMode == nil {
			continue
		}

		roles := map[string]*stats.DetailedStats{
			"attack":  gameMode.Attack,
			"defence": gameMode.Defence,
		}
		for roleName, roleStats := range roles {
			metrics := []struct {
				desc  *prometheus.Desc
				value int
			}{
				{descKills, roleStats.Kills},
				{descDeaths, roleStats.Deaths},
			}

			for _, metric := range metrics {
				ch <- prometheus.MustNewConstMetric(
					metric.desc,
					prometheus.GaugeValue,
					float64(metric.value),
					summarizedStats.SeasonSlug,
					username,
					gameModeName,
					roleName,
				)
			}
		}
	}
}

func (p ActionsMetricProvider) CollectErr(ch chan<- prometheus.Metric, err error) {
	for _, d := range allKillDescs {
		ch <- prometheus.NewInvalidMetric(d, err)
	}
}
