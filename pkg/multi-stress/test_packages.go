package multi_stress

import (
	"github.com/AlinaSecret/multi-stress-test/pkg/limited-stress"
	"github.com/AlinaSecret/multi-stress-test/pkg/packages-collector"
	"github.com/AlinaSecret/multi-stress-test/pkg/wait_group_max"
	"log"
	"os"
	"strings"
	"time"
)

func TestPackage(test packages_collector.PacakgeInfo, outputDirectory string, stressTime time.Duration, numProcesses int, logger *log.Logger) (limited_stress.StressResult, bool) {
	var compiledPath string
	var err error
	var result limited_stress.StressResult
	compiledPath = outputDirectory + "/" + strings.ReplaceAll(test.PackageName, "/", "-")
	if limited_stress.CompileTest(test.PackageName, test.Repo.Directory, compiledPath, logger) == true {
		result = limited_stress.ParseStressOutput(limited_stress.ExecStress(compiledPath, stressTime, test.Repo.Directory, numProcesses, logger))
		logger.Printf("Stress Testing is Finshed For Package Named: - %s \n", test.PackageName)
		err = os.Remove(compiledPath)
		if err != nil {
			logger.Println("Error deleting file:", err)
		}
		return result, true
	} else {
		logger.Printf("Found That Pacakge - %s Does Not Include Tests \n", test.PackageName)
	}
	return result, false
}

func TestPackages(jobs <-chan packages_collector.PacakgeInfo, outputDirectory string, stressTime time.Duration, numProcesses int, numWorkers int, logger *log.Logger) (<-chan TestSummary, *wait_group_max.WaitGroupMax) {
	var group = wait_group_max.CreateWorkGroupMax(numWorkers)
	results := make(chan TestSummary)
	group.Add(1)
	go func() {
		defer close(results)
		defer group.Done()
		var wg = wait_group_max.CreateWorkGroupMax(numWorkers)
		for pack := range jobs {
			wg.Add(1)
			currentPackage := pack
			go func() {
				var testSum TestSummary
				testSum.Results, testSum.IsTest = TestPackage(currentPackage, outputDirectory, stressTime, numProcesses, logger)
				testSum.TestSpecs = currentPackage
				testSum.IsSkipped = false
				results <- testSum
				wg.Done()
			}()
		}
		wg.Wait()
	}()
	return results, group
}
