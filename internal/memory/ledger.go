package memory

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Event struct {
	ID        int       `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Action    string    `json:"action"` // CREATE, EDIT, DELETE, JUMP
	Payload   string    `json:"payload"`
	Hash      string    `json:"hash"`
}

type Ledger struct {
	filePath string
	events   []Event
}

func NewLedger(filePath string) (*Ledger, error) {
	l := &Ledger{filePath: filePath, events: []Event{}}
	err := l.load()
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return l, nil
}

func (l *Ledger) Record(action string, payload string) (*Event, error) {
	prevHash := ""
	if len(l.events) > 0 {
		prevHash = l.events[len(l.events)-1].Hash
	}

	event := Event{
		ID:        len(l.events) + 1,
		Timestamp: time.Now(),
		Action:    action,
		Payload:   payload,
	}

	data := fmt.Sprintf("%d:%s:%s:%s:%s", event.ID, event.Timestamp.String(), action, payload, prevHash)
	sum := sha256.Sum256([]byte(data))
	event.Hash = hex.EncodeToString(sum[:8])

	l.events = append(l.events, event)
	return &event, l.persist(event)
}

func (l *Ledger) persist(event Event) error {
	file, err := os.OpenFile(l.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	bytes, err := json.Marshal(event)
	if err != nil {
		return err
	}

	_, err = file.WriteString(string(bytes) + "\n")
	return err
}

func (l *Ledger) load() error {
	file, err := os.Open(l.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var event Event
		if err := json.Unmarshal(scanner.Bytes(), &event); err == nil {
			l.events = append(l.events, event)
		}
	}
	return nil
}

func (l *Ledger) GetEvents() []Event {
	return l.events
}

func (l *Ledger) Undo() (*Event, error) {
	if len(l.events) == 0 {
		return nil, fmt.Errorf("no hay eventos para deshacer")
	}

	last := l.events[len(l.events)-1]
	l.events = l.events[:len(l.events)-1]

	// Reescribir archivo WAL
	file, err := os.Create(l.filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	for _, ev := range l.events {
		bytes, _ := json.Marshal(ev)
		file.WriteString(string(bytes) + "\n")
	}

	return &last, nil
}
