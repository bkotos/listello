package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"

	"github.com/bkotos/listello/internal/bootstrap"
	instanceadapter "github.com/bkotos/listello/internal/listello-instance-context/adapter"
	instanceapp "github.com/bkotos/listello/internal/listello-instance-context/application"
	ppadapter "github.com/bkotos/listello/internal/personal-productivity-context/adapter"
	"github.com/bkotos/listello/internal/sqlite"
)

func main() {
	var port int
	var host string

	cmd := &cobra.Command{
		Use:   "listello-server",
		Short: "Listello HTTP API server",
		RunE: func(cmd *cobra.Command, _ []string) error {
			workspace := sqlite.NewWorkspaceDB()
			defer workspace.Close()

			eventLog := bootstrap.MustOpenEventLog("domain_events.log")
			defer eventLog.Close()

			listService := bootstrap.NewListService(workspace, eventLog)
			itemService := bootstrap.NewItemService(workspace, eventLog)
			locatorPath, err := instanceadapter.ListelloInstanceLocatorPath()
			if err != nil {
				return err
			}
			repo, err := instanceadapter.NewListelloInstanceRepository(locatorPath)
			if err != nil {
				return err
			}
			instanceService := newWorkspaceOpeningInstanceService(
				instanceapp.NewListelloInstanceService(
					repo,
					instanceadapter.NewFilesystemPersistenceAdapter(),
					ppadapter.NewLoggingEventPublisher(eventLog),
				),
				workspace,
			)
			instance, err := instanceService.GetInstance()
			if err != nil {
				return err
			}
			if instance != nil && instance.Persistence.IsInitialized() {
				if err := openWorkspaceDB(workspace, instance.Persistence.Location); err != nil {
					return err
				}
			}
			addr := fmt.Sprintf("%s:%d", host, port)
			cmd.Printf("listening on http://%s\n", addr)
			return http.ListenAndServe(addr, newAPIServer(listService, itemService, instanceService))
		},
	}

	cmd.Flags().IntVarP(&port, "port", "p", 8080, "port to listen on")
	cmd.Flags().StringVar(&host, "host", "0.0.0.0", "host to bind to")

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
