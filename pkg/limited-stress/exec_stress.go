package limited_stress

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type StressResult struct {
	Time           string
	NumRuns        int
	NumFailures    int
	FailureMessage string
}

func (r StressResult) ToString() string {
	return fmt.Sprintf("Time: %s, Nums Runs: %d, Nums Failures: %d, Failure Message: %s", r.Time, r.NumRuns, r.NumFailures, r.FailureMessage)
}

func ExecStress(path string, timeout time.Duration, executionDirectory string, numProcesses int, logger *log.Logger) string {
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()
	outputBuff := new(bytes.Buffer)
	cmd := exec.CommandContext(ctx, "stress", "-p", strconv.Itoa(numProcesses), path)
	cmd.Stdout = outputBuff
	cmd.Dir = executionDirectory
	err := cmd.Run()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			str := outputBuff.String()
			logger.Printf("Executed Stress For , Path: %s \n", path)
			return str
		} else {
			logger.Printf("Failed To Execute Stress, Path: %s with error: %s\n", path, err.Error())
		}
		return ""
	}
	return ""
}

func CompileTest(packageName string, directory string, outputPath string, logger *log.Logger) bool {
	cmd := exec.Command("go", "test", "-race", "-c", "-o", outputPath, packageName)
	//check if its important to run in directory
	cmd.Dir = directory
	output, err := cmd.Output()
	if err != nil {
		logger.Println("Error Executing Compile Command For Package: %s Error:", packageName, err)
	}
	if strings.Contains(string(output), "no test files") {
		return false
	}
	return true
}

func parseStressError(output string) string {
	separatedError := strings.Split(output, errorIdentifier)
	//in cse of repetitive error, better to return 2 error
	if len(separatedError) <= 1 {
		return separatedError[0]
	}
	return separatedError[1]
}

const errorIdentifier = "/tmp/go-stress"

func ParseStressOutput(output string) StressResult {
	var result = StressResult{}
	pattern := `(?P<Time>.*?): (?P<Runs>\d+) runs so far, (?P<Failures>\d+) failures`
	regex := regexp.MustCompile(pattern)
	matches := regex.FindAllStringSubmatch(output, -1)
	if len(matches) != 0 {
		lastMatch := matches[len(matches)-1]
		result.Time = lastMatch[regex.SubexpIndex("Time")]
		result.NumRuns, _ = strconv.Atoi(lastMatch[regex.SubexpIndex("Runs")])
		result.NumFailures, _ = strconv.Atoi(lastMatch[regex.SubexpIndex("Failures")])
	}
	if result.NumFailures != 0 || len(matches) == 0 {
		result.FailureMessage = parseStressError(output)
	}
	return result
}
