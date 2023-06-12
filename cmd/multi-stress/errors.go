package main

import (
	"fmt"
)

type DirectoryDoesntExist struct {
	name string
}

func (e DirectoryDoesntExist) Error() string {
	return fmt.Sprintf("	%s Doesn't Exist", e.name)
}

type MultiStressYamlInputError struct {
	DirectoryErrors  []DirectoryDoesntExist
	RepositoryErrors []RepositoryDoesntExist
}

const new_line = "\n"
const ErrorMeesageDelimiter = new_line

func (e MultiStressYamlInputError) Error() string {
	dirMessage := concatErrors(e.DirectoryErrors, ErrorMeesageDelimiter)
	if dirMessage != "" {
		dirMessage = "For The Following Directories: \n" + dirMessage
	}
	repoMessage := concatErrors(e.RepositoryErrors, ErrorMeesageDelimiter)
	if repoMessage != "" {
		repoMessage = "For The Following Repositories: \n" + repoMessage
	}
	completeMessage := `
The Following Errors occurred:
  %s
  %s
    `
	return fmt.Sprintf(completeMessage, dirMessage, repoMessage)
}

type RepositoryDoesntExist struct {
	name string
}

func (e RepositoryDoesntExist) Error() string {
	return fmt.Sprintf("	%s Doesn't Exist", e.name)
}

func concatErrors[T error](errs []T, delimiter string) string {
	if len(errs) == 0 {
		return ""
	}
	completeMessage := errs[0].Error()
	for _, err := range errs[1:] {
		completeMessage += delimiter + err.Error()
	}
	return completeMessage
}
