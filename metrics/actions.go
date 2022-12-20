package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stnokott/r6api/types/stats"
)

var (
	labelsActions = []string{"season", "username", "gamemode", "role"}
	descKills     = prometheus.NewDesc(
		"kills",
		"Kills by [season,user,gamemode,role]",
		labelsActions,
		nil,
	)
	descDeaths = prometheus.NewDesc(
		"deaths",
		"Deaths by [season,user,gamemode,role]",
		labelsActions,
		nil,
	)
	descEntryKills = prometheus.NewDesc(
		"entry_kills",
		"Entry kills (first kill in a round) by [season,user,gamemode,role]",
		labelsActions,
		nil,
	)
	descEntryDeaths = prometheus.NewDesc(
		"entry_deaths",
		"Entry deaths (first death in a round) by [season,user,gamemode,role]",
		labelsActions,
		nil,
	)
	descHeadshotPercentage = prometheus.NewDesc(
		"headshot_percentage",
		"Headshot percentage by [season,user,gamemode,role]",
		labelsActions,
		nil,
	)
	descMatchesWon = prometheus.NewDesc(
		"matches_won",
		"Matches won by [season,user,gamemode,role]",
		labelsActions,
		nil,
	)
	descMatchesLost = prometheus.NewDesc(
		"matches_lost",
		"Matches lost by [season,user,gamemode,role]",
		labelsActions,
		nil,
	)
	descRoundsWon = prometheus.NewDesc(
		"rounds_won",
		"Rounds won by [season,user,gamemode,role]",
		labelsActions,
		nil,
	)
	descRoundsLost = prometheus.NewDesc(
		"rounds_lost",
		"Rounds lost by [season,user,gamemode,role]",
		labelsActions,
		nil,
	)
	descRoundsWithAce = prometheus.NewDesc(
		"rounds_with_ace",
		"Rounds with Ace by [season,user,gamemode,role]",
		labelsActions,
		nil,
	)
	descRoundsWithClutch = prometheus.NewDesc(
		"rounds_with_clutch",
		"Rounds with clutch by [season,user,gamemode,role]",
		labelsActions,
		nil,
	)
	descRoundsWithEntryDeath = prometheus.NewDesc(
		"rounds_with_entry_death",
		"Rounds with entry death by [season,user,gamemode,role]",
		labelsActions,
		nil,
	)
	descRoundsWithEntryKill = prometheus.NewDesc(
		"rounds_with_entry_kill",
		"Rounds with entry kill by [season,user,gamemode,role]",
		labelsActions,
		nil,
	)
	descRoundsWithKOST = prometheus.NewDesc(
		"rounds_with_kost",
		"Rounds with KOST (Kill,Objective,Survival,Trade) by [season,user,gamemode,role]",
		labelsActions,
		nil,
	)
	descRoundsWithKill = prometheus.NewDesc(
		"rounds_with_kill",
		"Rounds with kill by [season,user,gamemode,role]",
		labelsActions,
		nil,
	)
	descRoundsWithMultiKill = prometheus.NewDesc(
		"rounds_with_multikill",
		"Rounds with multi-kill by [season,user,gamemode,role]",
		labelsActions,
		nil,
	)
)

var allKillDescs = []*prometheus.Desc{
	descKills,
	descDeaths,
	descEntryKills,
	descEntryDeaths,
	descHeadshotPercentage,
	descMatchesWon,
	descMatchesLost,
	descRoundsWon,
	descRoundsLost,
	descRoundsWithAce,
	descRoundsWithClutch,
	descRoundsWithEntryDeath,
	descRoundsWithEntryKill,
	descRoundsWithKOST,
	descRoundsWithKill,
	descRoundsWithMultiKill,
}

type ActionsMetricProvider struct {
	Stats    *stats.SummarizedStats
	Username string
}

func (p ActionsMetricProvider) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(p, ch)
}

func (p ActionsMetricProvider) Collect(ch chan<- prometheus.Metric) {
	if p.Stats == nil {
		return
	}

	gameModes := map[string]*stats.SummarizedGameModeStats{
		"casual":   p.Stats.Casual,
		"unranked": p.Stats.Unranked,
		"ranked":   p.Stats.Ranked,
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
				value float64
			}{
				{descKills, float64(roleStats.Kills)},
				{descDeaths, float64(roleStats.Deaths)},
				{descEntryKills, float64(roleStats.EntryKills)},
				{descEntryDeaths, float64(roleStats.EntryDeaths)},
				{descHeadshotPercentage, roleStats.HeadshotPercentage},
				{descMatchesWon, float64(roleStats.MatchesWon)},
				{descMatchesLost, float64(roleStats.MatchesLost)},
				{descRoundsWon, float64(roleStats.RoundsWon)},
				{descRoundsLost, float64(roleStats.RoundsLost)},
				{descRoundsWithAce, roleStats.RoundsWithAce},
				{descRoundsWithClutch, roleStats.RoundsWithClutch},
				{descRoundsWithEntryDeath, roleStats.RoundsWithEntryDeath},
				{descRoundsWithEntryKill, roleStats.RoundsWithEntryKill},
				{descRoundsWithKOST, roleStats.RoundsWithKOST},
				{descRoundsWithKill, roleStats.RoundsWithKill},
				{descRoundsWithMultiKill, roleStats.RoundsWithMultikill},
			}

			for _, metric := range metrics {
				ch <- prometheus.MustNewConstMetric(
					metric.desc,
					prometheus.GaugeValue,
					metric.value,
					p.Stats.SeasonSlug,
					p.Username,
					gameModeName,
					roleName,
				)
			}
		}
	}
}

func ActionsErr(ch chan<- prometheus.Metric, err error) {
	for _, d := range allKillDescs {
		ch <- prometheus.NewInvalidMetric(d, err)
	}
}
