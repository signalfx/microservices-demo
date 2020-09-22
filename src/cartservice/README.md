# cartservice

## Environment Variables

`REDIS_SPAN_ERROR_RATE`: float [0, 1], Percentage of redis spans that will be flagged as an error

`EXTERNAL_DB_NAME`: string, Name of external database
`EXTERNAL_DB_ACCESS_RATE`: float [0, 1], Percentage of redis spans that will be turned into external database spans
`EXTERNAL_DB_MAX_DURATION_MILLIS`: int, Artificial delay added to external database spans (mock value)
