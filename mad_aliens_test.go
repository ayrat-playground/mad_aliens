package main

import (
	"os"
	"testing"
)

func TestParsesCityFromValidString(t *testing.T) {
	input := "Foo north=Bar west=Baz south=Qu-ux"

	name, city, error := parse_city(input)

	if name != "Foo" || city.north != "Bar" || city.west != "Baz" || city.south != "Qu-ux" || error != nil {
		t.Errorf("City was parsed wrong")
	}
}

func TestParsesWorldMapFromValidFile(t *testing.T) {
	input, _ := os.Open("./map")

	defer input.Close()

	world_map, _ := do_parse_world_map(input)

	if len(world_map.cities) != 2 {
		t.Errorf("Cities number should equal 2")
	}

	first_city := world_map.cities["Foo"]

	if first_city.north != "Bar" || first_city.west != "Baz" || first_city.south != "Qu-ux" {
		t.Errorf("City was parsed wrong")
	}

	second_city := world_map.cities["Bar"]

	if second_city.west != "Bee" || second_city.south != "Foo" {
		t.Errorf("City was parsed wrong")
	}

}

func TestDestroysCity(t *testing.T) {
	world_map := WorldMap{
		cities: map[string]City{
			"Foo": City{north: "Bar", west: "Bee"},
			"Bar": City{south: "Foo"},
		},
		alive_alien_count: 2,
	}

	destroy_city(&world_map, "Bar")

	if len(world_map.cities) != 1 {
		t.Errorf("Cities number should equal 2")
	}

	if world_map.alive_alien_count != 0 {
		t.Errorf("alive_alien_count should equal 0")
	}

	city := world_map.cities["Foo"]

	if city.north == "Bar" || city.west != "Bee" {
		t.Errorf("City was cleaned wrong")
	}
}

func TestAlienVisitsCity(t *testing.T) {
	world_map := WorldMap{
		cities: map[string]City{
			"Foo": City{north: "Bar", west: "Bee"},
			"Bar": City{south: "Foo"},
		},
		alive_alien_count: 1,
	}

	alien := Alien{number: 5}

	visit_city(&world_map, "Bar", alien)

	city := world_map.cities["Bar"]

	if city.south != "Foo" || city.alien != alien {
		t.Errorf("City was visited wrong")
	}
}

func TestAlienVisitsCityDestroyingIt(t *testing.T) {
	world_map := WorldMap{
		cities: map[string]City{
			"Foo": City{north: "Bar", west: "Bee"},
			"Bar": City{south: "Foo", alien: Alien{number: 5}},
		},
		alive_alien_count: 2,
	}

	alien := Alien{number: 7}

	visit_city(&world_map, "Bar", alien)

	if len(world_map.cities) != 1 {
		t.Errorf("Cities number should equal 2")
	}

	if world_map.alive_alien_count != 0 {
		t.Errorf("alive_alien_count should equal 0")
	}

	city := world_map.cities["Foo"]

	if city.north == "Bar" || city.west != "Bee" {
		t.Errorf("City was cleaned wrong")
	}
}

func TestInitialDistributionOfSingleAlien(t *testing.T) {
	world_map := WorldMap{
		cities: map[string]City{
			"Foo": City{north: "Bar", west: "Bee"},
			"Bar": City{south: "Foo"},
		},
		alive_alien_count: 1,
	}

	initial_distribution(&world_map, 1)

	foo := world_map.cities["Foo"]
	bar := world_map.cities["Bar"]

	if (foo.alien == Alien{} && bar.alien == Alien{}) || (foo.alien != Alien{} && bar.alien != Alien{}) {
		t.Errorf("Initial distribution is wrong")
	}
}

func TestInitialDistributionOfThreeAliens(t *testing.T) {
	world_map := WorldMap{
		cities: map[string]City{
			"Foo": City{north: "Bar", west: "Bee"},
			"Bar": City{south: "Foo"},
		},
		alive_alien_count: 3,
	}

	initial_distribution(&world_map, 3)

	if world_map.alive_alien_count != 1 {
		t.Errorf("alive_alien_count should equal 1")
	}

	if len(world_map.cities) != 1 {
		t.Errorf("Cities number should equal 1")
	}
}

func TestAvailableDirections(t *testing.T) {
	city := City{north: "Bar", west: "Bee"}

	directions := available_directions(city)

	if len(directions) != 2 || directions[0] != "Bar" || directions[1] != "Bee" {
		t.Errorf("directions are wrong")
	}
}

func TestOneIterationStep(t *testing.T) {
	world_map := WorldMap{
		cities: map[string]City{
			"Foo": City{north: "Bar", alien: Alien{number: 5}},
			"Bar": City{south: "Foo"},
		},
		alive_alien_count: 1,
	}

	iteration_step(&world_map)

	foo := world_map.cities["Foo"]
	bar := world_map.cities["Bar"]

	if foo.alien != (Alien{}) && world_map.alive_alien_count != 1 {
		t.Errorf("Alient wasn't moved")
	}

	if bar.alien != (Alien{number: 5}) {
		t.Errorf("Alien wasn't moved")
	}
}

func TestStartProcessing(t *testing.T) {
	world_map := WorldMap{
		cities: map[string]City{
			"Foo": City{north: "Bar"},
			"Bar": City{south: "Foo"},
		},
		alive_alien_count: 4,
	}

	start_processing(&world_map)

	if len(world_map.cities) != 0 || world_map.alive_alien_count != 0 {
		t.Errorf("Cities weren't destroyed")
	}
}

func TestStartProcessingLeavesOneCity(t *testing.T) {
	world_map := WorldMap{
		cities: map[string]City{
			"Foo": City{north: "Bar", west: "Bee"},
			"Bar": City{south: "Foo"},
			"Bee": City{east: "Foo", west: "Bar"},
		},
		alive_alien_count: 5,
	}

	start_processing(&world_map)

	if len(world_map.cities) != 1 || world_map.alive_alien_count != 1 {
		t.Errorf("Cities weren't destroyed")
	}
}
