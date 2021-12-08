package helpers

import "fmt"

// Message is used to arbitrary print or OK or NOK message
func Message(message string, err error) {
	if err != nil {
		MessageNOK(message)
		MessageError(err)
	} else {
		MessageOK(message)
	}
}

// MessageOK print a OK message
// in a specific format
func MessageOK(message string) {
	fmt.Printf("[%s] %s\n", ColorGreen("OK"), message)
}

// MessageNOK print a OK message
// in a specific format
func MessageNOK(message string) {
	fmt.Printf("[%s] %s\n", ColorRed("KO"), message)
}

// MessageProcessing print a message
// that indicate a in pending process
func MessageProcessing(message string) {
	fmt.Printf("[%s] %s\n", ColorWhite(".."), ColorCyan(message))
}

// MessageWarning print a warning message
func MessageWarning(message string) {
	fmt.Printf("[%s] %s\n", ColorYellow("!!"), ColorYellow(message))
}

// SpaceLine ...
func SpaceLine() {
	fmt.Print("--------\n")
}

// MessageError print an error message
func MessageError(err error) {
	fmt.Printf("%s\n", ColorRed(err.Error()))
}
