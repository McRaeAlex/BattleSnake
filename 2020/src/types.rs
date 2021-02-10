use std::fmt;

use serde::{Deserialize, Serialize};

#[serde(rename_all = "camelCase")]
#[derive(Serialize)]
pub(crate) struct StartResponse {
    pub color: String,
    pub head_type: String,
    pub tail_type: String,
}

#[serde(rename_all = "lowercase")]
#[derive(Deserialize, Serialize, Eq, PartialEq, Debug)]
pub(crate) enum Movement {
    Up,
    Down,
    Left,
    Right,
}

#[derive(Serialize, Deserialize)]
pub(crate) struct MoveResponse {
    #[serde(rename(serialize = "move", deserialize = "move"))]
    pub movement: Movement,
}

#[derive(Serialize, Deserialize, Eq, PartialEq, Ord, PartialOrd, Copy, Clone, Debug)]
pub(crate) struct Point {
    pub x: isize,
    pub y: isize,
}

impl Point {
    pub fn get_relative_direction_of(&self, p: &Point) -> Option<Movement> {
        match (
            p.x as isize - self.x as isize,
            p.y as isize - self.y as isize,
        ) {
            (1, 0) => Some(Movement::Right),
            (-1, 0) => Some(Movement::Left),
            (0, 1) => Some(Movement::Down),
            (0, -1) => Some(Movement::Up),
            _ => None,
        }
    }
}

#[derive(Serialize, Deserialize, Debug)]
pub(crate) struct Snake {
    pub id: String,
    pub name: String,
    pub health: usize,
    pub body: Vec<Point>,
}

#[derive(Serialize, Deserialize, Debug)]
pub(crate) struct Board {
    pub height: usize,
    pub width: usize,
    pub food: Vec<Point>,
    pub snakes: Vec<Snake>,
}

#[derive(Serialize, Deserialize, Debug)]
pub(crate) struct MoveRequest {
    pub turn: usize,
    pub board: Board,
    pub you: Snake,
}

#[derive(Debug, Eq, PartialEq, Clone, Copy)]
pub(crate) enum Tile {
    // Debating to add a wall type aswell...
    YouTail,
    YouHead,
    Snake,
    SnakeHead,
    Food,
    Empty,
    Tail,
}

impl std::fmt::Display for Tile {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        let symbol = match *self {
            Tile::Snake => "S",
            Tile::Tail => "T",
            Tile::Empty => "0",
            Tile::SnakeHead => "H",
            Tile::Food => "F",
            Tile::YouHead => "Y",
            Tile::YouTail => "L"
        };

        write!(f, "{}", symbol)?;
        Ok(())
    }
}

#[derive(Debug)]
pub(crate) struct GameBoard {
    height: usize,
    width: usize,
    tiles: Vec<Tile>,
}

impl GameBoard {
    pub fn set_val(&mut self, x: usize, y: usize, t: Tile) {
        self.tiles[self.height * y + x] = t;
    }

    pub fn from_move_req(req: &MoveRequest) -> GameBoard {
        let height = req.board.height;
        let width = req.board.width;
        let mut b = Vec::with_capacity(height * width);
        for _ in 0..(height * width) {
            b.push(Tile::Empty);
        }

        let mut board = GameBoard {
            height: height,
            width: width,
            tiles: b,
        };

        // Go through and add each snake to the board
        for snake in req.board.snakes.iter() {
            for bod in snake.body.iter() {
                board.set_val(bod.x as usize, bod.y as usize, Tile::Snake);
            }

            let tail: &Point = snake.body.last().unwrap();
            board.set_val(tail.x as usize, tail.y as usize, Tile::Tail);

            let head = snake.body.first().unwrap();
            board.set_val(head.x as usize, head.y as usize, Tile::SnakeHead);
        }

        // Go through the food and add to board
        for food in req.board.food.iter() {
            board.set_val(food.x as usize, food.y as usize, Tile::Food);
        }

        // Add self to board
        let tail: &Point = req.you.body.last().unwrap();
        board.set_val(tail.x as usize, tail.y as usize, Tile::YouTail);

        let head: &Point = req.you.body.first().unwrap();
        board.set_val(head.x as usize, head.y as usize, Tile::YouHead);

        return board;
    }

    pub fn get_val(&self, p: Point) -> Option<&Tile> {
        self.tiles.get(p.y as usize * self.height + p.x as usize)
    }

    pub fn in_bounds(&self, p: &Point) -> bool {
        let x = p.x as isize;
        let y = p.y as isize;
        if 0 <= x && 0 <= y && x < self.width as isize && y < self.height as isize {
            true
        } else {
            false
        }
    }
}

impl std::fmt::Display for GameBoard {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(f, "height: {} width: {}\n", self.height, self.width)?;
        for y in 0..self.height {
            for x in 0..self.width {
                write!(f, "{}", *(self.get_val(Point{x: x as isize, y: y as isize}).unwrap()))?;
            }
            writeln!(f, "")?;
        }
        Ok(())
    }
}
