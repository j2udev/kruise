package kruise

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/cqroot/prompt"
	"github.com/cqroot/prompt/input"
)

func normalInputPrompt(p string) string {
	return inputPrompt(p, input.EchoNormal)
}

func sensitiveInputPrompt(p string) string {
	return inputPrompt(p, input.EchoPassword)
}

func inputPrompt(p string, mode textinput.EchoMode) string {
	val, err := prompt.New().Ask(p).Input("", input.WithEchoMode(mode))
	if err != nil {
		Logger.Fatal(err)
	}
	return val
}
