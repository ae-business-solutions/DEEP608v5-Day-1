#![feature(proc_macro_hygiene, decl_macro)]
#[macro_use] extern crate rocket;
extern crate redis;

use rocket::Outcome;
use rocket::http::Status;
use rocket::request::{self, Request, FromRequest, Form};
use rocket_contrib::templates::Template;
use serde::Serialize;
use uuid::Uuid;
use std::env;
use std::collections::BTreeMap;
use std::collections::HashMap;

#[derive(FromForm, Debug)]
struct Req {
  user: String,
  url: String,
  reason: String
}

struct Username {
  username: String
}

impl<'a, 'r> FromRequest<'a, 'r> for Username {
  type Error = String;
  fn from_request(request: &'a Request<'r>) -> request::Outcome<Self, Self::Error> {
      let username = request.headers().get_one("x-username");
      match username {
        Some(username) => {
          Outcome::Success(Username{username: username.to_string()})
        },
        None => Outcome::Failure((Status::Forbidden, "Unauthorized".to_string()))
      }
  }
}

#[get("/health")]
fn health() -> &'static str {
  "OK"
}

#[get("/")]
fn index(username: Username) -> Template {
  #[derive(Serialize)]
  struct Context {
    user: String,
  }

  let context = Context {
    user: username.username,
  };
  Template::render("index", context)
}

#[catch(404)]
fn not_found(req: &Request) -> String {
    print!("{}", req);
    format!("Oh no! We couldn't find the requested path '{}'", req.uri())
}

fn connect() -> redis::Connection {
  let redis_host = env::var("REDIS_HOST").expect("missing environment variable REDIS_HOST");
  let redis_password = env::var("REDIS_PASSWORD").unwrap_or_default();
  let uri_scheme = match env::var("IS_TLS") {
      Ok(_) => "rediss",
      Err(_) => "redis",
  };
  let redis_conn_url = format!("{}://:{}@{}", uri_scheme, redis_password, redis_host);
  redis::Client::open(redis_conn_url)
      .expect("Invalid connection URL")
      .get_connection()
      .expect("failed to connect to Redis")
}

#[post("/request", data = "<request_form>")]
fn new_request(request_form: Form<Req>) -> Template {
  let request: Req = request_form.into_inner();
  let uuid = Uuid::new_v4().to_simple().to_string();
  let mut conn = connect();
  let mut request_data: BTreeMap<String, String> = BTreeMap::new();
  request_data.insert(String::from("user"), request.user.clone());
  request_data.insert(String::from("url"), request.url.clone());
  request_data.insert(String::from("reason"), request.reason.clone());
  let _: () = redis::cmd("HSET")
        .arg(format!("{}:{}", "request", uuid))
        .arg(request_data)
        .query(&mut conn)
        .expect("failed to execute HSET");
  let context: HashMap<&str, &str> = HashMap::with_capacity(0);
  Template::render("submitted", context)
}

fn main() {
  rocket::ignite()
    .register(catchers![not_found])
    .mount("/", routes![index, health, new_request])
    .attach(Template::fairing())
    .launch();
}