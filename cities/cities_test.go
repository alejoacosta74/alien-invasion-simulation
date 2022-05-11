package cities

import (
	"reflect"
	"testing"
	"testing/fstest"
)

func TestCityFileReader(t *testing.T) {

	fs := fstest.MapFS{
		"worldmap.txt": {Data: []byte(CitiesRawData)},
	}

	cityList, _ := LoadCitiesFS(fs)

	t.Run("Reads all the cities", func(t *testing.T) {
		if len(cityList) != 5 {
			t.Errorf("Expected 5 cities, got %d", len(cityList))
		}
	})

	t.Run("Reads all neighbours and directions", func(t *testing.T) {
		for i, tt := range TestCitylist {
			if !reflect.DeepEqual(tt, cityList[i]) {
				t.Errorf("Expected %v, got %v", tt, cityList[i])
			}

		}
	})
}

func TestLoadMap(t *testing.T) {
	t.Run("Should create a world structure from a file", func(t *testing.T) {

		want := map[string]*City{
			"Foo": {Name: "Foo",
				Neighbors: map[Directions]string{
					North: "Bar",
					South: "Qu-ux",
					West:  "Bazz",
				},
			},
			"Bar": {Name: "Bar",
				Neighbors: map[Directions]string{
					South: "Foo",
					West:  "Bee",
				},
			},
			"Qu-ux": {Name: "Qu-ux",
				Neighbors: map[Directions]string{
					North: "Foo",
				},
			},
			"Bazz": {Name: "Bazz",
				Neighbors: map[Directions]string{
					East: "Foo",
				},
			},
			"Bee": {Name: "Bee",
				Neighbors: map[Directions]string{
					East: "Bar",
				},
			},
		}

		got, err := LoadCityMapFromFile("../data/world_sample.txt")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(got) != len(want) {
			t.Errorf("Expected 5 cities, got %d", len(got))
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Expected %v, got %v", want, got)
		}
	})
}

const (
	CitiesRawData = `Foo north=Bar south=Qu-ux east=Buz west=Bazz
Bar north=Foo south=Qu-ux east=Buz west=Bazz
Qu-ux north=Bar south=Foo east=Buz west=Bazz
Buz north=Bar south=Foo east=Qu-ux west=Bazz
Bazz north=Bar south=Foo east=Qu-ux west=Buz
`
)

var TestCitylist = []*City{
	{
		Name: "Foo",
		Neighbors: map[Directions]string{
			North: "Bar",
			South: "Qu-ux",
			East:  "Buz",
			West:  "Bazz",
		},
	},
	{
		Name: "Bar",
		Neighbors: map[Directions]string{
			North: "Foo",
			South: "Qu-ux",
			East:  "Buz",
			West:  "Bazz",
		},
	},
	{
		Name: "Qu-ux",
		Neighbors: map[Directions]string{
			North: "Bar",
			South: "Foo",
			East:  "Buz",
			West:  "Bazz",
		},
	},
	{
		Name: "Buz",
		Neighbors: map[Directions]string{
			North: "Bar",
			South: "Foo",
			East:  "Qu-ux",
			West:  "Bazz",
		},
	},
	{
		Name: "Bazz",
		Neighbors: map[Directions]string{
			North: "Bar",
			South: "Foo",
			East:  "Qu-ux",
			West:  "Buz",
		},
	},
}
