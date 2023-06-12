package main

import (
	packagesCollector "github.com/AlinaSecret/multi-stress-test/pkg/packages-collector"
	"os"
	"os/exec"
	"sync"
)

type YAMLInput struct {
	Directories  []string `yaml:"Directories"`
	Repositories []string `yaml:"Repositories"`
}

func (yamlInput *YAMLInput) parseDirectories() []DirectoryDoesntExist {
	var directoryErrors []DirectoryDoesntExist
	for _, dir := range yamlInput.Directories {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			directoryErrors = append(directoryErrors, DirectoryDoesntExist{name: dir})
		}
	}
	return directoryErrors
}

const successExitCode = 0

func (yamlInput *YAMLInput) parseRepositories() []RepositoryDoesntExist {
	var repositoryErrors []RepositoryDoesntExist
	// wg := *wait_group_max.CreateWorkGroupMax(runtime.NumCPU())
	var wg sync.WaitGroup
	repositoryErrorsC := make(chan RepositoryDoesntExist)
	for _, repo := range yamlInput.Repositories {
		r := repo
		wg.Add(1)
		go func() {
			defer wg.Done()
			cmd := exec.Command("git", "ls-remote", r)
			err := cmd.Run()
			if exitError, ok := err.(*exec.ExitError); ok && exitError.ExitCode() != successExitCode {
				repositoryErrorsC <- RepositoryDoesntExist{name: r}
			}
		}()
	}
	go func() {
		wg.Wait()
		close(repositoryErrorsC)

	}()
	for invalidRepo := range repositoryErrorsC {
		repositoryErrors = append(repositoryErrors, invalidRepo)
	}
	return repositoryErrors
}

func (yamlInput *YAMLInput) ConvertToRepoInfo() []packagesCollector.RepoInfo {
	var repos []packagesCollector.RepoInfo
	for _, dir := range yamlInput.Directories {
		repos = append(repos, packagesCollector.CreateLocalRepo(dir))
	}
	for _, repo := range yamlInput.Repositories {
		repos = append(repos, packagesCollector.CreateRemoteRepo(repo))
	}
	return repos
}

func (yamlInput *YAMLInput) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// Create Alias so structs can unmarshal itself:
	type YAMLInputAlias YAMLInput
	data := (*YAMLInputAlias)(yamlInput)
	err := unmarshal(data)
	if err == nil {
		directoryErrors := yamlInput.parseDirectories()
		RepositoryErrors := yamlInput.parseRepositories()
		if len(directoryErrors) != 0 || len(RepositoryErrors) != 0 {
			return MultiStressYamlInputError{RepositoryErrors: RepositoryErrors, DirectoryErrors: directoryErrors}
		}
		return nil
	}
	return err
}
