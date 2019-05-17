package cmd

import (
  "github.com/spf13/cobra"
  "github.com/reactiveops/dd-manager/pkg/config"
  "github.com/reactiveops/dd-manager/pkg/controller"
  log "github.com/sirupsen/logrus"
  "os"
)


func RootCmd() *cobra.Command {
  root := &cobra.Command {
    Use: "dd-manager",
    Short: "Kubernetes datadog monitor manager",
    Long: "A kubernetes agent that manages datadog monitors.",
    Run: run,
  }
  return root
}


func loadConfig(cmd *cobra.Command)(*config.Config) {
  log.SetReportCaller(true)
  log.SetOutput(os.Stdout)

  config := config.New()
  return config
}


func run(cmd *cobra.Command, args []string) {
  conf := loadConfig(cmd)
  controller.NewController(conf)
}
