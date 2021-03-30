# enum - command options enumerator
## enumerate a command with options

Example:
enum -opt="{POST,PUT,GET}" -t=3 -cmd="curl -X %o https://google.com" -file="output.txt"
enum -opt-file="options.txt" -t=3 -cmd="curl -X %o https://google.com" -file="output.txt"
enum -threads=10 -repeat=100 -cmd="curl -X GET https://google.com" //output to stdout

you can even stack it i think?