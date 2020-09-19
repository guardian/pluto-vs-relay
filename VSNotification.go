package main

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

func (doc *VSNotificationDocument) GetAction() string {
	entry, didFind := doc.FindKey("action")
	if didFind {
		return entry.Value
	} else {
		return ""
	}
}
