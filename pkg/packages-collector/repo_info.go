package packages_collector

import (
	"path/filepath"
	"strings"
)

type RepoInfo struct {
	Name      string
	Directory string
	URL       string
	IsRemote  bool
}

const gitSuffix = ".git"

func CreateRemoteRepo(repositoryLink string) RepoInfo {
	repositoryLink = strings.TrimSuffix(repositoryLink, gitSuffix)
	_, repoName := filepath.Split(repositoryLink)
	return RepoInfo{URL: repositoryLink, Name: repoName, IsRemote: true}
}

func CreateLocalRepo(directory string) RepoInfo {
	return RepoInfo{Directory: directory, Name: filepath.Base(directory), IsRemote: false}
}
