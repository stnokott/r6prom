package store

import (
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/stnokott/r6api"
	"github.com/stnokott/r6api/types/metadata"
	"github.com/stnokott/r6api/types/stats"
)

func (s *Store) sendMatchStats(profile *r6api.Profile, meta *metadata.Metadata, t time.Time) error {
	currentSeason := meta.Seasons[len(meta.Seasons)-1]
	summarizedStats := new(stats.SummarizedStats)
	if err := s.api.GetStats(profile, currentSeason.Slug, summarizedStats); err != nil {
		return err
	}

	gameModes := map[string]*stats.SummarizedGameModeStats{
		"all":      summarizedStats.All,
		"casual":   summarizedStats.Casual,
		"unranked": summarizedStats.Unranked,
		"ranked":   summarizedStats.Ranked,
	}

	for gameModeName, gameModeStats := range gameModes {
		s.influxAPI.WritePoint(
			influxdb2.NewPoint(
				"matches",
				map[string]string{
					"season_slug": currentSeason.Slug,
					"season_name": currentSeason.Name,
					"username":    profile.Name,
					"gamemode":    gameModeName,
				},
				map[string]interface{}{
					"matches_played": gameModeStats.MatchesPlayed,
					"matches_won":    gameModeStats.MatchesWon,
					"matches_lost":   gameModeStats.MatchesLost,
				},
				t,
			),
		)
	}
	return nil
}
