package turbine

type App interface {
	Run(Turbine) error
}

type Turbine interface {
	Resources(string) (Resource, error)
	Process(Records, Function) Records
	RegisterSecret(string) error
}
