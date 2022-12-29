package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stnokott/r6api/types/stats"
)

var (
	labelsActionsWithRole = []string{"season", "username", "gamemode", "role", "operator"}
	metricKills           = metricDetails{
		desc: prometheus.NewDesc(
			"kills",
			"Kills by [season,user,gamemode,role,operator]",
			labelsActionsWithRole,
			nil,
		),
		metricType: prometheus.CounterValue,
	}
	metricDeaths = metricDetails{
		desc: prometheus.NewDesc(
			"deaths",
			"Deaths by [season,user,gamemode,role,operator]",
			labelsActionsWithRole,
			nil,
		),
		metricType: prometheus.CounterValue,
	}
	metricEntryKills = metricDetails{
		desc: prometheus.NewDesc(
			"entry_kills",
			"Entry kills (first kill in a round) by [season,user,gamemode,role,operator]",
			labelsActionsWithRole,
			nil,
		),
		metricType: prometheus.CounterValue,
	}
	metricEntryDeaths = metricDetails{
		desc: prometheus.NewDesc(
			"entry_deaths",
			"Entry deaths (first death in a round) by [season,user,gamemode,role,operator]",
			labelsActionsWithRole,
			nil,
		),
		metricType: prometheus.CounterValue,
	}
	metricHeadshotPercentage = metricDetails{
		desc: prometheus.NewDesc(
			"headshot_percentage",
			"Headshot percentage by [season,user,gamemode,role,operator]",
			labelsActionsWithRole,
			nil,
		),
		metricType: prometheus.GaugeValue,
	}
	metricRoundsPlayed = metricDetails{
		desc: prometheus.NewDesc(
			"rounds_played",
			"Rounds won by [season,user,gamemode,role,operator]",
			labelsActionsWithRole,
			nil,
		),
		metricType: prometheus.CounterValue,
	}
	metricRoundsWon = metricDetails{
		desc: prometheus.NewDesc(
			"rounds_won",
			"Rounds won by [season,user,gamemode,role,operator]",
			labelsActionsWithRole,
			nil,
		),
		metricType: prometheus.CounterValue,
	}
	metricRoundsLost = metricDetails{
		desc: prometheus.NewDesc(
			"rounds_lost",
			"Rounds lost by [season,user,gamemode,role,operator]",
			labelsActionsWithRole,
			nil,
		),
		metricType: prometheus.CounterValue,
	}
	metricRoundsWithAce = metricDetails{
		desc: prometheus.NewDesc(
			"rounds_with_ace",
			"Rounds with Ace by [season,user,gamemode,role,operator]",
			labelsActionsWithRole,
			nil,
		),
		metricType: prometheus.GaugeValue,
	}
	metricRoundsWithClutch = metricDetails{
		desc: prometheus.NewDesc(
			"rounds_with_clutch",
			"Rounds with clutch by [season,user,gamemode,role,operator]",
			labelsActionsWithRole,
			nil,
		),
		metricType: prometheus.GaugeValue,
	}
	metricRoundsWithEntryDeath = metricDetails{
		desc: prometheus.NewDesc(
			"rounds_with_entry_death",
			"Rounds with entry death by [season,user,gamemode,role,operator]",
			labelsActionsWithRole,
			nil,
		),
		metricType: prometheus.GaugeValue,
	}
	metricRoundsWithEntryKill = metricDetails{
		desc: prometheus.NewDesc(
			"rounds_with_entry_kill",
			"Rounds with entry kill by [season,user,gamemode,role,operator]",
			labelsActionsWithRole,
			nil,
		),
		metricType: prometheus.GaugeValue,
	}
	metricRoundsWithKOST = metricDetails{
		desc: prometheus.NewDesc(
			"rounds_with_kost",
			"Rounds with KOST (Kill,Objective,Survival,Trade) by [season,user,gamemode,role,operator]",
			labelsActionsWithRole,
			nil,
		),
		metricType: prometheus.GaugeValue,
	}
	metricRoundsWithKill = metricDetails{
		desc: prometheus.NewDesc(
			"rounds_with_kill",
			"Rounds with kill by [season,user,gamemode,role,operator]",
			labelsActionsWithRole,
			nil,
		),
		metricType: prometheus.GaugeValue,
	}
	metricRoundsWithMultiKill = metricDetails{
		desc: prometheus.NewDesc(
			"rounds_with_multikill",
			"Rounds with multi-kill by [season,user,gamemode,role,operator]",
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
	SummarizedStats *stats.SummarizedStats
	OperatorStats   *stats.OperatorStats
	Username        string
}

func (p ActionsMetricProvider) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(p, ch)
}

func (p ActionsMetricProvider) Collect(ch chan<- prometheus.Metric) {
	p.collectFromSummarizedStats(ch)
	p.collectFromOperatorStats(ch)
}

func (p ActionsMetricProvider) collectFromSummarizedStats(ch chan<- prometheus.Metric) {
	if p.SummarizedStats == nil {
		return
	}

	gameModes := map[string]*stats.SummarizedGameModeStats{
		"all":      p.SummarizedStats.All,
		"casual":   p.SummarizedStats.Casual,
		"unranked": p.SummarizedStats.Unranked,
		"ranked":   p.SummarizedStats.Ranked,
	}
	for gameModeName, gameMode := range gameModes {
		if gameMode == nil {
			continue
		}

		// can do general metrics with summarized data
		p.collectGeneralMetrics(ch, gameModeName, gameMode)
		p.collectSummarizedMetrics(ch, gameModeName, gameMode)
	}
}

func (p ActionsMetricProvider) collectGeneralMetrics(ch chan<- prometheus.Metric, gameModeName string, gameMode *stats.SummarizedGameModeStats) {
	metrics := []metricInstance{
		{metricMatchesPlayed, float64(gameMode.MatchesPlayed)},
		{metricMatchesWon, float64(gameMode.MatchesWon)},
		{metricMatchesLost, float64(gameMode.MatchesLost)},
	}
	for _, metric := range metrics {
		ch <- prometheus.MustNewConstMetric(
			metric.details.desc,
			metric.details.metricType,
			metric.value,
			p.SummarizedStats.SeasonSlug,
			p.Username,
			gameModeName,
		)
	}
}

func (p ActionsMetricProvider) collectSummarizedMetrics(ch chan<- prometheus.Metric, gameModeName string, gameMode *stats.SummarizedGameModeStats) {
	roles := map[string]*stats.DetailedStats{
		"all":     gameMode.All,
		"attack":  gameMode.Attack,
		"defence": gameMode.Defence,
	}
	for roleName, roleStats := range roles {
		metrics := []metricInstance{
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

		for _, metric := range metrics {
			ch <- prometheus.MustNewConstMetric(
				metric.details.desc,
				metric.details.metricType,
				metric.value,
				p.SummarizedStats.SeasonSlug,
				p.Username,
				gameModeName,
				roleName,
				"all", // operator name
			)
		}
	}
}

func (p ActionsMetricProvider) collectFromOperatorStats(ch chan<- prometheus.Metric) {
	if p.OperatorStats == nil {
		return
	}

	gameModes := map[string]*stats.NamedTeamRoles{
		"all":      p.OperatorStats.All,
		"casual":   p.OperatorStats.Casual,
		"unranked": p.OperatorStats.Unranked,
		"ranked":   p.OperatorStats.Ranked,
	}
	for gameModeName, gameMode := range gameModes {
		if gameMode == nil {
			continue
		}

		p.collectOperatorMetrics(ch, gameModeName, gameMode)
	}
}

func (p ActionsMetricProvider) collectOperatorMetrics(ch chan<- prometheus.Metric, gameModeName string, gameMode *stats.NamedTeamRoles) {
	roles := map[string][]stats.NamedTeamRoleStats{
		"all":     gameMode.All,
		"attack":  gameMode.Attack,
		"defence": gameMode.Defence,
	}
	for roleName, roleStats := range roles {
		for _, operatorStats := range roleStats {
			metrics := []metricInstance{
				{metricKills, float64(operatorStats.Kills)},
				{metricDeaths, float64(operatorStats.Deaths)},
				{metricEntryKills, float64(operatorStats.EntryKills)},
				{metricEntryDeaths, float64(operatorStats.EntryDeaths)},
				{metricHeadshotPercentage, operatorStats.HeadshotPercentage},
				{metricRoundsPlayed, float64(operatorStats.RoundsPlayed)},
				{metricRoundsWon, float64(operatorStats.RoundsWon)},
				{metricRoundsLost, float64(operatorStats.RoundsLost)},
				{metricRoundsWithAce, operatorStats.RoundsWithAce},
				{metricRoundsWithClutch, operatorStats.RoundsWithClutch},
				{metricRoundsWithEntryDeath, operatorStats.RoundsWithEntryDeath},
				{metricRoundsWithEntryKill, operatorStats.RoundsWithEntryKill},
				{metricRoundsWithKOST, operatorStats.RoundsWithKOST},
				{metricRoundsWithKill, operatorStats.RoundsWithKill},
				{metricRoundsWithMultiKill, operatorStats.RoundsWithMultikill},
			}

			for _, metric := range metrics {
				if operatorStats.Name != "" {
					ch <- prometheus.MustNewConstMetric(
						metric.details.desc,
						metric.details.metricType,
						metric.value,
						p.SummarizedStats.SeasonSlug,
						p.Username,
						gameModeName,
						roleName,
						operatorStats.Name,
					)
				}
			}
		}
	}
}

func ActionsErr(ch chan<- prometheus.Metric, err error) {
	for _, d := range allActionsDescs {
		ch <- prometheus.NewInvalidMetric(d.desc, err)
	}
}
