package main

import (
	"fmt"
	"log"

	"github.com/billdev1958/CSV2Sql.git/migrator"
)

func main() {
	// Rutas absolutas (ajusta según tu carpeta real)
	filePathMedicines := "/home/billy/Desktop/ProyectosUaemex/CSV2Sql/lista_medicamentos.csv"
	filePathDiagnoses := "/home/billy/Desktop/ProyectosUaemex/CSV2Sql/diagnosticos_cie_1.csv"

	// 1. Conectar a la base de datos
	db, err := migrator.ConnectDB()
	if err != nil {
		log.Fatalf("Error al conectar con la base de datos: %v\n", err)
	}
	defer db.Close()

	// 2. Procesar (migrar) Medicinas
	if err := migrator.ProcessMedicines(filePathMedicines, db); err != nil {
		log.Fatalf("Error al procesar Medicinas: %v\n", err)
	}

	// 3. Procesar (migrar) Diagnósticos
	if err := migrator.ProcessDiagnoses(filePathDiagnoses, db); err != nil {
		log.Fatalf("Error al procesar Diagnósticos: %v\n", err)
	}

	fmt.Println("¡Migración finalizada con éxito!")
}
