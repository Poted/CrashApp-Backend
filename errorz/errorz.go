package errorz

import (
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"strings"
)

type AppError struct {
	Error    error
	FuncInfo []FuncInfo
}

type FuncInfo struct {
	FuncLine int
	Function string
	FuncPath string
}

type IAppError interface {
	SendError(msg error) error
}

func SendError(errz error) error {

	var resetColor = "\033[0m"
	var color = "\033[33m"

	appError := &AppError{
		Error:    errz,
		FuncInfo: GetFuncInfo(),
	}

	errorJson, err := json.MarshalIndent(appError, "", "  ")
	if err != nil {
		r := fmt.Sprintf("unkown error occured; cannot parse error: (%s)", err)
		return errors.New(r)
	}

	return errors.New(color + string(errorJson) + resetColor + "\n")
}

func GetFuncInfo() []FuncInfo {

	lvl := 4
	stackTrace := make([]uintptr, 10)

	info := FuncInfo{}
	ret := []FuncInfo{}

	for {

		runtime.Callers(lvl, stackTrace)
		frame, _ := runtime.CallersFrames(stackTrace).Next()

		info = FuncInfo{
			FuncLine: getLineNumber(lvl),
			Function: getFunctionName(lvl),
			FuncPath: frame.File,
		}

		if info.Function == "" {
			break
		}

		lvl += 1

		if !strings.Contains(strings.ToLower(info.FuncPath), "/go/") {
			ret = append(ret, info)
		}
	}

	return ret
}

func getLineNumber(lvl int) int {
	var _, _, line, _ = runtime.Caller(lvl)
	return line
}

func getFunctionName(lvl int) string {
	var pc, _, _, _ = runtime.Caller(lvl)
	fullFuncName := runtime.FuncForPC(pc).Name()
	parts := strings.Split(fullFuncName, "/")
	return parts[len(parts)-1]
}
