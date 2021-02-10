#![feature(proc_macro_hygiene, decl_macro)]

#[macro_use]
extern crate rocket;


#[cfg(test)] mod test;

mod types;
mod logic;

use types::{MoveResponse, StartResponse, MoveRequest, GameBoard};
use logic::logic;

use rocket::Rocket;
use rocket_contrib::json::Json;


#[post("/start")]
fn start() -> Json<StartResponse> {
    Json(StartResponse {
        color: "#b7410e".to_string(),
        head_type: "tongue".to_string(),
        tail_type: "sharp".to_string(),
    })
}

#[post("/move", format = "json", data = "<data>")]
fn movement(data: Json<MoveRequest>) -> Json<MoveResponse> {
    // Parse into structure we want
    let board = data.into_inner();
    println!("board: {:?}", board);
    let board_with_objects = GameBoard::from_move_req(&board);

    // Do logic
    let mov = logic(&board_with_objects, &board);

    // Response
    Json(MoveResponse {
        movement: mov,
    })
}

#[post("/ping")]
fn ping() -> &'static str {
    "I am alive!"
}

#[get("/")]
fn index() -> &'static str {
    "Hello world"
}

fn rocket() -> Rocket {
    rocket::ignite().mount("/", routes![index, start, movement, ping])
}

fn main() {
    rocket().launch();
}
