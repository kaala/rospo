package cmd

import (
	"log"

	"github.com/ferama/rospo/pkg/conf"
	"github.com/ferama/rospo/pkg/pipe"
	"github.com/ferama/rospo/pkg/sshc"
	"github.com/ferama/rospo/pkg/sshd"
	"github.com/ferama/rospo/pkg/tun"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:  "rospo config_file_path.yaml",
	Long: "The tool to create relieable ssh tunnels.",
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := conf.LoadConfig(args[0])
		if err != nil {
			log.Fatalln(err)
		}

		if conf.Pipe != nil {
			for _, f := range conf.Pipe {
				go pipe.NewPipe(f).Start()
			}
		}

		if conf.SshD != nil {
			if conf.Tunnel != nil {
				go sshd.NewSshServer(conf.SshD).Start()
			} else {
				sshd.NewSshServer(conf.SshD).Start()
			}
		}

		if conf.Tunnel != nil && len(conf.Tunnel) > 0 {
			if conf.SshClient == nil {
				log.Fatalln("You need to configure sshclient section to support tunnels")
			}
			client := sshc.NewSshConnection(conf.SshClient)
			go client.Start()
			for idx, c := range conf.Tunnel {
				if idx == len(conf.Tunnel)-1 {
					tun.NewTunnel(client, c).Start()
				} else {
					go tun.NewTunnel(client, c).Start()
				}
			}
		}
	},
}

// Execute executes the root command
func Execute() error {
	return rootCmd.Execute()
}
