package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

// completionCmd represents the completion command
func CompletionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate completion script",
		Long: `To load completions:

	Bash:
	
	$ source <(meroxa completion bash)
	
	# To load completions for each session, execute once:
	Linux:
	  $ meroxa completion bash > /etc/bash_completion.d/meroxa
	MacOS:
	  $ meroxa completion bash > /usr/local/etc/bash_completion.d/meroxa
	
	Zsh:
	
	# If shell completion is not already enabled in your environment you will need
	# to enable it.  You can execute the following once:
	
	$ echo "autoload -U compinit; compinit" >> ~/.zshrc
	
	# To load completions for each session, execute once:
	$ meroxa completion zsh > "${fpath[1]}/_meroxa"
	
	# You will need to start a new shell for this setup to take effect.
	
	Fish:
	
	$ meroxa completion fish | source
	
	# To load completions for each session, execute once:
	$ meroxa completion fish > ~/.config/fish/completions/meroxa.fish
	`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.ExactValidArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			switch args[0] {
			case "bash":
				cmd.Root().GenBashCompletion(os.Stdout)
			case "zsh":
				cmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				cmd.Root().GenFishCompletion(os.Stdout, true)
			case "powershell":
				cmd.Root().GenPowerShellCompletion(os.Stdout)
			}
		},
	}
}
