package xwindows

import (
	"fmt"
	"testing"
	"time"
)

func TestNewWindowsManager(t *testing.T) {
	type args struct {
		windowDuration     time.Duration
		cleanupInterval    time.Duration
		inactivityDuration time.Duration
	}
	tests := []struct {
		name string
		args args
		want *WindowsManager
	}{
		// TODO: Add test cases.
		{
			name: "test-1",
			args: args{
				windowDuration:     10 * time.Second,
				cleanupInterval:    3 * time.Second,
				inactivityDuration: 1 * time.Minute,
			},
			want: nil,
		},
	}
	// pushFunc1 := func(event string) {
	// 	fmt.Println("Pushing 1 event:", event)
	// 	fmt.Println(time.Now().Unix())
	// }
	// pushFunc2 := func(event string) {
	// 	fmt.Println("Pushing 2 event:", event)
	// 	fmt.Println(time.Now().Unix())
	// }
	// pushFunc3 := func(event string) {
	// 	fmt.Println("Pushing 3 event:", event)
	// 	fmt.Println(time.Now().Unix())
	// }
	// pushFunc4 := func(event string) {
	// 	fmt.Println("Pushing 4 event:", event)
	// 	fmt.Println(time.Now().Unix())
	// }

	pushFuncAll := func(windowID string, event Event) {
		fmt.Println("Pushing all event:", event, windowID, time.Now().Format("2006-01-02 15:04:05"))
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			docManager := NewWindowsManager(tt.args.windowDuration, tt.args.cleanupInterval, tt.args.inactivityDuration, pushFuncAll)
			// Simulating events being added to different documents
			go func() {
				docManager.AddEvent("doc1", Event{
					Headers: nil,
					Value:   fmt.Sprintf("Edit 1 for Document 1, %d", time.Now().Unix()),
				})
				time.Sleep(8 * time.Second)
				docManager.AddEvent("doc1", Event{
					Headers: nil,
					Value:   fmt.Sprintf("Edit 2 for Document 1, %d", time.Now().Unix()),
				})
				time.Sleep(8 * time.Second)
				docManager.AddEvent("doc1", Event{
					Headers: nil,
					Value:   fmt.Sprintf("Edit 3 for Document 1, %d", time.Now().Unix()),
				})
			}()

			go func() {
				docManager.AddEvent("doc2", Event{
					Headers: nil,
					Value:   fmt.Sprintf("Edit 1 for Document 2, %d", time.Now().Unix()),
				})
				time.Sleep(11 * time.Second)
				docManager.AddEvent("doc2", Event{
					Headers: nil,
					Value:   fmt.Sprintf("Edit 2 for Document 2, %d", time.Now().Unix()),
				})
			}()

			go func() {
				docManager.AddEvent("doc3", Event{
					Headers: nil,
					Value:   fmt.Sprintf("Edit 1 for Document 3, %d", time.Now().Unix()),
				})
				time.Sleep(6 * time.Second)
				docManager.AddEvent("doc3", Event{
					Headers: nil,
					Value:   fmt.Sprintf("Edit 2 for Document 3, %d", time.Now().Unix()),
				})
				docManager.AddEvent("doc3", Event{
					Headers: nil,
					Value:   fmt.Sprintf("Edit 3 for Document 3, %d", time.Now().Unix()),
				})
				docManager.AddEvent("doc3", Event{
					Headers: nil,
					Value:   fmt.Sprintf("Edit 4 for Document 3, %d", time.Now().Unix()),
				})
				docManager.AddEvent("doc3", Event{
					Headers: nil,
					Value:   fmt.Sprintf("Edit 5 for Document 3, %d", time.Now().Unix()),
				})
				docManager.AddEvent("doc3", Event{
					Headers: nil,
					Value:   fmt.Sprintf("Edit 6 for Document 3, %d", time.Now().Unix()),
				})
				time.Sleep(6 * time.Second)
				docManager.AddEvent("doc3", Event{
					Headers: nil,
					Value:   fmt.Sprintf("Edit 7 for Document 3, %d", time.Now().Unix()),
				})
			}()

			go func() {
				docManager.AddEvent("doc4", Event{
					Headers: nil,
					Value:   fmt.Sprintf("Edit 1 for Document 4, %d", time.Now().Unix()),
				})
				time.Sleep(15 * time.Second)
				docManager.AddEvent("doc4", Event{
					Headers: nil,
					Value:   fmt.Sprintf("Edit 2 for Document 4, %d", time.Now().Unix()),
				})
			}()

			// Allow some time for all events to be processed
			time.Sleep(2 * time.Minute)

		})
	}
}
