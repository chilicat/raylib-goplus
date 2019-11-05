package raylib

/*
#include "raylib.h"
#include <stdlib.h>
#ifndef GO_TRACE
#define GO_TRACE

#define MAX_TRACELOG_BUFFER_SIZE   128  // Max length of one trace-log message. From raylib/utils.c

void Go_TraceLogWrapper(int logType, const char *text) {
  TraceLog(logType, text);
}

void Go_CustomCallbackHook(int logType, const char *text, va_list args) {
  char buffer[MAX_TRACELOG_BUFFER_SIZE] = { 0 };
  vsprintf (buffer, text, args);
  onTraceCallback(logType, buffer);
}

void Go_EnableCustomCallback() {
  SetTraceLogCallback(&Go_CustomCallbackHook);
}

void Go_DisableCustomCallback() {
  SetTraceLogCallback(NULL);
}


#endif
*/
import "C"
import (
	"fmt"
	"unsafe"
)

type TraceLogType int32

const (
	LogAll TraceLogType = iota
	LogTrace
	LogDebug
	LogInfo
	LogWarning
	LogError
	LogFatal
	LogNone
)

var logLevelType TraceLogType = LogInfo
var traceCallback func(logType TraceLogType, text string)

//SetTraceLogLevel : Set the current threshold (minimum) log level
func SetTraceLogLevel(logType TraceLogType) {
	logLevelType = logType
	C.SetTraceLogLevel(C.int(logType))
}

//SetTraceLogExit : Set the exit threshold (minimum) log level
func SetTraceLogExit(logType TraceLogType) {
	C.SetTraceLogExit(C.int(logType))
}

//SetTraceLogCallback : Set the callback for custom logging.
func SetTraceLogCallback(callback func(logType TraceLogType, text string)) {
	traceCallback = callback
	if traceCallback != nil {
		C.Go_EnableCustomCallback()
	} else {
		C.Go_DisableCustomCallback()
	}
}

//TraceLog shows log messages.
func TraceLog(logType TraceLogType, a ...interface{}) {
	line := fmt.Sprintln(a...)
	if traceCallback != nil {

		//Make sure we are the correct level. Since we are doing this ourselves,
		// we need to keep track of this
		if logType < logLevelType {
			return
		}

		//We have a custom tracer, so lets just use that and avoid the C end entirely.
		traceCallback(logType, line)

	} else {

		//We dont have a custom tracer yet, so we will just execute our wrapper over it.
		// We have a wrapper because cgo doesn't support variable length arguments.
		cline := C.CString(line)
		defer C.free(unsafe.Pointer(cline))
		C.Go_TraceLogWrapper(C.int(logType), cline)
	}
}