package main

import "strings"

type KeyValuePair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type VSNotificationDocument struct {
	Field []KeyValuePair `json:"field"`
}

/**
returns the given key and "true" or empty key and "false"
*/
func (doc *VSNotificationDocument) FindKey(key string) (KeyValuePair, bool) {
	for _, entry := range doc.Field {
		if entry.Key == key {
			return entry, true
		}
	}
	return KeyValuePair{}, false
}

/**
returns the value of the "action" field (i.e., create/stop/update) or an empty string
if none is set
*/
func (doc *VSNotificationDocument) GetAction() string {
	entry, didFind := doc.FindKey("action")
	if didFind {
		return strings.ToLower(entry.Value)
	} else {
		return ""
	}
}

/**
returns the value of the "type" field (i.e. ESSENCE_VERSION/TRANSCODE/etc.) lowercased or the
provided default value if none is set
*/
func (doc *VSNotificationDocument) GetType(defaultValue string) string {
	entry, didFind := doc.FindKey("type")
	if didFind {
		return strings.ToLower(entry.Value)
	} else {
		return defaultValue
	}
}
