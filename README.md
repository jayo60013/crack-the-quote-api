# Crack the Quote API

API for Crack the quote game

## Endpoints
[x] `GET /daily`: Fetches ciphered quote, author, date string and day number 
[x] `POST /daily/letter/check`: Checks if the user's guess for a given cipher letter is correct. Returns true or false
[x] `POST /daily/letter/solve`: Returns the correct letter for a given cipher letter
[x] `POST /daily/quote/check`: Checks if the user's cipher map is correct. Returns true or false.
[ ] `POST /daily/quote/solve`: TODO

### Example Request
```bash
curl http://localhost:9100/daily
```

### Example Response
```json
{
  "Author": "Megan Whalen Turner,  The Queen of Attolia",
  "CipherQuote": "gu g ws ltz dwfo bu ltz xbej, gl gj mzawqjz ltzn robf sz jb fzyy, obl mzawqjz ltzn swrz sn sgoe qd ubv sz.",
  "DateString": "December 24th, 2024",
  "DayNumber": 2
}
```
## License

This project is licensed under the MIT License - see the LICENSE file for details.
