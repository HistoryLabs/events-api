# Events API

An API that fetches historical events from Wikipedia. Built to provide data to "Today In History" style apps (such as the [Historian Discord Bot](https://github.com/HistoryLabs/historian-bot)).

## To use

You can find the most recent version of the API deployed at https://events.historylabs.io/.

## Endpoints

### GET /date

```http
GET /date
```

#### Parameters

|    Parameter    |       Type       |  Default |                                        Description                                       |
|:---------------:|:----------------:|:--------:|:----------------------------------------------------------------------------------------:|
|     `month`     |       1-12       | Required |                    The month to find events for (paired with `date`).                    |
|      `date`     |       1-31       | Required |                         The day of the month to find events for.                         |
| `minYear` (opt) | `-500` to `2022` |  `-500`  |  Only includes events that occurred after the given year (negative numbers mean BC/BCE). |
| `maxYear` (opt) | `-500` to `2022` |  `2022`  | Only includes events that occurred before the given year (negative numbers mean BC/BCE). |

#### Return DTO

```json
{
    "totalResults": 56,
    "sourceUrl": "https://en.wikipedia.org/wiki/March_2",
    "events": [
        {
            "year": "537",
            "yearInt": 537,
            "event": "Siege of Rome: The Ostrogoth army under king Vitiges begins the siege of the capital. Belisarius conducts a delaying action outside the Flaminian Gate; he and a detachment of his bucellarii are almost cut off."
        },
        {
            "year": "986",
            "yearInt": 986,
            "event": "Louis V becomes the last Carolingian king of West Francia after the death of his father, Lothaire."
        },
        {
            "year": "1476",
            "yearInt": 1476,
            "event": "Burgundian Wars: The Old Swiss Confederacy hands Charles the Bold, Duke of Burgundy, a major defeat in the Battle of Grandson in Canton of Neuch\âtel."
        },
        {
            "year": "1484",
            "yearInt": 1484,
            "event": "The College of Arms is formally incorporated by Royal Charter signed by King Richard III of England."
        },
        "..."
    ]
}
```

### GET /year

```http
GET /year
```

#### Parameters

|     Parameter    |        Type       | Default Value |               Description               |
|:----------------:|:-----------------:|:-------------:|:---------------------------------------:|
| `onlyDate` (opt) | `true` or `false` |    `false`    | Will only give events that have a date. |

#### Return DTO

```json
{
    "totalResults": 15,
    "sourceUrl": "https://en.wikipedia.org/wiki/AD_1500",
    "events": [
        {
            "date": "January 5",
            "event": "Duke Ludovico Sforza recaptures Milan, but is soon driven out again by the French."
        },
        {
            "date": "February 17",
            "event": "Battle of Hemmingstedt: The Danish army fails to conquer the peasants' republic of Dithmarschen."
        },
        {
            "date": "July 14",
            "event": "The Muscovites defeat the Lithuanians and the Poles in the Battle of Vedrosha."
        },
        {
            "date": "August 10",
            "event": "Diogo Dias discovers an island which he names St Lawrence (after the saint's day on which it was first sighted), later to be known as Madagascar."
        },
        {
            "date": "November 11",
            "event": "Treaty of Granada: Louis XII of France and Ferdinand II of Aragon agree to divide the Kingdom of Naples between them."
        },
        {
            "date": "December 31",
            "event": "The last incunable is printed in Venice."
        },
        {
            "date": "",
            "event": "Europe's population is estimated at 56.7 million people (Spielvogel)."
        },
        "..."
    ]
