package cli

import (
	"fmt"
	"io"
)

const (
	ExitOK    = 0
	ExitError = 1
)

func WriteError(w io.Writer, err error) {
	fmt.Fprintf(w, "%v", err.Error())
	fmt.Fprint(w, "\n--- stacktrace ---")
	//switch e := err.(type) {
	//case *errors.AnnotatedError:
	//	if e.OutputStackTrace() {
	//		fmt.Fprintf(w, "%+v\n", e.StackTrace())
	//	}
	//default:
	//	fmt.Fprintf(w, "%+v", err)
	//}
}
