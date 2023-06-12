package test_report

type ITestSummary interface {
	GetRepoName() string
	GetPackageName() string
	WasTested() bool
	WasSkipped() bool
	GetTestTime() string
	GetNumberOfRuns() int
	GetNumberOfFailures() int
	GetFailureMessage() string
	HasFailed() bool
}
