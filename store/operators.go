package store

import (
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/stnokott/r6api"
	"github.com/stnokott/r6api/types/metadata"
	"github.com/stnokott/r6api/types/stats"
)

func (s *Store) sendOperatorStats(profile *r6api.Profile, meta *metadata.Metadata, t time.Time) error {
	currentSeason := meta.Seasons[len(meta.Seasons)-1]
	operatorStats := new(stats.OperatorStats)
	if err := s.api.GetStats(profile, currentSeason.Slug, operatorStats); err != nil {
		return err
	}

	gameModes := map[string]*stats.NamedTeamRoles{
		"all":      operatorStats.All,
		"casual":   operatorStats.Casual,
		"unranked": operatorStats.Unranked,
		"ranked":   operatorStats.Ranked,
	}

	for gameModeName, gameModeStats := range gameModes {
		roles := map[string][]stats.NamedTeamRoleStats{
			"all":     gameModeStats.All,
			"attack":  gameModeStats.Attack,
			"defence": gameModeStats.Defence,
		}
		for roleName, roleStats := range roles {
			for _, operatorStats := range roleStats {
				s.influxAPI.WritePoint(
					influxdb2.NewPoint(
						"actions",
						map[string]string{
							"season_slug": currentSeason.Slug,
							"season_name": currentSeason.Name,
							"username":    profile.Name,
							"gamemode":    gameModeName,
							"role":        roleName,
							"operator":    operatorStats.Name,
						},
						map[string]interface{}{
							"kills":                   operatorStats.Kills,
							"deaths":                  operatorStats.Deaths,
							"assists":                 operatorStats.Assists,
							"melee_kills":             operatorStats.MeleeKills,
							"team_kills":              operatorStats.TeamKills,
							"trades":                  operatorStats.Trades,
							"revives":                 operatorStats.Revives,
							"headshots":               operatorStats.Headshots,
							"rounds_played":           operatorStats.RoundsPlayed,
							"rounds_won":              operatorStats.RoundsWon,
							"rounds_lost":             operatorStats.RoundsLost,
							"minutes_played":          operatorStats.MinutesPlayed,
							"kills_per_round":         operatorStats.KillsPerRound,
							"headshot_percentage":     operatorStats.HeadshotPercentage,
							"entry_deaths":            operatorStats.EntryDeaths,
							"entry_death_trades":      operatorStats.EntryDeathTrades,
							"entry_kills":             operatorStats.EntryKills,
							"entry_kill_trades":       operatorStats.EntryKillTrades,
							"rounds_survived":         operatorStats.RoundsSurvived,
							"rounds_with_kill":        operatorStats.RoundsWithKill,
							"rounds_with_multikill":   operatorStats.RoundsWithMultikill,
							"rounds_with_ace":         operatorStats.RoundsWithAce,
							"rounds_with_clutch":      operatorStats.RoundsWithClutch,
							"rounds_with_kost":        operatorStats.RoundsWithKOST,
							"rounds_with_entry_death": operatorStats.RoundsWithEntryDeath,
							"rounds_with_entry_kill":  operatorStats.RoundsWithEntryKill,
							"distance_per_round":      operatorStats.DistancePerRound,
							"distance_total":          operatorStats.DistanceTotal,
							"time_alive_per_match":    operatorStats.TimeAlivePerMatch,
							"time_dead_per_match":     operatorStats.TimeDeadPerMatch,
						},
						t,
					),
				)
			}
		}
	}
	return nil
}
