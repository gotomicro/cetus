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
				windowDuration:     30 * time.Second,
				cleanupInterval:    10 * time.Second,
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

	pushFuncAll := func(event string) {
		fmt.Println("Pushing all event:", event)
		fmt.Println(time.Now().Unix())
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			docManager := NewWindowsManager(tt.args.windowDuration, tt.args.cleanupInterval, tt.args.inactivityDuration, pushFuncAll)
			// Simulating events being added to different documents
			go func() {
				docManager.AddEvent("doc1", fmt.Sprintf("Edit 1 for Document 1, %d", time.Now().Unix()))
				time.Sleep(10 * time.Second)
				docManager.AddEvent("doc1", fmt.Sprintf("Edit 2 for Document 1, %d", time.Now().Unix()))

			}()

			go func() {
				docManager.AddEvent("doc2", fmt.Sprintf("Edit 1 for Document 2, %d", time.Now().Unix()))
				time.Sleep(20 * time.Second)
				docManager.AddEvent("doc2", fmt.Sprintf("Edit 2 for Document 2, %d", time.Now().Unix()))
			}()

			go func() {
				docManager.AddEvent("doc3", fmt.Sprintf("Edit 1 for Document 3, %d", time.Now().Unix()))
				time.Sleep(30 * time.Second)
				docManager.AddEvent("doc3", fmt.Sprintf("Edit 2 for Document 3, %d", time.Now().Unix()))
			}()

			go func() {
				docManager.AddEvent("doc4", fmt.Sprintf("Edit 1 for Document 4, %d", time.Now().Unix()))
				time.Sleep(60 * time.Second)
				docManager.AddEvent("doc4", fmt.Sprintf("Edit 2 for Document 4, %d", time.Now().Unix()))
			}()

			// Allow some time for all events to be processed
			time.Sleep(2 * time.Minute)

		})
	}
}
