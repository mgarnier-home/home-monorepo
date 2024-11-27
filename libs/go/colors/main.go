package colors

import (
	"fmt"
	"math"
)

// Convert HSL to RGB and return as a hex color code
func HSLToHex(h, s, l float64) string {
	c := (1 - math.Abs(2*l-1)) * s
	x := c * (1 - math.Abs(math.Mod(h/60.0, 2)-1))
	m := l - c/2

	var r, g, b float64
	switch {
	case h >= 0 && h < 60:
		r, g, b = c, x, 0
	case h >= 60 && h < 120:
		r, g, b = x, c, 0
	case h >= 120 && h < 180:
		r, g, b = 0, c, x
	case h >= 180 && h < 240:
		r, g, b = 0, x, c
	case h >= 240 && h < 300:
		r, g, b = x, 0, c
	case h >= 300 && h < 360:
		r, g, b = c, 0, x
	}

	// Convert RGB to 0-255 and apply offset
	r = (r + m) * 255
	g = (g + m) * 255
	b = (b + m) * 255

	// Format as hexadecimal
	return fmt.Sprintf("#%02X%02X%02X", int(r), int(g), int(b))
}

// Generate an HSL-based hex color from a string
func GenerateHexColor(input string) string {
	hash := fnv32a(input)

	// Map hash to HSL values
	hue := float64(hash % 360)                 // Hue: [0, 360)
	saturation := 0.6 + float64(hash%40)/100.0 // Saturation: [(0.x * 100)%, 100%)
	lightness := 0.5 + float64(hash%40)/100.0  // Lightness: [(0.x * 100)%, 70%)

	// Convert HSL to hex
	return HSLToHex(hue, saturation, lightness)
}

// A simple non-cryptographic hash function
func fnv32a(input string) uint32 {
	const prime = 16777619
	var hash uint32 = 2166136261
	for i := 0; i < len(input); i++ {
		hash ^= uint32(input[i])
		hash *= prime
	}
	return hash
}
