package notification

import (
	"time"

	"charm.land/lipgloss/v2"
)

type Level int

const (
	Error int = iota
	Info
)

var errorStyle = lipgloss.NewStyle().
	Background(lipgloss.Red)

type Notification struct {
	startTime time.Time
	msg       string
	level     Level
	open      bool
}

func NewNotification(msg string, level Level) Notification {
	return Notification{
		startTime: time.Now(),
		msg:       msg,
		level:     level,
		open:      false,
	}
}

type Model struct {
	notifications []Notification
	bool
}

func New(notifications []Notification) Model {
	return Model{notifications: notifications}
}

func (n Notification) Show() {
	n.open = true
	time.Sleep(10 * time.Second)
	n.open = false
}

// func (m Model) Show() {
// 	for _, n := range m.notifications {
//
// 	}
// }
