# libpostal

The Mapzen Libpostal API is a ... around [the Libpostal C library](https://github.com/openvenues/libpostal) for parsing and normalizing street addresses. The Libpostal API has two endpoints that are exposed as HTTP `GET` requests, returning results as JSON data structures.

The 'parse' endpoint analyzes an address string and returns its component parts (street number, street name, city and so on). The 'expand' endpoint analyzes an address string and returns a set of normalized equivalent strings.

## GET /expand _?address=ADDRESS&apikey=APIKEY_

This endpoint accepts a single address parameter and expands it into one or more normalized forms suitable for geocoder queries.

|parameter|value|
| :--- | :--- |
| `api_key` | [get yours here](https://mapzen.com/developers) |
| `address` | 475+Sansome+St+San+Francisco+CA |

```
curl -s -X GET 'https://libpostal.mapzen.com/expand?address=475+Sansome+St+San+Francisco+CA&apikey=APIKEY' | python -mjson.tool
[
    "475 sansome saint san francisco california",
    "475 sansome saint san francisco ca",
    "475 sansome street san francisco california",
    "475 sansome street san francisco ca"
]
```

## GET /parse _?address=ADDRESS&apikey=APIKEY_

This endpoint accepts a single `address` parameter and parses it in to its components.

|parameter|value|
| :--- | :--- |
| `api_key` | [get yours here](https://mapzen.com/developers) |
| `address` | 475+Sansome+St+San+Francisco+CA |

```
curl -s -X GET 'https://libpostal.mapzen.com/parse?address=475+Sansome+St+San+Francisco+CA&apikey=APIKEY' | python -mjson.tool
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

By default both Libpostal and the Libpostal API return results a list of dictionaries, each containing a `label` and `value` key. This is because...

If you would prefer to return a simple dictionary with labels as keys and values as lists of possible strings you should append the `format=keys` parameter.

|parameter|value|
| :--- | :--- |
| `api_key` | [get yours here](https://mapzen.com/developers) |
| `address` | 475+Sansome+St+San+Francisco+CA |
| `format` | keys |

_Remember: all parameters should be URL encoded._

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
