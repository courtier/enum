# enum - command options enumerator

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