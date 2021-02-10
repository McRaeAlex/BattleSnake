use crate::rocket;
use crate::types::{MoveResponse, Movement};

use rocket::http::ContentType;
use rocket::http::Status;
use rocket::local::Client;
use serde_json;

#[test]
fn test_start() {
  let client = Client::new(rocket()).expect("Valid rocket instance");
  let resp = client.post("/start").dispatch();

  assert_eq!(resp.status(), Status::Ok);
}

#[test]
fn test_ping() {
  let client = Client::new(rocket()).expect("Valid rocket instance");
  let resp = client.post("/ping").dispatch();

  assert_eq!(resp.status(), Status::Ok);
}

#[test]
fn movement() {
  println!("Testing different situations");
  let client = Client::new(rocket()).expect("Valid rocket instance");
  // Go through tests and parse each test case.
  for test in TEST_CASES.iter() {
    println!("Test description: {}", test.description);

    // Then create a request for each test case and run it!
    let req = client
      .post("/move")
      .header(ContentType::JSON)
      .body(test.data);
    let mut resp = req.dispatch();

    // Then check to makesure the output is correct
    assert_eq!(resp.status(), Status::Ok);
    assert_eq!(resp.content_type(), Some(ContentType::JSON));

    let body: MoveResponse = serde_json::from_reader(
      resp
        .body()
        .expect("Response did not have a body")
        .into_inner(),
    )
    .expect("Could not deserialize json");

    // Check our test case
    match &test.correct {
      Tests::Eq(mov) => assert_eq!(body.movement, *mov),
      Tests::Neq(mov) => assert_ne!(body.movement, *mov),
      Tests::Eq2(move1, move2) => assert!([move1, move2].iter().any(|x| **x == body.movement)),
      Tests::Neq2(move1, move2) => assert!([move1, move2].iter().all(|x| **x != body.movement)),
    }
  }
}

#[allow(dead_code)]
enum Tests {
  Eq(Movement),
  Eq2(Movement, Movement),
  Neq(Movement),
  Neq2(Movement, Movement),
}

struct MoveTestCase<'a> {
  description: &'a str,
  data: &'a str,
  correct: Tests,
}

// Remember to increment the number each time a test case is added
const TEST_CASES: [MoveTestCase; 3] = [
  MoveTestCase {
    // Testing chase tail when no food or other snake on board
    description: "Should chase tail when no food on board or other snakes",
    data: r#"{
      "game": {
        "id": "game-id-string"
      },
      "turn": 4,
      "board": {
        "height": 15,
        "width": 15,
        "food": [
        ],
        "snakes": [
          {
            "id": "snake-id-string",
            "name": "Sneky Snek",
            "health": 90,
            "body": [
              { "x": 3, "y": 2 },
              { "x": 3, "y": 3 },
              { "x": 4, "y": 3 },
              { "x": 4, "y": 2 }
            ]
          }
        ]
      },
      "you": {
        "id": "snake-id-string",
        "name": "Sneky Snek",
        "health": 90,
        "body": [
          { "x": 3, "y": 2 },
          { "x": 3, "y": 3 },
          { "x": 4, "y": 3 },
          { "x": 4, "y": 2 }
        ]
      }
    }"#,
    correct: Tests::Neq(Movement::Down),
  },
  MoveTestCase {
    description: "Should avoid killing self when in circle",
    data: r#"{
    "game": {
      "id": "game-id-string"
    },
    "turn": 4,
    "board": {
      "height": 15,
      "width": 15,
      "food": [
      ],
      "snakes": [
        {
          "id": "snake-id-string",
          "name": "Sneky Snek",
          "health": 90,
          "body": [
            { "x": 3, "y": 2 },
            { "x": 3, "y": 3 }
          ]
        }
      ]
    },
    "you": {
      "id": "snake-id-string",
      "name": "Sneky Snek",
      "health": 90,
      "body": [
        { "x": 3, "y": 2 },
        { "x": 3, "y": 3 }
      ]
    }
  }"#,
    correct: Tests::Neq(Movement::Down),
  },
  MoveTestCase {
    description: "Should avoid hitting others heads when smaller", // TODO: this
    data: r#"{
    "game": {
      "id": "game-id-string"
    },
    "turn": 4,
    "board": {
      "height": 15,
      "width": 15,
      "food": [
      ],
      "snakes": [
        {
          "id": "snake-id-string",
          "name": "Sneky Snek",
          "health": 90,
          "body": [
            { "x": 3, "y": 2 },
            { "x": 3, "y": 3 }
          ]
        }
      ]
    },
    "you": {
      "id": "snake-id-string",
      "name": "Sneky Snek",
      "health": 90,
      "body": [
        { "x": 3, "y": 2 },
        { "x": 3, "y": 3 }
      ]
    }
  }"#,
    correct: Tests::Neq(Movement::Down),
  },
  
];
