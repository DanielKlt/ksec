package secrets

import (
	"context"
	"fmt"
	"os"

	k8s "github.com/DanielKlt/ksec/pkg/clients"
	"github.com/DanielKlt/ksec/pkg/environment"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	core_v1 "k8s.io/api/core/v1"
	"k8s.io/cli-runtime/pkg/printers"
)

var GetCmd = &cobra.Command{
	Use:     "get",
	Aliases: []string{"g"},
	Short:   "get secret cmd",
	Long:    "Get a decoded secret from k8s",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := environment.GetEnv().GetViper().BindPFlag("namespace", cmd.PersistentFlags().Lookup("namespace")); err != nil {
			return err
		}
		if err := environment.GetEnv().GetViper().BindPFlag("kubeconfig", cmd.PersistentFlags().Lookup("kubeconfig")); err != nil {
			return err
		}
		if err := environment.GetEnv().GetViper().BindPFlag("name", cmd.PersistentFlags().Lookup("name")); err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return showSecret(context.Background())
	},
}

func showSecret(ctx context.Context) error {
	var kcfgPath string
	kcfgPath = environment.GetEnv().GetViper().GetString("kubeconfig")
	environment.GetEnv().GetLogger().Debug("kubeconfig from viper", zap.String("kubeconfig", kcfgPath))
	if kcfgPath == "" {
		if v, ok := os.LookupEnv("KUBECONFIG"); ok {
			kcfgPath = v
			environment.GetEnv().GetLogger().Debug("using kubeconfig from env", zap.String("kubeconfig", kcfgPath))
		} else {
			return fmt.Errorf("kubeconfig not set")
		}
	}

	environment.GetEnv().GetLogger().Debug("creating kubernetes client", zap.String("kubeconfig", kcfgPath))
	c, err := k8s.NewK8sClient(kcfgPath)
	if err != nil {
		return err
	}

	if environment.GetEnv().GetViper().GetString("namespace") == "" {
		environment.GetEnv().GetLogger().Debug("namespace not set, using default namespace")
		environment.GetEnv().GetViper().Set("namespace", "default")
	}

	if environment.GetEnv().GetViper().GetString("name") == "" {
		environment.GetEnv().GetLogger().Debug("secret name not set")
		return fmt.Errorf("secret name not set")
	}

	s, err := c.GetSecret(ctx, environment.GetEnv().GetViper().GetString("namespace"), environment.GetEnv().GetViper().GetString("name"))
	if err != nil {
		return err
	}

	return prettyPrint(s)
}

func init() {
	GetCmd.PersistentFlags().String("namespace", "", "namespace of the secret")
	GetCmd.PersistentFlags().String("kubeconfig", "", "path to kubeconfig")
	GetCmd.PersistentFlags().String("name", "", "name of the secret")
}

func prettyPrint(s *core_v1.Secret) error {
	var tmpl = []byte(`apiVersion: {{ (index .metadata.managedFields 0).apiVersion }}
metadata:
  name: {{ .metadata.name }}
  namespace: {{ .metadata.namespace }}
  creationTimestamp: {{ .metadata.creationTimestamp }}
  resourceVersion: {{ .metadata.resourceVersion }}
  uid: {{ .metadata.uid }}
{{- if .data }}
data:
{{- range $k, $v := .data }}
  {{ $k }}: {{ $v | base64decode }}
{{- end }}
{{- else if .stringData }}
stringData:
{{- range $k, $v := .stringData }}
  {{ $k }}: {{ $v }}
{{- end }}
{{- end }}
type: {{ .type }}
`)
	p, err := printers.NewGoTemplatePrinter(tmpl)
	if err != nil {
		return err
	}
	p.PrintObj(s, os.Stdout)

	return nil
}
