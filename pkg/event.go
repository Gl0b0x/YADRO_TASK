package pkg

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

type Event struct {
	ID          int
	timeEvent   time.Time
	clientName  string
	tableNumber int
}

func parseEvent(eventStr []string) Event {
	var tableNumber int
	timeEvent, _ := parseTime(eventStr[0])
	if timeEvent.Hour() == 0 && timeEvent.Minute() == 0 {
		timeEvent = timeEvent.AddDate(0, 0, 1)
	}
	eventID, _ := strconv.Atoi(eventStr[1])
	clientName := eventStr[2]
	if eventID == 2 && len(eventStr) == 4 {
		tableNumber, _ = strconv.Atoi(eventStr[3])
	}
	return Event{ID: eventID, timeEvent: timeEvent, clientName: clientName, tableNumber: tableNumber}
}

func eventToString(t time.Time, ID int, clientName string, tableNumber int) string {
	if tableNumber == 0 {
		return fmt.Sprintf("%s %d %s", t.Format("15:04"), ID, clientName)
	}
	return fmt.Sprintf("%s %d %s %d", t.Format("15:04"), ID, clientName, tableNumber)
}

func parseTime(timeString string) (time.Time, bool) {
	eventTime, err := time.Parse("15:04", timeString)
	if err != nil {
		return time.Time{}, false
	}
	if eventTime.Format("15:04") != timeString {
		return time.Time{}, false
	}
	return eventTime, true
}

func validEvent(eventStr []string, countComputers int) bool {
	if len(eventStr) != 3 && len(eventStr) != 4 {
		return false
	}
	_, ok := parseTime(eventStr[0])
	if !ok {
		return false
	}
	if eventStr[1] != "1" && eventStr[1] != "2" && eventStr[1] != "3" && eventStr[1] != "4" {
		return false
	}
	if !validClientName(eventStr[2]) {
		return false
	}
	if eventStr[1] == "2" {
		if len(eventStr) != 4 {
			return false
		}
		n, err := strconv.Atoi(eventStr[3])
		if err != nil || n < 1 || n > countComputers {
			return false
		}
	} else {
		if len(eventStr) != 3 {
			return false
		}
	}
	return true
}

func validClientName(clientName string) bool {
	expression := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	return expression.MatchString(clientName)
}
