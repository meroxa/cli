module github.com/meroxa/cli

go 1.13

require (
	github.com/alexeyco/simpletable v0.0.0-20200730140406-5bb24159ccfb
	github.com/danieljoos/wincred v1.1.0 // indirect
	github.com/docker/docker-credential-helpers v0.6.3
	github.com/fatih/color v1.9.0
	github.com/google/gops v0.3.10 // indirect
	github.com/gorilla/mux v1.7.3
	github.com/manifoldco/promptui v0.7.0
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/meroxa/meroxa-go v0.0.0-20210107153739-b88de43297a9
	github.com/mitchellh/go-homedir v1.1.0
	github.com/nmrshll/rndm-go v0.0.0-20170430161430-8da3024e53de
	github.com/pkg/browser v0.0.0-20180916011732-0a3d74bf9ce4 // indirect
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.6.2
	golang.org/x/oauth2 v0.0.0-20180821212333-d2e6202438be
	golang.org/x/sys v0.0.0-20200515095857-1151b9dac4a9 // indirect
)

replace github.com/meroxa/meroxa-go => ../meroxa-go
