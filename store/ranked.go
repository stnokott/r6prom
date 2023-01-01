package store

import (
	"fmt"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/stnokott/r6api"
	"github.com/stnokott/r6api/types/metadata"
)

func (s *Store) sendRankedStats(profile *r6api.Profile, meta *metadata.Metadata, t time.Time) error {
	seasons, err := s.api.GetRankedHistory(profile, 1)
	if err != nil {
		return err
	}
	if len(seasons) == 0 {
		return fmt.Errorf("got no ranked history for user %s", profile.Name)
	}
	stats := seasons[0]

	s.influxAPI.WritePoint(
		influxdb2.NewPoint(
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
	)
	return nil
}
