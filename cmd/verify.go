package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/urfave/cli"

	"github.com/arigatomachine/cli/api"
	"github.com/arigatomachine/cli/config"
	"github.com/arigatomachine/cli/errs"
)

func init() {
	version := cli.Command{
		Name:      "verify",
		ArgsUsage: "<code>",
		Usage:     "Verify the email address for your account",
		Category:  "ACCOUNT",
		Action: chain(
			ensureDaemon, ensureSession, verifyEmailCmd,
		),
	}
	Cmds = append(Cmds, version)
}

func verifyEmailCmd(ctx *cli.Context) error {
	return verifyEmail(ctx, nil, false)
}

func verifyEmail(ctx *cli.Context, code *string, subCommand bool) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}
	var verifyCode string
	if code != nil {
		verifyCode = *code
	}

	client := api.NewClient(cfg)
	c := context.Background()

	// Code is nil, check args for value
	if !subCommand {
		args := ctx.Args()
		if len(args) > 1 {
			return errs.NewUsageExitError("Too many arguments", ctx)
		}
		if len(args) != 1 {
			return errs.NewUsageExitError("Missing verification code", ctx)
		}

		verifyCode = args[0]
	}

	err = client.Users.VerifyEmail(c, verifyCode)
	if err != nil {
		if strings.Contains(err.Error(), "wrong user state: active") {
			return errs.NewExitError("Email already verified :)")
		}
		return errs.NewExitError("Email verification failed, please try again.")
	}

	fmt.Println("Your email is now verified.")
	return nil
}
