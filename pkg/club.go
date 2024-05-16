package pkg

import (
	"bufio"
	"fmt"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"
)

type ComputerClub struct {
	CountComputers       int
	countAvailablePlaces int
	Computers            []*Computer
	timeStart            time.Time
	timeEnd              time.Time
	price                int
	WaitingList          []string
	Clients              map[string]*Computer
	Events               []Event
}

func InitClub(scanner *bufio.Scanner) (*ComputerClub, bool) {
	var (
		countComputers, price int
		timeStart, timeEnd    time.Time
		ok                    bool
		err                   error
	)
	if !scanner.Scan() {
		return nil, false
	}
	countComputers, err = strconv.Atoi(scanner.Text())
	if err != nil || countComputers < 1 {
		return nil, false
	}
	if !scanner.Scan() {
		return nil, false
	}
	workingTimeStr := strings.Split(scanner.Text(), " ")
	if len(workingTimeStr) != 2 {
		return nil, false
	}
	timeStart, ok = parseTime(workingTimeStr[0])
	if !ok {
		return nil, false
	}
	timeEnd, ok = parseTime(workingTimeStr[1])
	if !ok {
		return nil, false
	}
	if timeEnd.Hour() == 0 && timeStart.Minute() == 0 {
		timeEnd = timeEnd.AddDate(0, 0, 1)
	}
	if timeEnd.Before(timeStart) {
		return nil, false
	}
	if !scanner.Scan() {
		return nil, false
	}
	price, err = strconv.Atoi(scanner.Text())
	if err != nil || price < 0 {
		return nil, false
	}
	computers := make([]*Computer, countComputers+1)
	for i := 1; i <= countComputers; i++ {
		computers[i] = &Computer{Number: i}
	}
	return &ComputerClub{CountComputers: countComputers,
			Computers:            computers,
			timeStart:            timeStart,
			timeEnd:              timeEnd,
			Events:               make([]Event, 0),
			countAvailablePlaces: countComputers,
			Clients:              make(map[string]*Computer, countComputers),
			price:                price},
		true
}

func (c *ComputerClub) DoWork() {
	var (
		eventStr   string
		isOutEvent bool
	)
	fmt.Println(c.timeStart.Format("15:04"))
	for i, event := range c.Events {
		fmt.Println(eventToString(event.timeEvent, event.ID, event.clientName, event.tableNumber))
		switch event.ID {
		case 1:
			eventStr, isOutEvent = c.clientCome(i)
		case 2:
			eventStr, isOutEvent = c.clientSit(i)
		case 3:
			eventStr, isOutEvent = c.clientWaiting(i)
		case 4:
			eventStr, isOutEvent = c.clientLeft(i)
		}
		if isOutEvent {
			fmt.Println(eventStr)
		}
	}
	if len(c.Clients) != 0 {
		clients := make([]string, len(c.Clients))
		count := 0
		for clientName := range c.Clients {
			clients[count] = clientName
			count++
		}
		sort.Strings(clients)
		for _, clientName := range clients {
			fmt.Println(eventToString(c.timeEnd, EventLeft, clientName, 0))
			if c.Clients[clientName] != nil {
				c.Clients[clientName].CountProfit(c.timeEnd, c.price)
			}
			delete(c.Clients, clientName)
			c.countAvailablePlaces++
		}
	}
	fmt.Println(c.timeEnd.Format("15:04"))
	for i := 1; i <= c.CountComputers; i++ {
		fmt.Printf("%d %d %s\n", i, c.Computers[i].Profit, c.Computers[i].TotalTime.Format("15:04"))
	}
}

func (c *ComputerClub) ParseEvents(scanner *bufio.Scanner, countComputers int) bool {
	var predEventTime time.Time
	for scanner.Scan() {
		eventStr := strings.Split(scanner.Text(), " ")
		if !validEvent(eventStr, countComputers) {
			return false
		}
		event := parseEvent(eventStr)
		if predEventTime.IsZero() {
			predEventTime = event.timeEvent
		}
		if event.timeEvent.Before(predEventTime) {
			return false
		}
		predEventTime = event.timeEvent
		c.Events = append(c.Events, event)
	}
	return true
}

func (c *ComputerClub) clientCome(i int) (string, bool) {
	event := c.Events[i]
	if event.timeEvent.Before(c.timeStart) || event.timeEvent.After(c.timeEnd) {
		return eventToString(event.timeEvent, EventError, ErrorClubClose, 0), true
	}
	_, ok := c.Clients[event.clientName]
	if ok {
		return eventToString(event.timeEvent, EventError, ErrorAlready, 0), true
	}
	c.Clients[event.clientName] = nil
	return "", false
}

func (c *ComputerClub) clientSit(i int) (string, bool) {
	event := c.Events[i]
	computer, ok := c.Clients[event.clientName]
	if !ok {
		return eventToString(event.timeEvent, EventError, ErrorClientUnknown, 0), true
	}
	if c.Computers[event.tableNumber].IsBusy {
		return eventToString(event.timeEvent, EventError, ErrorPlaceBusy, 0), true
	}
	if computer != nil {
		computer.IsBusy = false
		c.countAvailablePlaces++
	}
	c.Computers[event.tableNumber].IsBusy = true
	c.Computers[event.tableNumber].TimeStart = event.timeEvent
	c.Clients[event.clientName] = c.Computers[event.tableNumber]
	c.countAvailablePlaces--
	return "", false
}

func (c *ComputerClub) clientWaiting(i int) (string, bool) {
	event := c.Events[i]
	if c.countAvailablePlaces != 0 {
		return eventToString(event.timeEvent, EventError, ErrorAvailablePlaces, 0), true
	}
	if len(c.WaitingList) == c.CountComputers {
		return eventToString(event.timeEvent, EventLeft, event.clientName, 0), true
	}
	c.WaitingList = append(c.WaitingList, event.clientName)
	return "", false
}

func (c *ComputerClub) clientLeft(i int) (string, bool) {
	event := c.Events[i]
	computer, ok := c.Clients[event.clientName]
	if !ok {
		return eventToString(event.timeEvent, EventError, ErrorClientUnknown, 0), true
	}
	delete(c.Clients, event.clientName)
	if computer != nil {
		computer.CountProfit(event.timeEvent, c.price)
		c.countAvailablePlaces++
		if len(c.WaitingList) != 0 {
			clientName := c.WaitingList[0]
			computer.IsBusy = true
			computer.TimeStart = event.timeEvent
			c.Clients[clientName] = computer
			c.countAvailablePlaces--
			c.WaitingList = slices.Delete(c.WaitingList, 0, 1)
			return eventToString(event.timeEvent, EventSit, clientName, computer.Number), true
		}
	}
	return "", false
}
