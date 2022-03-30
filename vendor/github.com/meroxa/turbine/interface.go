package turbine

type App interface {
	Run(Turbine) error
}

type Turbine interface {
	Resources(string) (Resource, error)
	Process(Records, Function) (Records, RecordsWithErrors)
	RegisterSecret(string) error
}
