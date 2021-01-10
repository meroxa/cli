module github.com/meroxa/cli

go 1.15

require (
	github.com/alexeyco/simpletable v0.0.0-20200730140406-5bb24159ccfb
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/fatih/color v1.9.0
	github.com/gorilla/mux v1.7.3
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/meroxa/meroxa-go v0.0.0-20210107153739-b88de43297a9
	github.com/mitchellh/go-homedir v1.1.0
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.3
	github.com/spf13/viper v1.6.2
	github.com/stretchr/testify v1.5.1 // indirect
	golang.org/x/sys v0.0.0-20200515095857-1151b9dac4a9 // indirect
)

replace github.com/meroxa/meroxa-go => ../meroxa-go
