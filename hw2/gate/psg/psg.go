package psg

import (
	"context"
	"fmt"
	"hw2/models/dto"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Psg struct {
	Conn *pgxpool.Pool
}

func NewPsg(psgAddr string, login, password string) (*Psg, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s/postgres", login, password, psgAddr)
	poolConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}
	poolConfig.MaxConns = 5 // Установите максимальное количество соединений, если это необходимо
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, err
	}
	return &Psg{Conn: pool}, nil
}

// RecordAdd добавляет новую запись в базу данных.
func (p *Psg) RecordAdd(record dto.Record) (int64, error) {
	query := `
		INSERT INTO records(name, last_name, middle_name, phone, address)
		VALUES($1, $2, $3, $4, $5)
		RETURNING id
	`
	var id int64
	err := p.Conn.QueryRow(context.Background(), query, record.Name, record.LastName, record.MiddleName, record.Phone, record.Address).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// RecordsGet возвращает записи из базы данных на основе предоставленных полей Record.
func (p *Psg) RecordsGet(record dto.Record) ([]dto.Record, error) {
	query := `SELECT id, name, last_name, middle_name, phone, address FROM records WHERE 1=1`
	var args []interface{}

	if record.Name != "" {
		query += " AND name = $" + fmt.Sprint(len(args)+1)
		args = append(args, record.Name)
	}
	if record.LastName != "" {
		query += " AND last_name = $" + fmt.Sprint(len(args)+1)
		args = append(args, record.LastName)
	}
	if record.MiddleName != "" {
		query += " AND middle_name = $" + fmt.Sprint(len(args)+1)
		args = append(args, record.MiddleName)
	}
	if record.Phone != "" {
		query += " AND phone = $" + fmt.Sprint(len(args)+1)
		args = append(args, record.Phone)
	}
	if record.Address != "" {
		query += " AND address = $" + fmt.Sprint(len(args)+1)
		args = append(args, record.Address)
	}

	rows, err := p.Conn.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []dto.Record
	for rows.Next() {
		var r dto.Record
		err := rows.Scan(&r.ID, &r.Name, &r.LastName, &r.MiddleName, &r.Phone, &r.Address)
		if err != nil {
			return nil, err
		}
		records = append(records, r)
	}
	
	return records, nil
}

// RecordUpdate обновляет существующую запись в базе данных по номеру телефона.
func (p *Psg) RecordUpdate(record dto.Record) error {
	query := `
		UPDATE records
		SET name = $1, last_name = $2, middle_name = $3, address = $4
		WHERE phone = $5
	`
	_, err := p.Conn.Exec(context.Background(), query, record.Name, record.LastName, record.MiddleName, record.Address, record.Phone)
	return err
}

// RecordDeleteByPhone удаляет запись из базы данных по номеру телефона.
func (p *Psg) RecordDeleteByPhone(phone string) error {
	query := `DELETE FROM records WHERE phone = $1`
	_, err := p.Conn.Exec(context.Background(), query, phone)
	return err
}

