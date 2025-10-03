package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	config := &GeneratorConfig{}

	flag.StringVar(&config.Kind, "kind", "", "Entity name to generate (e.g., User, Product)")
	flag.Parse()

	if config.Kind == "" {
		fmt.Fprintf(os.Stderr, "Usage: %s --kind <EntityName>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nExample:\n")
		fmt.Fprintf(os.Stderr, "  %s --kind User\n", os.Args[0])
		os.Exit(1)
	}

	if err := generateEntity(config); err != nil {
		fmt.Fprintf(os.Stderr, "Generate failed: %v\n", err)
		os.Exit(1)
	}
}