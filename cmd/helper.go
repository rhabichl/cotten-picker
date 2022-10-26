package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func getClassIndexByName(name string) int {
	for i, v := range classes {
		if v.Name == name {
			return i
		}
	}
	return -1
}

func getPackageFomUsers(path string) {
	var packages []string
	err := filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			f := strings.Split(info.Name(), ".")
			if f[len(f)-1] == "java" {
				p, err := getPackageName(path)
				if err != nil {
					fmt.Println(err)
					return nil
				}
				if len(p) == 0 {
					return nil
				}
				packages = append(packages, p)

			}
		}
		return nil
	})

	if err != nil {
		fmt.Println(err)
	}
	pc := promptContent{
		"Select the package you want to use or use a custom one",
		"Select the package",
	}

	packageName = promptGetSelect(pc, packages)

}

func getPackageName(path string) (string, error) {
	if len(path) == 0 {
		return "", fmt.Errorf("path can't be empty")
	}
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	sContent := string(content)

	lines := strings.Split(sContent, "\n")
	if !strings.HasPrefix(lines[0], "package") {
		return "", nil
	}
	name := strings.Split(lines[0], " ")[1]
	name = strings.ReplaceAll(name, ";", "")

	return name, nil
}
