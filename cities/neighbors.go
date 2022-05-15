package cities

type Directions int

const (
	North Directions = iota
	South
	East
	West
)

func (d Directions) String() string {
	switch d {
	case North:
		return "North"
	case South:
		return "South"
	case East:
		return "East"
	case West:
		return "West"
	}
	return ""
}

type neighbors map[Directions]string

func (n neighbors) GetNeighbor(d Directions) string {
	return n[d]
}

func (n neighbors) GetNeighbors() map[Directions]string {
	return n
}

func (n neighbors) NumNeighbors() int {
	return len(n)
}
