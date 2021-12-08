package helpers

var (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
)

// ColorRed wrap a text in red color
func ColorRed(s string) string {
	return string(colorRed) + s + string(colorReset)
}

// ColorWhite wrap a text in white color
func ColorWhite(s string) string {
	return string(colorWhite) + s + string(colorReset)
}

// ColorCyan wrap a text in cyan color
func ColorCyan(s string) string {
	return string(colorCyan) + s + string(colorReset)
}

// ColorPurple wrap a text in purple color
func ColorPurple(s string) string {
	return string(colorPurple) + s + string(colorReset)
}

// ColorGreen wrap a text in green color
func ColorGreen(s string) string {
	return string(colorGreen) + s + string(colorReset)
}

// ColorYellow wrap a text in yellow color
func ColorYellow(s string) string {
	return string(colorYellow) + s + string(colorReset)
}
