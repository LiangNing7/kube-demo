/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"

	"github.com/spf13/cobra"
	gserver "k8s.io/apiserver/pkg/server"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kube-demo",
	Short: "An aggregated API Server",
	Long:  `This is an aggregated API Server, wrote manually`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

func NewCommandStartServer(stopCh <-chan struct{}) *cobra.Command {
	options := NewServerOptions()
	rootCmd.RunE = func(c *cobra.Command, args []string) error {
		if err := options.Complete(); err != nil {
			return err
		}
		if err := options.Validate(); err != nil {
			return err
		}
		if err := run(options, stopCh); err != nil {
			return err
		}
		return nil
	}
	flags := rootCmd.Flags()
	options.RecommendedOptions.AddFlags(flags)
	return rootCmd
}

func run(o *ServerOptions, stopCh <-chan struct{}) error {
	c, err := o.Config()
	if err != nil {
		return err
	}

	s, err := c.Complete().NewServer()
	if err != nil {
		return err
	}

	s.GenericAPIServer.AddPostStartHook("start-provision-server-informers",
		func(context gserver.PostStartHookContext) error {
			c.GenericConfig.SharedInformerFactory.Start(context.Done())
			o.SharedInformerFactory.Start(context.Done())
			return nil
		})
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-stopCh
		cancel()
	}()
	return s.GenericAPIServer.PrepareRun().PrepareRun().PrepareRun().RunWithContext(ctx)
}
