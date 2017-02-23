package logger

import "runtime"

type wrapperError struct {
	Message    string
	Context    Fields
	Frame      *wrapperFrame
	InnerError error
}

type wrapperFrame struct {
	Function string
	File     string
	Line     int
}

func (e *wrapperError) Error() string {
	return e.Message
}

var _ Fielder = (*wrapperError)(nil)

func (err *wrapperError) ToFields() Fields {
	fields := Fields{
		"message": err.Message,
		"context": err.Context,
	}

	if innerError := errorToField(err.InnerError); innerError != nil {
		fields["innerError"] = innerError
	}

	if frame := err.Frame; frame != nil {
		if frame.Function != "" {
			fields["function"] = frame.Function
		}
		fields["file"] = frame.File
		fields["line"] = frame.Line
	}

	return fields
}

func WrapError(inner error, message string) error {
	return &wrapperError{
		Message:    message,
		Frame:      getFrame(2),
		InnerError: inner,
	}
}

func WrapErrorWithContext(inner error, message string, context Fields) error {
	return &wrapperError{
		Message:    message,
		Context:    context,
		Frame:      getFrame(2),
		InnerError: inner,
	}
}

func WrapErrorWithContextAndStack(inner error, message string, context Fields, stackSkip int) error {
	return &wrapperError{
		Message:    message,
		Context:    context,
		Frame:      getFrame(2 + stackSkip),
		InnerError: inner,
	}
}

func getFrame(skip int) *wrapperFrame {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return nil
	}

	var funcName string
	if f := runtime.FuncForPC(pc); f != nil {
		funcName = f.Name()
	}

	return &wrapperFrame{
		Function: funcName,
		File:     file,
		Line:     line,
	}
}
