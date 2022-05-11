package app_test

import (
	"context"
	"strings"
	"sync"
	"testing"

	"github.com/alejoacosta74/allien_invasion/aliens"
	"github.com/alejoacosta74/allien_invasion/app"
	"github.com/alejoacosta74/allien_invasion/internal"
	"github.com/alejoacosta74/allien_invasion/log"
	"github.com/alejoacosta74/allien_invasion/world"
)

func TestAlienInvasion(t *testing.T) {
	var buf = internal.NewSafeBuffer()
	logger, _ := log.GetLogger(log.WithWriter(buf), log.WithDebugLevel(true))

	t.Run("With odd number of aliens, at least 1 alien should stay alive and max iterations reached", func(t *testing.T) {
		N := 3
		I := 20
		ctx, cancelFunc := context.WithCancel(context.Background())
		alienDoneChan := make(chan int, N)
		startWanderChan := make(chan struct{}, 1)
		completedWanderChan := make(chan struct{}, 1)
		wg := sync.WaitGroup{}
		worldMap, err := world.NewWorldFromFile("../data/world_sample.txt")
		if err != nil {
			t.Fatalf("Error creating world: %s", err)
		}
		aliens.Start(ctx, N, alienDoneChan, worldMap, startWanderChan, completedWanderChan, &wg)
		done := app.HandleInvasionSignals(logger, cancelFunc, N, I, startWanderChan, completedWanderChan, alienDoneChan, &wg)
		<-done

		// fmt.Println("buffer:\n", buf.String())

		invasionLog := buf.String()
		if strings.Count(invasionLog, "Killing remaining 1 aliens") != 1 {
			t.Errorf("Expected 1 alien to be alive, but got %d", strings.Count(invasionLog, "Killing remaining 1 aliens"))
		}
	})

	buf.Reset()

	t.Run("With aliens (even) >> 5, all cities will likely be destroyed", func(t *testing.T) {
		N := 50
		I := 1000
		ctx, cancelFunc := context.WithCancel(context.Background())
		alienDoneChan := make(chan int, N)
		startWanderChan := make(chan struct{}, 1)
		completedWanderChan := make(chan struct{}, 1)
		wg := sync.WaitGroup{}
		worldMap, err := world.NewWorldFromFile("../data/world_sample.txt")
		if err != nil {
			t.Fatalf("Error creating world: %s", err)
		}
		aliens.Start(ctx, N, alienDoneChan, worldMap, startWanderChan, completedWanderChan, &wg)
		done := app.HandleInvasionSignals(logger, cancelFunc, N, I, startWanderChan, completedWanderChan, alienDoneChan, &wg)
		<-done

		// fmt.Println("buffer:\n", buf.String())
		invasionLog := buf.String()
		if strings.Count(invasionLog, "no more cities available. Exiting") == 0 {
			t.Errorf("Expected 0 cities to remain undestroyed")
		}
	})
}
