package cmd

import (
	"io"
	"os"

	argo "github.com/argoproj/argo/cmd/argo/commands"
	"github.com/spf13/cobra"
	kubeadm "k8s.io/kubernetes/cmd/kubeadm/app/cmd"
	kubectl "k8s.io/kubernetes/pkg/kubectl/cmd"

	// Import to initialize client auth plugins.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"github.com/ian-howell/airshipctl/cmd/bootstrap"
	"github.com/ian-howell/airshipctl/cmd/completion"
	"github.com/ian-howell/airshipctl/pkg/environment"
	"github.com/ian-howell/airshipctl/pkg/log"
)

// NewAirshipCTLCommand creates a root `airshipctl` command with the default commands attached
func NewAirshipCTLCommand(out io.Writer) (*cobra.Command, *environment.AirshipCTLSettings, error) {
	rootCmd, settings, err := NewRootCmd(out)
	return AddDefaultAirshipCTLCommands(rootCmd, settings), settings, err
}

// NewRootCmd creates the root `airshipctl` command. All other commands are
// subcommands branching from this one
func NewRootCmd(out io.Writer) (*cobra.Command, *environment.AirshipCTLSettings, error) {
	settings := &environment.AirshipCTLSettings{}
	rootCmd := &cobra.Command{
		Use:           "airshipctl",
		Short:         "airshipctl is a unified entrypoint to various airship components",
		SilenceErrors: true,
		SilenceUsage:  true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			log.Init(settings.Debug, cmd.OutOrStderr())
		},
	}
	rootCmd.SetOutput(out)
	rootCmd.AddCommand(NewVersionCommand())

	settings.InitFlags(rootCmd)

	return rootCmd, settings, nil
}

// AddDefaultAirshipCTLCommands is a convenience function for adding all of the
// default commands to airshipctl
func AddDefaultAirshipCTLCommands(cmd *cobra.Command, settings *environment.AirshipCTLSettings) *cobra.Command {
	cmd.AddCommand(argo.NewCommand())
	cmd.AddCommand(bootstrap.NewBootstrapCommand(settings))
	cmd.AddCommand(completion.NewCompletionCommand())
	cmd.AddCommand(kubectl.NewDefaultKubectlCommand())
	// Should we use cmd.OutOrStdout?
	cmd.AddCommand(kubeadm.NewKubeadmCommand(os.Stdin, os.Stdout, os.Stderr))

	kustomizeCmd, _, err := cmd.Find([]string{"kubectl", "kustomize"})
	if err != nil {
		log.Fatalf("Unable to find subcommand '%s'", err.Error())
	}

	cmd.AddCommand(kustomizeCmd)

	return cmd
}
