package turbine

type CLI interface {
	Upgrade(appPath string, vendor bool) error
}
