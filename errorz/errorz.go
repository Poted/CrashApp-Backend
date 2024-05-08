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
	ErrMsg   string
	FuncLine int
	Function string
	FuncPath string
}

type IAppError interface {
	SendError(msg error) error
}

func SendError(errz error) error {

	stackTrace := make([]uintptr, 10)
	runtime.Callers(2, stackTrace)
	frame, _ := runtime.CallersFrames(stackTrace).Next()

	var resetColor = "\033[0m"
	var firstColor = "\033[31m"
	// var startIcon = "\021"
	var secondColor = "\033[34m"

	appError := &AppError{
		Error: errz,
		ErrMsg: func() string {

			if errors.Unwrap(errz) != nil { /// HANDLE PANIC
				err := errors.Unwrap(errz)
				if err != nil {
					return "1"
				}
				return errors.Unwrap(errz).Error()
			}

			return errz.Error()
		}(),
		FuncLine: getLineNumber(),
		Function: getFunctionName(),
		FuncPath: frame.File,
	}

	errorJson, err := json.MarshalIndent(appError, "", " ")
	if err != nil {
		r := fmt.Sprintf("unkown error occured; cannot parse error: (%s)", err)
		return errors.New(r)
	}

	for _, line := range strings.Split(string(errorJson), "\n") {

		parts := strings.SplitN(line, ":", 2)

		if !strings.Contains(line, ":") {
			fmt.Printf(secondColor + line + "\n")
		} else if len(parts) == 2 {
			fieldName := strings.TrimSpace(strings.Trim(parts[0], ``))
			fieldValue := strings.TrimSpace(strings.Trim(parts[1], ``))
			// fmt.Printf("%s.%v%s: %v%s\n", startIcon, firstColor, fieldName, secondColor, fieldValue)
			fmt.Printf("%s.%v%s: %v%s\n", firstColor, fieldName, secondColor, fieldValue)
		} else {
			fmt.Print(resetColor)
		}

	}

	return errors.New(string(errorJson))
}

func getLineNumber() int {
	var _, _, line, _ = runtime.Caller(2)
	return line
}

func getFunctionName() string {
	var pc, _, _, _ = runtime.Caller(2)
	fullFuncName := runtime.FuncForPC(pc).Name()
	parts := strings.Split(fullFuncName, "/")
	return parts[len(parts)-1]
}
