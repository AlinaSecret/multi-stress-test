package main

import (
	"fmt"
	"github.com/AlinaSecret/multi-stress-test/pkg/multi-stress"
	"github.com/akamensky/argparse"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"runtime"
	"time"
)

func main() {
	var logger *log.Logger
	parser := argparse.NewParser("multi-stress", "Command-line application")
	yamlPath := parser.String("y", "yaml", &argparse.Options{Required: true, Help: "Path to the YAML file containing the list of repositories and directories to test"})
	outputPath := parser.String("o", "output", &argparse.Options{Help: "Path to the output file where the report will be saved, should end with xlsx"})
	stressTestTime := parser.Int("t", "time", &argparse.Options{Help: "Time duration for stress test per package (in seconds)", Default: 10})
	parallelProcesses := parser.Int("p", "parallel", &argparse.Options{Help: "Number of parallel processes to run stress tests per package", Default: runtime.NumCPU()})
	numWorkers := parser.Int("w", "workers", &argparse.Options{Help: "Number of packages to run stress tests in parallel (should be synchronized with the number of parallel processes per package)", Default: 3})
	cloneDirectory := parser.String("d", "directory", &argparse.Options{Help: "Directory to clone all repositories. If not specified, repositories will be cloned in a temporary directory that will be deleted"})
	verbose := parser.Flag("v", "verbose", &argparse.Options{Help: "Enable verbose logging", Default: false})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}

	yamlFile, err := os.ReadFile(*yamlPath)
	if err != nil {
		fmt.Println("Error reading YAML file:", err)
		return
	}
	if *verbose {
		logger = log.New(os.Stdout, "", 5)
	} else {
		file, _ := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		defer file.Close()
		logger = log.New(file, "", log.Ldate|log.Ltime)
	}
	var yamlInput YAMLInput
	err = yaml.Unmarshal(yamlFile, &yamlInput)
	if err != nil {
		fmt.Println("Validation Error in YAML:\n", err)
		return
	}
	duration := time.Duration(*stressTestTime)
	multi_stress.RunMultiStress(*outputPath, yamlInput.ConvertToRepoInfo(), *cloneDirectory, duration, *parallelProcesses, *numWorkers, logger)
}
