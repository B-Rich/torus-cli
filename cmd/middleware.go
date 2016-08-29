package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/urfave/cli"

	"github.com/arigatomachine/cli/api"
	"github.com/arigatomachine/cli/apitypes"
)

// Chain allows easy sequential calling of BeforeFuncs and AfterFuncs.
// Chain will exit on the first error seen.
// XXX Chain is only public while we need it for passthrough.go
func Chain(funcs ...func(*cli.Context) error) func(*cli.Context) error {
	return func(ctx *cli.Context) error {

		for _, f := range funcs {
			err := f(ctx)
			if err != nil {
				return err
			}
		}

		return nil
	}
}

// EnsureDaemon ensures that the daemon is running, and is the correct version,
// before a command is exeucted.
// the daemon will be started/restarted once, to try and launch the latest
// version.
// XXX EnsureDaemon is only public while we need it for passthrough.go
func EnsureDaemon(ctx *cli.Context) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	proc, err := findDaemon(cfg)
	if err != nil {
		return err
	}

	spawned := false

	if proc == nil {
		err := spawnDaemon()
		if err != nil {
			return err
		}

		spawned = true
	}

	client := api.NewClient(cfg)

	var v *apitypes.Version
	increment := 5 * time.Millisecond
	for d := increment; d < 1*time.Second; d += increment {
		v, err = client.Version.Get(context.Background())
		if err == nil {
			break
		}
		time.Sleep(d)
	}

	if err != nil {
		return cli.NewExitError("Error communicating with the daemon: "+err.Error(), -1)
	}

	if v.Version == cfg.Version {
		return nil
	}

	if spawned {
		return cli.NewExitError("The daemon version is incorrect. Check for stale processes", -1)
	}

	fmt.Println("The daemon version is out of date and is being restarted.")
	fmt.Println("You will need to login again.")

	_, err = stopDaemon(proc)
	if err != nil {
		return err
	}

	return EnsureDaemon(ctx)
}

// EnsureSession ensures that the user is logged in with the daemon and has a
// valid session. If not, it will attempt to log the user in via environment
// variables. If they do not exist, of the login fails, it will abort the
// command.
// XXX EnsureSession is only public while we need it for passthrough.go
func EnsureSession(ctx *cli.Context) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	client := api.NewClient(cfg)
	_, err = client.Session.Get(context.Background())

	hasSession := true
	if err != nil {
		if cerr, ok := err.(*apitypes.Error); ok {
			if cerr.Type == apitypes.UnauthorizedError {
				hasSession = false
			}
		}
		if hasSession {
			return cli.NewExitError("Error communicating with the daemon: "+err.Error(), -1)
		}
	}

	if hasSession {
		return nil
	}

	email, hasEmail := os.LookupEnv("AG_EMAIL")
	password, hasPassword := os.LookupEnv("AG_PASSWORD")

	if hasEmail && hasPassword {
		fmt.Println("Attempting to login with email: " + email)

		err := client.Session.Login(context.Background(), email, password)
		if err != nil {
			fmt.Println("Could not log in: " + err.Error())
		} else {
			return nil
		}
	}

	msg := "You must be logged in to run '" + ctx.Command.FullName() + "'.\n" +
		"Login using 'login' or create an account using 'signup'."
	return cli.NewExitError(msg, -1)
}