package metrics

import (
	"time"

	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/stnokott/r6api"
	"github.com/stnokott/r6api/types/metadata"
)

type StatResponse struct {
	P    *write.Point
	Done bool
	Err  error
}

type StatSenderFunc func(*r6api.R6API, *r6api.Profile, *metadata.Metadata, time.Time, chan<- StatResponse)
