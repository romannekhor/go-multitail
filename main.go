package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var commands arrayFlags
var outputColors arrayFlags
var labels arrayFlags

// KILLSIG is a kill signal for process execution go-routines
const KILLSIG string = "KILLSIG"

// ProcessSignal is a type to represent instances of signal (e.g. sigOk and sigKill)
type ProcessSignal int

const (
	sigOk ProcessSignal = iota
	sigKill
)

// OutputLine is a unit of process output
type OutputLine struct {
	signal     ProcessSignal
	line       string
	processLbl string
}

func execute(label string, command []string, output chan<- OutputLine) {
	tailCmd := exec.Command(command[0], command[1:]...)

	tailOut, _ := tailCmd.StdoutPipe()
	tailCmd.Start()

	scanner := bufio.NewScanner(tailOut)
	for scanner.Scan() {
		line := scanner.Text()

		lineStruct := new(OutputLine)
		lineStruct.line = line
		lineStruct.processLbl = label
		lineStruct.signal = sigOk

		output <- *lineStruct
	}
	tailCmd.Wait()

	lineStruct := new(OutputLine)
	lineStruct.processLbl = label
	lineStruct.signal = sigKill
	output <- *lineStruct
}

func main() {
	colorNameToColor := map[string]color.Attribute{
		"red":     color.FgRed,
		"green":   color.FgGreen,
		"yellow":  color.FgYellow,
		"blue":    color.FgBlue,
		"magenta": color.FgMagenta,
		"cyan":    color.FgCyan,
	}

	var lblToColor = make(map[string]*color.Color)

	flag.Var(&commands, "cmd", "Command to execute")
	flag.Var(&labels, "l", "Command labels")
	flag.Var(&outputColors, "color", "Output color")
	flag.Parse()

	fmt.Printf("Commands: %#v\n", commands)
	fmt.Printf("Labels: %#v\n", labels)

	if len(commands) != len(labels) {
		fmt.Println("Error. Number of commands doesn't match the number of labels provided")
		os.Exit(1)
	}

	output := make(chan OutputLine)

	for i := 0; i < len(commands); i++ {
		lbl := labels[i]
		cmdStr := commands[i]
		cmd := strings.Split(cmdStr, " ")

		colorName := outputColors[i]

		outputColor, ok := colorNameToColor[colorName]

		if !ok {
			fmt.Printf("Error. Unknown color '%s'", colorName)
			os.Exit(1)
		}

		lblToColor[lbl] = color.New(outputColor)

		go execute(lbl, cmd, output)
	}

	var procsRunning = len(commands)
	var prevLabel string

	for lineStruct := range output {
		if prevLabel != lineStruct.processLbl {
			color.Red("======================================================")
		}

		if lineStruct.signal == sigKill {
			procsRunning--

			if procsRunning <= 0 {
				break
			}

		} else {
			colorObj, _ := lblToColor[lineStruct.processLbl]

			colorObj.Printf("%s :: %s\n", lineStruct.processLbl, lineStruct.line)
		}

		prevLabel = lineStruct.processLbl
	}

	color.Blue("BYE!")

}
