package user

import (
	"math"
	"strings"
)

func compass(direction float64) string {
	direction = math.Mod(direction+360, 361)
	direction = direction * 2 / 10

	comp := []string{
		"§4S §r", "| ", "| ", "| ", "| ", "| ", "| ", "| ", "| ", "§4S", "W §r",
		"| ", "| ", "| ", "| ", "| ", "| ", "| ", "§4W §r", "| ", "| ", "| ", "| ",
		"| ", "| ", "| ", "| ", "§4N", "§4W §r", "| ", "| ", "| ", "| ", "| ", "| ",
		"| ", "§4N §r", "| ", "| ", "| ", "| ", "| ", "| ", "| ", "§4N", "§4E §r", "| ",
		"| ", "| ", "| ", "| ", "| ", "| ", "§4E §r", "| ", "| ", "| ", "| ", "| ",
		"| ", "| ", "§4S", "§4E §r", "| ", "| ", "| ", "| ", "| ", "| ", "| ", "§4S §r",
		"| ", "| ", "| ", "| ", "| ", "| ", "| ", "| ", "§4S", "§4W §r", "| ", "| ",
		"| ", "| ", "| ", "| ", "| ", "§4W §r", "| ", "| ", "| ", "| ", "| ", "| ",
		"| ", "| ", "§4N", "§4W §r", "| ", "| ", "| ", "| ", "| ", "| ", "| ", "§4N §r",
		"| ", "| ", "| ", "| ", "| ", "| ", "| ", "§4N", "§4E §r", "| ", "| ", "| ", "| ",
		"| ", "| ", "| ", "§4E §r", "| ", "| ", "| ", "| ", "| ", "| ", "| ", "§4S",
		"§4E §r", "| ", "| ", "| ", "| ", "| ", "| ", "| ", "§4S §r", "| ", "| ", "| ",
		"| ", "| ", "| ", "| ", "| ", "§4S", "W §r", "| ", "| ", "| ", "| ", "| ", "| ",
		"| ", "§4W §r", "| ", "| ", "| ", "| ", "| ", "| ", "| ", "§4N", "§4W §r", "| ",
		"| ", "| ", "| ", "| ", "| ", "| ", "§4N §r", "| ", "| ", "| ", "| ", "| ", "| ",
		"| ", "§4N", "§4E §r", "| ", "| ", "| ", "| ", "| ", "| ", "| ", "§4E §r", "| ",
		"| ", "| ", "| ", "| ", "| ", "| ", "§4S", "§4E §r", "| ", "| ", "| ", "| ", "| ", "| ",
	}
	direction += 70

	start := int(direction) - int(math.Floor(float64(25)/2))

	if start < 0 {
		start = 0
	}

	end := start + 25
	if end > len(comp) {
		end = len(comp)
	}

	slice := comp[start:end]
	return strings.Join(slice, "")
}
