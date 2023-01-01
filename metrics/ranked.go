package metrics

import (
	"fmt"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/stnokott/r6api"
	"github.com/stnokott/r6api/types/metadata"
)

func SendRankedStats(api *r6api.R6API, profile *r6api.Profile, meta *metadata.Metadata, t time.Time, chData chan<- StatResponse) {
	seasons, err := api.GetRankedHistory(profile, 1)
	if err != nil {
		chData <- StatResponse{Err: err}
		return
	}
	if len(seasons) == 0 {
		chData <- StatResponse{Err: fmt.Errorf("got no ranked history for user %s", profile.Name)}
		return
	}
	stats := seasons[0]

	chData <- StatResponse{
		P: influxdb2.NewPoint(
			"ranked",
			map[string]string{
				"season_slug": meta.SeasonSlugFromID(stats.SeasonID),
				"season_name": meta.SeasonNameFromID(stats.SeasonID),
				"username":    profile.Name,
			},
			map[string]interface{}{
				"mmr":         stats.MMR,
				"rank":        stats.Rank,
				"skill_mean":  stats.SkillMean,
				"skill_stdev": stats.SkillStdev,
			},
			t,
		),
	}
	chData <- StatResponse{Done: true}
}
