package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"

	"github.com/bkotos/listello/internal/bootstrap"
)

func main() {
	var port int
	var host string

	cmd := &cobra.Command{
		Use:   "listello-server",
		Short: "Listello HTTP API server",
		RunE: func(cmd *cobra.Command, _ []string) error {
			db := bootstrap.MustOpenDB("listello.db")
			defer db.Close()

			eventLog := bootstrap.MustOpenEventLog("domain_events.log")
			defer eventLog.Close()

			listService := bootstrap.NewListService(db, eventLog)
			itemService := bootstrap.NewItemService(db, eventLog)
			addr := fmt.Sprintf("%s:%d", host, port)
			cmd.Printf("listening on http://%s\n", addr)
			return http.ListenAndServe(addr, newAPIServer(listService, itemService))
		},
	}

	cmd.Flags().IntVarP(&port, "port", "p", 8080, "port to listen on")
	cmd.Flags().StringVar(&host, "host", "0.0.0.0", "host to bind to")

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
