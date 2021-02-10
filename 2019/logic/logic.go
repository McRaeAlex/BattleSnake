package logic

import (
	"fmt"

	"github.com/mcraealex/BattleSnake2019/structs"
)

/*
	Representation:
	0 - open space
	1 - snake
	2 - food
	3 - tail
	4 - our tail
*/

// CreateBoard creates a BoardR which has all th relavent data on it
func CreateBoard(g structs.GeneralRequest) structs.BoardR {
	// make the board dynamicly
	board := structs.BoardR{
		Board:  make([]int, g.Board.Width*g.Board.Height),
		Width:  g.Board.Width,
		Height: g.Board.Height,
	}

	// place the foo on the board
	for _, point := range g.Board.Food {
		*board.GetRef(point.X, point.Y) = 2
	}
	// place the snakes on the board
	for _, snake := range g.Board.Snakes {
		for i := 0; i < len(snake.Body); i++ {
			point := snake.Body[i]
			if i == (len(snake.Body) - 1) {

				*board.GetRef(point.X, point.Y) = 3
			} else {
				*board.GetRef(point.X, point.Y) = 1
			}
		}
	}
	// place our snakes tail
	ourTail := g.You.Body[len(g.You.Body)-1]
	*board.GetRef(ourTail.X, ourTail.Y) = 4

	return board
}

// Logic wow
// 0 - right
// 1 - down
// 2 - left
// 3 - up
func Logic(g structs.GeneralRequest) int {
	board := CreateBoard(g)
	board.Print()
	fmt.Printf("state: %v\n", g)
	// run bfs and get the closest tail, food, and own tail
	food, tail, ownTail := getMoves(g, board)
	if g.You.Health < 50+(g.Turn/5) || len(g.You.Body) < 5 {
		// chase food
		if food != -1 {
			return food
		} else if ownTail != -1 {
			// if cannot find food
			return ownTail
		} else if tail != -1 {
			// if cannot find tail
			return tail
		}
	} else {
		if ownTail != -1 {
			// if cannot find food
			return ownTail
		} else if tail != -1 {
			// if cannot find tail
			return tail
		} else if food != -1 {
			return food
		}
	}

	return 0
}

func getMoves(g structs.GeneralRequest, b structs.BoardR) (int, int, int) {
	visited := make(map[structs.Point]bool)
	parents := make(map[structs.Point]structs.Point)
	tail := -1
	myTail := -1
	food := -1
	q := structs.Queue{Data: make([]structs.Point, 0)}
	head := g.You.Body[0]

	// enqueue all the spaces around the head that are not a snake body part
	enqueueAroundHead(head, &q, b, g)

	for q.Length() != 0 {
		//fmt.Printf("Queue: %v\n", q)
		node := q.Dequeue()

		// if we have seen this node before next node
		if visited[node] {
			continue
		}

		//fmt.Printf("Node: %v\n", node)

		// enqueue around
		enqueueAroundPoint(node, &q, b, g, &parents, &visited)

		// if node is a tail
		if b.GetValP(node) == 3 {
			fmt.Println("Found tail")
			temp := rebuild(parents, node)
			tail = getDirection(head, temp)
			fmt.Printf("tail: %v\n", tail)
		}
		// if node is a food
		if b.GetValP(node) == 2 {
			fmt.Println("Found food")
			temp := rebuild(parents, node)
			food = getDirection(head, temp)
			fmt.Printf("food: %v\n", food)
		}
		// if node is myTail
		if b.GetValP(node) == 4 {
			fmt.Println("Found my tail")
			temp := rebuild(parents, node)
			myTail = getDirection(head, temp)
			fmt.Printf("myTail: %v\n", myTail)
		}

		// make node visited
		visited[node] = true

		// if we have found paths to food tail and myTail return
		if myTail != -1 && tail != -1 && food != -1 {
			break
		}
	}

	return food, tail, myTail
}

func rebuild(m map[structs.Point]structs.Point, p structs.Point) structs.Point {
	temp := p
	for {
		val, ok := m[temp]
		if ok {
			temp = val
		} else {
			break
		}
	}
	return temp
}

func getDirection(g structs.Point, p structs.Point) int {
	diffX := p.X - g.X
	diffY := p.Y - g.Y
	if diffX > 0 {
		return 0
	}
	if diffX < 0 {
		return 2
	}
	if diffY > 0 {
		return 1
	}
	if diffY < 0 {
		return 3
	}
	return 0
}

func inBounds(node structs.Point, g structs.GeneralRequest) bool {
	if node.Y < 0 || node.X < 0 || node.Y >= g.Board.Height || node.X >= g.Board.Width {
		return false
	}
	return true
}

func enqueueAroundHead(head structs.Point, q *structs.Queue, b structs.BoardR, g structs.GeneralRequest) {
	rightHead := structs.Point{
		X: head.X + 1,
		Y: head.Y,
	}
	// if it is in bounds and not a snake body
	if inBounds(rightHead, g) && (b.GetValP(rightHead) != 1 && b.GetValP(rightHead) != 4) {
		q.Enqueue(rightHead)
		// we don't set its parent
	}

	leftHead := structs.Point{
		X: head.X - 1,
		Y: head.Y,
	}
	// if it is in bounds and not a snake body
	if inBounds(leftHead, g) && (b.GetValP(leftHead) != 1 && b.GetValP(leftHead) != 4) {
		q.Enqueue(leftHead)
		// we don't set its parent
	}

	upHead := structs.Point{
		X: head.X,
		Y: head.Y - 1,
	}
	// if it is in bounds and not a snake body
	if inBounds(upHead, g) && (b.GetValP(upHead) != 1 && b.GetValP(upHead) != 4) {
		q.Enqueue(upHead)
		// we don't set its parent
	}

	downHead := structs.Point{
		X: head.X,
		Y: head.Y + 1,
	}
	// if it is in bounds and not a snake body
	if inBounds(downHead, g) && (b.GetValP(downHead) != 1 && b.GetValP(downHead) != 4) {
		q.Enqueue(downHead)
		// we don't set its parent
	}
}

func enqueueAroundPoint(p structs.Point, q *structs.Queue, b structs.BoardR, g structs.GeneralRequest, pm *map[structs.Point]structs.Point, vm *map[structs.Point]bool) {
	// only if they are in bounds and not a snake body
	leftPoint := structs.Point{
		X: p.X - 1,
		Y: p.Y,
	}
	if inBounds(leftPoint, g) && b.GetValP(leftPoint) != 1 && !(*vm)[leftPoint] {
		q.Enqueue(leftPoint)
		(*pm)[leftPoint] = p
	}

	rightPoint := structs.Point{
		X: p.X + 1,
		Y: p.Y,
	}
	if inBounds(rightPoint, g) && b.GetValP(rightPoint) != 1 && !(*vm)[rightPoint] {
		q.Enqueue(rightPoint)
		(*pm)[rightPoint] = p
	}

	downPoint := structs.Point{
		X: p.X,
		Y: p.Y + 1,
	}
	if inBounds(downPoint, g) && b.GetValP(downPoint) != 1 && !(*vm)[downPoint] {
		q.Enqueue(downPoint)
		(*pm)[downPoint] = p
	}

	upPoint := structs.Point{
		X: p.X,
		Y: p.Y - 1,
	}
	if inBounds(upPoint, g) && b.GetValP(upPoint) != 1 && !(*vm)[upPoint] {
		q.Enqueue(upPoint)
		(*pm)[upPoint] = p
	}
}
