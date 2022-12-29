package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stnokott/r6api/types/stats"
)

var (
	labelsActionsWithRole = []string{"season", "username", "gamemode", "role"}
	metricKills           = metricDetails{
		desc: prometheus.NewDesc(
			"kills",
			"Kills by [season,user,gamemode,role]",
			labelsActionsWithRole,
			nil,
		),
		metricType: prometheus.CounterValue,
	}
	metricDeaths = metricDetails{
		desc: prometheus.NewDesc(
			"deaths",
			"Deaths by [season,user,gamemode,role]",
			labelsActionsWithRole,
			nil,
		),
		metricType: prometheus.CounterValue,
	}
	metricEntryKills = metricDetails{
		desc: prometheus.NewDesc(
			"entry_kills",
			"Entry kills (first kill in a round) by [season,user,gamemode,role]",
			labelsActionsWithRole,
			nil,
		),
		metricType: prometheus.CounterValue,
	}
	metricEntryDeaths = metricDetails{
		desc: prometheus.NewDesc(
			"entry_deaths",
			"Entry deaths (first death in a round) by [season,user,gamemode,role]",
			labelsActionsWithRole,
			nil,
		),
		metricType: prometheus.CounterValue,
	}
	metricHeadshotPercentage = metricDetails{
		desc: prometheus.NewDesc(
			"headshot_percentage",
			"Headshot percentage by [season,user,gamemode,role]",
			labelsActionsWithRole,
			nil,
		),
		metricType: prometheus.GaugeValue,
	}
	metricRoundsPlayed = metricDetails{
		desc: prometheus.NewDesc(
			"rounds_played",
			"Rounds won by [season,user,gamemode,role]",
			labelsActionsWithRole,
			nil,
		),
		metricType: prometheus.CounterValue,
	}
	metricRoundsWon = metricDetails{
		desc: prometheus.NewDesc(
			"rounds_won",
			"Rounds won by [season,user,gamemode,role]",
			labelsActionsWithRole,
			nil,
		),
		metricType: prometheus.CounterValue,
	}
	metricRoundsLost = metricDetails{
		desc: prometheus.NewDesc(
			"rounds_lost",
			"Rounds lost by [season,user,gamemode,role]",
			labelsActionsWithRole,
			nil,
		),
		metricType: prometheus.CounterValue,
	}
	metricRoundsWithAce = metricDetails{
		desc: prometheus.NewDesc(
			"rounds_with_ace",
			"Rounds with Ace by [season,user,gamemode,role]",
			labelsActionsWithRole,
			nil,
		),
		metricType: prometheus.GaugeValue,
	}
	metricRoundsWithClutch = metricDetails{
		desc: prometheus.NewDesc(
			"rounds_with_clutch",
			"Rounds with clutch by [season,user,gamemode,role]",
			labelsActionsWithRole,
			nil,
		),
		metricType: prometheus.GaugeValue,
	}
	metricRoundsWithEntryDeath = metricDetails{
		desc: prometheus.NewDesc(
			"rounds_with_entry_death",
			"Rounds with entry death by [season,user,gamemode,role]",
			labelsActionsWithRole,
			nil,
		),
		metricType: prometheus.GaugeValue,
	}
	metricRoundsWithEntryKill = metricDetails{
		desc: prometheus.NewDesc(
			"rounds_with_entry_kill",
			"Rounds with entry kill by [season,user,gamemode,role]",
			labelsActionsWithRole,
			nil,
		),
		metricType: prometheus.GaugeValue,
	}
	metricRoundsWithKOST = metricDetails{
		desc: prometheus.NewDesc(
			"rounds_with_kost",
			"Rounds with KOST (Kill,Objective,Survival,Trade) by [season,user,gamemode,role]",
			labelsActionsWithRole,
			nil,
		),
		metricType: prometheus.GaugeValue,
	}
	metricRoundsWithKill = metricDetails{
		desc: prometheus.NewDesc(
			"rounds_with_kill",
			"Rounds with kill by [season,user,gamemode,role]",
			labelsActionsWithRole,
			nil,
		),
		metricType: prometheus.GaugeValue,
	}
	metricRoundsWithMultiKill = metricDetails{
		desc: prometheus.NewDesc(
			"rounds_with_multikill",
			"Rounds with multi-kill by [season,user,gamemode,role]",
			labelsActionsWithRole,
			nil,
		),
		metricType: prometheus.GaugeValue,
	}

	labelsActions       = []string{"season", "username", "gamemode"}
	metricMatchesPlayed = metricDetails{
		desc: prometheus.NewDesc(
			"matches_played",
			"Matches won by [season,user,gamemode]",
			labelsActions,
			nil,
		),
		metricType: prometheus.CounterValue,
	}
	metricMatchesWon = metricDetails{
		desc: prometheus.NewDesc(
			"matches_won",
			"Matches won by [season,user,gamemode]",
			labelsActions,
			nil,
		),
		metricType: prometheus.CounterValue,
	}
	metricMatchesLost = metricDetails{
		desc: prometheus.NewDesc(
			"matches_lost",
			"Matches lost by [season,user,gamemode]",
			labelsActions,
			nil,
		),
		metricType: prometheus.CounterValue,
	}
)

var allActionsDescs = []metricDetails{
	metricKills,
	metricDeaths,
	metricEntryKills,
	metricEntryDeaths,
	metricHeadshotPercentage,
	metricRoundsPlayed,
	metricRoundsWon,
	metricRoundsLost,
	metricRoundsWithAce,
	metricRoundsWithClutch,
	metricRoundsWithEntryDeath,
	metricRoundsWithEntryKill,
	metricRoundsWithKOST,
	metricRoundsWithKill,
	metricRoundsWithMultiKill,
	metricMatchesPlayed,
	metricMatchesWon,
	metricMatchesLost,
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
	metricsNoRole := []metricInstance{
		{metricMatchesPlayed, float64(gameMode.MatchesPlayed)},
		{metricMatchesWon, float64(gameMode.MatchesWon)},
		{metricMatchesLost, float64(gameMode.MatchesLost)},
	}
	for _, metric := range metricsNoRole {
		ch <- prometheus.MustNewConstMetric(
			metric.details.desc,
			metric.details.metricType,
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
		metricsWithRole := []metricInstance{
			{metricKills, float64(roleStats.Kills)},
			{metricDeaths, float64(roleStats.Deaths)},
			{metricEntryKills, float64(roleStats.EntryKills)},
			{metricEntryDeaths, float64(roleStats.EntryDeaths)},
			{metricHeadshotPercentage, roleStats.HeadshotPercentage},
			{metricRoundsPlayed, float64(roleStats.RoundsPlayed)},
			{metricRoundsWon, float64(roleStats.RoundsWon)},
			{metricRoundsLost, float64(roleStats.RoundsLost)},
			{metricRoundsWithAce, roleStats.RoundsWithAce},
			{metricRoundsWithClutch, roleStats.RoundsWithClutch},
			{metricRoundsWithEntryDeath, roleStats.RoundsWithEntryDeath},
			{metricRoundsWithEntryKill, roleStats.RoundsWithEntryKill},
			{metricRoundsWithKOST, roleStats.RoundsWithKOST},
			{metricRoundsWithKill, roleStats.RoundsWithKill},
			{metricRoundsWithMultiKill, roleStats.RoundsWithMultikill},
		}

		for _, metric := range metricsWithRole {
			ch <- prometheus.MustNewConstMetric(
				metric.details.desc,
				metric.details.metricType,
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
	for _, d := range allActionsDescs {
		ch <- prometheus.NewInvalidMetric(d.desc, err)
	}
}
