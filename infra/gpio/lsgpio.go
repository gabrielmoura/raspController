package gpio

import (
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// Regular Expressions for parsing
var (
	chipRegex = regexp.MustCompile(`GPIO chip: (\S+), "(.+)", (\d+) GPIO lines`)
	lineRegex = regexp.MustCompile(`\s*line\s*(\d+):\s*"(.+)"\s*(used|unused)\s*\[([^\[\]]*)\]\s*`)
)

// GPIOInfo represents the information of a GPIO chip
type GPIOInfo struct {
	DeviceName string     `json:"deviceName"`
	Name       string     `json:"name"`
	Lines      []LineInfo `json:"lines"`
}

// LineInfo represents the information of a GPIO line
type LineInfo struct {
	Number    int      `json:"number"`
	Name      string   `json:"name"`
	Function  string   `json:"function"`
	Used      bool     `json:"used"`
	Direction string   `json:"direction"`
	Flags     []string `json:"flags"`
}

// GetGPIOInfo retrieves information about available GPIOs
func GetGPIOInfo() ([]GPIOInfo, error) {
	if _, err := exec.LookPath("lsgpio"); err != nil {
		return nil, errors.New("lsgpio command not found")
	}

	output, err := exec.Command("lsgpio").Output()
	if err != nil {
		return nil, fmt.Errorf("error executing lsgpio: %v", err)
	}

	return parseGPIOOutput(string(output))
}

// parseGPIOOutput processes the output of the lsgpio command
func parseGPIOOutput(output string) ([]GPIOInfo, error) {
	var gpioInfo []GPIOInfo
	var currentChip *GPIOInfo

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if chipMatch := chipRegex.FindStringSubmatch(line); chipMatch != nil {
			if currentChip != nil {
				gpioInfo = append(gpioInfo, *currentChip)
			}
			currentChip = &GPIOInfo{
				DeviceName: chipMatch[1],
				Name:       chipMatch[2],
			}
		} else if lineMatch := lineRegex.FindStringSubmatch(line); lineMatch != nil {
			lineInfo := LineInfo{
				Number:   parseInt(lineMatch[1]),
				Name:     lineMatch[2],
				Function: lineMatch[3],
				Used:     strings.Contains(lineMatch[4], "used"),
			}

			if strings.Contains(lineMatch[4], "output") {
				lineInfo.Direction = "output"
			} else {
				lineInfo.Direction = "input"
			}

			lineInfo.Flags = strings.Split(lineMatch[4], ", ")[1:]
			if currentChip != nil {
				currentChip.Lines = append(currentChip.Lines, lineInfo)
			}
		} else {
			// Handle non-matching lines (potentially gpiochip0)
			if currentChip != nil {
				// Create a simple LineInfo assuming it's part of the current chip
				lineInfo := LineInfo{
					Function:  "unused",
					Used:      false,
					Direction: "unknown",
					Flags:     []string{},
				}

				// You could attempt to extract at least the number and name here if possible

				currentChip.Lines = append(currentChip.Lines, lineInfo)
			}
		}
	}

	if currentChip != nil {
		gpioInfo = append(gpioInfo, *currentChip)
	}

	return gpioInfo, nil
}

// parseInt is a helper to parse strings as integers safely
func parseInt(s string) int {
	num, _ := strconv.Atoi(s)
	return num
}
