package migrator

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ReaderMedicine(pathCsv string) ([]Medicine, error) {

	file, err := os.Open(pathCsv)
	if err != nil {
		return nil, fmt.Errorf("error al abrir el archivo .csv: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error al leer el csv: %w", err)
	}

	medicines := make([]Medicine, len(records))

	for i, record := range records {
		if len(record) < 6 {
			fmt.Printf("Registro incompleto en la línea %d: %v\n", i+1, record)
			continue
		}
		medicines[i] = Medicine{
			Substance:           record[0],
			Presentation:        record[1],
			RouteAdministration: record[2],
			Dose:                record[3],
			Quantity:            record[4],
			Frequency:           record[5],
		}
	}

	return medicines, nil
}

func ReaderDiagnosis(pathCsv string) ([]Diagnoses, error) {
	file, err := os.Open(pathCsv)
	if err != nil {
		return nil, fmt.Errorf("error al abrir el archivo .csv: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error al leer el csv: %w", err)
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("el archivo CSV está vacío")
	}

	startIndex := 1
	if len(records[0]) < 2 || records[0][0] != "Key" {
		startIndex = 0
	}

	diagnoses := make([]Diagnoses, 0, len(records)-startIndex)

	for i := startIndex; i < len(records); i++ {
		record := records[i]

		if len(record) < 2 {
			fmt.Printf("Registro incompleto en la línea %d: %v\n", i+1, record)
			continue
		}

		diagnoses = append(diagnoses, Diagnoses{
			Key:       record[0],
			Diagnosis: record[1],
		})
	}

	return diagnoses, nil
}

func ConnectDB() (*pgxpool.Pool, error) {
	dsn := "postgres://root:secret@localhost:5432/university_db?sslmode=disable"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("error al conectar con la base de datos: %w", err)
	}
	return pool, nil
}

func ProcessMedicines(pathCsv string, db *pgxpool.Pool) error {
	medicines, err := ReaderMedicine(pathCsv)
	if err != nil {
		return fmt.Errorf("error al leer medicinas: %w", err)
	}

	ctx := context.Background()
	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error al iniciar transacción: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `
        INSERT INTO medicines (
            substance, 
            presentation, 
            route_administration, 
            dose, 
            quantity, 
            frequency
        ) VALUES ($1, $2, $3, $4, $5, $6)
    `

	for _, med := range medicines {
		_, err := tx.Exec(ctx, query,
			med.Substance,
			med.Presentation,
			med.RouteAdministration,
			med.Dose,
			med.Quantity,
			med.Frequency,
		)
		if err != nil {
			return fmt.Errorf("error al insertar medicina: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("error al confirmar transacción: %w", err)
	}

	fmt.Println("[OK] Medicinas insertadas correctamente.")
	return nil
}

func ProcessDiagnoses(pathCsv string, db *pgxpool.Pool) error {
	diagnoses, err := ReaderDiagnosis(pathCsv)
	if err != nil {
		return fmt.Errorf("error al leer diagnósticos: %w", err)
	}

	ctx := context.Background()
	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error al iniciar transacción: %w", err)
	}
	defer tx.Rollback(ctx)

	// Definir columnas de la tabla
	columns := []string{"key", "diagnosis"}

	// Preparar los datos en un formato compatible con pgx.CopyFrom
	copyData := make([][]interface{}, len(diagnoses))
	for i, diag := range diagnoses {
		copyData[i] = []interface{}{diag.Key, diag.Diagnosis}
	}

	// Ejecutar inserción en lotes con CopyFrom
	_, err = tx.CopyFrom(ctx, pgx.Identifier{"diagnoses"}, columns, pgx.CopyFromRows(copyData))
	if err != nil {
		return fmt.Errorf("error al insertar diagnósticos: %w", err)
	}

	// Confirmar la transacción
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("error al confirmar transacción: %w", err)
	}

	fmt.Println("[OK] Diagnósticos insertados correctamente.")
	return nil
}
