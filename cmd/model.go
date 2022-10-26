/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/rhabichl/cotten-picker/models"
	"github.com/spf13/cobra"
)

var packageName string

var classes []models.Class

var path string

// modelCmd represents the model command
var modelCmd = &cobra.Command{
	Use:   "model",
	Short: "To ineract with models",
	Long:  "Please give a path where the models should be made",
	Run: func(cmd *cobra.Command, args []string) {
		getPackageFomUsers(args[0])
		fmt.Println(packageName)
		prompt := promptContent{
			"What do you want to do?:",
			"choose what to do:",
		}
		path = args[0]
		category := getActionPromt(prompt, []string{"create class", "connect two classes (1 - n, n -m)", "write class"})
		fmt.Println(packageName, category)
	},
}

func init() {
	rootCmd.AddCommand(modelCmd)
}

func createClasse() {
	var c models.Class
	p := promptContent{
		"What is the class name:",
		"What is the class name:",
	}
	c.Name = promptGetInput(p)
	c.IsBetween = false

	c.Vars = append(c.Vars, models.Variable{
		Name:        "Id",
		Type:        "Long",
		Annotations: []string{"@Id", "@GeneratedValue(strategy = GenerationType.IDENTITY)"},
		Security:    "private",
	})

	// a loop which allows for multiple variable creations
	var r bool = true
	for r {
		// check if the user wants to make another variable
		prompt := promptui.Prompt{
			Label:     "New variable",
			IsConfirm: true,
		}
		result, _ := prompt.Run()
		r = false
		if result == "y" {
			r = true
			// ceate new var
			v, i := createVariable()
			c.Vars = append(c.Vars, v)
			if len(i.Name) != 0 {
				c.Import = append(c.Import, i)
			}
		}

	}

	var f []models.Printable
	for _, v := range c.Vars {
		f = append(f, v)
	}
	// always add jpa as import
	c.Import = append(c.Import, models.Import{Name: "javax.persistence.*"})
	var m []models.Printable
	for _, v := range c.Import {
		m = append(m, v)
	}

	fmt.Println("\n", models.GenerateClass(c.Name, packageName, f, m))
	classes = append(classes, c)
}

// function to create a new variable for the class
func createVariable() (models.Variable, models.Import) {
	var m models.Variable
	p := promptContent{
		"What is the variable name (Date, Firstname, ...):",
		"What is the variable name:",
	}
	m.Name = promptGetInput(p)
	// get the type of the variable
	var i models.Import
	m.Type, i = getVarType()
	// set the security private
	m.Security = "private"
	// add the annotations for jpa
	m.Annotations = []string{fmt.Sprintf("@Column(name = \"%s\")", strings.ToLower(m.Name))}
	return m, i
}

// function to let the user select the variable type
func getVarType() (string, models.Import) {

	p := promptContent{
		"What is the variable type (String, int, ...):",
		"Select your variable Type",
	}

	// availibe types
	i := []string{"String", "Integer", "Long", "Date"}
	// if the type requires an extra import just add it here
	imports := []models.Import{
		{Name: ""},
		{Name: ""},
		{Name: ""},
		{Name: "java.sql.Date"},
	}

	ps := promptui.Select{
		Label: p.label,
		Items: i,
	}

	index, result, _ := ps.Run()
	return result, imports[index]
}

func updateClasses() {
	// check which two classes should be selected
	p := promptContent{
		"Select the first class",
		"Select the first class",
	}
	var c []string
	for _, v := range classes {
		c = append(c, v.Name)
	}

	prompt1 := promptui.Select{
		Label: p.label,
		Items: c,
	}

	_, result1, _ := prompt1.Run()

	var c2 []string
	for _, v := range classes {
		if result1 != v.Name {
			c2 = append(c2, v.Name)
		}
	}

	p2 := promptContent{
		"Select the second class",
		"Select the second class",
	}

	prompt2 := promptui.Select{
		Label: p2.label,
		Items: c2,
	}

	_, result2, _ := prompt2.Run()

	// select what to do n - 1 or n - m

	p3 := promptContent{
		"Select the connection methode (1 - n, n - m)",
		"Select the connection methode (1 - n, n - m)",
	}

	c3 := []string{fmt.Sprintf("1 - n (%s - %s)", result1, result2),
		fmt.Sprintf("n - 1 (%s - %s)", result1, result2),
		fmt.Sprintf("n - m (%s - %s)", result1, result2)}

	prompt3 := promptui.Select{
		Label: p3.label,
		Items: c3,
	}

	i, _, _ := prompt3.Run()
	var i1, i2 int
	switch i {
	case 0:
		i1 = getClassIndexByName(result1)
		if i1 == -1 {
			return
		}
		i2 = getClassIndexByName(result2)
		if i2 == -1 {
			return
		}
	case 1:
		i2 = getClassIndexByName(result1)
		if i2 == -1 {
			return
		}
		i1 = getClassIndexByName(result2)
		if i1 == -1 {
			return
		}

	case 2:
		i1 = getClassIndexByName(result1)
		if i1 == -1 {
			return
		}
		i2 = getClassIndexByName(result2)
		if i2 == -1 {
			return
		}
		NtoM(i1, i2)
		return
	}
	v1 := models.Variable{
		Name:        strings.ToLower(classes[i2].Name),
		Security:    "private",
		Type:        fmt.Sprintf("Set<%s>", classes[i2].Name),
		Annotations: []string{"@JsonIgnore", fmt.Sprintf("@OneToMany(mappedBy = \"%s\", fetch = FetchType.LAZY)", strings.ToLower(classes[i1].Name))},
	}
	classes[i1].Vars = append(classes[i1].Vars, v1)
	classes[i1].Import = append(classes[i1].Import, models.Import{Name: "java.util.Set"})

	v2 := models.Variable{
		Name:        strings.ToLower(classes[i1].Name),
		Type:        classes[i1].Name,
		Security:    "private",
		Annotations: []string{"@ManyToOne(fetch = FetchType.LAZY)", fmt.Sprintf("@JoinColumn(name = \"%s_id\")", strings.ToLower(classes[i1].Name))},
	}

	classes[i2].Vars = append(classes[i2].Vars, v2)
}

func NtoM(i1, i2 int) {
	pc := promptContent{
		"What name should the betweentable have",
		"What name should the betweentable have",
	}

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
		Default:   fmt.Sprintf("%s_to_%s", classes[i1].Name, classes[i2].Name),
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}

	c := models.Class{
		Name:   result,
		Import: []models.Import{{Name: "javax.persistence.*"}},
	}
	// make new id
	c.Vars = append(c.Vars, models.Variable{
		Name:        "Id",
		Type:        "Long",
		Annotations: []string{"@Id", "@GeneratedValue(strategy = GenerationType.IDENTITY)"},
		Security:    "private",
	})

	c.Vars = append(c.Vars, models.Variable{
		Name:        strings.ToLower(classes[i1].Name),
		Type:        classes[i1].Name,
		Security:    "private",
		Annotations: []string{"@ManyToOne(fetch = FetchType.LAZY)", fmt.Sprintf("@JoinColumn(name = \"%s_id\")", strings.ToLower(classes[i1].Name))},
	})

	c.Vars = append(c.Vars, models.Variable{
		Name:        strings.ToLower(classes[i2].Name),
		Type:        classes[i2].Name,
		Security:    "private",
		Annotations: []string{"@ManyToOne(fetch = FetchType.LAZY)", fmt.Sprintf("@JoinColumn(name = \"%s_id\")", strings.ToLower(classes[i2].Name))},
	})
	c.IsBetween = true

	classes[i1].Vars = append(classes[i1].Vars, models.Variable{
		Name:        strings.ToLower(c.Name),
		Security:    "private",
		Type:        fmt.Sprintf("Set<%s>", c.Name),
		Annotations: []string{"@JsonIgnore", fmt.Sprintf("@OneToMany(mappedBy = \"%s\", fetch = FetchType.LAZY)", strings.ToLower(classes[i1].Name))},
	})

	classes[i2].Vars = append(classes[i2].Vars, models.Variable{
		Name:        strings.ToLower(c.Name),
		Security:    "private",
		Type:        fmt.Sprintf("Set<%s>", c.Name),
		Annotations: []string{"@JsonIgnore", fmt.Sprintf("@OneToMany(mappedBy = \"%s\", fetch = FetchType.LAZY)", strings.ToLower(classes[i2].Name))},
	})

	classes[i1].Import = append(classes[i1].Import, models.Import{Name: "java.util.Set"})
	classes[i2].Import = append(classes[i2].Import, models.Import{Name: "java.util.Set"})

	var r bool = true
	for r {
		// check if the user wants to make another variable
		prompt := promptui.Prompt{
			Label:     "New variable",
			IsConfirm: true,
		}
		result, _ := prompt.Run()
		r = false
		if result == "y" {
			r = true
			// ceate new var
			v, i := createVariable()
			c.Vars = append(c.Vars, v)
			if len(i.Name) != 0 {
				c.Import = append(c.Import, i)
			}
		}

	}
	classes = append(classes, c)
}

func writeClasses() {

	if err := os.Mkdir(fmt.Sprintf("%s/models", path), os.ModePerm); err != nil {
		fmt.Println(err)
	}

	prompt := promptui.Prompt{
		Label:     "You want the Repo:",
		IsConfirm: true,
	}
	result, _ := prompt.Run()
	if result == "y" {
		if err := os.Mkdir(fmt.Sprintf("%s/models/repository", path), os.ModePerm); err != nil {
			fmt.Println(err)
		}
		for _, v := range classes {
			var sb strings.Builder
			sb.WriteString(fmt.Sprintf("package %s.models.repository;\n\n", packageName))
			sb.WriteString(fmt.Sprintf("import %s.models.%s;\n", packageName, v.Name))
			sb.WriteString("import org.springframework.data.jpa.repository.JpaRepository;\n\n")
			sb.WriteString(fmt.Sprintf("public interface %sRepository extends JpaRepository<%s, Long> {\n}", v.Name, v.Name))
			os.WriteFile(fmt.Sprintf("%s/models/repository/%sRepository.java", path, v.Name), []byte(sb.String()), 0664)
		}
	}
	for _, c := range classes {
		var f2 []models.Printable
		for _, v := range c.Vars {
			f2 = append(f2, v)
		}

		var m2 []models.Printable
		for _, v := range c.Import {
			m2 = append(m2, v)
		}
		os.WriteFile(fmt.Sprintf("%s/models/%s.java", path, c.Name), []byte(models.GenerateClass(c.Name, packageName, f2, m2)), 0664)
	}

	prompt1 := promptui.Prompt{
		Label:     "You want apis 2:",
		IsConfirm: true,
	}

	result1, _ := prompt1.Run()
	if result1 == "n" {
		return
	}
	os.MkdirAll(fmt.Sprintf("%s/controllers/api", path), os.ModePerm)

	for _, v := range classes {
		if !v.IsBetween {
			result = apiController(v)
			os.WriteFile(fmt.Sprintf("%s/controllers/api/Rest%sController.java", path, v.Name), []byte(result), 0664)
		}
	}
}

func apiController(c models.Class) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("package %s.controllers.api;\n\n", packageName))
	sb.WriteString("import org.springframework.beans.factory.annotation.Autowired;\n")
	sb.WriteString("import org.springframework.web.bind.annotation.*;\n")
	sb.WriteString("import java.util.List;\n")
	sb.WriteString(fmt.Sprintf("import %s.models.%s;\n", packageName, c.Name))
	sb.WriteString(fmt.Sprintf("import %s.models.repository.%sRepository;\n\n", packageName, c.Name))

	sb.WriteString(fmt.Sprintf("@RestController\n@RequestMapping(\"/%s\")\npublic class Rest%sController {\n", strings.ToLower(c.Name), c.Name))
	sb.WriteString(fmt.Sprintf("\n@Autowired\n%sRepository r;\n\n", c.Name))
	sb.WriteString(fmt.Sprintf("@GetMapping(\"/\")\npublic List<%s> getAll(){\nreturn r.findAll();\n}", c.Name))

	sb.WriteString("\n}")
	return sb.String()
}
