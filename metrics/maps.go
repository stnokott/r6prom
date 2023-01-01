package metrics

import (
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/stnokott/r6api"
	"github.com/stnokott/r6api/types/metadata"
	"github.com/stnokott/r6api/types/stats"
)

func SendMapStats(api *r6api.R6API, profile *r6api.Profile, meta *metadata.Metadata, t time.Time, chData chan<- StatResponse) {
	currentSeason := meta.Seasons[len(meta.Seasons)-1]
	mapStats := new(stats.MapStats)
	if err := api.GetStats(profile, currentSeason.Slug, mapStats); err != nil {
		chData <- StatResponse{Err: err}
		return
	}

	gameModes := map[string]*stats.NamedTeamRoles{
		"all":      mapStats.All,
		"casual":   mapStats.Casual,
		"unranked": mapStats.Unranked,
		"ranked":   mapStats.Ranked,
	}

	for gameModeName, gameModeStats := range gameModes {
		roles := map[string][]stats.NamedTeamRoleStats{
			"all":     gameModeStats.All,
			"attack":  gameModeStats.Attack,
			"defence": gameModeStats.Defence,
		}
		for roleName, roleStats := range roles {
			for _, mapStats := range roleStats {
				chData <- StatResponse{
					P: influxdb2.NewPoint(
						"actions",
						map[string]string{
							"season_slug": currentSeason.Slug,
							"season_name": currentSeason.Name,
							"username":    profile.Name,
							"gamemode":    gameModeName,
							"role":        roleName,
							"map":         mapStats.Name,
						},
						map[string]interface{}{
							"kills":                   mapStats.Kills,
							"deaths":                  mapStats.Deaths,
							"assists":                 mapStats.Assists,
							"melee_kills":             mapStats.MeleeKills,
							"team_kills":              mapStats.TeamKills,
							"trades":                  mapStats.Trades,
							"revives":                 mapStats.Revives,
							"headshots":               mapStats.Headshots,
							"rounds_played":           mapStats.RoundsPlayed,
							"rounds_won":              mapStats.RoundsWon,
							"rounds_lost":             mapStats.RoundsLost,
							"minutes_played":          mapStats.MinutesPlayed,
							"kills_per_round":         mapStats.KillsPerRound,
							"headshot_percentage":     mapStats.HeadshotPercentage,
							"entry_deaths":            mapStats.EntryDeaths,
							"entry_death_trades":      mapStats.EntryDeathTrades,
							"entry_kills":             mapStats.EntryKills,
							"entry_kill_trades":       mapStats.EntryKillTrades,
							"rounds_survived":         mapStats.RoundsSurvived,
							"rounds_with_kill":        mapStats.RoundsWithKill,
							"rounds_with_multikill":   mapStats.RoundsWithMultikill,
							"rounds_with_ace":         mapStats.RoundsWithAce,
							"rounds_with_clutch":      mapStats.RoundsWithClutch,
							"rounds_with_kost":        mapStats.RoundsWithKOST,
							"rounds_with_entry_death": mapStats.RoundsWithEntryDeath,
							"rounds_with_entry_kill":  mapStats.RoundsWithEntryKill,
							"distance_per_round":      mapStats.DistancePerRound,
							"distance_total":          mapStats.DistanceTotal,
							"time_alive_per_match":    mapStats.TimeAlivePerMatch,
							"time_dead_per_match":     mapStats.TimeDeadPerMatch,
						},
						t,
					),
				}

			}
		}
	}
	chData <- StatResponse{Done: true}
}
