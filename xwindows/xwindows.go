package xwindows

import (
	"fmt"
	"sync"
	"time"
)

const maxWindowDuration = 30 * time.Minute

// Debounce struct manages the debounce logic for a single document
type Debounce struct {
	windowDuration time.Duration
	timer          *time.Timer
	latestEvent    string
	windowID       string
	lock           sync.Mutex
	lastActiveTime time.Time
	pushFunc       func(string, string)
}

// NewDebounce creates a new Debounce
func NewDebounce(windowDuration time.Duration, pushFunc func(string, string), windowID string) *Debounce {
	return &Debounce{
		windowDuration: windowDuration,
		lastActiveTime: time.Now(),
		pushFunc:       pushFunc,
		windowID:       windowID,
	}
}

// pushEvent is the actual logic to push the event
func (d *Debounce) pushEvent() {
	d.lock.Lock()
	defer d.lock.Unlock()

	if d.latestEvent != "" {
		d.pushFunc(d.latestEvent, d.windowID)
		d.latestEvent = ""
	}

	// 清空定时器，允许下一个事件创建新窗口
	d.timer = nil
}

// addEvent adds a new event and resets the timer
func (d *Debounce) addEvent(event string) {
	d.lock.Lock()
	defer d.lock.Unlock()

	d.latestEvent = event
	d.lastActiveTime = time.Now()

	// 只有在没有定时器时才创建新的
	if d.timer == nil {
		d.timer = time.AfterFunc(d.windowDuration, d.pushEvent)
	}
	// 如果定时器已存在，让它继续运行至窗口期结束
}

// WindowsManager manages multiple documents and their debounceMap
type WindowsManager struct {
	debounceMap        sync.Map
	windowDuration     time.Duration
	cleanupInterval    time.Duration
	inactivityDuration time.Duration
	pushFunc           func(string, string)
	stopCh             chan struct{}
}

// NewWindowsManager creates a new windowsManager
func NewWindowsManager(windowDuration, cleanupInterval, inactivityDuration time.Duration, pushFunc func(string, string)) *WindowsManager {
	if windowDuration > maxWindowDuration {
		windowDuration = maxWindowDuration
	}
	dm := &WindowsManager{
		windowDuration:     windowDuration,
		cleanupInterval:    cleanupInterval,
		inactivityDuration: inactivityDuration,
		pushFunc:           pushFunc,
		stopCh:             make(chan struct{}),
	}
	go dm.startCleanupRoutine()
	return dm
}

// Close stops the cleanup routine and releases resources
func (dm *WindowsManager) Close() {
	close(dm.stopCh)
}

// RemoveWindow removes the document with the given ID
func (dm *WindowsManager) RemoveWindow(windowID string) {
	if value, ok := dm.debounceMap.Load(windowID); ok {
		debounce := value.(*Debounce)
		debounce.lock.Lock()
		if debounce.timer != nil {
			debounce.timer.Stop()
			debounce.timer = nil
		}
		debounce.lock.Unlock()
		dm.debounceMap.Delete(windowID)
	}
}

// AddEvent adds an event to the specified document
func (dm *WindowsManager) AddEvent(windowID string, event string) {
	if debounce, ok := dm.debounceMap.Load(windowID); ok {
		debounce.(*Debounce).addEvent(event)
	} else {
		nd := NewDebounce(dm.windowDuration, dm.pushFunc, windowID)
		nd.addEvent(event)
		dm.debounceMap.Store(windowID, nd)
	}
}

// startCleanupRoutine starts a goroutine that periodically cleans up inactive documents
func (dm *WindowsManager) startCleanupRoutine() {
	ticker := time.NewTicker(dm.cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			now := time.Now()
			dm.debounceMap.Range(func(key, value interface{}) bool {
				debounce := value.(*Debounce)
				debounce.lock.Lock()

				// 如果有定时器等待中，说明窗口仍在等待推送，不应该清理
				if debounce.timer != nil {
					debounce.lock.Unlock()
					return true
				}

				inactive := now.Sub(debounce.lastActiveTime) > dm.inactivityDuration
				debounce.lock.Unlock()

				if inactive {
					fmt.Printf("Removing inactive document: %s\n", key.(string))
					dm.RemoveWindow(key.(string))
				}
				return true
			})
		case <-dm.stopCh:
			return
		}
	}
}
