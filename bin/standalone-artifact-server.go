package main

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	// "path/filepath"

	// "github.com/nektos/act/cmd"
	"github.com/nektos/act/pkg/artifacts"
	"github.com/nektos/act/pkg/common"
	"github.com/spf13/cobra"
)

var version string = "0.0.0"

var input Input

func main() {
	ctx, cancel := common.CreateGracefulJobCancellationContext()
	defer cancel()
	// input = new(Input)

	// input := new(Input)
	rootCmd := createRootCommand(ctx, version)

	var exitFunc = os.Exit
	if err := rootCmd.Execute(); err != nil {
		exitFunc(1)
	}
}

// Input contains the input for the root command
type Input struct {
	artifactServerPath                 string
	artifactServerAddr                 string
	artifactServerPort                 string
	noCacheServer                      bool
	cacheServerPath                    string
	cacheServerExternalURL             string
	cacheServerAddr                    string
	cacheServerPort                    uint16
	actionCachePath                    string
	networkName                        string
}

func createActionRuntimeVars() {
	actionsRuntimeURL := fmt.Sprintf("http://%s:%s/", input.artifactServerAddr, input.artifactServerPort)
	runID := int64(1)
	actionsRuntimeToken, _ := common.CreateAuthorizationToken(runID, runID, runID)

	fmt.Printf("Now provide these arguments to act:\n")
	fmt.Printf("--env ACTIONS_RUNTIME_URL=%s --env ACTIONS_RESULTS_URL=%s --env ACTIONS_RUNTIME_TOKEN=%s\n", actionsRuntimeURL, actionsRuntimeURL, actionsRuntimeToken)

	fmt.Printf("\nAlternatively, you can put these environment variables into a environment file and use the --env-file option for act.\n")
	fmt.Printf("ACTIONS_RUNTIME_URL=%s\n", actionsRuntimeURL)
	fmt.Printf("ACTIONS_RESULTS_URL=%s\n", actionsRuntimeURL)
	fmt.Printf("ACTIONS_RUNTIME_TOKEN=%s\n", actionsRuntimeToken)
	fmt.Printf("\n")
}
//nolint:gocyclo
func newRunCommand(ctx context.Context) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		createActionRuntimeVars()

		cancel := artifacts.Serve(ctx, input.artifactServerPath, input.artifactServerAddr, input.artifactServerPort)
		defer cancel()
		done := make(chan os.Signal, 1)
		signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
		fmt.Println("Blocking, press ctrl+c to continue...")
		<-done  // Will block here until user hits ctrl+c
	}
}

func createRootCommand(ctx context.Context, version string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:               "standalone-artifacts-server [flags]",
		Short:             "Run the act artifacts server standalone",
		Args:              cobra.MaximumNArgs(0),
		Version:           version,
		SilenceUsage:      true,
		Run: newRunCommand(ctx),
	}

	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringVarP(&input.artifactServerPath, "artifact-server-path", "", "", "Defines the path where the artifact server stores uploads and retrieves downloads from. If not specified the artifact server will not start.")
	rootCmd.PersistentFlags().StringVarP(&input.artifactServerAddr, "artifact-server-addr", "", common.GetOutboundIP().String(), "Defines the address to which the artifact server binds.")
	rootCmd.PersistentFlags().StringVarP(&input.artifactServerPort, "artifact-server-port", "", "34567", "Defines the port where the artifact server listens.")
	// rootCmd.PersistentFlags().BoolVarP(&input.noCacheServer, "no-cache-server", "", false, "Disable cache server")
	// rootCmd.PersistentFlags().StringVarP(&input.cacheServerPath, "cache-server-path", "", filepath.Join(cmd.CacheHomeDir, "actcache"), "Defines the path where the cache server stores caches.")
	// rootCmd.PersistentFlags().StringVarP(&input.cacheServerExternalURL, "cache-server-external-url", "", "", "Defines the external URL for if the cache server is behind a proxy. e.g.: https://act-cache-server.example.com. Be careful that there is no trailing slash.")
	// rootCmd.PersistentFlags().StringVarP(&input.cacheServerAddr, "cache-server-addr", "", common.GetOutboundIP().String(), "Defines the address to which the cache server binds.")
	// rootCmd.PersistentFlags().Uint16VarP(&input.cacheServerPort, "cache-server-port", "", 0, "Defines the port where the artifact server listens. 0 means a randomly available port.")
	// rootCmd.PersistentFlags().StringVarP(&input.actionCachePath, "action-cache-path", "", filepath.Join(cmd.CacheHomeDir, "act"), "Defines the path where the actions get cached and host workspaces created.")
	// rootCmd.PersistentFlags().StringVarP(&input.networkName, "network", "", "host", "Sets a docker network name. Defaults to host.")
	return rootCmd
}
