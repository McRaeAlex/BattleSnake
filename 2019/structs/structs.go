package structs

import (
	"fmt"
)

// MoveResponse is a struct to send back to the battlesnake server
type MoveResponse struct {
	Move string `json:"move"`
}

// StartResponse is a struct to send back if /start is hit
type StartResponse struct {
	Color string `json:"color"`
}

// GeneralRequest is a struct for all JSON requests coming to the server
type GeneralRequest struct {
	Game struct {
		ID string `json:"id"`
	} `json:"game"`
	Turn  int   `json:"turn"`
	Board Board `json:"board"`
	You   Snake `json:"you"`
}

// Board is the JSON struct for the board and not the actual board itself
type Board struct {
	Height int     `json:"height"`
	Width  int     `json:"width"`
	Food   []Point `json:"food"`
	Snakes []Snake `json:"snakes"`
}

// Point is just a location on the board
type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Snake is the JSON structure for the snake
type Snake struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Health int     `json:"health"`
	Body   []Point `json:"body"`
}

// BoardR is the respresentation of the board with all of the snakes on it.
type BoardR struct {
	Board  []int
	Height int
	Width  int
}

func (b *BoardR) GetVal(x int, y int) int {
	return b.Board[y*b.Width+x]
}

func (b *BoardR) GetValP(p Point) int {
	return b.GetVal(p.X, p.Y)
}

func (b *BoardR) GetRef(x int, y int) *int {
	return &b.Board[y*b.Width+x]
}

func (b *BoardR) Print() {
	for i := 0; i < len(b.Board); i++ {
		fmt.Print(b.Board[i])
		if (i+1)%b.Width == 0 {
			fmt.Println()
		}
	}
}

type Queue struct {
	Data []Point
}

func (q *Queue) Dequeue() Point {
	val := q.Data[0]
	q.Data = q.Data[1:]
	return val
}

func (q *Queue) Enqueue(val Point) {
	q.Data = append(q.Data, val)
}

func (q *Queue) Length() int {
	return len(q.Data)
}

type GeneralRequestR struct {
	Game struct {
		ID string `json:"id"`
	} `json:"game"`
	Turn  int `json:"turn"`
	Board struct {
		Height int `json:"height"`
		Width  int `json:"width"`
		Food   []struct {
			X int `json:"x"`
			Y int `json:"y"`
		} `json:"food"`
		Snakes []struct {
			ID     string `json:"id"`
			Name   string `json:"name"`
			Health int    `json:"health"`
			Body   []struct {
				X int `json:"x"`
				Y int `json:"y"`
			} `json:"body"`
		} `json:"snakes"`
	} `json:"board"`
	You struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Health int    `json:"health"`
		Body   []struct {
			X int `json:"x"`
			Y int `json:"y"`
		} `json:"body"`
	} `json:"you"`
}
