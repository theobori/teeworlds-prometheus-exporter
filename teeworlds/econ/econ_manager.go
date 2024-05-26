package econ

import (
	"fmt"
	"sync"

	twecon "github.com/theobori/teeworlds-econ"
)

type econEventMetricsKey struct {
	Name string
	Regex string
}

var (
	EconEventMetrics = map[econEventMetricsKey]func(econMetrics *EconMetrics){
		// Teeworlds 0.7 events metrics
		{
			Name: "message",
			Regex: `\[chat\]: .*`,
		}: func(econMetrics *EconMetrics) {
			econMetrics.MessagesTotal++
		},
		{
			Name: "kill",
			Regex: `\[game\]: kill killer=.*`,
		}: func(econMetrics *EconMetrics) {
			econMetrics.KillsTotal++
		},
		{
			Name: "captured_flag",
			Regex: `\[game\]: flag_capture player=.*`,
		}: func(econMetrics *EconMetrics) {
			econMetrics.CapturedFlagsTotal++
		},
	}
)

type EconMetrics struct {
	MessagesTotal      uint
	CapturedFlagsTotal uint
	KillsTotal         uint
	VotesTotal         uint
}

type EconMananagerEntry struct {
	Econ       *twecon.Econ
	Metrics    EconMetrics
	IsHandling bool
}

type EconMananagerKey struct {
	// Server IP address
	Host string
	// Server port
	Port uint16
}

func NewEconManagerEntry(e *twecon.Econ) *EconMananagerEntry {
	return &EconMananagerEntry{
		Econ:       e,
		Metrics:    EconMetrics{},
		IsHandling: false,
	}
}

type EconManager struct {
	econs map[EconMananagerKey]*EconMananagerEntry
	mu sync.Mutex
}

func NewEconManager() *EconManager {
	return &EconManager{
		econs: make(map[EconMananagerKey]*EconMananagerEntry),
	}
}

func (em *EconManager) Register(e *twecon.Econ) error {
	if e == nil {
		return fmt.Errorf("nil econ")
	}

	c := e.Config()

	entry := EconMananagerEntry{
		Econ:    e,
		Metrics: EconMetrics{},
	}

	k := EconMananagerKey{
		Host: c.Host,
		Port: c.Port,
	}

	em.econs[k] = &entry

	return nil
}

func (em *EconManager) Delete(k EconMananagerKey) {
	delete(em.econs, k)
}

func (em *EconManager) RegisterEconEvents() error {
	em.mu.Lock()
	defer em.mu.Unlock()

	for _, entry := range em.econs {
		if entry == nil {
			continue
		}

		err := em.registerMetricEvents(entry.Econ, &entry.Metrics)
		if err != nil {
			return err
		}
	}

	return nil
}

func (em *EconManager) EconServersMetrics() map[EconMananagerKey]EconMetrics {
	ret := make(map[EconMananagerKey]EconMetrics)

	em.mu.Lock()
	defer em.mu.Unlock()

	for k, e := range em.econs {
		ret[k] = e.Metrics
	}

	return ret
}

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

func (em *EconManager) registerMetricEvents(e *twecon.Econ, metrics *EconMetrics) error {
	if e == nil || metrics == nil {
		return fmt.Errorf("nil econ or metrics")
	}

	for k, f := range(EconEventMetrics) {
		event := twecon.EconEvent{
			Name: k.Name,
			Regex: k.Regex,
			Func: func(econ *twecon.Econ, eventPayload string) any {
				f(metrics)

				return nil
			},
		}

		if err := e.EventManager.Register(&event); err != nil {
			return err
		}
	}

	return nil
}
