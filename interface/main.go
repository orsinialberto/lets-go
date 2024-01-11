package main

import "fmt"

type bot interface {
	getGreating() string
}

type englishBot struct{}

type spanishBot struct{}

func main() {
	eb := englishBot{}
	sb := spanishBot{}

	printGreeting(eb)
	printGreeting(sb)
}

func printGreeting(b bot) {
	fmt.Println(b.getGreating())
}

func (englishBot) getGreating() string {
	return "Hi There!"
}

func (spanishBot) getGreating() string {
	return "Hola!"
}
