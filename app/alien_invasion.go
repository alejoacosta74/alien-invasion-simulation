package app

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/alejoacosta74/allien_invasion/aliens"
	"github.com/alejoacosta74/allien_invasion/log"
	"github.com/alejoacosta74/allien_invasion/world"
	"github.com/schollz/progressbar/v3"
	"github.com/sirupsen/logrus"
)

// Starts aliens and coordinates their movement until the aliens are destroyed or the iterations are completed
func StartInvation(ctx context.Context, cancelFunc context.CancelFunc, N int, I int, sigs chan os.Signal, worldMap *world.World) chan struct{} {

	// get logger returns a singleton logger
	logger, _ := log.GetLogger()

	// channel to receive communication when aliens are done
	alienDoneChan := make(chan int, N)

	// channel to SEND 'start wandering' signal to alien
	startWanderChan := make(chan struct{}, 1)

	// channel to RECEIVE 'wandering completed' signal from alien
	completedWanderChan := make(chan struct{}, N)

	// Waitgroup to wait for all aliens to complete wandering
	var wg sync.WaitGroup

	// create and start N aliens
	aliens.Start(ctx, N, alienDoneChan, worldMap, startWanderChan, completedWanderChan, &wg)

	// start go rountine to dispatch and handle signals
	done := HandleInvasionSignals(logger, cancelFunc, N, I, startWanderChan, completedWanderChan, alienDoneChan, &wg)

	return done
}

func HandleInvasionSignals(logger *logrus.Logger, cancelFunc context.CancelFunc, N, I int, startWanderChan, completedWanderChan chan struct{}, alienDoneChan chan int, wg *sync.WaitGroup) chan struct{} {
	done := make(chan struct{})
	go func() {

		// put the aliens to start wandering and interact with each other
		startWanderChan <- struct{}{}
		logger.Infof("Alien invasion started with %d aliens and %d iterations", N, I)

		// used for displaying progress when only 1 alien left
		progBar := getBar(I)

		// keep progress of alien moves/iteractions
		moves := 1

		// keep track of alien that is exiting
		var alienId int
		// total iterations = N * I (iterations per alien)
		iterations := N * I
		for moves < iterations && N > 0 {
			select {
			// handle done signal from aliens
			case alienId = <-alienDoneChan:
				N--
				iterations = updateIterations(I, N, moves)
				logger.Debugf("Alien id %d destroyed,  remaining aliens: N = %d, moves: %d, updated iterations: %d", alienId, N, moves, iterations)
				// receives wandered completed signal from alien and sends the
				// order to start wandering to next alien
			case <-completedWanderChan:
				startWanderChan <- struct{}{}
				moves++
				logger.Debugf("Sent move order to alien (N = %d). Total moves: %d", N, moves)
			}
			if N == 1 && logger.GetLevel() != logrus.DebugLevel {
				progBar.Add(1)
				time.Sleep(5 * time.Millisecond)
			}
		}
		println()
		// if all iterations are completed, destroy the remaining aliens, if any
		if N > 0 {
			cancelFunc()
			logger.Infof("Killing remaining %d aliens after %d iterations", N, I)
			for N > 0 {
				alienId = <-alienDoneChan
				logger.Debugf("Alien id %d stopped\n", alienId)
				N--
				if N == 0 {
					break
				}
				logger.Debug("Waiting for remaining aliens to stop...")
				// time.Sleep(1 * time.Second)
			}
			logger.Debug("All aliens stopped")
			wg.Wait()
		}
		close(startWanderChan)
		close(alienDoneChan)
		done <- struct{}{}
	}()
	return done

}

// creates a progress bar used to display progress when only 1 alien is left
func getBar(I int) *progressbar.ProgressBar {
	bar := progressbar.NewOptions(I,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetWidth(15),
		progressbar.OptionSetDescription("[cyan][reset] 1 alien left. Waiting for it to exit or iterations to complete..."),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))
	return bar
}

func updateIterations(I, N, moves int) int {
	previousIterationsPerAlien := int(moves / (N + 1))
	remainingIterationsPerAlien := I - previousIterationsPerAlien
	return remainingIterationsPerAlien * N
}
