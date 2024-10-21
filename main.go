package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"github.com/nsf/termbox-go"
)

type Sink struct {
	ID        string
	Name      string
	Volume    int
	IsDefault bool
}

func getDefaultSink() (string, error) {
	cmd := exec.Command("pactl", "get-default-sink")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	defaultSinkID := strings.TrimSpace(string(output))
	if defaultSinkID != "" {
		return defaultSinkID, nil
	}
	return "", fmt.Errorf("default sink not found")
}

func getSinks() ([]Sink, error) {
	cmd := exec.Command("pactl", "list", "short", "sinks")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(bytes.NewReader(output))
	var sinks []Sink
	defaultSinkID, err := getDefaultSink()
	if err != nil {
		return nil, err
	}

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			sink := Sink{
				ID:   parts[0],
				Name: parts[1],
			}
			if sink.Name == defaultSinkID {
				sink.IsDefault = true
			}
			volume, err := getVolume(sink.ID)
			if err != nil {
				sink.Volume = 0
			} else {
				sink.Volume = volume
			}
			sinks = append(sinks, sink)
		}
	}
	return sinks, nil
}

func getVolume(sinkID string) (int, error) {
	cmd := exec.Command("pactl", "get-sink-volume", sinkID)
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}
	
	scanner := bufio.NewScanner(bytes.NewReader(output))
	if scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "/")
		if len(parts) >= 2 {
			volStr := strings.TrimSpace(parts[1])
			volStr = strings.TrimSuffix(volStr, "%")
			vol, err := strconv.Atoi(volStr)
			if err != nil {
				return 0, err
			}
			return vol, nil
		}
	}
	return 0, fmt.Errorf("volume not found")
}

func setVolume(sinkID string, newVolume int) error {
	volumeStr := fmt.Sprintf("%d%%", newVolume)
	cmd := exec.Command("pactl", "set-sink-volume", sinkID, volumeStr)
	return cmd.Run()
}

func setDefaultSink(sinkID string) error {
	cmd := exec.Command("pactl", "set-default-sink", sinkID)
	if err := cmd.Run(); err != nil {
		return err
	}
	
	cmd = exec.Command("pactl", "list", "short", "sink-inputs")
	inputsOutput, err := cmd.Output()
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(bytes.NewReader(inputsOutput))
	for scanner.Scan() {
		inputID := strings.Fields(scanner.Text())[0]
		moveCmd := exec.Command("pactl", "move-sink-input", inputID, sinkID)
		if err := moveCmd.Run(); err != nil {
			return err
		}
	}
	return nil
}

func visualizeVolume(volume int) string {
	barLength := 20
	numBars := (volume * barLength) / 100
	bar := strings.Repeat("#", numBars)
	spaces := strings.Repeat(" ", barLength-numBars)
	return fmt.Sprintf("Volume : [%s%s] %d%%", bar, spaces, volume)
}

func displayList(sinks []Sink, selected int, maxNameLen int) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	
	header := "Select an audio output device (Arrow keys or j/k: navigate, Arrow keys or h/l: adjust volume, number: select, Enter: change, q: quit):"
	for i, ch := range header {
		termbox.SetCell(i, 0, ch, termbox.ColorYellow, termbox.ColorDefault)
	}
	
	for idx, sink := range sinks {
		y := idx + 2 
		var line string
		var color termbox.Attribute = termbox.ColorDefault
		
		if sink.IsDefault {
			line = "    * " 
			color = termbox.ColorBlue 
		} else {
			line = "      " 
		}

		if idx == selected {
			if sink.IsDefault {
				line = "-> * "
			} else { 
				line = "->   "
				color = termbox.ColorGreen
			}
		}
			
		formatName := fmt.Sprintf("%d %s", idx, sink.Name)
		if len(formatName) < maxNameLen {
			formatName += strings.Repeat(" ", maxNameLen-len(formatName))
		}
		line += formatName
		
		for x, ch := range line {
			termbox.SetCell(x, y, ch, color, termbox.ColorDefault)
		}
		
		volumeBar := visualizeVolume(sink.Volume)
		for x, ch := range volumeBar {
			termbox.SetCell(maxNameLen+5+x, y, ch, termbox.ColorDefault, termbox.ColorDefault)
		}
	}
	termbox.Flush()
}

func main() {
	err := termbox.Init()
	if err != nil {
		fmt.Println("Failed to initialize termbox:", err)
		return
	}
	defer termbox.Close()
	
	sinks, err := getSinks()
	if err != nil {
		fmt.Println("Error fetching sinks:", err)
		return
	}
	if len(sinks) == 0 {
		fmt.Println("No audio output devices found.")
		return
	}
	
	maxNameLen := 0
	for _, sink := range sinks {
		nameLen := len(fmt.Sprintf("%d %s", sinks[0].ID, sink.Name))
		if nameLen > maxNameLen {
			maxNameLen = nameLen
		}
	}
	selected := 0
	displayList(sinks, selected, maxNameLen)
	
mainloop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyCtrlC, termbox.KeyEsc:
				break mainloop
case termbox.KeyArrowDown:
				if selected < len(sinks)-1 {
					selected++
				}
			case termbox.KeyArrowUp:
				if selected > 0 {
					selected--
				}
			case termbox.KeyArrowLeft:
				if sinks[selected].Volume > 0 {
					newVolume := sinks[selected].Volume - 5
					if newVolume < 0 {
						newVolume = 0
					}
					err := setVolume(sinks[selected].ID, newVolume)
					if err == nil {
						sinks[selected].Volume = newVolume
					}
				}
			case termbox.KeyArrowRight:
				if sinks[selected].Volume < 100 {
					newVolume := sinks[selected].Volume + 5
					if newVolume > 100 {
						newVolume = 100
					}
					err := setVolume(sinks[selected].ID, newVolume)
					if err == nil {
						sinks[selected].Volume = newVolume
					}
				}
				
			default:
				switch ev.Ch {
			case 'j', 'J':
					if selected < len(sinks)-1 {
						selected++
					}
				case 'k', 'K':
					if selected > 0 {
						selected--
					}
				case 'h', 'H':
					if sinks[selected].Volume > 0 {
						newVolume := sinks[selected].Volume - 5
						if newVolume < 0 {
							newVolume = 0
						}
						err := setVolume(sinks[selected].ID, newVolume)
						if err == nil {
							sinks[selected].Volume = newVolume
						}
					}
				case 'l', 'L':
					if sinks[selected].Volume < 100 { 
						newVolume := sinks[selected].Volume + 5
						if newVolume > 100 {
							newVolume = 100
						}
						err := setVolume(sinks[selected].ID, newVolume)
						if err == nil {
							sinks[selected].Volume = newVolume
						}
					}
				case 'q', 'Q':
				break mainloop
default:
					if ev.Ch >= '0' && ev.Ch <= '9' {
						num := int(ev.Ch - '0')
						if num < len(sinks) {
							selected = num
						}
					}
				}
				
				if ev.Key == termbox.KeyEnter {
					err := setDefaultSink(sinks[selected].ID)
					if err == nil {
						
						for i := range sinks {
							if sinks[i].ID == sinks[selected].ID {
								sinks[i].IsDefault = true
							} else {
								sinks[i].IsDefault = false
							}
						}
					}
				}
			}
		case termbox.EventError:
			fmt.Println("Termbox event error:", ev.Err)
			
		}
		
		sinks, err = getSinks()
		if err != nil {
			fmt.Println("Error fetching sinks:", err)
			break
		}

		maxNameLen = 0
		for _, sink := range sinks {
			nameLen := len(fmt.Sprintf("%d %s", sinks[0].ID, sink.Name))
			if nameLen > maxNameLen {
				maxNameLen = nameLen
			}
		}

		if selected >= len(sinks) {
			selected = len(sinks) - 1
		}

		displayList(sinks, selected, maxNameLen)
	}
}

