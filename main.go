package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/billdev1958/CSV2Sql.git/migrator"
)

func main() {
	// 1. Mostrar menú y obtener DSN del usuario
	fmt.Println("--- Migrador de CSV a PostgreSQL ---")
	fmt.Print("Por favor, introduce tu DSN (cadena de conexión a la base de datos) y presiona Enter:\n> ")

	reader := bufio.NewReader(os.Stdin)
	dsn, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Error al leer la entrada del usuario: %v", err)
	}
	dsn = strings.TrimSpace(dsn) // Limpiar espacios en blanco y saltos de línea

	// Rutas a los archivos CSV (ajusta si es necesario)
	filePathMedicines := "./lista_medicamentos.csv"
	filePathDiagnoses := "./diagnosticos_cie_1.csv"

	// 2. Conectar a la base de datos con el DSN proporcionado
	db, err := migrator.ConnectDB(dsn)
	if err != nil {
		log.Fatalf("Error al conectar con la base de datos: %v\n", err)
	}
	defer db.Close()

	// 3. Procesar (migrar) Medicinas
	fmt.Printf("\nProcesando archivo de medicinas: %s\n", filePathMedicines)
	if err := migrator.ProcessMedicines(filePathMedicines, db); err != nil {
		log.Fatalf("Error al procesar Medicinas: %v\n", err)
	}

	// 4. Procesar (migrar) Diagnósticos
	fmt.Printf("\nProcesando archivo de diagnósticos: %s\n", filePathDiagnoses)
	if err := migrator.ProcessDiagnoses(filePathDiagnoses, db); err != nil {
		log.Fatalf("Error al procesar Diagnósticos: %v\n", err)
	}

	fmt.Println("\n¡Migración finalizada con éxito!")
}
