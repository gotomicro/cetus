package xwindows

import (
	"fmt"
	"sync"
	"time"
)

// Debounce struct manages the debounce logic for a single document
type Debounce struct {
	windowDuration time.Duration
	timer          *time.Timer
	latestEvent    string
	lock           sync.Mutex
	lastActiveTime time.Time
	pushFunc       func(string)
}

// NewDebounce creates a new Debounce
func NewDebounce(windowDuration time.Duration, pushFunc func(string)) *Debounce {
	return &Debounce{
		windowDuration: windowDuration,
		lastActiveTime: time.Now(),
		pushFunc:       pushFunc,
	}
}

// pushEvent is the actual logic to push the event
func (d *Debounce) pushEvent() {
	d.lock.Lock()
	defer d.lock.Unlock()

	if d.latestEvent != "" {
		d.pushFunc(d.latestEvent)
		d.latestEvent = ""
	}
}

// addEvent adds a new event and resets the timer
func (d *Debounce) addEvent(event string) {
	d.lock.Lock()
	defer d.lock.Unlock()

	d.latestEvent = event
	d.lastActiveTime = time.Now()

	if d.timer != nil {
		d.timer.Stop()
	}

	d.timer = time.AfterFunc(d.windowDuration, d.pushEvent)
}

// WindowsManager manages multiple documents and their debounceMap
type WindowsManager struct {
	debounceMap        sync.Map
	windowDuration     time.Duration
	cleanupInterval    time.Duration
	inactivityDuration time.Duration
	pushFunc           func(string)
}

// NewWindowsManager creates a new windowsManager
func NewWindowsManager(windowDuration, cleanupInterval, inactivityDuration time.Duration) *WindowsManager {
	dm := &WindowsManager{
		windowDuration:     windowDuration,
		cleanupInterval:    cleanupInterval,
		inactivityDuration: inactivityDuration,
	}
	go dm.startCleanupRoutine()
	return dm
}

// AddWindow adds a new document with the given ID
func (dm *WindowsManager) AddWindow(docID string, pushFunc func(string)) {
	dm.debounceMap.Store(docID, NewDebounce(dm.windowDuration, pushFunc))
}

// RemoveWindow removes the document with the given ID
func (dm *WindowsManager) RemoveWindow(docID string) {
	dm.debounceMap.Delete(docID)
}

// AddEvent adds an event to the specified document
func (dm *WindowsManager) AddEvent(windowID string, event string) {
	if debounce, ok := dm.debounceMap.Load(windowID); ok {
		debounce.(*Debounce).addEvent(event)
	}
}

// startCleanupRoutine starts a goroutine that periodically cleans up inactive documents
func (dm *WindowsManager) startCleanupRoutine() {
	for {
		time.Sleep(dm.cleanupInterval)
		now := time.Now()
		dm.debounceMap.Range(func(key, value interface{}) bool {
			debounce := value.(*Debounce)
			debounce.lock.Lock()
			inactive := now.Sub(debounce.lastActiveTime) > dm.inactivityDuration
			debounce.lock.Unlock()
			if inactive {
				fmt.Printf("Removing inactive document: %s\n", key.(string))
				dm.RemoveWindow(key.(string))
			}
			return true
		})
	}
}
