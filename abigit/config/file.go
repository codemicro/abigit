package config

import (
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const configFileName = "config.yml"

var (
	rawConfigFileContents map[string]any
	lastKey               string
)

func loadConfigFileFromDisk() {
	if rawConfigFileContents != nil {
		return
	}

	fcont, err := os.ReadFile(configFileName)
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to load file %s", configFileName)
	}
	rawConfigFileContents = make(map[string]any)
	if err := yaml.Unmarshal(fcont, &rawConfigFileContents); err != nil {
		log.Fatal().Err(err).Msg("could not unmarshal config file")
	}
}

type optionalItem struct {
	item  any
	found bool
}

var indexedPartRegexp = regexp.MustCompile(`(?m)([a-zA-Z]+)(?:\[(\d+)\])?`)

func fetchFromFile(key string) optionalItem {
	// http[2].bananas
	loadConfigFileFromDisk()
	lastKey = key

	parts := strings.Split(key, ".")
	var cursor any = rawConfigFileContents
	for _, part := range parts {
		components := indexedPartRegexp.FindStringSubmatch(part)
		key := components[1]
		index, _ := strconv.ParseInt(components[2], 10, 32)
		isIndexed := components[2] != ""

		item, found := cursor.(map[string]any)[key]
		if !found {
			return optionalItem{nil, false}
		}

		if isIndexed {
			arr, conversionOk := item.([]any)
			if !conversionOk {
				log.Fatal().Msgf("attempted to index non-indexable item %s", key)
			}
			cursor = arr[index]
		} else {
			cursor = item
		}
	}
	return optionalItem{cursor, true}
}

func required(key string) optionalItem {
	opt := fetchFromFile(key)
	if !opt.found {
		log.Fatal().Msgf("required key %s not found", lastKey)
	}
	return opt
}

func withDefault(key string, defaultValue any) optionalItem {
	opt := fetchFromFile(key)
	if !opt.found {
		return optionalItem{item: defaultValue, found: true}
	}
	return opt
}

func asInt(x optionalItem) int {
	if !x.found {
		return 0
	}
	return x.item.(int)
}

func asString(x optionalItem) string {
	if !x.found {
		return ""
	}
	return x.item.(string)
}

func asBool(x optionalItem) bool {
	if !x.found {
		return false
	}
	return x.item.(bool)
}
