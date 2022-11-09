/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// sceneCmd represents the scene command
var sceneCmd = &cobra.Command{
	Use:   "scene",
	Short: "Create simple scenes for CRUD operations for your entities",
	Long: `Create simple CRUD fxml scenes for java fx.
	`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("scene called")
		path, err := cmd.Flags().GetString("path")
		if err != nil {
			fmt.Println(err)
			return
		}
		if path == "." {
			path, err = cmd.Flags().GetString("root")
			if err != nil {
				fmt.Println(err)
				return
			}
		}

		fmt.Println(getResourceFolder(path))
	},
}

// tries to find the folder with the resources and the scenes
func getResourceFolder(rootPath string) string {
	var root string
	filepath.Walk(rootPath, func(path string, info fs.FileInfo, err error) error {
		if strings.HasSuffix(path, "fxml") {
			root = path
			return nil
		}
		return nil
	})

	fmt.Println(root)
	return rootPath
}

func init() {
	rootCmd.AddCommand(sceneCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sceneCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sceneCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	sceneCmd.Flags().StringP("path", "p", ".", "The path where the resources are located")
	sceneCmd.Flags().String("root", ".", "the root folder to start the search")
}
