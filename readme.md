# enum - command options enumerator

Download: `go install github.com/courtier/enum@0.0.1` or you can download one of the binaries

```
usage: enum [-h|--help] [-o|--options "<value>" [-o|--options "<value>" ...]]
            [-i|--option-file "<value>"] [-t|--threads <integer>] -c|--cmd
            "<value>" [-r|--repeat <integer>] [-f|--file "<value>"]

            Enumerate a command

Arguments:

  -h  --help         Print help information
  -o  --options      Options to enumerate through
  -i  --option-file  File to load options from
  -t  --threads      Number of threads to use when enumerating. Default: 1
  -c  --cmd          Command to enumerate, options will replace %o, multiple %o
                     are allowed, %-o will be replaced with %o
  -r  --repeat       Repeat command this many times, won't work if there are
                     options defined. Default: 1
  -f  --file         Output file, if not defined command outputs will be
                     printed to stdout
```

## Formatting
Options should be input through cli like so:

`-o POST,PUT,GET`

If you want commas in the options: `-o {"POST,", "PUT,", "GET,"}`

If you want to enumerate a through z or 0 through 9 with length of 2: `-o [az]{2}`

When inputting options through a file one set of options should be a single line, and three dashes must separate sets of options like so:
```
POST,PUT,GET
---
http://httpstat.us/200,http://google.com
---
{"POST,", "PUT,", "GET,"}
```
And of course you must have the same amount of option sets and %o's in your command.

## Example
`enum -c "echo %o" -o "[az]{2}"`
This will echo all 676 combinations of letters a through z.