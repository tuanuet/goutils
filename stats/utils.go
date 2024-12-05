package stats

import (
	"sort"
	"strings"
)

const (
	prometheusKeySeparator = "_"
)

func buildHashKey(contextTag, key string, sortedTagKeys ...string) string {
	var hashKey strings.Builder

	if len(contextTag) > 0 {
		hashKey.WriteString(contextTag)
		hashKey.WriteString(prometheusKeySeparator)
	}

	hashKey.WriteString(key)

	if len(sortedTagKeys) > 0 {
		hashKey.WriteString(prometheusKeySeparator)
		hashKey.WriteString(strings.Join(sortedTagKeys, prometheusKeySeparator))
	}
	return hashKey.String()
}

func getSortedTagKeys(tags ...Tag) []string {
	tagKeys := make([]string, 0, len(tags))
	for _, tag := range tags {
		tagKeys = append(tagKeys, tag.Key)
	}
	sort.Strings(tagKeys)
	return tagKeys
}

func getTagMap(tags ...Tag) map[string]string {
	tagMap := make(map[string]string)
	for _, tag := range tags {
		tagMap[tag.Key] = tag.Value
	}
	return tagMap
}
