module github.com/meroxa/cli

go 1.13

require (
	github.com/gorilla/mux v1.7.3
	github.com/manifoldco/promptui v0.7.0
	github.com/meroxa/meroxa-go v0.0.0-20200407000008-ab97e83a1bf2
	github.com/mitchellh/go-homedir v1.1.0
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.6.2
)

replace github.com/meroxa/meroxa-go => ../meroxa-go
