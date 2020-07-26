package land

import (
	"bufio"
	"fmt"
	"hash/adler32"
	"math"
	"os"
	"regexp"
	"strings"

	"github.com/gosimple/slug"
	log "github.com/sirupsen/logrus"
)

// Define the directions
const (
	North = int8(1)
	South = int8(-1)
	East  = int8(2)
	West  = int8(-2)
)

// map file line spec validation
var specValidator = regexp.MustCompile(`^([a-zA-Z][a-zA-Z1-9-]+)(\s+(north|west|south|east)=([a-zA-Z][a-zA-Z0-9-]*))*\s*$`)

// hash calculate the hash of a string
// we use adler since it is fast and produces short hashes
func hash(in string) uint32 {
	return adler32.Checksum([]byte(in))
}

// compute the hash of a route
// the hash is direction insensitive
func routeHash(from, to uint32, d int8) (h uint32, dir int8) {
	if from < to {
		h, dir = hash(fmt.Sprint(from, "-", to)), d
		return
	}
	h, dir = hash(fmt.Sprint(to, "-", from)), d*-1
	return

}

// Assumption
// 1. a city can only have 4 neighbors max (north, east, south, west)
// 2. if A is north of B then B is south of A
// 3. input string may or may not have bi-directional infos
// 4. if there are bi-directional info then they must be valid for 2.

// Land represent the map of a land with cities and routes between them
type Land struct {
	// keeps the association City Name -> Hash(city Name)
	cities map[uint32]string
	// keep the routes in the form cityFrom[direction] = cityTo
	// links are bi-directional
	routes map[uint32]map[int8]uint32
	// quick way to know if a route exists
	routeH map[uint32]int8
}

// NewLand creates a new land
func NewLand() *Land {
	return &Land{
		cities: make(map[uint32]string),
		routes: make(map[uint32]map[int8]uint32),
		routeH: make(map[uint32]int8),
	}
}

// Size returns the number of cities and number of routes
func Size(land *Land) (cities int, routes int) {
	cities = len(land.cities)
	routes = len(land.routeH)
	return
}

// AddCity add a city and retrieve the city ID
// if the city is already there then the existing uid
// will returned
func AddCity(l *Land, name string) (id uint32) {
	id = hash(slug.Make(name))
	if _, found := l.cities[id]; !found {
		l.cities[id] = name
	}
	return
}

// AddRoute adds a route from city A to city B and back
func AddRoute(l *Land, from, to uint32, direction int8) (err error) {
	if from == to {
		err = fmt.Errorf("From and to are the same (cannot add a route): %v", from)
		return
	}
	// check the direction
	if direction == 0 || math.Abs(float64(direction)) > 2 {
		err = fmt.Errorf("Unknown direction %v", direction)
		return
	}
	// test whenever a route already exists
	_h, _d := routeHash(from, to, direction)
	if d, ex := l.routeH[_h]; ex && d != _d {
		// route exists with different direction
		err = fmt.Errorf("route from %v to %v exists with a different direction: %v", from, to, d)
		return
	}

	// init a route, create the direction map if it does not exists
	ir := func(x uint32) (v map[int8]uint32) {
		v, ex := l.routes[x]
		if !ex {
			v = make(map[int8]uint32)
			l.routes[x] = v
		}
		return v
	}
	// get the "from city" routes
	fromRoutes := ir(from)
	// if there is already an existing target for that direction
	// check that it is the same of the one requested, otherwise return error
	if t, ex := fromRoutes[direction]; ex && t != to {
		err = fmt.Errorf("target direction %v to %v exists to a different location: %v", direction, to, from)
		return
	}
	// get the "to city" routes
	toRoutes := ir(to)
	// if there were no errors add both routes
	fromRoutes[direction] = to
	toRoutes[direction*-1] = from
	// save the edge index
	l.routeH[_h] = _d
	return
}

// DestroyCity remove a city and all the routes from/to it
func DestroyCity(land *Land, id uint32) {
	if routes, ex := land.routes[id]; ex {
		for d, t := range routes {
			delete(land.routes[t], d*-1)
			_h, _ := routeHash(id, t, d)
			delete(land.routeH, _h)
		}
		delete(land.routes, id)
	}
	delete(land.cities, id)
}

// Optimize remove nodes without edges
func Optimize(land *Land) {
	for id, routes := range land.routes {
		if len(routes) == 0 {
			delete(land.routes, id)
			delete(land.cities, id)
		}
	}
}

// LoadFromFile load a land definition from a file
func LoadFromFile(land *Land, mapFile string) (err error) {
	file, err := os.Open(mapFile)
	if err != nil {
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if e := parseMapLine(land, scanner.Text()); e != nil {
			log.Error(e)
		}
	}
	Optimize(land)
	err = scanner.Err()
	return
}

func parseMapLine(land *Land, specs string) (err error) {
	// validate the string according to the regexp (above)
	if !specValidator.Match([]byte(specs)) {
		err = fmt.Errorf("Invalid input string %s", specs)
		return
	}
	// get the city name
	pieces := strings.Split(specs, " ")
	srcName := pieces[0]
	src := AddCity(land, srcName)
	// parse neighbors
	for _, neighbor := range pieces[1:] {
		// shortest direction is east/west 4 characters
		x := strings.Index(neighbor, "=")
		if x < 0 {
			continue
		}
		targetName := neighbor[x+1:]
		target := AddCity(land, targetName)
		// parse the direction
		dir := North // default direction
		dirName := neighbor[0:x]
		switch strings.ToLower(dirName) {
		case "east":
			dir = East
		case "south":
			dir = South
		case "west":
			dir = West
		}
		if err := AddRoute(land, src, target, dir); err != nil {
			fmt.Println("SKIP invalid route from", srcName, "to", targetName, "via", dirName)
			log.Debug(err)
		}
	}
	return
}
