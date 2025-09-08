package migrator

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// --- PROCESADOR PARA MEDICINAS (REFACTORIZADO) ---

// medicineCopySource es un adaptador que enseña a CopyFrom cómo leer
// el archivo CSV línea por línea, sin cargarlo todo a memoria.
type medicineCopySource struct {
	reader *csv.Reader
	record []string
	err    error
}

// Next avanza al siguiente registro. Retorna false cuando no hay más.
func (mcs *medicineCopySource) Next() bool {
	mcs.record, mcs.err = mcs.reader.Read()
	return mcs.err == nil
}

// Values retorna el registro actual como un slice de `any`.
func (mcs *medicineCopySource) Values() ([]any, error) {
	values := make([]any, len(mcs.record))
	for i, v := range mcs.record {
		values[i] = v
	}
	return values, nil
}

// Err retorna cualquier error que haya ocurrido durante la lectura.
func (mcs *medicineCopySource) Err() error {
	if mcs.err == io.EOF {
		return nil // EOF no es un error para CopyFrom, es la señal de que se terminó.
	}
	return mcs.err
}

// ProcessMedicines ahora procesa el CSV como un stream, sin límite de memoria.
func ProcessMedicines(pathCsv string, db *pgxpool.Pool) error {
	file, err := os.Open(pathCsv)
	if err != nil {
		return fmt.Errorf("error al abrir el archivo .csv: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	ctx := context.Background()
	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error al iniciar transacción: %w", err)
	}
	defer tx.Rollback(ctx)

	// Nombres de las columnas en la tabla `medicines` donde se insertarán los datos.
	columns := []string{
		"substance",
		"presentation",
		"route_administration",
		"dose",
		"quantity",
		"frequency",
	}

	// Creamos nuestra fuente de datos en streaming.
	copySource := &medicineCopySource{reader: reader}

	// Ejecutamos la inserción masiva. CopyFrom es el método más rápido.
	_, err = tx.CopyFrom(ctx, pgx.Identifier{"medicines"}, columns, copySource)
	if err != nil {
		return fmt.Errorf("error durante la inserción masiva (CopyFrom): %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("error al confirmar transacción: %w", err)
	}

	fmt.Println("[OK] Medicinas insertadas correctamente.")
	return nil
}

// --- PROCESADOR PARA DIAGNÓSTICOS (REFACTORIZADO) ---

// diagnosisCopySource es el adaptador para el archivo de diagnósticos.
type diagnosisCopySource struct {
	reader *csv.Reader
	record []string
	err    error
}

func (dcs *diagnosisCopySource) Next() bool {
	dcs.record, dcs.err = dcs.reader.Read()
	return dcs.err == nil
}

func (dcs *diagnosisCopySource) Values() ([]any, error) {
	values := make([]any, len(dcs.record))
	for i, v := range dcs.record {
		values[i] = v
	}
	return values, nil
}

func (dcs *diagnosisCopySource) Err() error {
	if dcs.err == io.EOF {
		return nil
	}
	return dcs.err
}

// ProcessDiagnoses ahora también procesa el CSV como un stream.
func ProcessDiagnoses(pathCsv string, db *pgxpool.Pool) error {
	file, err := os.Open(pathCsv)
	if err != nil {
		return fmt.Errorf("error al abrir el archivo .csv: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	ctx := context.Background()
	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error al iniciar transacción: %w", err)
	}
	defer tx.Rollback(ctx)

	columns := []string{"key", "diagnosis"}

	copySource := &diagnosisCopySource{reader: reader}

	_, err = tx.CopyFrom(ctx, pgx.Identifier{"diagnoses"}, columns, copySource)
	if err != nil {
		return fmt.Errorf("error al insertar diagnósticos (CopyFrom): %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("error al confirmar transacción: %w", err)
	}

	fmt.Println("[OK] Diagnósticos insertados correctamente.")
	return nil
}
