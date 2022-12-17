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
	NumSeasons        uint8
}

const (
	envEmail             string = "UBI_EMAIL"
	envPassword          string = "UBI_PASSWORD"
	envObservedUsernames string = "UBI_OBSERVED_USERNAMES"
)

var requiredEnvs = []string{envEmail, envPassword, envObservedUsernames}

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

	return
}
