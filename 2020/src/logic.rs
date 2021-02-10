use crate::types::{GameBoard, MoveRequest, Movement, Point, Tile};

use std::collections::{BTreeMap, BTreeSet, VecDeque};

/**
 * This is the logic function which makes sure the snake works
 */
pub(crate) fn logic(board: &GameBoard, req: &MoveRequest) -> Movement {
    let you = &req.you;
    let mut mov;

    // TODO: we need to check if the move we make will kill us if another snake moves into the space
    

    // Logic
    if you.body.len() < 10 || you.health < 50 {
        println!("Going for food!");
        mov = bfs(board, req, Tile::Food);
    } else {
        println!("Going for own tail!");
        mov = bfs(board, req, Tile::YouTail);
        if mov.is_none() {
            println!("Did not find a tail");
            mov = bfs(board, req, Tile::Tail);
        }
    }
    
    let movement = match mov {
        Some(val) => val,
        None => {
            println!("We did not make a move defaulting to up");
            Movement::Up
        },
    }; // When in doubt go up
    println!("Made the move: {:?}", movement);

    movement
}

// Returns the closest to tail and food
fn bfs(board: &GameBoard, req: &MoveRequest, tile_type: Tile) -> Option<Movement> {
    println!("{}", board);
    let mut tiles: VecDeque<Point> = VecDeque::new();
    let mut visited: BTreeSet<Point> = BTreeSet::new();
    let mut parent_graph: BTreeMap<Point, Point> = BTreeMap::new();


    // Add the tiles around the head
    let head = req.you.body.first().unwrap(); // There should always be a head in a valid request
    enqueue_around_head(head, &board, &mut tiles);

    while !tiles.is_empty() {

        let curr = tiles.pop_front().unwrap();

        visited.insert(curr);

        let t = match board.get_val(curr) {
            Some(v) => v,
            None => continue,
        };
        
        if *t == tile_type {
            return head.get_relative_direction_of(&rebuild(&parent_graph, &curr));
        }
        if *t == Tile::Empty || *t == Tile::Tail {
            enqueue_around(&curr, &board, &mut tiles, &mut visited, &mut parent_graph)
        }
    }

    None
}

fn enqueue_around_head(
    p: &Point,
    b: &GameBoard,
    q: &mut VecDeque<Point>
) {
    let around = [
        Point { x: p.x, y: p.y - 1 },
        Point { x: p.x, y: p.y + 1 },
        Point { x: p.x - 1, y: p.y },
        Point { x: p.x + 1, y: p.y },
    ];
    for a in around.iter() {
        let val = match b.get_val(*a) {
            Some(val) => val,
            None => continue,
        };
        if *val != Tile::Snake && *val != Tile::SnakeHead && b.in_bounds(a) {
            q.push_back(*a);
        }
    }
}

fn enqueue_around(
    p: &Point,
    b: &GameBoard,
    q: &mut VecDeque<Point>,
    visited: &mut BTreeSet<Point>,
    parent_graph: &mut BTreeMap<Point, Point>,
) {
    let around = [
        Point { x: p.x, y: p.y - 1 },
        Point { x: p.x, y: p.y + 1 },
        Point { x: p.x - 1, y: p.y },
        Point { x: p.x + 1, y: p.y },
    ];
    for a in around.iter() {
        if !visited.contains(a) && b.in_bounds(a) {
            q.push_back(*a);
            parent_graph.insert(*a, *p);
            visited.insert(*a);
        }
    }
}

/// follows the graph of tiles to find the tile that the one point was enqueued by
fn rebuild(parents: &BTreeMap<Point, Point>, p: &Point) -> Point {
    let mut temp = p;
    loop {
        match parents.get(temp) {
            Some(val) => temp = val,
            None => break,
        }
    }

    *temp
}
