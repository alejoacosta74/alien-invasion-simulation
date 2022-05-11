package cities

type City struct {
	Name            string
	Neighbors       neighbors
	visited         bool
	visitor         int
	visitorQuitChan chan struct{}
}

func NewCity(name string, nbrs neighbors) *City {
	city := &City{
		Name:      name,
		Neighbors: nbrs,
		visited:   false,
	}
	return city
}

func (c *City) IsVisited() bool {
	return c.visited
}

func (c *City) Visit(id int, quitChan chan struct{}) {
	c.visited = true
	c.visitor = id
	c.visitorQuitChan = quitChan
}

func (c *City) Leave(id int, quitChan chan struct{}) {
	c.visited = false
	c.visitor = 0
	c.visitorQuitChan = nil
}

func (c *City) GetVisitor() int {
	return c.visitor
}

func (c *City) GetVisitorQuitChannel() chan struct{} {
	return c.visitorQuitChan
}

func (c *City) GetNeighbor(d Directions) string {
	return c.Neighbors[d]
}

func (c *City) GetNeighborsMap() map[Directions]string {
	// fmt.Printf("City: %+v\n", c)
	// fmt.Printf("CityNeigborhs: %+v\n", c.Neighbors)
	return c.Neighbors
}

func (c *City) NumNeighbors() int {
	return len(c.Neighbors)
}
