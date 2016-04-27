[![Build Status](https://travis-ci.org/yukithm/json2csv.svg?branch=master)](https://travis-ci.org/yukithm/json2csv)

json2csv
========

Convert JSON into CSV. CSV header is generated from each key path as a JSON Pointer.
json2csv can be used as a library and command line tool.


Install
-------

Use `go get` or just download [binary releases](https://github.com/yukithm/json2csv/releases).

```
go get github.com/yukithm/json2csv/cmd/json2csv
```


Usage
-----

json2csv reads JSON content from STDIN or the file specified by the argument.

```sh
json2csv example.json
```

```sh
cat example.json | json2csv
```

Convert object array:

```json
[
    {
        "id": 1,
        "name": "foo",
        "favorites": {
            "fruits": "apple",
            "color": "red"
        }
    },
    {
        "id": 2,
        "name": "bar",
        "favorites": {
            "fruits": "orange"
        }
    },
    {
        "id": 3,
        "name": "baz",
        "favorites": {
            "fruits": "banana",
            "color": "yellow"
        }
    }
]
```

```sh
$ json2csv example1.json

/id,/name,/favorites/color,/favorites/fruits
1,foo,red,apple
2,bar,,orange
3,baz,yellow,banana
```

Convert object:

```json
{
    "status": 200,
    "result": [
        {
            "id": 1,
            "name": "foo"
        },
        {
            "id": 2,
            "name": "bar"
        }
    ]
}
```

```sh
$ json2csv example2.json

/status,/result/0/id,/result/0/name,/result/1/id,/result/1/name
200,1,foo,2,bar
```

Convert an array of the object:

Use `--path=<JSON Pointer>` option.

```sh
$ json2csv --path=/result example2.json

/id,/name
1,foo
2,bar
```

Transpose rows and columns:

```sh
$ json2csv --transpose example1.json

/id,1,2,3
/name,foo,bar,baz
/favorites/color,red,,yellow
/favorites/fruits,apple,orange,banana
```

### Header styles

By default, header is represented with JSON Pointer.
`--header-style=STYLE` option can change styles.

| style       | example        |
|-------------|----------------|
| jsonpointer | /foo/bar/0/baz |
| slash       | foo/bar/0/baz  |
| dot         | foo.bar.0.baz  |
| dot-bracket | foo.bar[0].baz |

Note: `slash` style similar to `jsonpointer` style, but `slash` style doesn't start with '/' and doesn't escape special characters ('/' and '~') defined in [RFC 6901](https://tools.ietf.org/html/rfc6901).

Note: `dot-bracket` style similar to `dot` style, but `dot-bracket` style uses square brackets for array indexes.


License
-------

MIT


Author
------

Yuki (@yukithm)
