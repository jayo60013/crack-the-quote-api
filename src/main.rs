mod models;
mod routes;
mod utils;

use std::{sync::Arc, time::Instant};

use actix_cors::Cors;
use actix_web::{
    App, HttpServer, http,
    middleware::{Compress, Logger},
    web,
};
use env_logger::Env;
use log::{debug, error, info};
use tokio::sync::RwLock;

use crate::{
    models::{daily_puzzle::DailyPuzzle, daily_puzzle_response::DailyPuzzleResponse},
    utils::{
        constants::PORT_NUMBER, daily_puzzle_utils::get_daily_puzzle_entity,
        db_utils::connect_pool, init_quotes_utils::initialise_quotes_table,
    },
};

type DailyPuzzleCache = Arc<RwLock<DailyPuzzle>>;
type DailyPuzzleResponseCache = Arc<RwLock<DailyPuzzleResponse>>;

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    env_logger::Builder::from_env(
        Env::default().default_filter_or("actix_web=debug,crack_the_quote_api=debug"),
    )
    .init();

    let pool = match connect_pool().await {
        Ok(v) => v,
        Err(e) => {
            error!("Failed to connect to postgres db with error: {e}");
            std::process::exit(1);
        }
    };

    let start = Instant::now();
    match initialise_quotes_table(&pool).await {
        Ok(v) => {
            let duration = start.elapsed();
            info!("Added {v} quotes to table in {:.1}ms", duration.as_millis())
        }
        Err(e) => {
            error!("Unable to init quotes table: {e}");
            std::process::exit(1);
        }
    };

    let daily_puzzle = match get_daily_puzzle_entity(&pool).await {
        Ok(v) => v,
        Err(e) => {
            error!("Unable to init quotes table: {e}");
            std::process::exit(1);
        }
    };
    debug!("{:?}", daily_puzzle);
    let daily_puzzle_response = DailyPuzzleResponse {
        author: daily_puzzle.author.clone(),
        cipher_quote: daily_puzzle.cipher_quote.clone(),
        date_string: daily_puzzle.date_string.clone(),
        day_number: daily_puzzle.day_number,
    };

    let daily_puzzle_cache: DailyPuzzleCache = Arc::new(RwLock::new(daily_puzzle));
    let daily_puzzle_response_cache: DailyPuzzleResponseCache =
        Arc::new(RwLock::new(daily_puzzle_response));

    info!("Starting HTTP server on port {}", PORT_NUMBER);
    HttpServer::new(move || {
        let cors_config = || {
            Cors::default()
                .allowed_origin("http://localhost:5173")
                .allowed_methods(vec!["GET", "POST"])
                .allowed_headers(vec![http::header::CONTENT_TYPE])
                .max_age(3600)
        };

        App::new()
            .wrap(Logger::default())
            .wrap(cors_config())
            .wrap(Compress::default())
            .service(
                web::scope("/daily")
                    .app_data(web::Data::new(daily_puzzle_response_cache.clone()))
                    .wrap(cors_config())
                    .configure(routes::daily_routes::init)
                    .service(
                        web::scope("/letter")
                            .app_data(web::Data::new(daily_puzzle_cache.clone()))
                            .wrap(cors_config())
                            .configure(routes::daily_letter_routes::init),
                    )
                    .service(
                        web::scope("/quote")
                            .app_data(web::Data::new(daily_puzzle_cache.clone()))
                            .wrap(cors_config())
                            .configure(routes::daily_quote_routes::init),
                    ),
            )
    })
    .bind(("0.0.0.0", PORT_NUMBER))?
    .run()
    .await
}
