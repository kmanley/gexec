package gexec

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
)

var ErrTimeout = errors.New("Timeout elapsed")

type Done chan error

type GracefulCmd struct {
	*exec.Cmd
	done Done
	log  io.Writer
}

func GracefulCommand(c *exec.Cmd, log io.Writer) *GracefulCmd {
	return &GracefulCmd{c, nil, log}
}

func (this *GracefulCmd) Log(msg string) {
	if this.log != nil {
		this.log.Write([]byte(msg))
	}
}

func (this *GracefulCmd) Start() error {
	err := this.Cmd.Start()
	if err != nil {
		return err
	}
	this.done = make(Done, 1)
	go func() {
		this.done <- this.Cmd.Wait()
	}()
	return nil
}

func (this *GracefulCmd) Wait(timeout time.Duration) error {
	var err error
	if timeout == 0 {
		err = <-this.done
		return err
	} else {
		select {
		case <-time.After(timeout):
			{
				return ErrTimeout
			}
		case err = <-this.done:
			{
				return err
			}
		}
	}
}

func (this *GracefulCmd) Kill(sigIntWaitTime time.Duration) error {
	err := this.Cmd.Process.Signal(os.Interrupt)
	if err != nil {
		this.Log(fmt.Sprintf("sending SIGINT to %s failed: %s", this.Cmd.Process.Pid, err))
		// fall through to sending SIGKILL below
	} else {
		this.Log(fmt.Sprintf("sent SIGINT to %d; waiting up to %s for exit...", this.Cmd.Process.Pid, sigIntWaitTime.String()))
		select {
		case err = <-this.done:
			{
				// process exited
				this.Log(fmt.Sprintf("pid %d exited", this.Cmd.Process.Pid))
				return err
			}
		case <-time.After(sigIntWaitTime):
			{
				// SIGINT timeout elapsed
				this.Log(fmt.Sprintf("timed out waiting for pid %d to handle SIGINT", this.Cmd.Process.Pid))
				break
			}
		}
	}
	this.Log(fmt.Sprintf("pid %d still running; sending SIGKILL", this.Cmd.Process.Pid))
	err = this.Cmd.Process.Kill()
	if err != nil {
		return err
	}
	err = <-this.done
	return err
}
