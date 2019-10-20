package main

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"strings"
)

const (
	NORTH           string = "north"
	SOUTH           string = "south"
	WEST            string = "west"
	EAST            string = "east"
	ITERATION_COUNT int    = 10_000
)

type WorldMap struct {
	alive_alien_count int
	cities            map[string]City
}

type Alien struct {
	number int
}

type City struct {
	north, south, west, east string
	alien                    Alien
}

func (p City) String() string {
	return fmt.Sprintf("north: %v, south: %v, west: %v, east: %v", p.north, p.south, p.west, p.east)
}

func main() {
	file, err := open_file()

	defer file.Close()

	if err != nil {
		fmt.Println("Could not open the file", err)
		return
	}

	world_map, err := parse_world_map(file)

	if err != nil {
		fmt.Println("Could not parse the world map", err)
		return
	}

	start_processing(&world_map)

	print_result(world_map)
}

func open_file() (*os.File, error) {
	args := os.Args

	return os.Open(args[1])
}

func parse_world_map(file *os.File) (WorldMap, error) {
	if len(os.Args) < 3 {
		return WorldMap{}, errors.New("Can't parse an alien number")
	}

	alien_arg := os.Args[2]

	var alien_number int64
	var err error

	if alien_arg != "" {
		alien_number, err = strconv.ParseInt(alien_arg, 10, 64)

		if err != nil {
			return WorldMap{}, err
		}

	}

	world_map, err := do_parse_world_map(file)

	if err != nil {
		return world_map, err
	}

	world_map.alive_alien_count = int(alien_number)

	return world_map, nil
}

func do_parse_world_map(file *os.File) (WorldMap, error) {
	scanner := bufio.NewScanner(file)
	cities := make(map[string]City)

	for scanner.Scan() {
		line := scanner.Text()

		name, city, error := parse_city(line)

		if error != nil {
			return WorldMap{}, error
		}

		cities[name] = city
	}

	return WorldMap{0, cities}, nil
}

func parse_city(line string) (string, City, error) {
	line_parts := strings.Split(line, " ")

	name := line_parts[0]

	direction_pairs := line_parts[1:]

	var city City

	for _, dir_pair := range direction_pairs {
		dir_pair_list := strings.Split(dir_pair, "=")

		if len(dir_pair_list) != 2 {
			return "", City{}, errors.New("Can't parse a city")
		}

		direction := dir_pair_list[0]
		city_name := dir_pair_list[1]

		switch direction {
		case NORTH:
			city.north = city_name
		case SOUTH:
			city.south = city_name
		case WEST:
			city.west = city_name
		case EAST:
			city.east = city_name
		default:
			return "", City{}, errors.New("Unknown direction")
		}
	}

	return name, city, nil
}

func start_processing(world_map *WorldMap) {
	initial_distribution(world_map, world_map.alive_alien_count)

	for i := 0; i < ITERATION_COUNT && world_map.alive_alien_count != 0; i++ {
		iteration_step(world_map)
	}
}

func initial_distribution(world_map *WorldMap, aliens_count int) {
	for i := 1; i <= aliens_count; i++ {
		new_alien := Alien{number: i}

		if len(world_map.cities) > 0 {
			random_city := random_key(world_map.cities).(string)

			visit_city(world_map, random_city, new_alien)
		}
	}
}

func iteration_step(world_map *WorldMap) {
	var touched_aliens []Alien

	for k, v := range world_map.cities {
		if v.alien != (Alien{}) && !contains(touched_aliens, v.alien) {
			available_cities := available_directions(v)

			if len(available_cities) > 0 {
				random_city_idx := rand.Int() % len(available_cities)
				random_city := available_cities[random_city_idx]

				visit_city(world_map, random_city, v.alien)

				touched_aliens = append(touched_aliens, v.alien)

				v.alien = Alien{}

				world_map.cities[k] = v
			}
		}
	}
}

func available_directions(city City) []string {
	var directions []string

	if city.north != "" {
		directions = append(directions, city.north)
	}

	if city.south != "" {
		directions = append(directions, city.south)
	}

	if city.west != "" {
		directions = append(directions, city.west)
	}

	if city.east != "" {
		directions = append(directions, city.east)
	}

	return directions

}

func visit_city(world_map *WorldMap, city_name string, alien Alien) {
	city := world_map.cities[city_name]

	if city.alien != (Alien{}) {
		fmt.Printf("%v has been destroyed by alien %v and alien %v\n", city_name, alien.number, city.alien.number)
		destroy_city(world_map, city_name)
	} else {
		city.alien = alien
		world_map.cities[city_name] = city
	}
}

func destroy_city(world_map *WorldMap, city_name string) {
	delete(world_map.cities, city_name)

	for k, v := range world_map.cities {
		switch city_name {
		case v.north:
			v.north = ""
		case v.south:
			v.south = ""
		case v.west:
			v.west = ""
		case v.east:
			v.east = ""
		}

		world_map.cities[k] = v
	}

	world_map.alive_alien_count -= 2
}

func print_result(world_map WorldMap) {
	for k, v := range world_map.cities {
		fmt.Printf("%v", k)

		if v.north != "" {
			fmt.Printf(" north=%v", v.north)
		}

		if v.south != "" {
			fmt.Printf(" south=%v", v.south)
		}

		if v.west != "" {
			fmt.Printf(" west=%v", v.west)
		}

		if v.east != "" {
			fmt.Printf(" east=%v", v.east)
		}

		fmt.Printf("\n")
	}
}

func random_key(mapI interface{}) interface{} {
	keys := reflect.ValueOf(mapI).MapKeys()

	return keys[rand.Intn(len(keys))].Interface()
}

func contains(aliens []Alien, alien Alien) bool {
	for _, a := range aliens {
		if a == alien {
			return true
		}
	}
	return false
}
