use serde::Serialize;

#[derive(Serialize)]
pub struct DailyPuzzleResponse {
    pub author: String,
    pub cipher_quote: String,
    pub date_string: String,
    pub day_number: u16,
}
