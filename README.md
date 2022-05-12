
# Alien invasion

<html><img src="./aliens.png" height="300px"/></html>

## Description

Golang based app that simulates a earth invasion of outerspace mad aliens.

These aliens start out at random places on the given world map, and wander around randomly, following links (directions). 

In each iteration, the aliens can travel in any of the directions leading out of a city.

When two aliens end up in the same place, they fight, and in the process kill each other and destroy the city. When a city is destroyed, it is removed from the map, and so are any roads that lead into or out of it.

The program runs until all the aliens have been destroyed, or each alien has moved at least 10,000 times. 

When two aliens fight, the program prints out a message like:

`City Bee destroyed by Aliens 2 and 4 `

## Features
- Concurrency model based on `channel` type (no `mutex` used)
- Debug mode available
- Progress bar that shows execution progress when program is waiting for iterations to complete:
```
1 alien left. Waiting for it to exit or iterations to complete...  94% [=============> ]  [5s:0s]
```

- Sample maps provided in `data` folder:
  
  `world_sample.txt` (5 cities)

  `world_big.txt` (13 cities)

- Command line options
```bash
usage: main --file=FILE [<flags>]

Flags:
      --help              Show context-sensitive help (also try --help-long and --help-man).
  -a, --aliens=12         Number of aliens. Defaults to system's number of CPUs.
  -i, --iterations=10000  Number of iterations. Defaults to 10.000.
  -d, --debug             debug mode
  -f, --file=FILE         input source file
      --version           Show application version.
```
- Test suite available
  
## Assumptions

- If an alien wants to "land" but no cities are available (all cities have already been destroyed), it will "die" (i.e. exit) automatically
- Cities can only have a maximun of 1 neighbor city *per* direction (i.e, north, south, east, west)
- At the moment, the `world` is only supported as a `map` data structure
- When a city is destroyed, the worldmap is updated without bridging (i.e. connecting) remaining cities. This behaviour could potentially lead to cities being isolated
- Input file format: The city and each of the pairs should be separated by a single space, and the directions are separated from their respective cities with an equals (=) sign, i.e:

`world_sample.txt`
```
Foo north=Bar west=Baz south=Qu-ux
Bar south=Foo west=Bee
```
  

## Usage examples

### Example 1

- Input:
  
  `aliens` = 4

  `iterations` = 500

  `file` = ./data/world_sample.txt

```bash
❯ go run main.go -a 4 -i 500  -f ./data/world_sample.txt
```

- Reference ouput
```bash
INFO[0000] World loaded from file                       

Printing world:
----------------------------------------
Number of cities:  5
        City:  Bee
                Neighbor:  Bar  ->  East
        City:  Foo
                Neighbor:  Bazz  ->  West
                Neighbor:  Bar  ->  North
                Neighbor:  Qu-ux  ->  South
        City:  Bar
                Neighbor:  Foo  ->  South
                Neighbor:  Bee  ->  West
        City:  Qu-ux
                Neighbor:  Foo  ->  North
        City:  Bazz
                Neighbor:  Foo  ->  East
----------------------------------------

INFO[0000] Alien invasion started with 4 aliens and 500 iterations 
WARN[0000] City Bee destroyed by Aliens 2 and 4          alien=Alien2 id=2
WARN[0000] City Qu-ux destroyed by Aliens 3 and 1        alien=Alien3 id=3

INFO[0000] Invasion terminated ...                      

Printing world:
----------------------------------------
Number of cities:  3
        City:  Foo
                Neighbor:  Bazz  ->  West
                Neighbor:  Bar  ->  North
        City:  Bar
                Neighbor:  Foo  ->  South
        City:  Bazz
                Neighbor:  Foo  ->  East
----------------------------------------

INFO[0000] Program finished  
```

### Example 2

- Input:
  
  `aliens` = 11

  `iterations` = 5000

  `file` = ./data/world_big.txt

```bash
❯ go run main.go -a 11 -i 5000  -f ./data/world_big.txt
```

- Reference ouput

```bash
INFO[0000] World loaded from file                       

Printing world:
----------------------------------------
Number of cities:  13
        City:  City1
                Neighbor:  City2  ->  North
                Neighbor:  City4  ->  East
                Neighbor:  City6  ->  South
                Neighbor:  City8  ->  West
        City:  City10
                Neighbor:  City2  ->  South
        City:  City9
                Neighbor:  City2  ->  East
                Neighbor:  City8  ->  South
        City:  City12
                Neighbor:  City6  ->  North
        City:  City13
                Neighbor:  City8  ->  East
        City:  City7
                Neighbor:  City6  ->  East
                Neighbor:  City8  ->  North
        City:  City11
                Neighbor:  City4  ->  West
        City:  City2
                Neighbor:  City10  ->  North
                Neighbor:  City3  ->  East
                Neighbor:  City9  ->  West
                Neighbor:  City1  ->  South
        City:  City4
                Neighbor:  City1  ->  West
                Neighbor:  City3  ->  North
                Neighbor:  City11  ->  East
                Neighbor:  City5  ->  South
        City:  City6
                Neighbor:  City1  ->  North
                Neighbor:  City5  ->  East
                Neighbor:  City12  ->  South
                Neighbor:  City7  ->  West
        City:  City8
                Neighbor:  City13  ->  West
                Neighbor:  City1  ->  East
                Neighbor:  City9  ->  North
                Neighbor:  City7  ->  South
        City:  City3
                Neighbor:  City2  ->  West
                Neighbor:  City4  ->  South
        City:  City5
                Neighbor:  City6  ->  West
                Neighbor:  City4  ->  North
----------------------------------------

INFO[0000] Alien invasion started with 11 aliens and 5000 iterations 
WARN[0000] City City12 destroyed by Aliens 5 and 4       alien=Alien5 id=5
WARN[0000] City City1 destroyed by Aliens 3 and 9        alien=Alien3 id=3
WARN[0000] City City4 destroyed by Aliens 8 and 10       alien=Alien8 id=8
WARN[0000] City City6 destroyed by Aliens 11 and 1       alien=Alien11 id=11
WARN[0000] City City2 destroyed by Aliens 7 and 2        alien=Alien7 id=7
 1 alien left. Waiting for it to exit or iterations to complete...  99% [=============> ]  [27s:0s] 
INFO[0027] Killing remaining 1 aliens after 5000 iterations 
INFO[0027] Invasion terminated ...                      

Printing world:
----------------------------------------
Number of cities:  8
        City:  City8
                Neighbor:  City9  ->  North
                Neighbor:  City7  ->  South
                Neighbor:  City13  ->  West
        City:  City3
        City:  City5
        City:  City7
                Neighbor:  City8  ->  North
        City:  City11
        City:  City10
        City:  City9
                Neighbor:  City8  ->  South
        City:  City13
                Neighbor:  City8  ->  East
----------------------------------------

INFO[0027] Program finished  
```

### Example 3

- Input:
  
  `aliens` = 30

  `iterations` = 10000

  `file` = ./data/world_big.txt

```bash
❯ go run main.go -a 11 -i 5000  -f ./data/world_big.txt
```

- Reference ouput

```bash
❯ go run -race main.go -a 30 -i 10000  -f ./data/world_big.txt
INFO[0000] World loaded from file                       

INFO[0000] Alien invasion started with 30 aliens and 10000 iterations 
WARN[0000] City City7 destroyed by Aliens 12 and 14      alien=Alien12 id=12
WARN[0000] City City2 destroyed by Aliens 16 and 2       alien=Alien16 id=16
WARN[0000] City City6 destroyed by Aliens 11 and 4       alien=Alien11 id=11
WARN[0000] City City11 destroyed by Aliens 7 and 5       alien=Alien7 id=7
WARN[0000] City City1 destroyed by Aliens 13 and 3       alien=Alien13 id=13
WARN[0000] City City8 destroyed by Aliens 9 and 6        alien=Alien9 id=9
WARN[0000] City City10 destroyed by Aliens 15 and 19     alien=Alien15 id=15
WARN[0000] City City13 destroyed by Aliens 25 and 1      alien=Alien25 id=25
WARN[0000] City City9 destroyed by Aliens 17 and 28      alien=Alien17 id=17
WARN[0000] City City4 destroyed by Aliens 10 and 21      alien=Alien10 id=10
WARN[0000] City City3 destroyed by Aliens 26 and 8       alien=Alien26 id=26
WARN[0000] City City12 destroyed by Aliens 24 and 22     alien=Alien24 id=24
WARN[0000] City City5 destroyed by Aliens 30 and 23      alien=Alien30 id=30
INFO[0000] no more cities available. Exiting             alien=Alien20 id=20
INFO[0000] no more cities available. Exiting             alien=Alien18 id=18
INFO[0000] no more cities available. Exiting             alien=Alien27 id=27
INFO[0000] no more cities available. Exiting             alien=Alien29 id=29

INFO[0000] Invasion terminated ...                      

Printing world:
----------------------------------------
Number of cities:  0
----------------------------------------

INFO[0000] Program finished          
```

## Test coverage

```
❯ go test -v ./..
=== RUN   TestAliensLifeCycleBasic
=== RUN   TestAliensLifeCycleBasic/Should_start_N_aliens
=== RUN   TestAliensLifeCycleBasic/Should_stop_N_aliens_on_cancel
--- PASS: TestAliensLifeCycleBasic (0.30s)
    --- PASS: TestAliensLifeCycleBasic/Should_start_N_aliens (0.30s)
    --- PASS: TestAliensLifeCycleBasic/Should_stop_N_aliens_on_cancel (0.00s)
PASS
ok      github.com/alejoacosta74/allien_invasion/aliens 1.172s
=== RUN   TestAlienInvasion
=== RUN   TestAlienInvasion/With_odd_number_of_aliens,_at_least_1_alien_should_stay_alive_and_max_iterations_reached

=== RUN   TestAlienInvasion/With_aliens_(even)_>>_5,_all_cities_will_likely_be_destroyed

--- PASS: TestAlienInvasion (0.01s)
    --- PASS: TestAlienInvasion/With_odd_number_of_aliens,_at_least_1_alien_should_stay_alive_and_max_iterations_reached (0.00s)
    --- PASS: TestAlienInvasion/With_aliens_(even)_>>_5,_all_cities_will_likely_be_destroyed (0.01s)
PASS
ok      github.com/alejoacosta74/allien_invasion/app    0.721s
=== RUN   TestCityFileReader
=== RUN   TestCityFileReader/Reads_all_the_cities
=== RUN   TestCityFileReader/Reads_all_neighbours_and_directions
--- PASS: TestCityFileReader (0.00s)
    --- PASS: TestCityFileReader/Reads_all_the_cities (0.00s)
    --- PASS: TestCityFileReader/Reads_all_neighbours_and_directions (0.00s)
=== RUN   TestLoadMap
=== RUN   TestLoadMap/Should_create_a_world_structure_from_a_file
--- PASS: TestLoadMap (0.00s)
    --- PASS: TestLoadMap/Should_create_a_world_structure_from_a_file (0.00s)
PASS
```


## To Do

- Add support for graph structure
- Increase test coverage