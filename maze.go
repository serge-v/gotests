/*
Create a random maze
*/

package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	maze_width         = 20
	maze_height        = 20
	highlight_solution = false
)

type room struct {
	x, y int
}

func (r room) String() string {
	return fmt.Sprintf("(%d,%d)", r.x, r.y)
}

func (r room) id() int {
	return (r.y * maze_width) + r.x
}

// whetwher walls are open or not. (open if true)
// There are (num_rooms * 2) walls. Some are on borders, but nevermind them ;)
type wall_register [maze_width * maze_height * 2]bool

var wr = wall_register{}

// rooms are visited or not
type room_register [maze_width * maze_height]bool

var rr = room_register{}

// very subjective entrance/exit...
var entrance = room{0, 0}
var exit = room{maze_width - 1, maze_height - 1}

func main() {
	seed := time.Now().UnixNano()
	rand.Seed(seed)

	stack := make([]room, 0, maze_width*maze_height)
	current_room := room{rand.Intn(maze_width), rand.Intn(maze_height)}

	// mark start position visited
	rr[current_room.id()] = true

	for {
		// Slice of neighbors we can move
		available_neighbors := current_room.nonvisited_neighbors()

		// cannot move. Go back!
		if len(available_neighbors) < 1 {
			if len(stack) == 0 { // everything is visited
				break
			}

			current_room = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			continue
		}

		// pick a random room to move.
		next := available_neighbors[rand.Intn(len(available_neighbors))]

		// mark next visited
		rr[next.id()] = true

		// open wall between current_room and next:
		first, second := order_rooms(current_room, next)

		// second is either at the right or bottom of first.
		if second.x == first.x+1 {
			wr[first.id()*2] = true
		} else if second.y == first.y+1 {
			wr[first.id()*2+1] = true
		} else { // probably impossible or maybe not...
			panic("Wot?!?")
		}
		// push next to stack
		stack = append(stack, next)
		current_room = next
	}

	// try to solve
	find_solution(entrance)
	//	fmt.Printf("\ni think the solution is: %v\n\n", solution)

	// print maze

	// print upper border
	for x := 0; x < maze_width; x++ {
		if x == 0 {
			// maze entrance
			if highlight_solution {
				fmt.Printf("\x1b[41m \x1b[49mV") // background red
			} else {
				fmt.Printf(" V")
			}
		} else {
			if highlight_solution {
				fmt.Printf("\x1b[41m _\x1b[49m") // background red
			} else {
				fmt.Printf(" _")
			}
		}
	}
	if highlight_solution {
		fmt.Println("\x1b[41m \x1b[49m") // background red
	} else {
		fmt.Println(" ")
	}

	for y := 0; y < maze_height; y++ {

		// left border
		if highlight_solution {
			fmt.Printf("\x1b[41m|\x1b[49m") // background red
		} else {
			fmt.Printf("|")
		}

		for x := 0; x < maze_width; x++ {
			draw_room := room{x, y}

			part_of_solution := room_in_rooms_list(draw_room, solution)

			id := draw_room.id()
			right := "|"
			bottom := "_"
			if wr[id*2] {
				right = " "
			}
			if wr[id*2+1] {
				bottom = " "
			}
			if x == exit.x && y == exit.y {
				right = " >" // maze exit
			}

			if highlight_solution {
				if part_of_solution {
					//if bottom == "_" {
					//	bottom = "\x1b[31m_\x1b[39m" // foreground red
					//}
					if right == "|" {
						right = "\x1b[41m|\x1b[49m" // background red
					}
					fmt.Printf("%s%s", bottom, right)
				} else {
					fmt.Printf("\x1b[41m%s%s\x1b[49m", bottom, right) // background red
				}
			} else {
				fmt.Printf("%s%s", bottom, right)
			}
		}
		fmt.Println()
	}
	//	fmt.Printf("seed: %v\n", seed)
}

// return slice of neighbor rooms
func (r room) neighbors() []room {
	rslice := make([]room, 0, 4)
	if r.x < maze_width-1 {
		// right
		rslice = append(rslice, room{r.x + 1, r.y})
	}
	if r.x > 0 {
		// left
		rslice = append(rslice, room{r.x - 1, r.y})
	}
	if r.y < maze_height-1 {
		// bottom
		rslice = append(rslice, room{r.x, r.y + 1})
	}
	if r.y > 0 {
		// top
		rslice = append(rslice, room{r.x, r.y - 1})
	}
	return rslice
}

// return rooms that are not visited yet
func (r room) nonvisited_neighbors() []room {
	rslice := make([]room, 0, 4)
	for _, r := range r.neighbors() {
		if rr[r.id()] == false {
			rslice = append(rslice, r)
		}
	}
	return rslice
}

// order to rooms by closeness to origin (upperleft)
func order_rooms(room1, room2 room) (room, room) {
	dist1 := room1.x*room1.x + room1.y*room1.y
	dist2 := room2.x*room2.x + room2.y*room2.y
	if dist1 < dist2 {
		return room1, room2
	}
	return room2, room1
}

func room_in_rooms_list(r room, rooms []room) bool {
	for i := 0; i < len(rooms); i++ {
		if r.x == rooms[i].x && r.y == rooms[i].y {
			return true
		}
	}
	return false
}

var solution = make([]room, 0, maze_width*maze_height)
var solution_visited = make([]room, 0, maze_width*maze_height)

func find_solution(entrance room) []room {
	find_path(entrance)
	return solution
}
func find_path(r room) bool {
	if r.x == exit.x && r.y == exit.y {
		solution = append([]room{r}, solution...)
		return true
	}
	if room_in_rooms_list(r, solution_visited) {
		return false
	}
	solution_visited = append(solution_visited, r)

	if r.x < maze_width-1 {
		// right
		next_room := room{r.x + 1, r.y}
		open_wall := wr[r.id()*2]
		if open_wall && find_path(next_room) {
			solution = append([]room{r}, solution...)
			return true
		}
	}
	if r.x > 0 {
		// left
		next_room := room{r.x - 1, r.y}
		open_wall := wr[next_room.id()*2]
		if open_wall && find_path(next_room) {
			solution = append([]room{r}, solution...)
			return true
		}
	}
	if r.y < maze_height-1 {
		// bottom
		next_room := room{r.x, r.y + 1}
		open_wall := wr[r.id()*2+1]
		if open_wall && find_path(next_room) {
			solution = append([]room{r}, solution...)
			return true
		}
	}
	if r.y > 0 {
		// top
		next_room := room{r.x, r.y - 1}
		open_wall := wr[next_room.id()*2+1]
		if open_wall && find_path(next_room) {
			solution = append([]room{r}, solution...)
			return true
		}
	}

	return false
}
