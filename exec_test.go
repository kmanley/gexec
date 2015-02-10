package gexec

import (
	"bytes"
	"fmt"
	_ "github.com/davecgh/go-spew/spew"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
	"time"
)

// thank you https://github.com/benbjohnson/testing for these helper functions
// assert fails the test if the condition is false.
func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.Fail()
	}
}

// ok fails the test if an err is not nil.
func assertnil(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.Fail()
	}
}

// err fails the test if err is nil.
func asserterr(tb testing.TB, err error) {
	if err == nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: expected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.Fail()
	}
}

// equals fails the test if exp is not equal to act.
func assertequals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.Fail()
	}
}

type Logger struct {
}

func (g *Logger) Write(s []byte) (int, error) {
	fmt.Println(string(s))
	return len(s), nil
}

func TestExitSuccess(t *testing.T) {
	cmd := GracefulCommand(exec.Command("python", "test/exitcode_0.py"), &Logger{})
	cmd.Start()
	err := cmd.Wait(0)
	//spew.Dump(res)
	assertnil(t, err)
	assertequals(t, "exit status 0", cmd.ProcessState.String())
}

func TestExitFail(t *testing.T) {
	cmd := GracefulCommand(exec.Command("python", "test/exitcode_123.py"), &Logger{})
	cmd.Start()
	err := cmd.Wait(0)
	//spew.Dump(err)
	asserterr(t, err)
	assertequals(t, "exit status 123", cmd.ProcessState.String())
}

func TestWaitThenKill1(t *testing.T) {
	cmd := GracefulCommand(exec.Command("python", "test/exits_on_sigint.py"), &Logger{})
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Start()
	err := cmd.Wait(1 * time.Second)
	assertequals(t, ErrTimeout, err)
	err = cmd.Kill(5 * time.Second)
	assertequals(t, "exit status 42", err.Error())
	assertequals(t, "exit status 42", cmd.ProcessState.String())
	assertequals(t, "got SIGINT; exiting...\n", stdout.String())
}

func TestWaitThenKill2(t *testing.T) {
	cmd := GracefulCommand(exec.Command("python", "test/ignores_sigint.py"), &Logger{})
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Start()
	err := cmd.Wait(1 * time.Second)
	assertequals(t, ErrTimeout, err)
	err = cmd.Kill(5 * time.Second)
	assertequals(t, "signal: killed", err.Error())
	assertequals(t, "signal: killed", cmd.ProcessState.String())
	assertequals(t, "got SIGINT; ignoring it...\n", stdout.String())
}
