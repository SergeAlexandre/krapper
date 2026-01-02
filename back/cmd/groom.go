package cmd

import (
	"encoding/json"
	"fmt"
	"krapper/internal/wrap"
	"log"

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
	Short: "Load and groom a wrap file",
	Run: func(cmd *cobra.Command, args []string) {
		for _, fileName := range args {
			fmt.Printf("\n---------------------------------------------- Processing file %s\n", fileName)
			w, err := wrap.Load(fileName)
			if err != nil {
				log.Fatal(err)
			}
			if w != nil {
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
			} else {
				fmt.Printf("Not a wrap file")
			}
		}
	},
}
