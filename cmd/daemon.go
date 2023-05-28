/*
Copyright © 2023 Andreas Pfurtscheller <andreas@fruits.co>
*/
package cmd

import (
	"context"

	"github.com/aplr/pubsub-emulator/app"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// daemonCmd represents the serve command
var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Start the operator daemon.",
	Run:   runDaemon,
}

func runDaemon(cmd *cobra.Command, args []string) {
	log.Infof("PubSub operator version %s starting", cmd.Root().Version)

	app := app.NewApp()

	err := app.Run(context.Background())

	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	rootCmd.AddCommand(daemonCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}