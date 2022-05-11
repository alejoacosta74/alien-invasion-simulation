package cities

import (
	"fmt"
	"strings"
)

func parseCityToMap(cityMap map[string]*City, line string) error {
	parts := strings.Split(line, " ")
	if len(parts) == 0 {
		return fmt.Errorf("no cities found to parse")
	}
	if _, ok := cityMap[parts[0]]; !ok {
		cityMap[parts[0]] = &City{
			Name:      parts[0],
			Neighbors: make(map[Directions]string),
		}
	}
	for _, part := range parts[1:] {
		if len(part) == 0 {
			continue
		}
		neighbour := strings.Split(part, "=")
		if len(neighbour) != 2 {
			return fmt.Errorf("wrong format for neighbour: %s", part)
		}
		neighbourCityName := neighbour[1]
		if _, ok := cityMap[neighbourCityName]; !ok {
			cityMap[neighbourCityName] = &City{
				Name:      neighbourCityName,
				Neighbors: make(map[Directions]string),
			}
		}
	}
	return nil
}

func parseNeighborToMap(cityMap map[string]*City, line string) error {
	parts := strings.Split(line, " ")
	if len(parts) == 0 {
		return fmt.Errorf("no cities found to parse")
	}
	cityName := parts[0]

	for _, part := range parts[1:] {
		if len(part) == 0 {
			continue
		}
		neighbour := strings.Split(part, "=")
		if len(neighbour) != 2 {
			return fmt.Errorf("wrong format for neighbour: %s", part)
		}
		direction, err := parseDirection(neighbour[0])
		if err != nil {
			return fmt.Errorf("eror parsing direction for neighbor: %s", neighbour[0])
		}
		neighbourCityName := neighbour[1]
		city := cityMap[cityName]
		if _, ok := city.Neighbors[direction]; !ok {
			city.Neighbors[direction] = neighbourCityName
		}
		oppositeDirection := getOppositeDirection(direction)
		if _, ok := cityMap[neighbourCityName].Neighbors[oppositeDirection]; !ok {
			cityMap[neighbourCityName].Neighbors[oppositeDirection] = cityName
		}
	}
	return nil
}

func parseDirection(direction string) (Directions, error) {
	switch direction {
	case "north":
		return North, nil
	case "south":
		return South, nil
	case "east":
		return East, nil
	case "west":
		return West, nil
	}
	return 0, fmt.Errorf("unknown direction: %s", direction)
}

func getOppositeDirection(direction Directions) Directions {
	switch direction {
	case North:
		return South
	case South:
		return North
	case East:
		return West
	case West:
		return East
	}
	return 0
}

func parseCity(line string) (*City, error) {
	parts := strings.Split(line, " ")
	if len(parts) == 0 {
		return nil, fmt.Errorf("no cities found to parse")
	}
	c := &City{
		Name:      parts[0],
		Neighbors: make(map[Directions]string),
	}
	for _, part := range parts[1:] {
		if len(part) == 0 {
			continue
		}
		neighbour := strings.Split(part, "=")
		if len(neighbour) != 2 {
			return nil, fmt.Errorf("wrong format for neighbour: %s", part)
		}
		direction, err := parseDirection(neighbour[0])
		if err != nil {
			return nil, fmt.Errorf("eror parsin direction for neighbor: %s", neighbour[0])
		}
		c.Neighbors[direction] = neighbour[1]
	}
	return c, nil
}
