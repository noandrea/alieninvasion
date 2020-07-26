package land

import (
	"fmt"
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestAddCity(t *testing.T) {

	l := NewLand()

	type args struct {
		l    *Land
		name string
	}
	tests := []struct {
		name   string
		args   args
		wantID uint32
	}{
		{"1", args{l, "berlin"}, 144114301},
		{"1", args{l, "BerLIN"}, 144114301},
		{"2", args{l, "london"}, 149488267},
		{"2", args{l, " London "}, 149488267},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotID := AddCity(tt.args.l, tt.args.name); gotID != tt.wantID {
				t.Errorf("AddCity() = %v, want %v", gotID, tt.wantID)
			}
		})
	}
	//
	citiesExpected, routesExpected := Size(l)
	citiesActual, routesActual := 2, 0

	if citiesExpected != citiesActual {
		t.Errorf("AddCity() = %v, want %v", citiesActual, citiesExpected)
	}

	if routesExpected != routesActual {
		t.Errorf("AddCity() = %v, want %v", routesActual, routesExpected)
	}

}

func TestAddRoute(t *testing.T) {

	l := NewLand()

	type args struct {
		l         *Land
		from      uint32
		to        uint32
		direction int8
	}
	tests := []struct {
		args    args
		wantErr bool
	}{
		{args{l, 1, 2, North}, false}, // 1n2 ; 2s1
		{args{l, 1, 3, North}, true},  // duplicated route
		{args{l, 1, 4, East}, false},  // 1e4; 4w1
		{args{l, 2, 3, East}, false},  // 2e3; 3w2
		{args{l, 3, 2, West}, false},  // 2e3; 3w2 ** back route
		{args{l, 3, 2, North}, true},  // 3w2 is the only possibility
		{args{l, 3, 2, South}, true},  // 3w2 is the only possibility
		{args{l, 3, 2, East}, true},   // 3w2 is the only possibility
		{args{l, 1, 1, West}, true},   // route to itself
		{args{l, 4, 6, -3}, true},     // invalid direction
		{args{l, 4, 7, 0}, true},      // invalid direction
		{args{l, 4, 8, 3}, true},      // invalid direction
	}
	expectedRoutes := 3 // 3 valid + 1 back route

	for x, tt := range tests {
		t.Run(fmt.Sprint("test_", tt.args), func(t *testing.T) {
			if err := AddRoute(tt.args.l, tt.args.from, tt.args.to, tt.args.direction); (err != nil) != tt.wantErr {
				t.Errorf("%v AddRoute() error = %v, wantErr %v", x, err, tt.wantErr)
			}
		})
	}

	for s, rs := range l.routes {
		for d, e := range rs {
			log.Printf("%2v >> %2v > %2v", s, d, e)
		}
	}

	_, r := Size(l)
	if r != expectedRoutes {
		t.Errorf("count AddRoute() expected = %v, got %v", r, expectedRoutes)
	}
}

func TestDestroyCity(t *testing.T) {

	l := NewLand()

	c0 := AddCity(l, "city_0") // 143524425
	c1 := AddCity(l, "city_1") // 143589962
	c2 := AddCity(l, "city_2") // 143655499
	c3 := AddCity(l, "city_3") // 143721036
	c4 := AddCity(l, "city_4") // 143786573
	c5 := AddCity(l, "city_5") // 143852110
	c6 := AddCity(l, "city_6") // 143917647
	c7 := AddCity(l, "city_7") // 143983184
	c8 := AddCity(l, "city_8") // 144048721
	c9 := AddCity(l, "city_9") // 144114258
	// the routes are arranged like in a numpad
	AddRoute(l, c1, c2, East)
	AddRoute(l, c1, c4, South)
	AddRoute(l, c2, c3, East)
	AddRoute(l, c2, c5, South)
	AddRoute(l, c3, c6, South)
	AddRoute(l, c4, c5, East)
	AddRoute(l, c4, c7, South)
	AddRoute(l, c5, c6, East)
	AddRoute(l, c5, c8, South)
	AddRoute(l, c6, c9, South)
	AddRoute(l, c7, c8, East)
	AddRoute(l, c8, c9, East)
	AddRoute(l, c8, c0, South)

	type args struct {
		land  *Land
		id    uint32
		nLeft int // number of nodes left
		eLeft int // number of edges left
	}
	tests := []struct {
		name string
		args args
	}{
		{"c0", args{l, c0, 9, 12}},
		{"c5", args{l, c5, 8, 8}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			DestroyCity(tt.args.land, tt.args.id)
			n, e := Size(l)
			if n != tt.args.nLeft {
				t.Errorf("DestroyCity() = nodes(%v), want %v", n, tt.args.nLeft)
			}
			if e != tt.args.eLeft {
				t.Errorf("DestroyCity() = edges(%v), want %v", e, tt.args.eLeft)
			}
		})
	}
}
