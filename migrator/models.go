package migrator

type Medicine struct {
	Substance           string
	Presentation        string
	RouteAdministration string
	Dose                string
	Quantity            string
	Frequency           string
}

type Diagnoses struct {
	Key       string
	Diagnosis string
}
