package aliens

import (
	"context"
	"fmt"
	"math/rand"
	"reflect"
	"sync"
	"time"

	"github.com/alejoacosta74/allien_invasion/cities"
	"github.com/alejoacosta74/allien_invasion/log"
	"github.com/alejoacosta74/allien_invasion/world"
	"github.com/sirupsen/logrus"
)

type Alien struct {
	Name                string
	id                  int
	wg                  *sync.WaitGroup
	logger              *logrus.Entry
	worldMap            world.IWorld
	city                *cities.City
	alienDoneChan       chan int      //
	startWanderChan     chan struct{} // channel to receive move signal
	completedWanderChan chan struct{} // channel to send signal when alien has completed wandering
	quitChan            chan struct{} // channel to receive quit signal
}

// starts n gorountines (aliens) to wander the world
func Start(ctx context.Context, n int, alienDoneChan chan int, worldMap world.IWorld, startWanderChan, completedWanderChan chan struct{}, wg *sync.WaitGroup) {
	for i := 1; i <= n; i++ {
		wg.Add(1)
		go func(i int) {
			alien := NewAlien("Alien"+fmt.Sprintf("%d", i), i, alienDoneChan, worldMap, startWanderChan, completedWanderChan, wg)
			alien.Run(ctx)
		}(i)
	}
}

// Creates a new alien
func NewAlien(name string, id int, alienDoneChan chan int, worldMap world.IWorld, startWanderChan, completedWanderChan chan struct{}, wg *sync.WaitGroup) *Alien {
	alienLogger, _ := log.GetLogger()
	alien := &Alien{
		Name:          "Alien",
		id:            id,
		wg:            wg,
		alienDoneChan: alienDoneChan,
		logger: alienLogger.WithFields(logrus.Fields{
			"alien": name,
			"id":    id,
		}),
		worldMap:            worldMap,
		startWanderChan:     startWanderChan,
		completedWanderChan: completedWanderChan,
		quitChan:            make(chan struct{}, 1),
	}
	return alien
}

// Run the alien, by waiting for it turn to move and wander, and also waits for cancelation signal or exit signal from other aliens
func (a *Alien) Run(ctx context.Context) {
	a.logger.Debug("...starting")
	for {
		select {
		case <-ctx.Done():
			a.alienDoneChan <- a.id
			a.wg.Done()
			a.logger.Debugf("received cancel signal...exiting")
			return
		case <-a.startWanderChan:
			a.logger.Debugf("received order to wander")
			destroyer, killed := a.move(ctx)
			a.completedWanderChan <- struct{}{}
			if killed || destroyer {
				a.alienDoneChan <- a.id
				a.wg.Done()
				switch {
				case killed:
					a.logger.Debugf("...got killed. Exiting")
				case destroyer:
					a.logger.Debugf("... just destroyed a city and alien. Exiting")
				}
				return
			}
			continue
		case <-a.quitChan:
			a.alienDoneChan <- a.id
			a.wg.Done()
			a.logger.Debugf("... just received quit (killed) signal. Exiting")
			return
		default:
			continue
		}
	}
}

// move() requests a city to the world map for the alien
// if the city is not visited, alien will move into it and wait for its next turn to move again
// if the city is visited, alien will destroy the city, the alien visitor and will also destroy itself (exit)
func (a *Alien) move(ctx context.Context) (destroyer bool, killed bool) {
	var newCity *cities.City
	var neighborName string
	// If alien has no ciy attached, it needs to land on a city
	if a.city == nil {
		newCity = a.worldMap.GetRandomCity()
		//if no more cities available, exit
		if newCity == nil {
			a.logger.Info("no more cities available. Exiting")
			destroyer = false
			killed = true
			return
		}
		neighborName = newCity.Name

	} else {

		// get current city
		city := a.city
		// get neighbor cities
		n := city.NumNeighbors()
		if n == 0 {
			a.logger.Debug("locked. No neighbors cities available")
			return false, false
		}

		// get a random direction for a neighbor city
		getDirection := func() cities.Directions {
			rand.Seed(time.Now().UnixNano())
			keys := reflect.ValueOf(city.GetNeighborsMap()).MapKeys()
			d := keys[rand.Intn(len(keys))].Interface().(cities.Directions)
			return d
		}
		d := getDirection()
		neighborName = city.GetNeighbor(d)
		newCity = a.worldMap.GetCityByName(neighborName)
		a.logger.Debugf("going to move to %s", neighborName)
	}
	// if alien is currently in a city, it needs to leave it first before moving to another city
	if a.city != nil {
		a.leaveCity()
	}
	a.logger.WithField("toCity", neighborName).Debug("Alien moving...")
	if !newCity.IsVisited() {
		// if the candidate city is available, move into it and
		// return to main loop and wait for next turn to move
		a.moveTo(newCity)
		destroyer = false
		killed = false
		return
	} else {
		// if the candidate city is already occupied, destroy it
		a.destroyCity(newCity)
		destroyer = true
		killed = false
		return
	}
}

func (a *Alien) moveTo(c *cities.City) {
	c.Visit(a.id, a.quitChan)
	a.city = c
}

func (a *Alien) destroyCity(c *cities.City) {
	visitorId := c.GetVisitor()
	visitorCh := c.GetVisitorQuitChannel()
	a.worldMap.DestroyCity(c)
	a.logger.Warnf("City %s destroyed by Aliens %d and %d", c.Name, a.id, visitorId)
	visitorCh <- struct{}{}
	a.logger.Debugf("sent quit signal to visitor %d", visitorId)
}

func (a *Alien) leaveCity() {
	a.city.Leave(a.id, a.quitChan)
	a.city = nil
}
