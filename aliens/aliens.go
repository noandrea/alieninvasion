package aliens

import (
	"fmt"
	"sync"

	"github.com/noandrea/alieninvasion/land"
	log "github.com/sirupsen/logrus"
)

const (
	ghost = -1
)

// Invasion simulates the invasion
type Invasion struct {
	aliens       int
	alienBase    map[int]uint32
	alienBaseM   sync.RWMutex
	maxRounds    int
	currentRound int
}

// NewInvasion prepare new invasion
func NewInvasion(l *land.Land, aliens int, maxRounds int) *Invasion {
	invasion := &Invasion{
		aliens:    aliens,
		maxRounds: maxRounds,
		alienBase: make(map[int]uint32),
	}
	// populate the alien bases
	for i := 0; i < invasion.aliens; i++ {
		invasion.alienBase[i] = land.Space
	}
	log.Debug("Added ", len(invasion.alienBase), " aliens ready to roll")
	return invasion
}

func killAliens(i *Invasion, aliens ...int) {
	i.alienBaseM.Lock()
	defer i.alienBaseM.Unlock()
	for _, aID := range aliens {
		delete(i.alienBase, aID)
	}
}

// Run an iteration of the invasion
func Run(world *land.Land, invasion *Invasion) (r int, err error) {
	// begin a new round
	r = invasion.currentRound
	log.Debug("> > Begin round ", r)
	if invasion.currentRound >= invasion.maxRounds {
		err = fmt.Errorf("Time is over")
	}
	// check if there is still place
	if n, _ := land.Size(world); n == 0 {
		err = fmt.Errorf("The invasion is over, the world has been destroyed")
		log.Debug("all cities have been destroyed")
		return
	}

	// check how many aliens are left
	if len(invasion.alienBase) <= 1 {
		err = fmt.Errorf("The invasion is over")
		log.Debug("all aliens are dead")
		return
	}
	// save movements for current round
	cityOccupant := make(map[uint32]int)
	for aID, currentCityID := range invasion.alienBase {
		// move the alien to a new place
		targetCityID, err := land.MoveFrom(world, currentCityID)
		log.Debug("Alien ", aID, " moving from ", land.GetCityName(world, currentCityID), " - ", currentCityID, " to ", land.GetCityName(world, targetCityID), " - ", targetCityID)
		if err != nil {
			fmt.Println("There is nowhere to go!")
			log.Debug(err)
			killAliens(invasion, aID)
			continue
		}
		// here is the conflict and mutual destruction
		if occupantID, es := cityOccupant[targetCityID]; es {
			// if the city was already visited
			if occupantID >= 0 {
				// if there is a live occupant then destroy the city
				// will be better to push the message in a channel
				// instead of printing directly
				// TODO: return the message
				fmt.Printf("%v has​ been​ destroyed​ by​ alien​ %v and​ alien​ %v!\n", land.GetCityName(world, targetCityID), aID, occupantID)
				land.DestroyCity(world, targetCityID)
				killAliens(invasion, aID, occupantID)
				cityOccupant[targetCityID] = ghost // set the city as destroyed
			} else {
				// if there is no one then the alien will die anyway (will be useful if using concurency)
				fmt.Printf("%v got caught in the distruction of the city and died as well\n", occupantID)
				killAliens(invasion, occupantID)
			}
			continue
		}
		// here is a successful occupation
		cityOccupant[targetCityID] = aID
		// update alien location
		invasion.alienBase[aID] = targetCityID
	}
	invasion.currentRound++
	log.Debug("Alien population: ", len(invasion.alienBase))
	log.Debug("> > End round ", r)
	return
}
