package packages_collector

import (
	"log"
	"os/exec"
	"strings"
	"sync"
)

type PacakgeInfo struct {
	PackageName string
	Repo        RepoInfo
}

func CollectPackages(repos <-chan RepoInfo, logger *log.Logger) <-chan PacakgeInfo {
	var wg sync.WaitGroup
	pacakgeInfos := make(chan PacakgeInfo)
	for repo := range repos {
		wg.Add(1)
		go getRepositoryPackages(pacakgeInfos, repo, &wg, logger)
	}
	go func() {
		wg.Wait()
		close(pacakgeInfos)
	}()
	return pacakgeInfos
}

func getRepositoryPackages(jobs chan<- PacakgeInfo, repo RepoInfo, wg *sync.WaitGroup, logger *log.Logger) {
	defer wg.Done()
	cmd := exec.Command("go", "list", "./...")
	cmd.Dir = repo.Directory
	output, err := cmd.Output()
	if err != nil {
		logger.Println("Error executing command:", err)
	}
	packages := strings.Split(string(output), "\n")
	for _, packageName := range packages {
		if len(packageName) > 0 {
			jobs <- PacakgeInfo{PackageName: packageName, Repo: repo}
		}
	}
	logger.Println("Got List Of Packages For: ", repo.Name)
	return
}
