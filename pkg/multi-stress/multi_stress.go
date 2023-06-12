package multi_stress

import (
	packages_collector "github.com/AlinaSecret/multi-stress-test/pkg/packages-collector"
	test_report "github.com/AlinaSecret/multi-stress-test/pkg/test-report"
	"log"
	"os"
	"time"
)

func RunMultiStress(reportPath string, repos []packages_collector.RepoInfo, cloneDirectory string, stressTime time.Duration, numProcesses int, numWorkers int, logger *log.Logger) {
	report := test_report.New(reportPath, "Sheet1", logger)
	report.AddHeaders()
	if len(cloneDirectory) == 0 {
		cloneDirectory, err := os.MkdirTemp("", "repos")
		if err != nil {
			logger.Println("Error creating temporary directory:", err)
			return
		}

		defer removeDirectory(cloneDirectory, logger)

	}

	readyRepos := SetupRepos(repos, cloneDirectory, logger)
	packages := packages_collector.CollectPackages(readyRepos, logger)

	tempDir, err := os.MkdirTemp("", "tests")
	if err != nil {
		logger.Println("Error creating temporary directory:", err)
		return
	}

	defer removeDirectory(tempDir, logger)

	results, wg := TestPackages(packages, tempDir, stressTime, numProcesses, numWorkers, logger)
	go func() {
		for res := range results {
			report.AddTest(res)
		}
	}()
	wg.Wait()
	report.Save()
}

func removeDirectory(dir string, logger *log.Logger) {
	err := os.RemoveAll(dir)
	if err != nil {
		logger.Println("Error: Failed To Remove Directory %s %v", dir, err)
	}
}
