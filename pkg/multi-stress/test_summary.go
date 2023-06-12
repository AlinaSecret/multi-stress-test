package multi_stress

import (
	"fmt"
	limited_stress "github.com/AlinaSecret/multi-stress-test/pkg/limited-stress"
	packages_collector "github.com/AlinaSecret/multi-stress-test/pkg/packages-collector"
)

type TestSummary struct {
	TestSpecs packages_collector.PacakgeInfo
	Results   limited_stress.StressResult
	IsTest    bool
	IsSkipped bool
}

func (ts TestSummary) ToString() string {
	multiline := `
    Test Summary:
	Repo: %s
	Package Name: %s
	Was Tested: %v
	Was Skipped: %v
	Stress Results: %s

    `
	return fmt.Sprintf(multiline, ts.TestSpecs.Repo.Name, ts.TestSpecs.PackageName, ts.IsTest, ts.IsSkipped, ts.Results.ToString())
}

func (ts TestSummary) GetRepoName() string {
	return ts.TestSpecs.Repo.Name
}

func (ts TestSummary) GetPackageName() string {
	return ts.TestSpecs.PackageName
}

func (ts TestSummary) WasTested() bool {
	return ts.IsTest
}

func (ts TestSummary) WasSkipped() bool {
	return ts.IsSkipped
}

func (ts TestSummary) GetTestTime() string {
	return ts.Results.Time
}

func (ts TestSummary) GetNumberOfRuns() int {
	return ts.Results.NumRuns
}

func (ts TestSummary) GetNumberOfFailures() int {
	return ts.Results.NumFailures
}

func (ts TestSummary) GetFailureMessage() string {
	return ts.Results.FailureMessage
}

func (ts TestSummary) HasFailed() bool {
	if ts.GetNumberOfFailures() > 0 || len(ts.GetFailureMessage()) != 0 {
		return true
	}
	return false
}
