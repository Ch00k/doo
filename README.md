[![test](https://github.com/Ch00k/doo/workflows/tests/badge.svg)](https://github.com/Ch00k/doo/actions)
[![codecov](https://codecov.io/gh/Ch00k/doo/branch/master/graphs/badge.svg)](https://codecov.io/github/Ch00k/doo)
[![Docker Build](https://img.shields.io/docker/cloud/build/ch00k/doo.svg)](https://hub.docker.com/r/ch00k/doo)


* [Running](#running)
* [Usage](#usage)
  * [List all entries](#list-all-entries)
  * [Show an entry](#show-an-entry)
  * [Create an entry](#create-an-entry)
  * [Update an entry](#update-an-entry)
  * [Complete an entry](#complete-an-entry)
  * [Add a comment to an entry](#add-a-comment-to-an-entry)
  * [Delete a comment from an entry](#delete-a-comment-from-an-entry)
  * [List all tags](#list-all-tags)
  * [Show a tag](#show-a-tag)
  * [Tag an entry](#tag-an-entry)
  * [Untag an entry](#untag-an-entry)
  * [Delete an entry](#delete-an-entry)
* [Running tests](#running-tests)
* [License](#license)


# doo

*doo* is a yet another ToDo app. It implements a REST API that allows creating entries, tagging them and commenting on
them, as well as marking them as completed.

*doo* is implemented in Go and uses [Gin](https://github.com/gin-gonic/gin) web framework and [GORM](https://gorm.io/)
ORM. Database backend of choice is [PostgreSQL](https://www.postgresql.org).

*doo* can run in Docker. Images are available [here](https://hub.docker.com/r/ch00k/doo).

## Running

The easiest way to run *doo* locally is with [docker-compose](https://docs.docker.com/compose). Make sure you have it
installed, and execute

```
$ make startall
```

This will start the database and the REST API server, which will be listening on http://localhost:8080.

The REST API server itself can also run outside of Docker, which can be used for debugging purposes. The database server
can be started with

```
$ make startdb
```

after which executing

```
$ make run
```

will start the REST API server.

## Usage

The examples below use [HTTPie](https://httpie.io). Check out [here](https://httpie.io/docs#json) for more info on how
it works with JSON.

### List all entries

```
$ http localhost:8080/entries
HTTP/1.1 200 OK
Content-Length: 348
Content-Type: application/json; charset=utf-8
Date: Wed, 17 Feb 2021 18:21:42 GMT

[
    {
        "Comments": [],
        "CompletedAt": 1613585629,
        "CreatedAt": 1613585401656,
        "ID": 2,
        "Tags": [
            {
                "Entries": null,
                "Name": "food"
            },
            {
                "Entries": null,
                "Name": "cooking"
            }
        ],
        "Text": "Buy Szechuan Sauce. Real quick!",
        "UpdatedAt": 1613585629348
    },
    {
        "Comments": [
            {
                "CreatedAt": 1613586152801,
                "EntryID": 1,
                "ID": 2,
                "Text": "Can't wait to taste it again!",
                "UpdatedAt": 1613586152801
            }
        ],
        "CompletedAt": null,
        "CreatedAt": 1613585231995,
        "ID": 1,
        "Tags": [],
        "Text": "Buy Szechuan Sauce",
        "UpdatedAt": 1613586152800
    }
]
```

### Show an entry

```
$ http localhost:8080/entries/1
HTTP/1.1 200 OK
Content-Length: 242
Content-Type: application/json; charset=utf-8
Date: Wed, 17 Feb 2021 18:23:52 GMT

{
    "Comments": [
        {
            "CreatedAt": 1613586152801,
            "EntryID": 1,
            "ID": 2,
            "Text": "Can't wait to taste it again!",
            "UpdatedAt": 1613586152801
        }
    ],
    "CompletedAt": null,
    "CreatedAt": 1613585231995,
    "ID": 1,
    "Tags": [],
    "Text": "Buy Szechuan Sauce",
    "UpdatedAt": 1613586152800
}
```

### Create an entry

```
$ http POST localhost:8080/entries Text="Buy Szechuan Sauce"
HTTP/1.1 201 Created
Content-Length: 135
Content-Type: application/json; charset=utf-8
Date: Wed, 17 Feb 2021 18:07:11 GMT

{
    "Comments": null,
    "CompletedAt": null,
    "CreatedAt": 1613585231995,
    "ID": 1,
    "Tags": null,
    "Text": "Buy Szechuan Sauce",
    "UpdatedAt": 1613585231995
}
```

Entries can also be created with tags

```
$ http POST localhost:8080/entries Text="Buy Szechuan Sauce" Tags:='[{"Name": "food"}, {"Name": "cooking"}]'
HTTP/1.1 201 Created
Content-Length: 197
Content-Type: application/json; charset=utf-8
Date: Wed, 17 Feb 2021 18:10:01 GMT

{
    "Comments": null,
    "CompletedAt": null,
    "CreatedAt": 1613585401656,
    "ID": 2,
    "Tags": [
        {
            "Entries": null,
            "Name": "food"
        },
        {
            "Entries": null,
            "Name": "cooking"
        }
    ],
    "Text": "Buy Szechuan Sauce",
    "UpdatedAt": 1613585401656
}
```

### Update an entry

```
$ http PUT localhost:8080/entries/2 Text="Buy Szechuan Sauce. Real quick!"
HTTP/1.1 200 OK
Content-Length: 148
Content-Type: application/json; charset=utf-8
Date: Wed, 17 Feb 2021 18:11:22 GMT

{
    "Comments": null,
    "CompletedAt": null,
    "CreatedAt": 1613585401656,
    "ID": 2,
    "Tags": null,
    "Text": "Buy Szechuan Sauce. Real quick!",
    "UpdatedAt": 1613585482084
}
```

### Complete an entry

```
$ http POST localhost:8080/entries/2/complete
HTTP/1.1 200 OK
Content-Length: 214
Content-Type: application/json; charset=utf-8
Date: Wed, 17 Feb 2021 18:13:49 GMT

{
    "Comments": [],
    "CompletedAt": 1613585629,
    "CreatedAt": 1613585401656,
    "ID": 2,
    "Tags": [
        {
            "Entries": null,
            "Name": "food"
        },
        {
            "Entries": null,
            "Name": "cooking"
        }
    ],
    "Text": "Buy Szechuan Sauce. Real quick!",
    "UpdatedAt": 1613585629348
}
```

### Add a comment to an entry

```
$ http POST localhost:8080/entries/1/comments Text="Can't wait to taste it again!"
HTTP/1.1 201 Created
Content-Length: 111
Content-Type: application/json; charset=utf-8
Date: Wed, 17 Feb 2021 18:16:26 GMT

{
    "CreatedAt": 1613585786227,
    "EntryID": 1,
    "ID": 1,
    "Text": "Can't wait to taste it again!",
    "UpdatedAt": 1613585786227
}
```

### Delete a comment from an entry

```
$ http DELETE localhost:8080/entries/1/comments/1
HTTP/1.1 200 OK
Content-Length: 0
Date: Wed, 17 Feb 2021 18:18:03 GMT
```

### List all tags

```
$ http localhost:8080/tags
HTTP/1.1 200 OK
Content-Length: 539
Content-Type: application/json; charset=utf-8
Date: Wed, 17 Feb 2021 18:26:53 GMT

[
    {
        "Entries": [
            {
                "Comments": null,
                "CompletedAt": 1613585629,
                "CreatedAt": 1613585401656,
                "ID": 2,
                "Tags": null,
                "Text": "Buy Szechuan Sauce. Real quick!",
                "UpdatedAt": 1613585629348
            }
        ],
        "Name": "food"
    },
    {
        "Entries": [
            {
                "Comments": null,
                "CompletedAt": 1613585629,
                "CreatedAt": 1613585401656,
                "ID": 2,
                "Tags": null,
                "Text": "Buy Szechuan Sauce. Real quick!",
                "UpdatedAt": 1613585629348
            }
        ],
        "Name": "cooking"
    },
    {
        "Entries": [
            {
                "Comments": null,
                "CompletedAt": null,
                "CreatedAt": 1613586403406,
                "ID": 3,
                "Tags": null,
                "Text": "Go to Gazorpazorp",
                "UpdatedAt": 1613586403406
            }
        ],
        "Name": "adventures"
    }
]
```

### Show a tag
```
$ http localhost:8080/tags/adventures
HTTP/1.1 200 OK
Content-Length: 168
Content-Type: application/json; charset=utf-8
Date: Wed, 17 Feb 2021 18:27:43 GMT

{
    "Entries": [
        {
            "Comments": null,
            "CompletedAt": null,
            "CreatedAt": 1613586403406,
            "ID": 3,
            "Tags": null,
            "Text": "Go to Gazorpazorp",
            "UpdatedAt": 1613586403406
        }
    ],
    "Name": "adventures"
}
```

### Tag an entry

```
$ echo '[{"Name": "spacetravel"}]' | http POST localhost:8080/entries/3/tag
HTTP/1.1 200 OK
Content-Length: 0
Date: Wed, 17 Feb 2021 18:32:32 GMT
```

### Untag an entry

```
$ http DELETE localhost:8080/entries/3/tags/spacetravel
HTTP/1.1 200 OK
Content-Length: 0
Date: Wed, 17 Feb 2021 18:33:41 GMT
```

### Delete an entry

```
$ http DELETE localhost:8080/entries/3
HTTP/1.1 200 OK
Content-Length: 0
Date: Wed, 17 Feb 2021 18:34:19 GMT
```

## Running tests

To run all tests execute

```
$ make test
```

## License

*doo* is UNLICENSEd. See [UNLICENSE](./UNLICENSE)
