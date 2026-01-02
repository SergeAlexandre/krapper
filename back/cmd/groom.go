package cmd

import (
	"encoding/json"
	"fmt"
	"krapper/internal/misc"
	"krapper/internal/wrap"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var groomParams struct {
	yaml bool
	json bool
}

func init() {
	groomCmd.PersistentFlags().BoolVarP(&groomParams.yaml, "yaml", "y", false, "dump in yaml format")
	groomCmd.PersistentFlags().BoolVarP(&groomParams.json, "json", "j", false, "dump in json format")
}

var groomCmd = &cobra.Command{
	Use:   "groom",
	Short: "Load a groom a wrap file",
	Run: func(cmd *cobra.Command, args []string) {
		for _, fileName := range args {
			fmt.Printf("\n---------------------------------------------- Processing file %s\n", fileName)
			w, err := load(fileName)
			if err != nil {
				log.Fatal(err)
			}
			if groomParams.json {
				jsonData, err := json.MarshalIndent(w, "", "  ")
				if err != nil {
					log.Fatalf("Error marshalling to JSON: %v", err)
				}
				fmt.Println(string(jsonData))
			}

			if groomParams.yaml {
				yamlData, err := yaml.Marshal(w)
				if err != nil {
					log.Fatalf("Error marshalling to YAML: %v", err)
				}
				fmt.Println(string(yamlData))
			}
		}

	},
}

func load(fName string) (*wrap.Wrap, error) {

	filename, err := filepath.Abs(fName)
	if err != nil {
		log.Fatalf("Error getting absolute path of yaml file: %v", err)
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file %s: %v", filename, err)
	}
	defer func() { _ = file.Close() }()

	var w wrap.Wrap
	decoder := yaml.NewDecoder(file)
	decoder.KnownFields(true)
	if err := decoder.Decode(&w); err != nil {
		return nil, fmt.Errorf("error decoding file %s: %v", filename, err)
	}

	// Template logic
	if w.Template != "" && w.TemplateFile != "" {
		return nil, fmt.Errorf("template and templateFile are mutually exclusive")
	}

	if w.TemplateFile != "" {
		templatePath := misc.AdjustPath(filepath.Dir(filename), w.TemplateFile)

		content, err := os.ReadFile(templatePath)
		if err != nil {
			return nil, fmt.Errorf("error reading template file %s: %v", templatePath, err)
		}

		w.Template = wrap.WrTemplate(content)
		w.TemplateFile = ""
	}

	err = w.Groom()
	if err != nil {
		return nil, err
	}
	return &w, nil
}
