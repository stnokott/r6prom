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

	gameModes := map[string]*map[string]stats.NamedMapStatDetails{
		"all":      mapStats.All,
		"casual":   mapStats.Casual,
		"unranked": mapStats.Unranked,
		"ranked":   mapStats.Ranked,
	}

	for gameModeName, gameModeStats := range gameModes {
		if gameModeStats == nil {
			continue
		}
		for mapName, mapStats := range *gameModeStats {
			labels := map[string]string{
				"season_slug": currentSeason.Slug,
				"season_name": currentSeason.Name,
				"username":    profile.Name,
				"gamemode":    gameModeName,
				"map":         mapName,
			}
			chData <- StatResponse{
				P: influxdb2.NewPoint(
					"maps",
					labels,
					map[string]interface{}{
						"matches_played":          mapStats.MatchesPlayed,
						"matches_won":             mapStats.MatchesWon,
						"matches_lost":            mapStats.MatchesLost,
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
			if mapStats.Bombsites != nil {
				sendMapBombsiteStats(mapStats.Bombsites, chData, labels, t)
			}
		}
	}
	chData <- StatResponse{Done: true}
}

func sendMapBombsiteStats(s *stats.BombsiteGamemodeStats, chData chan<- StatResponse, srcLabels map[string]string, t time.Time) {
	teamRoles := map[string][]stats.BombsiteTeamRoleStats{
		"all":     s.All,
		"attack":  s.Attack,
		"defence": s.Defence,
	}

	for teamRoleName, teamRole := range teamRoles {
		for _, bombsiteStats := range teamRole {
			labels := make(map[string]string, len(srcLabels))
			for k, v := range srcLabels {
				labels[k] = v
			}
			labels["role"] = teamRoleName
			labels["bombsite"] = bombsiteStats.Name

			chData <- StatResponse{
				P: influxdb2.NewPoint(
					"bombsites",
					labels,
					map[string]interface{}{
						"kills":                   bombsiteStats.Kills,
						"deaths":                  bombsiteStats.Deaths,
						"assists":                 bombsiteStats.Assists,
						"melee_kills":             bombsiteStats.MeleeKills,
						"team_kills":              bombsiteStats.TeamKills,
						"trades":                  bombsiteStats.Trades,
						"revives":                 bombsiteStats.Revives,
						"headshots":               bombsiteStats.Headshots,
						"rounds_played":           bombsiteStats.RoundsPlayed,
						"rounds_won":              bombsiteStats.RoundsWon,
						"rounds_lost":             bombsiteStats.RoundsLost,
						"minutes_played":          bombsiteStats.MinutesPlayed,
						"kills_per_round":         bombsiteStats.KillsPerRound,
						"headshot_percentage":     bombsiteStats.HeadshotPercentage,
						"entry_deaths":            bombsiteStats.EntryDeaths,
						"entry_death_trades":      bombsiteStats.EntryDeathTrades,
						"entry_kills":             bombsiteStats.EntryKills,
						"entry_kill_trades":       bombsiteStats.EntryKillTrades,
						"rounds_survived":         bombsiteStats.RoundsSurvived,
						"rounds_with_kill":        bombsiteStats.RoundsWithKill,
						"rounds_with_multikill":   bombsiteStats.RoundsWithMultikill,
						"rounds_with_ace":         bombsiteStats.RoundsWithAce,
						"rounds_with_clutch":      bombsiteStats.RoundsWithClutch,
						"rounds_with_kost":        bombsiteStats.RoundsWithKOST,
						"rounds_with_entry_death": bombsiteStats.RoundsWithEntryDeath,
						"rounds_with_entry_kill":  bombsiteStats.RoundsWithEntryKill,
						"distance_per_round":      bombsiteStats.DistancePerRound,
						"distance_total":          bombsiteStats.DistanceTotal,
						"time_alive_per_match":    bombsiteStats.TimeAlivePerMatch,
						"time_dead_per_match":     bombsiteStats.TimeDeadPerMatch,
					},
					t,
				),
			}
		}
	}
}
