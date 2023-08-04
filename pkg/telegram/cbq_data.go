package telegram

import (
	"fmt"
	"strconv"
	"strings"
)

func cbqParseDataGetAdEventId(data string) (adEventId int64, err error) {
	dataSlice := strings.Split(data, ";")
	if len(dataSlice) != 1 {
		return 0, fmt.Errorf("dataSlice incorrect. dataSlice: %v", dataSlice)
	}

	id, err := strconv.ParseInt(dataSlice[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("error pasge AdEventId: %w", err)
	}

	return id, nil
}
