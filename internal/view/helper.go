package view

func centerText(text string, lenght int) string {
	if len(text) >= lenght {
		return text
	}

	textRune := []rune(text)
	whitespaces := (lenght - len(text)) / 2
	whitespace := " "
	for i := 0; i < whitespaces-1; i++ {
		textRune = append([]rune(whitespace), textRune...)
	}

	return string(textRune)
}
