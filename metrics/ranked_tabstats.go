package metrics

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/stnokott/r6api"
	"github.com/stnokott/r6api/constants"
	"github.com/stnokott/r6api/types/metadata"
)

type rankedTabStats struct {
	CurrentSeason struct {
		Ranked struct {
			MMR      int    `json:"mmr"`
			RealMMR  int    `json:"real_mmr"`
			RankSlug string `json:"rank_slug"`
		} `json:"ranked"`
	} `json:"current_season_records"`
}

const tabStatsBaseURL = "https://r6.apitab.net/website/profiles/"

func getRankedTabStats(profile *r6api.Profile) (result *rankedTabStats, err error) {
	requestURL := tabStatsBaseURL + profile.ProfileID + "?update=false"
	req, _ := http.NewRequest("GET", requestURL, nil)
	req.Header.Add("User-Agent", constants.USER_AGENT)
	req.Header.Add("Accept", "application/json")
	var resp *http.Response
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	defer func() {
		innerErr := resp.Body.Close()
		if err == nil {
			err = innerErr
		}
	}()

	result = new(rankedTabStats)
	err = json.NewDecoder(resp.Body).Decode(result)
	return
}

func SendRankedTabStatsStats(_ *r6api.R6API, profile *r6api.Profile, meta *metadata.Metadata, t time.Time, chData chan<- StatResponse) {
	tabStats, err := getRankedTabStats(profile)
	if err != nil {
		chData <- StatResponse{Err: err}
		return
	}

	currentSeason := meta.Seasons[len(meta.Seasons)-1]
	rankSlugSplit := strings.SplitN(tabStats.CurrentSeason.Ranked.RankSlug, "-", 2)
	seasonID, err := strconv.Atoi(rankSlugSplit[0])
	if err != nil {
		chData <- StatResponse{Err: err}
		return
	}

	chData <- StatResponse{
		P: influxdb2.NewPoint(
			"ranked_tabstats",
			map[string]string{
				"season_slug": currentSeason.Slug,
				"season_name": currentSeason.Name,
				"season_id":   strconv.Itoa(seasonID),
				"username":    profile.Name,
			},
			map[string]interface{}{
				"mmr":       tabStats.CurrentSeason.Ranked.MMR,
				"real_mmr":  tabStats.CurrentSeason.Ranked.RealMMR,
				"rank_slug": rankSlugSplit[1],
			},
			t,
		),
	}
	chData <- StatResponse{Done: true}
}
