use std::collections::HashMap;

#[derive(Debug)]
pub struct DailyPuzzle {
    pub cipher_quote: String,
    pub author: String,
    pub date_string: String,
    pub day_number: u16,
    pub cipher_map: HashMap<char, char>,
}
