package turbine

type Function interface {
	Process(r []Record) []Record
}
