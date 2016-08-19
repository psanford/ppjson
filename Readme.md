ppjson
======

My replacement for json_pp.

The main impetus for writing this was to have a tool that produced non-randomized ordering of hashes. Occasionally I want to diff json files and that is quite painful when using json_pp.

I added some flags to allow from reading and writing to things other than std{in,out}. And I added a gnu sed like '-i' flag to do in place updates (although its not really like sed because it doesn't take an optional suffix).

Example
-------

```
# help
$ ppjson --help
Usage of ppjson:
  -i    update file inplace
  -in string
        input file (defaults to stdin)
  -out string
        output file (defaults to stdout)
  -replace
        update file inplace


# stdin array
$ echo '[1,2,3]' | ppjson
[
  1,
  2,
  3
]


# stdin hash
$ echo '{"a":1,"b":2,"c":3}' | ppjson
{
  "a": 1,
  "b": 2,
  "c": 3
}


# file input
$ echo '{"a":1,"b":2,"c":3}' > /tmp/input
$ ppjson /tmp/input
{
  "a": 1,
  "b": 2,
  "c": 3
}


# in-place update
$ echo '{"a":1,"b":2,"c":3}' > /tmp/input
$ ppjson -replace /tmp/input
$ cat /tmp/input
{
  "a": 1,
  "b": 2,
  "c": 3
}
```
