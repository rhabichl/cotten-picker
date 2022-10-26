package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
)

func promptGetInput(pc promptContent) string {
	validate := func(input string) error {
		if len(input) <= 0 {
			return errors.New(pc.errorMsg)
		}
		return nil
	}

	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }} ",
		Valid:   "{{ . | green }} ",
		Invalid: "{{ . | red }} ",
		Success: "{{ . | bold }} ",
	}

	prompt := promptui.Prompt{
		Label:     pc.label,
		Templates: templates,
		Validate:  validate,
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}
	return result
}

func getActionPromt(pc promptContent, items []string) string {
	promt := promptui.Select{
		Label: pc.label,
		Items: items,
	}
	for {
		index, result, err := promt.Run()
		if err != nil {
			return ""
		}
		fmt.Println(result)
		if index == 0 {
			createClasse()
		} else if index == 1 {
			updateClasses()
		} else if index == 2 {
			writeClasses()
			break
		}
	}

	return ""
}

func promptGetSelect(pc promptContent, items []string) string {
	index := -1
	var result string
	var err error

	for index < 0 {
		prompt := promptui.SelectWithAdd{
			Label:    pc.label,
			Items:    items,
			AddLabel: "Other",
		}

		index, result, err = prompt.Run()

		if index == -1 {
			items = append(items, result)
		}
	}

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Input: %s\n", result)

	return result
}
