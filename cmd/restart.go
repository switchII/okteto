package cmd

import (
	"github.com/okteto/okteto/pkg/k8s/pods"
	"github.com/okteto/okteto/pkg/log"

	k8Client "github.com/okteto/okteto/pkg/k8s/client"
	"github.com/okteto/okteto/pkg/model"

	"github.com/spf13/cobra"
)

//Restart restarts the pods of a given dev mode deployment
func Restart() *cobra.Command {
	var namespace string
	// get the name from it
	var devPath string

	cmd := &cobra.Command{
		Use:   "restart",
		Short: "Restarts the pods of your Okteto Environment",
		RunE: func(cmd *cobra.Command, args []string) error {
			dev, err := loadDev(devPath)
			if err != nil {
				return err
			}
			if namespace != "" {
				dev.Namespace = namespace
			}

			if err := executeRestart(dev); err != nil {
				return err
			}
			log.Success("Okteto Environment restarted")

			return nil
		},
	}

	cmd.Flags().StringVarP(&devPath, "file", "f", defaultManifest, "path to the manifest file")
	cmd.Flags().StringVarP(&namespace, "namespace", "n", "", "namespace where the exec command is executed")

	return cmd
}

func executeRestart(dev *model.Dev) error {
	client, _, namespace, err := k8Client.GetLocal()
	if err != nil {
		return err
	}

	if dev.Namespace == "" {
		dev.Namespace = namespace
	}

	restarted := map[string]bool{}
	for _, s := range dev.Services {
		if _, ok := restarted[s.Name]; ok {
			continue
		}
		if err := pods.Restart(s.Name, dev.Namespace, client); err != nil {
			return err
		}
		restarted[s.Name] = true
	}

	return nil
}