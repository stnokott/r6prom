package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Email             string
	Password          string
	ObservedUsernames []string
	RefreshCron       string
	InfluxURL         string
	InfluxAuthToken   string
	InfluxOrg         string
	InfluxBucket      string
}

const (
	envEmail             string = "UBI_EMAIL"
	envPassword          string = "UBI_PASSWORD"
	envObservedUsernames string = "UBI_OBSERVED_USERNAMES"
	envRefreshCron       string = "REFRESH_CRON"
	envInfluxURL         string = "INFLUX_URL"
	envInfluxAuthToken   string = "INFLUX_AUTH_TOKEN"
	envInfluxOrg         string = "INFLUX_ORGANIZATION"
	envInfluxBucket      string = "INFLUX_BUCKET"
)

var requiredEnvs = []string{
	envEmail,
	envPassword,
	envObservedUsernames,
	envRefreshCron,
	envInfluxURL,
	envInfluxAuthToken,
	envInfluxOrg,
	envInfluxBucket,
}

func Load() (c Config, err error) {
	vals := map[string]string{}

	for _, envKey := range requiredEnvs {
		val, exists := os.LookupEnv(envKey)
		if !exists {
			err = fmt.Errorf("environment variable %s missing", envKey)
			return
		}
		vals[envKey] = val
	}

	c.Email = vals[envEmail]
	c.Password = vals[envPassword]
	c.ObservedUsernames = strings.Split(vals[envObservedUsernames], ",")
	if len(c.ObservedUsernames) == 0 {
		err = fmt.Errorf("environment variable %s needs to contain at least one item", envObservedUsernames)
		return
	}
	c.RefreshCron = vals[envRefreshCron]
	c.InfluxURL = vals[envInfluxURL]
	c.InfluxAuthToken = vals[envInfluxAuthToken]
	c.InfluxOrg = vals[envInfluxOrg]
	c.InfluxBucket = vals[envInfluxBucket]

	return
}
