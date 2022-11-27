package logger

import (
	"bufio"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func getLineCount(file *os.File) int {
	fileScanner := bufio.NewScanner(file)
	lineCount := 0
	for fileScanner.Scan() {
		lineCount++
	}
	return lineCount
}

func TestLogger(t *testing.T) {
	fileName := "1_tst.log"
	q := New(fileName, "debug")
	q.Trace("trace")  // не должен напечатать
	q.Warning("warn") // должен
	q.Error("er")     // должен
	file, err := os.Open(fileName)
	require.NoError(t, err)
	defer file.Close()
	defer os.Remove(fileName)
	require.Equal(t, 2, getLineCount(file))

	fmt.Print("test2")
	fileName = "2_tst.log"
	q2 := New(fileName, "ERROR")
	q2.Trace("trace")      // не должен напечатать
	q2.Info("informatuom") // не должен
	q.Warning("Warning")   // не должен
	q2.Error("er")         // должен
	file, err = os.Open(fileName)
	require.NoError(t, err)
	defer file.Close()
	defer os.Remove(fileName)
	require.Equal(t, 1, getLineCount(file))
}
