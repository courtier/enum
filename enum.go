package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/akamensky/argparse"
	cartesianLibrary "github.com/schwarmco/go-cartesian-product"
)

//TODO
//add generative functions such as --opt "[a-z]{3}" -> would generate all combinations of a-z with 3 chars

var (
	smallAZ         = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "y", "z", "x"}
	capitalAZ       = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "Y", "Z", "X"}
	numbersZeroNine = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
)

var (
	options       []string
	actualOptions [][]string
	optionFile    string
	threads       int
	command       string
	repeat        int
	outputFile    string
)

func main() {
	parser := argparse.NewParser("enum", "Enumerate a command")

	optionsPointer := parser.StringList("o", "options", &argparse.Options{Required: false, Help: "Options to enumerate through"})
	optionFilePointer := parser.String("i", "option-file", &argparse.Options{Required: false, Help: "File to load options from"})
	threadsPointer := parser.Int("t", "threads", &argparse.Options{Required: false, Help: "Number of threads to use when enumerating", Default: 1})
	commandPointer := parser.String("c", "cmd", &argparse.Options{Required: true, Help: "Command to enumerate, options will replace %o, multiple %o are allowed, %-o will be replaced with %o"})
	repeatPointer := parser.Int("r", "repeat", &argparse.Options{Required: false, Help: "Repeat command this many times, won't work if there are options defined", Default: 1})
	outputFilePointer := parser.String("f", "file", &argparse.Options{Required: false, Help: "Output file, if not defined command outputs will be printed to stdout"})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}

	//dereference da pointers into actual variables
	options = *optionsPointer
	optionFile = *optionFilePointer
	threads = *threadsPointer
	command = *commandPointer
	repeat = *repeatPointer
	outputFile = *outputFilePointer

	if len(optionFile) > 0 && len(options) > 0 {
		log.Fatal("Options must only be defined once, either through file or as argument.")
	}
	if threads < 1 || repeat < 1 {
		log.Fatal("Threads or repeat cannot be lower than 1.")
	}

	if len(optionFile) > 0 {
		//read file
		fileContent, err := os.ReadFile(optionFile)
		if err != nil {
			log.Fatal("Error reading options file. ", err)
		}
		fileStr := string(fileContent)
		if len(strings.Split(fileStr, "\n")) < 1 {
			log.Fatal("Empty options file.")
		}
		options = []string{}
		//read file content line by line
		for _, line := range strings.Split(fileStr, "\n") {
			//--- splits options from one another
			if line != "---" {
				options = append(options, line)
			}
		}
	}

	//make sure we have options laoaded
	if strings.Contains(command, "%o") && len(options) < 1 {
		log.Fatal("No options were loaded, but command requires options.")
	} else if !strings.Contains(command, "%o") && len(options) > 0 {
		log.Fatal("Options were loaded, but command requires no options.")
	} else if strings.Count(command, "%o") != len(options) {
		log.Fatal("Uneven number of options loaded and options in command found.")
	}

	workAmount := 0
	workCommands := []string{}

	if len(options) > 0 {
		//extract options into 2d array
		for _, option := range options {
			if option[0] == '{' {
				//argparse does some weird things and we do some weird things
				option = option[1 : len(option)-1]
				newArr := strings.Split(option, ", ")
				actualOptions = append(actualOptions, newArr)
			} else if option[0] == '[' && option[len(option)-1] == '}' {
				varLength := strings.Split(strings.Split(option, "{")[1], "}")[0]
				lengths := []int{}
				if strings.Contains(varLength, ",") {
					for _, newLength := range strings.Split(varLength, ",") {
						length, err := strconv.Atoi(newLength)
						if err != nil {
							log.Fatal("Error processing option.")
						}
						lengths = append(lengths, length)
					}
				} else {
					length, err := strconv.Atoi(varLength)
					if err != nil {
						log.Fatal("Error processing option.")
					}
					lengths = append(lengths, length)
				}
				if len(lengths) > 0 {
					for _, setLength := range lengths {
						if strings.Contains(option, "az") {
							actualOptions = append(actualOptions, generateSetLength(setLength, smallAZ))
						} else if strings.Contains(option, "AZ") {
							actualOptions = append(actualOptions, generateSetLength(setLength, capitalAZ))
						} else if strings.Contains(option, "09") {
							actualOptions = append(actualOptions, generateSetLength(setLength, numbersZeroNine))
						}
					}
				}
			} else {
				actualOptions = append(actualOptions, strings.Split(option, ","))
			}
		}

		if len(actualOptions) > 1 {
			interfaceArr := make([][]interface{}, 0)
			for _, arr := range actualOptions {
				iArr := make([]interface{}, 0)
				for _, el := range arr {
					iArr = append(iArr, el)
				}
				interfaceArr = append(interfaceArr, iArr)
			}
			cartesianChan := cartesianLibrary.Iter(interfaceArr...)
			for product := range cartesianChan {
				workAmount++
				formattedCommand := command
				for _, option := range product {
					formattedCommand = strings.Replace(formattedCommand, "%o", option.(string), 1)
				}
				workCommands = append(workCommands, formattedCommand)
			}
		} else if len(actualOptions) == 1 {
			workAmount = len(actualOptions[0])
			for _, option := range actualOptions[0] {
				formattedCommand := command
				formattedCommand = strings.Replace(formattedCommand, "%o", option, 1)
				workCommands = append(workCommands, formattedCommand)
			}
		}
	} else if repeat > 0 {
		workAmount = repeat
		for i := 0; i < repeat; i++ {
			workCommands = append(workCommands, command)
		}
	}

	jobs := make(chan string, workAmount)
	results := make(chan string, workAmount)

	for i := 0; i < threads; i++ {
		go worker(jobs, results)
	}

	for i := 0; i < workAmount; i++ {
		jobs <- workCommands[i]
	}
	close(jobs)

	doneCounter := 0
	var saveResults strings.Builder
	for result := range results {
		if outputFile == "" {
			fmt.Println(result)
		} else {
			saveResults.WriteString(result)
		}
		doneCounter++
		if doneCounter == workAmount {
			break
		}
	}
	close(results)
	if outputFile != "" {
		err = appendStringToFile(outputFile, saveResults.String(), true)
		if err != nil {
			log.Fatal("Error saving results to file. ", err)
		}
	}
}

func worker(jobs, results chan string) {
	for job := range jobs {
		executeCommand(job, results)
	}
}

func executeCommand(command string, results chan string) {
	name, args := strings.Fields(command)[0], strings.Fields(command)[1:]
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal("Couldn't run command. ", err)
	}
	results <- string(out)
}

func appendStringToFile(fileName string, line string, overwrite bool) error {
	sBytes := []byte(line + "\n")

	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := file.Write(sBytes); err != nil {
		return err
	}
	file.Sync()
	return nil
}

func generateSetLength(length int, source []string) []string {
	interfaceArr := make([][]interface{}, 0)
	for i := 0; i < length; i++ {
		iArr := make([]interface{}, 0)
		for _, el := range source {
			iArr = append(iArr, el)
		}
		interfaceArr = append(interfaceArr, iArr)

	}
	result := []string{}
	cartesianChan := cartesianLibrary.Iter(interfaceArr...)
	for product := range cartesianChan {
		var pro strings.Builder
		for _, str := range product {
			pro.WriteString(str.(string))
		}
		result = append(result, pro.String())
	}
	return result
}
