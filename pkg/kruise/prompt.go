package kruise

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

func credentialPrompt(usernamePrompt string, passwordPrompt string) (username string, password string, err error) {
	validate := func(input string) error {
		return nil
	}
	u := promptui.Prompt{
		Label:    usernamePrompt,
		Validate: validate,
	}
	p := promptui.Prompt{
		Label:    passwordPrompt,
		Validate: validate,
		Mask:     '*',
	}
	c := promptui.Prompt{
		Label:    "Please confirm your password",
		Validate: validate,
		Mask:     '*',
	}
	username, err = u.Run()
	Fatal(err)
	password, err = p.Run()
	Fatal(err)
	validationPassword, err := c.Run()
	Fatal(err)
	if password != validationPassword {
		return "", "", fmt.Errorf("passwords do not match")
	}
	return username, password, nil
}
