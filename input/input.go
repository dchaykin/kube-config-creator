package input

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

const RegExpNames = `^[a-zA-Z0-9-_]+$`

type UserInput struct {
	serviceAccountName string
	roleName           string
	namespace          string
}

func (ui *UserInput) GetServiceAccountName() string {
	return ui.serviceAccountName
}

func (ui *UserInput) GetNamespace() string {
	if ui.namespace == "" {
		return "default"
	}
	return ui.namespace
}

func (ui *UserInput) GetRoleName() string {
	return ui.roleName
}

func (ui *UserInput) ReadServiceAccountName() error {
	return getInputFromUser("Service account", &ui.serviceAccountName, true)
}

func (ui *UserInput) ReadRoleName() error {
	return getInputFromUser("Role", &ui.roleName, true)
}

func (ui *UserInput) ReadNamespace() error {
	return getInputFromUser("Namespace", &ui.namespace, false)
}

func getInputFromUser(subject string, value *string, mandatory bool) (err error) {
	text := ""
	for matched := false; !matched; {
		scanner := bufio.NewScanner(os.Stdin)
		text = ""
		if scanner.Scan() {
			text = scanner.Text()
		}
		if len(text) == 0 {
			if mandatory {
				fmt.Printf("%s cannot be empty\n", subject)
				continue
			}
			break
		}
		matched, err = regexp.MatchString(RegExpNames, text)
		if err != nil {
			return err
		}
		if !matched {
			fmt.Printf("Names can contain letters, numbers, hyphens and underscores only.\nTry again: ")
			continue
		}
	}
	*value = text
	return nil
}
