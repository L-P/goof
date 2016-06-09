# goof REST API
This API is designed for consumption by the HTML/JS client. I don't
want to handle CalDAV client-side so a simple REST mapping will do.

All response follow the same format:
```json
{
    "Data": {},
    "Meta": {
        "Errors": []
    }
}
```

When `Meta.Errors` is not empty, `Data` will most likely be `null` and should
be discarded anyway.  
Empty arrays and objects may be `null` instead.


## GET /calendar/:id
`:id` is a calendar name.

### Query parameters
  * `start`: optional string, date eg. "2006-01-02"
  * `end`: optional string, date eg. "2006-02-30"

One of the two dates can be omitted. An empty range will return all events.

### Response
```json
{
    "Data": {
        "Calendar": {
            "Events": [ ]
        }
    },
    "Meta": {
        "Errors": null
    }
}
```
