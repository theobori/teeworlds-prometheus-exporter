package econ

import (
	"fmt"
	"sync"

	twecon "github.com/theobori/teeworlds-econ"
)

// Represents an event
type EconEventEntry struct {
	// Event name
	Name string
	// Event regex
	Regex string
}

var (
	// Econ events used for metrics
	EconEvents = []EconEventEntry{
		// Teeworlds 0.7 events metrics
		{
			Name:  "message",
			Regex: `\[chat\]: .*`,
		},
		{
			Name:  "kill",
			Regex: `\[game\]: kill killer=.*`,
		},
		{
			Name:  "captured_flag",
			Regex: `\[game\]: flag_capture player=.*`,
		},
	}
)

// Econ metrics storage format. (name, count)
type EconMetrics map[string]uint

// Econ manager map value
type EconMananagerEntry struct {
	// Econ client controller
	Econ *twecon.Econ
	// Metrics associated with a econ server
	Metrics EconMetrics
	// Indicating is the econ client is handling events
	IsHandling bool
}

// Econ manager map key
type EconMananagerKey struct {
	// Server IP address
	Host string
	// Server port
	Port uint16
}

// Create a new EconMananagerEntry struct
func NewEconManagerEntry(e *twecon.Econ) *EconMananagerEntry {
	return &EconMananagerEntry{
		Econ:       e,
		Metrics:    EconMetrics{},
		IsHandling: false,
	}
}

// Econs manager
type EconManager struct {
	econs map[EconMananagerKey]*EconMananagerEntry
	mu    sync.Mutex
}

// Create a new econ manager
func NewEconManager() *EconManager {
	return &EconManager{
		econs: make(map[EconMananagerKey]*EconMananagerEntry),
	}
}

// Register a econ client
func (em *EconManager) Register(e *twecon.Econ) error {
	if e == nil {
		return fmt.Errorf("nil econ")
	}

	c := e.Config()

	// init metrics with zeros
	metrics := EconMetrics{}

	for _, econEvent := range EconEvents {
		metrics[econEvent.Name] = 0
	}

	entry := EconMananagerEntry{
		Econ:    e,
		Metrics: metrics,
	}

	k := EconMananagerKey{
		Host: c.Host,
		Port: c.Port,
	}

	em.econs[k] = &entry

	return nil
}

// Delete a econ client
func (em *EconManager) Delete(k EconMananagerKey) {
	delete(em.econs, k)
}

// Register a econ events
func (em *EconManager) RegisterEconEvents() error {
	em.mu.Lock()
	defer em.mu.Unlock()

	for _, entry := range em.econs {
		if entry == nil {
			continue
		}

		err := registerMetricEvents(entry.Econ, &entry.Metrics)
		if err != nil {
			return err
		}
	}

	return nil
}

// Return metrics per econ server
func (em *EconManager) EconServersMetrics() map[EconMananagerKey]EconMetrics {
	ret := make(map[EconMananagerKey]EconMetrics)

	em.mu.Lock()
	defer em.mu.Unlock()

	for k, e := range em.econs {
		ret[k] = e.Metrics
	}

	return ret
}

// Start handling event for every econ client
func (em *EconManager) StartHandle() error {
	for _, entry := range em.econs {
		if entry == nil {
			continue
		}

		if !entry.IsHandling && entry.Econ != nil {
			entry.IsHandling = true

			go (*entry.Econ).HandleEvents()
		}
	}

	return nil
}

// Register the events for metrics
func registerMetricEvents(e *twecon.Econ, metrics *EconMetrics) error {
	if e == nil || metrics == nil {
		return fmt.Errorf("nil econ or metrics")
	}

	for _, event := range EconEvents {
		event := twecon.EconEvent{
			Name:  event.Name,
			Regex: event.Regex,
			Func: func(econ *twecon.Econ, eventPayload string) any {
				(*metrics)[event.Name]++

				return nil
			},
		}

		if err := e.EventManager.Register(&event); err != nil {
			return err
		}
	}

	return nil
}
