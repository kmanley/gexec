# gexec

gexec is a little library containing a GracefulCommand object that adds some useful functionality to [Go](http://golang.org)'s [os/exec](http://golang.org/pkg/os/exec/) package:

 * Wait with timeout
 * Graceful Kill
 
## Creating a Command

Create a GracefulCommand by passing it an os.exec.Command object

```go
package main
import (
	"os/exec"
	"github.com/kmanley/gexec"
)
cmd := gexec.GracefulCommand(exec.Command("python", "path/app.py"), nil)
```

If you want to integrate gexec with your logging system you can pass an adapter
for whatever logger you're using. It just has to conform to io.Writer

```go
type FmtLogger struct {
}

func (g *FmtLogger) Write(s []byte) (int, error) {
	fmt.Println(string(s))
	return len(s), nil
}

cmd := gexec.GracefulCommand(exec.Command("python", "path/app.py"), &FmtLogger{})
```


## Wait with timeout

GracefulCommand.Wait accepts a time.Duration timeout. If the subprocess doesn't exit
before the timeout elapses, Wait returns gexec.ErrTimeout. Otherwise, the return 
value is the same as Cmd.Wait.

```go
cmd := gexec.GracefulCommand(exec.Command("python", "path/app.py"), nil)
err := cmd.Wait(5 * time.Second)
if err == gexec.ErrTimeout {
	// ...handle timeout
} else {
	// err is the same object returned by Cmd.Wait
}
```

You can pass 0 to wait forever. 

```go
cmd := gexec.GracefulCommand(exec.Command("python", "path/app.py"), nil)
err := cmd.Wait(0)
```

## Graceful Kill

GracefulCommand.Kill first attempts to send SIGINT to your subprocess. You can specify
how long to wait for the process to exit in response to SIGINT. If it fails to exit
in time, GracefulCommand.Kill will then send SIGKILL to terminate the process. 

If the process exits in time in response to SIGINT, you can get its exit code the same way
you would with Wait. 

```go
cmd := gexec.GracefulCommand(exec.Command("python", "path/app.py"), nil)
err := cmd.Kill(5 * time.Second)
// If process exited in response to SIGINT in < 5s, then 
// err.Error() == "exit status <N>"
```

If Kill times out waiting for the process to exit in response to SIGINT and subsequently sends SIGKILL,
then the returned error will indicate the process was killed

```go
cmd := gexec.GracefulCommand(exec.Command("python", "path/ignores_signals.py"), nil)
err := cmd.Kill(5 * time.Second)
// err.Error() == "signal: killed"
```

If you specified a log adapter when creating your GracefulCommand, Kill will write some
diagnostic information to it during the Kill process. 

### License

MIT License
