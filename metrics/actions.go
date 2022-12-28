package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stnokott/r6api/types/stats"
)

var (
	labelsActionsWithRole = []string{"season", "username", "gamemode", "role"}
	descKills             = prometheus.NewDesc(
		"kills",
		"Kills by [season,user,gamemode,role]",
		labelsActionsWithRole,
		nil,
	)
	descDeaths = prometheus.NewDesc(
		"deaths",
		"Deaths by [season,user,gamemode,role]",
		labelsActionsWithRole,
		nil,
	)
	descEntryKills = prometheus.NewDesc(
		"entry_kills",
		"Entry kills (first kill in a round) by [season,user,gamemode,role]",
		labelsActionsWithRole,
		nil,
	)
	descEntryDeaths = prometheus.NewDesc(
		"entry_deaths",
		"Entry deaths (first death in a round) by [season,user,gamemode,role]",
		labelsActionsWithRole,
		nil,
	)
	descHeadshotPercentage = prometheus.NewDesc(
		"headshot_percentage",
		"Headshot percentage by [season,user,gamemode,role]",
		labelsActionsWithRole,
		nil,
	)
	descRoundsPlayed = prometheus.NewDesc(
		"rounds_played",
		"Rounds won by [season,user,gamemode,role]",
		labelsActionsWithRole,
		nil,
	)
	descRoundsWon = prometheus.NewDesc(
		"rounds_won",
		"Rounds won by [season,user,gamemode,role]",
		labelsActionsWithRole,
		nil,
	)
	descRoundsLost = prometheus.NewDesc(
		"rounds_lost",
		"Rounds lost by [season,user,gamemode,role]",
		labelsActionsWithRole,
		nil,
	)
	descRoundsWithAce = prometheus.NewDesc(
		"rounds_with_ace",
		"Rounds with Ace by [season,user,gamemode,role]",
		labelsActionsWithRole,
		nil,
	)
	descRoundsWithClutch = prometheus.NewDesc(
		"rounds_with_clutch",
		"Rounds with clutch by [season,user,gamemode,role]",
		labelsActionsWithRole,
		nil,
	)
	descRoundsWithEntryDeath = prometheus.NewDesc(
		"rounds_with_entry_death",
		"Rounds with entry death by [season,user,gamemode,role]",
		labelsActionsWithRole,
		nil,
	)
	descRoundsWithEntryKill = prometheus.NewDesc(
		"rounds_with_entry_kill",
		"Rounds with entry kill by [season,user,gamemode,role]",
		labelsActionsWithRole,
		nil,
	)
	descRoundsWithKOST = prometheus.NewDesc(
		"rounds_with_kost",
		"Rounds with KOST (Kill,Objective,Survival,Trade) by [season,user,gamemode,role]",
		labelsActionsWithRole,
		nil,
	)
	descRoundsWithKill = prometheus.NewDesc(
		"rounds_with_kill",
		"Rounds with kill by [season,user,gamemode,role]",
		labelsActionsWithRole,
		nil,
	)
	descRoundsWithMultiKill = prometheus.NewDesc(
		"rounds_with_multikill",
		"Rounds with multi-kill by [season,user,gamemode,role]",
		labelsActionsWithRole,
		nil,
	)

	labelsActions     = []string{"season", "username", "gamemode"}
	descMatchesPlayed = prometheus.NewDesc(
		"matches_played",
		"Matches won by [season,user,gamemode]",
		labelsActions,
		nil,
	)
	descMatchesWon = prometheus.NewDesc(
		"matches_won",
		"Matches won by [season,user,gamemode]",
		labelsActions,
		nil,
	)
	descMatchesLost = prometheus.NewDesc(
		"matches_lost",
		"Matches lost by [season,user,gamemode]",
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
	descRoundsPlayed,
	descRoundsWon,
	descRoundsLost,
	descRoundsWithAce,
	descRoundsWithClutch,
	descRoundsWithEntryDeath,
	descRoundsWithEntryKill,
	descRoundsWithKOST,
	descRoundsWithKill,
	descRoundsWithMultiKill,
	descMatchesPlayed,
	descMatchesWon,
	descMatchesLost,
}

type ActionsMetricProvider struct {
	Stats    *stats.SummarizedStats
	Username string
}

func (p ActionsMetricProvider) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(p, ch)
}

type metric struct {
	desc  *prometheus.Desc
	value float64
}

func (p ActionsMetricProvider) Collect(ch chan<- prometheus.Metric) {
	if p.Stats == nil {
		return
	}

	gameModes := map[string]*stats.SummarizedGameModeStats{
		"all":      p.Stats.All,
		"casual":   p.Stats.Casual,
		"unranked": p.Stats.Unranked,
		"ranked":   p.Stats.Ranked,
	}
	for gameModeName, gameMode := range gameModes {
		if gameMode == nil {
			continue
		}

		p.collectMetricsWithoutRole(ch, gameModeName, gameMode)
		p.collectMetricsWithRole(ch, gameModeName, gameMode)
	}
}

func (p ActionsMetricProvider) collectMetricsWithoutRole(ch chan<- prometheus.Metric, gameModeName string, gameMode *stats.SummarizedGameModeStats) {
	metricsNoRole := []metric{
		{descMatchesPlayed, float64(gameMode.MatchesPlayed)},
		{descMatchesWon, float64(gameMode.MatchesWon)},
		{descMatchesLost, float64(gameMode.MatchesLost)},
	}
	for _, metric := range metricsNoRole {
		ch <- prometheus.MustNewConstMetric(
			metric.desc,
			prometheus.GaugeValue,
			metric.value,
			p.Stats.SeasonSlug,
			p.Username,
			gameModeName,
		)
	}
}

func (p ActionsMetricProvider) collectMetricsWithRole(ch chan<- prometheus.Metric, gameModeName string, gameMode *stats.SummarizedGameModeStats) {
	roles := map[string]*stats.DetailedStats{
		"all":     gameMode.All,
		"attack":  gameMode.Attack,
		"defence": gameMode.Defence,
	}
	for roleName, roleStats := range roles {
		metricsWithRole := []metric{
			{descKills, float64(roleStats.Kills)},
			{descDeaths, float64(roleStats.Deaths)},
			{descEntryKills, float64(roleStats.EntryKills)},
			{descEntryDeaths, float64(roleStats.EntryDeaths)},
			{descHeadshotPercentage, roleStats.HeadshotPercentage},
			{descRoundsPlayed, float64(roleStats.RoundsPlayed)},
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

		for _, metric := range metricsWithRole {
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

func ActionsErr(ch chan<- prometheus.Metric, err error) {
	for _, d := range allKillDescs {
		ch <- prometheus.NewInvalidMetric(d, err)
	}
}
