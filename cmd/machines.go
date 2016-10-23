package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/urfave/cli"

	"github.com/arigatomachine/cli/api"
	"github.com/arigatomachine/cli/config"
	"github.com/arigatomachine/cli/errs"
	"github.com/arigatomachine/cli/identity"
)

func init() {
	machines := cli.Command{
		Name:     "machines",
		Usage:    "Manage machine for an organization",
		Category: "MACHINES",
		Subcommands: []cli.Command{
			{
				Name:  "create",
				Usage: "Create a machine for an organization",
				Flags: []cli.Flag{
					orgFlag("org to generate keypairs for", true),
				},
				Action: chain(
					ensureDaemon, ensureSession, loadDirPrefs, loadPrefDefaults,
					setUserEnv, checkRequiredFlags, generateKeypairs,
				),
			},
		},
	}
	Cmds = append(Cmds, machines)
}

func createMachine(ctx *cli.Context) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	client := api.NewClient(cfg)
	c := context.Background()

	org, orgName, newOrg, err := SelectCreateOrg(c, client, ctx.String("org"))

	var orgID *identity.ID
	if !newOrg {
		if org == nil {
			return errs.NewExitError("Org not found.")
		}
		orgID = org.ID
	}

	args := ctx.Args()
	name := ""
	if len(args) > 0 {
		name = args[0]
	}

	label := "Machine name"
	autoAccept := name != ""
	name, err = NamePrompt(&label, name, autoAccept)
	if err != nil {
		return handleSelectError(err, projectCreateFailed)
	}

	if newOrg {
		org, err := client.Orgs.Create(c, orgName)
		orgID = org.ID
		if err != nil {
			return errs.NewErrorExitError("Could not create org", err)
		}

		err = generateKeypairsForOrg(c, ctx, client, org.ID, false)
		if err != nil {
			return err
		}

		fmt.Printf("Org %s created.\n\n", orgName)
	}

	_, err = createMachineByName(c, client, orgID, name)
	return err
}

func createMachineByName(c context.Context, client *api.Client, orgID *identity.ID, name string) (*api.ProjectResult, error) {
	machine, err := client.Machines.Create(c, orgID, name)
	if orgID == nil {
		return nil, errs.NewExitError("Org not found")
	}
	if err != nil {
		if strings.Contains(err.Error(), "resource exists") {
			return nil, errs.NewExitError("Machine already exists")
		}
		return nil, errs.NewErrorExitError(projectCreateFailed, err)
	}
	fmt.Printf("Machine %s created.\n", name)
	return project, nil
}
