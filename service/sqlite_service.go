package service

import (
    "context"
    "database/sql"
    "encoding/json"
    "fmt"
    "log"

    _ "github.com/mattn/go-sqlite3"
    "github.com/rainbowmga/timetravel/entity"
)

type SQLiteRecordService struct {
    db *sql.DB
}

// NewSQLiteRecordService creates a new instance of SQLiteRecordService
func NewSQLiteRecordService() *SQLiteRecordService {
    db, err := sql.Open("sqlite3", "./records.db")
    if err != nil {
        log.Fatal(err)
    }

    // The `ver` field is used by v2 of the service
    createTable := `CREATE TABLE IF NOT EXISTS records (
        id INTEGER,
        data TEXT,
        ver INTEGER,
        PRIMARY KEY (id, ver)
    );`
    _, err = db.Exec(createTable)
    if err != nil {
        log.Fatal(err)
    }

    return &SQLiteRecordService{db: db}
}

// GetLatestVersion gets the latest version for the given id
func (s *SQLiteRecordService) GetLatestVersion(ctx context.Context, id int) (int, error) {
    var err error
    var verRead int
    query := "SELECT ver FROM records WHERE id = ? ORDER BY ver desc LIMIT 1"
    err = s.db.QueryRowContext(ctx, query, id).Scan(&verRead)
    if err != nil {
        if err == sql.ErrNoRows {
            return 1, ErrRecordDoesNotExist
        }
        return 1, err
    }
    return verRead, nil
}

// DeleteRecordForVersion deletes the specified version for the given id
func (s *SQLiteRecordService) DeleteRecordForVersion(ctx context.Context, id int, ver int) (error) {
    stmt, err := s.db.Prepare("DELETE FROM records WHERE id = ? AND ver = ?")
    if err != nil {
        return fmt.Errorf("failed to prepare statement: %v", err)
    }
    defer stmt.Close()

    // Execute the DELETE statement with the id and ver as parameters
    res, err := stmt.Exec(id, ver)
    if err != nil {
        return fmt.Errorf("failed to execute delete: %v", err)
    }

    // Check how many rows were affected by the DELETE
    rowsAffected, err := res.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to retrieve rows affected: %v", err)
    }

    if rowsAffected == 0 {
        return fmt.Errorf("no record found with id %d and version %d", id, ver)
    }
    return nil
}
// GetRecord retrieves a record by id
func (s *SQLiteRecordService) GetRecord(ctx context.Context, id int, ver int) (entity.Record, error) {
    var jsonData string
    var err error
    var verRead int
    if ver < 0 {
        query := "SELECT data, ver FROM records WHERE id = ? ORDER BY ver desc LIMIT 1"
        err = s.db.QueryRowContext(ctx, query, id).Scan(&jsonData, &verRead)
    } else {
        query := "SELECT data, ver FROM records WHERE id = ? AND ver = ?"
        err = s.db.QueryRowContext(ctx, query, id, ver).Scan(&jsonData, &verRead)
    }
    if err != nil {
        if err == sql.ErrNoRows {
            return entity.Record{}, ErrRecordDoesNotExist
        }
        return entity.Record{}, err
    }

    var data map[string]string
    err = json.Unmarshal([]byte(jsonData), &data)
    if err != nil {
        return entity.Record{}, err
    }

    return entity.Record{ID: id, Data: data, Ver: verRead}, nil
}

// CreateRecord inserts a new record into the database
func (s *SQLiteRecordService) CreateRecord(ctx context.Context, record entity.Record) error {
    // Convert map to JSON
    jsonData, err := json.Marshal(record.Data)
    if err != nil {
        return err
    }

    // Insert the record
    query := `INSERT INTO records (id, data, ver) VALUES (?, ?, 1)`
    _, err = s.db.ExecContext(ctx, query, record.ID, string(jsonData))
    if err != nil {
        return err
    }

    return nil
}

// UpdateRecord updates the existing record with new data or deletes keys if values are null
func (s *SQLiteRecordService) UpdateRecord(ctx context.Context, id int, updates map[string]*string) (entity.Record, error) {
    // Retrieve the existing record first
    record, err := s.GetRecord(ctx, id, -1)
    if err != nil {
        return entity.Record{}, err
    }

    // Update the record with new values or delete keys if the value is null
    for key, value := range updates {
        if value == nil {
            delete(record.Data, key)
        } else {
            record.Data[key] = *value
        }
    }

    // Convert updated map to JSON
    jsonData, err := json.Marshal(record.Data)
    if err != nil {
        return entity.Record{}, err
    }

    // Update the record in the database
    query := `UPDATE records SET data = ?, ver = ? WHERE id = ?`
    ver := record.Ver
    if src, ok := ctx.Value("src").(string); ok {
        if src == "v2" {
            // increment the version
            ver++
            query = `INSERT INTO records (id, data, ver) VALUES (?, ?, ?)`
            _, err = s.db.ExecContext(ctx, query, record.ID, string(jsonData), ver)
        } else {
            _, err = s.db.ExecContext(ctx, query, string(jsonData), ver, id)
        } 
    } else {
        _, err = s.db.ExecContext(ctx, query, string(jsonData), ver, id)
    }
    if err != nil {
        return entity.Record{}, err
    }

    record.Ver = ver
    return record, nil
}
