package service

import (
    "context"
    "database/sql"
    "encoding/json"
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

    // Create table if not exists
    createTable := `CREATE TABLE IF NOT EXISTS records (
        id INTEGER PRIMARY KEY,
        data TEXT
    );`
    _, err = db.Exec(createTable)
    if err != nil {
        log.Fatal(err)
    }

    return &SQLiteRecordService{db: db}
}

// GetRecord retrieves a record by id
func (s *SQLiteRecordService) GetRecord(ctx context.Context, id int) (entity.Record, error) {
    var jsonData string
    query := "SELECT data FROM records WHERE id = ?"
    err := s.db.QueryRowContext(ctx, query, id).Scan(&jsonData)
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

    return entity.Record{ID: id, Data: data}, nil
}

// CreateRecord inserts a new record into the database
func (s *SQLiteRecordService) CreateRecord(ctx context.Context, record entity.Record) error {
    // Convert map to JSON
    jsonData, err := json.Marshal(record.Data)
    if err != nil {
        return err
    }

    // Insert the record
    query := `INSERT INTO records (id, data) VALUES (?, ?)`
    _, err = s.db.ExecContext(ctx, query, record.ID, string(jsonData))
    if err != nil {
        return err
    }

    return nil
}

// UpdateRecord updates the existing record with new data or deletes keys if values are null
func (s *SQLiteRecordService) UpdateRecord(ctx context.Context, id int, updates map[string]*string) (entity.Record, error) {
    // Retrieve the existing record first
    record, err := s.GetRecord(ctx, id)
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
    query := `UPDATE records SET data = ? WHERE id = ?`
    _, err = s.db.ExecContext(ctx, query, string(jsonData), id)
    if err != nil {
        return entity.Record{}, err
    }

    return record, nil
}
