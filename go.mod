module github.com/meroxa/cli

go 1.13

require (
	github.com/alexeyco/simpletable v0.0.0-20200730140406-5bb24159ccfb
	github.com/google/gops v0.3.10 // indirect
	github.com/gorilla/mux v1.7.3
	github.com/manifoldco/promptui v0.7.0
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/meroxa/meroxa-go v0.0.0-20200807224254-6d5da351027d
	github.com/mitchellh/go-homedir v1.1.0
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.6.2
)

replace github.com/meroxa/meroxa-go => ../meroxa-go
