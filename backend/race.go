// https://gamedev.stackexchange.com/questions/84387/how-should-i-store-game-data-in-a-game-server

package main

// represents a quote and its respective progress to completion
import "strings"

type Race struct {
	quote            string
	quoteArray       []string
	typedWord        string
	currentWordIndex int
}

func newRace(quote string) *Race {
	return &Race{
		quote:            quote,
		quoteArray:       strings.SplitAfter(quote, " "),
		typedWord:        "",
		currentWordIndex: 0,
	}
}
