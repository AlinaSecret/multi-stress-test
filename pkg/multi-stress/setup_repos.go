package multi_stress

import (
	"github.com/AlinaSecret/multi-stress-test/pkg/packages-collector"
	"github.com/go-git/go-git/v5"
	"log"
	"os"
	"sync"
)

func SetupRepos(repos []packages_collector.RepoInfo, directory string, logger *log.Logger) chan packages_collector.RepoInfo {
	var wg sync.WaitGroup
	readyRepos := make(chan packages_collector.RepoInfo)
	for _, repo := range repos {
		r := packages_collector.RepoInfo{Name: repo.Name, Directory: repo.Directory, URL: repo.URL, IsRemote: repo.IsRemote}
		wg.Add(1)
		go func() {
			defer wg.Done()
			if r.IsRemote {
				r.Directory, _ = os.MkdirTemp(directory, r.Name)
				_, err := git.PlainClone(r.Directory, false, &git.CloneOptions{
					URL: r.URL,
				})
				if err != nil {
					logger.Printf("Error cloning repository: %v\n", err)
				} else {
					logger.Printf("Cloned Repository to: %s\n", r.Directory)
					readyRepos <- r
				}

			} else {
				readyRepos <- r
			}
		}()
	}
	go func() {
		wg.Wait()
		close(readyRepos)
	}()
	return readyRepos
}
