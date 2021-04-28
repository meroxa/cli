package deprecated

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// CompletionCmd represents the completion command.
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
		RunE: func(cmd *cobra.Command, args []string) error {
			switch args[0] {
			case "bash":
				return cmd.Root().GenBashCompletion(os.Stdout)
			case "zsh":
				return cmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				return cmd.Root().GenFishCompletion(os.Stdout, true)
			case "powershell":
				return cmd.Root().GenPowerShellCompletion(os.Stdout)
			default:
				return fmt.Errorf("unexpected argument %q", args[0])
			}
		},
	}
}
