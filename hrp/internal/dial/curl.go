package dial

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/test-instructor/yangfan/hrp/internal/builtin"
	"github.com/test-instructor/yangfan/hrp/internal/myexec"
)

const (
	normalResult = "STDOUT"
	errorResult  = "STDERR"
	failedResult = "FAILED"
)

type CurlResult struct {
	Result     string `json:"result"`
	ErrorMsg   string `json:"errorMsg"`
	ResultType string `json:"resultType"`
}

func DoCurl(args []string) (err error) {
	var saveTests bool
	for i, arg := range args {
		if arg == "--save-tests" {
			args = append(args[:i], args[i+1:]...)
			saveTests = true
		}
	}
	var curlResult CurlResult
	defer func() {
		if saveTests {
			dir, _ := os.Getwd()
			curlResultName := fmt.Sprintf("curl_result_%v.json", time.Now().Format("20060102150405"))
			curlResultPath := filepath.Join(dir, curlResultName)
			err = builtin.Dump2JSON(curlResult, curlResultPath)
			if err != nil {
				log.Error().Err(err).Msg("save dns resolution result failed")
			}
		}
	}()

	cmd := myexec.Command("curl", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		log.Error().Err(err).Msgf("fail to run curl command")
		curlResult.ErrorMsg = err.Error()
		curlResult.Result = stderr.String()
		curlResult.ResultType = errorResult
		return
	}
	if stdout.String() != "" {
		fmt.Println(stdout.String())
		curlResult.Result = stdout.String()
		curlResult.ResultType = normalResult
	} else if stderr.String() != "" {
		fmt.Println(stderr.String())
		curlResult.ErrorMsg = stderr.String()
		curlResult.ResultType = errorResult
	}
	return
}
