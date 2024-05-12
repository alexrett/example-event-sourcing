package event_listener

import (
	"strconv"
	"strings"
)

func ExtractIdAndTypeFromPayload(payload string) (string, int, string) {
	data := strings.Split(payload, ",")
	if len(data) == 3 {
		eventTypeID, err := strconv.Atoi(data[1])
		if err != nil {
			return "", 0, ""
		}
		return data[0], eventTypeID, data[2]
	}
	return "", 0, ""
}
