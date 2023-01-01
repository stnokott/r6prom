package metrics

import (
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/stnokott/r6api"
	"github.com/stnokott/r6api/types/metadata"
	"github.com/stnokott/r6api/types/stats"
)

func SendMatchStats(api *r6api.R6API, profile *r6api.Profile, meta *metadata.Metadata, t time.Time, chData chan<- StatResponse) {
	currentSeason := meta.Seasons[len(meta.Seasons)-1]
	summarizedStats := new(stats.SummarizedStats)
	if err := api.GetStats(profile, currentSeason.Slug, summarizedStats); err != nil {
		chData <- StatResponse{Err: err}
		return
	}

	gameModes := map[string]*stats.SummarizedGameModeStats{
		"all":      summarizedStats.All,
		"casual":   summarizedStats.Casual,
		"unranked": summarizedStats.Unranked,
		"ranked":   summarizedStats.Ranked,
	}

	for gameModeName, gameModeStats := range gameModes {
		chData <- StatResponse{
			P: influxdb2.NewPoint(
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
		}
	}
	chData <- StatResponse{Done: true}
}
