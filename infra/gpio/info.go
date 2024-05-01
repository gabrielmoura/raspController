package gpio

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
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
	if _, err := exec.LookPath("lsgpio"); err == nil {

		output, err := exec.Command("lsgpio").Output()
		if err != nil {
			return nil, fmt.Errorf("error executing lsgpio: %v", err)
		}
		return parseLsGPIO(string(output))
	}

	if _, err := exec.LookPath("gpioinfo"); err == nil {
		output, err := exec.Command("gpioinfo").Output()
		if err != nil {
			return nil, fmt.Errorf("error executing lsgpio: %v", err)
		}
		return parseGPIOInfo(string(output))

	}
	log.Println("gpiod is required to get GPIO information")
	return nil, errors.New("lsgpio and gpioinfo command not found")
}

// parseGPIOInfo processes the output of the gpioinfo command
func parseGPIOInfo(output string) ([]GPIOInfo, error) {
	chipRegex := regexp.MustCompile(`(gpiochip\d+) - (\d+) lines:`)
	lineRegex := regexp.MustCompile(`\s*line\s*(\d+):\s*["']?(\w+)["']?\s*(\S+)\s*(input|output)\s*(active-high|active-low)\s*(?:\s*\[([\w\s]+)\])?`)

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
				Name:       chipMatch[1],
			}
		} else if lineMatch := lineRegex.FindStringSubmatch(line); lineMatch != nil {
			lineInfo := LineInfo{
				Number:    parseInt(lineMatch[1]),
				Name:      lineMatch[2],
				Function:  parseCleanString(lineMatch[3]),
				Direction: lineMatch[4],
			}

			lineInfo.Used = strings.Contains(lineMatch[6], "used")

			lineInfo.Flags = []string{lineMatch[5]}

			if currentChip != nil {
				currentChip.Lines = append(currentChip.Lines, lineInfo)
			}
		}
	}

	if currentChip != nil {
		gpioInfo = append(gpioInfo, *currentChip)
	}

	return gpioInfo, nil
}

// parseGPIOOutput processes the output of the lsgpio command
func parseLsGPIO(output string) ([]GPIOInfo, error) {
	chipRegex := regexp.MustCompile(`GPIO chip: (\S+), "(.+)", (\d+) GPIO lines`)
	lineRegex := regexp.MustCompile(`\s*line\s*(\d+):\s*"(.+)"\s*(used|unused)\s*\[([^\[\]]*)\]\s*`)

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

func parseCleanString(s string) string {
	return strings.Trim(s, "\"")
}
