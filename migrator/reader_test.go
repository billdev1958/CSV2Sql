package migrator_test

import (
	"testing"

	"github.com/billdev1958/CSV2Sql.git/migrator"
)

func BenchmarkReaderMedicine(b *testing.B) {
	filePath := "/home/billy/Desktop/ProyectosUaemex/CSV2Sql/lista_medicamentos.csv"

	b.ResetTimer()

	for i := 0; i < b.N; i++ { // Go controla automáticamente b.N
		_, err := migrator.ReaderMedicine(filePath)
		if err != nil {
			b.Errorf("error al leer el CSV: %v", err)
		}
	}
}

func BenchmarkReaderDiagnose(b *testing.B) {
	filePath := "/home/billy/Desktop/ProyectosUaemex/CSV2Sql/diagnosticos_cie_1.csv"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := migrator.ReaderDiagnosis(filePath)
		if err != nil {
			b.Errorf("error al leer el CSV: %v", err)
		}
	}
}
