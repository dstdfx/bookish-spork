package log

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testFile = "test.log"

func writeTmpFile(t *testing.T, filename, contents string) (string, error) {
	tmpDir, err := ioutil.TempDir("/tmp", "tests")
	if err != nil {
		return "", err
	}

	t.Logf("created temporary dir: %s", tmpDir)

	filePath := strings.Join([]string{tmpDir, filename}, "/")

	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}

	t.Logf("will write to temporary file: %s", filePath)

	defer file.Close()

	bytes, err := file.WriteString(contents)
	if err != nil {
		return "", err
	}

	t.Logf("wrote %d bytes to temporary file: %s", bytes, filePath)

	return tmpDir, nil
}

func TestOutputConfigFileAndStdout(t *testing.T) {
	// Prepare test logging file.
	logFile := testFile
	tmpDir, err := writeTmpFile(t, testFile, "")
	assert.NoError(t, err)

	defer func() {
		err := os.RemoveAll(tmpDir)
		assert.NoError(t, err)
	}()

	filePath := strings.Join([]string{tmpDir, logFile}, "/")

	expectedOutputPaths := []string{filePath, stdout}
	expectedErrPaths := []string{filePath, stderr}

	actualOutputPaths, actualErrPaths, err := outputConfig(filePath, true)

	assert.NoError(t, err)
	assert.ElementsMatch(t, expectedOutputPaths, actualOutputPaths)
	assert.ElementsMatch(t, expectedErrPaths, actualErrPaths)
}

func TestOutputConfigFile(t *testing.T) {
	// Prepare test logging file.
	logFile := testFile
	tmpDir, err := writeTmpFile(t, testFile, "")
	assert.NoError(t, err)

	defer func() {
		err := os.RemoveAll(tmpDir)
		assert.NoError(t, err)
	}()

	filePath := strings.Join([]string{tmpDir, logFile}, "/")

	expectedOutputPaths := []string{filePath}
	expectedErrPaths := []string{filePath, stderr}

	actualOutputPaths, actualErrPaths, err := outputConfig(filePath, false)

	assert.NoError(t, err)
	assert.ElementsMatch(t, expectedOutputPaths, actualOutputPaths)
	assert.ElementsMatch(t, expectedErrPaths, actualErrPaths)
}

func TestOutputConfigNoFile(t *testing.T) {
	expectedOutputPaths := []string{stdout}
	expectedErrPaths := []string{stderr}

	actualOutputPaths, actualErrPaths, err := outputConfig("", true)

	assert.NoError(t, err)
	assert.ElementsMatch(t, expectedOutputPaths, actualOutputPaths)
	assert.ElementsMatch(t, expectedErrPaths, actualErrPaths)
}
