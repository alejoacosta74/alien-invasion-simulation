package aliens

import (
	"context"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/alejoacosta74/allien_invasion/internal"
	"github.com/alejoacosta74/allien_invasion/log"
	"github.com/alejoacosta74/allien_invasion/world"
)

func TestAliensLifeCycleBasic(t *testing.T) {
	var N = 4
	var buf = internal.NewSafeBuffer()
	log.GetLogger(log.WithWriter(buf), log.WithDebugLevel(true))

	w, err := world.NewWorldFromFile("../data/world_sample.txt")
	if err != nil {
		t.Fatalf("Error creating world: %s", err)
	}
	ctx, cancelFunc := context.WithCancel(context.Background())
	alienDoneChan := make(chan int, N)
	wg := sync.WaitGroup{}
	Start(ctx, N, alienDoneChan, w, nil, nil, &wg)
	t.Run("Should start N aliens", func(t *testing.T) {
		time.Sleep(time.Millisecond * 300)
		log := buf.String()
		if strings.Count(log, "starting") != N {
			t.Errorf("Expected %d aliens to start, but got %d", N, strings.Count(log, "starting"))
		}
	})
	buf.Reset()
	var aliens = N

	t.Run("Should stop N aliens on cancel", func(t *testing.T) {
		cancelFunc()
		for {
			select {
			case <-alienDoneChan:
				aliens--
				if aliens == 0 {
					close(alienDoneChan)
					return
				}
			case <-time.After(time.Millisecond * 500):
				t.Errorf("Expected %d aliens to be alive, but got %d", 0, aliens)
				return
			}
		}

	})
}
