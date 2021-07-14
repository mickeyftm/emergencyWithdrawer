package view

func centerText(text string, length int) string {
	if len(text) >= length {
		return text
	}

	textRune := []rune(text)
	whitespaces := (length - len(text)) / 2
	whitespace := " "
	for i := 0; i < whitespaces-1; i++ {
		textRune = append([]rune(whitespace), textRune...)
	}

	return string(textRune)
}
