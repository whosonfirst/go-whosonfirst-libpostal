# The Mapzen Libpostal API

The Mapzen Libpostal API is a ... around [the Libpostal C library](https://github.com/openvenues/libpostal) for parsing and normalizing street addresses in to data structures. The Libpostal API has two endpoints that are exposed as HTTP `GET` requests, returning results as JSON data structures.

## GET /expand _?address=ADDRESS&apikey=APIKEY_

|parameter|value|
| :--- | :--- |
| `address` | 475+Sansome+St+San+Francisco+CA |
| `api_key` | [get yours here](https://mapzen.com/developers) |

The `/expand` endpoint analyzes an address string and returns a set of normalized equivalent strings.

```
curl -s -X GET 'https://libpostal.mapzen.com/expand?address=475+Sansome+St+San+Francisco+CA&api_key=APIKEY' | python -mjson.tool
[
    "475 sansome saint san francisco california",
    "475 sansome saint san francisco ca",
    "475 sansome street san francisco california",
    "475 sansome street san francisco ca"
]
```

## GET /parse _?address=ADDRESS&apikey=APIKEY_

|parameter|value|
| :--- | :--- |
| `address` | 475+Sansome+St+San+Francisco+CA |
| `api_key` | [get yours here](https://mapzen.com/developers) |

The `/parse` endpoint analyzes an address string and returns its component parts (street number, street name, city and so on). 

```
curl -s -X GET 'https://libpostal.mapzen.com/parse?address=475+Sansome+St+San+Francisco+CA&api_key=APIKEY' | python -mjson.tool
[
    {
        "label": "house_number",
        "value": "475"
    },
    {
        "label": "road",
        "value": "sansome st"
    },
    {
        "label": "city",
        "value": "san francisco"
    },
    {
        "label": "state",
        "value": "ca"
    }
]
```

By default both [Libpostal](https://github.com/openvenues/libpostal) and the Libpostal API return results a list of dictionaries, each containing a `label` and `value` key. This is because...

|parameter|value|
| :--- | :--- |
| `address` | 475+Sansome+St+San+Francisco+CA |
| `api_key` | [get yours here](https://mapzen.com/developers) |
| `format` | keys |

If you would prefer to haves results returned instead as a simple dictionary with labels as keys and values as lists of possible strings you should append the `format=keys` parameter.

```
curl -s -X GET 'https://libpostal.mapzen.com/parse?address=475+Sansome+St+San+Francisco+CA&format=keys' | python -mjson.tool
{
    "city": [
        "san francisco"
    ],
    "house_number": [
        "475"
    ],
    "road": [
        "sansome st"
    ],
    "state": [
        "ca"
    ]
}
```
