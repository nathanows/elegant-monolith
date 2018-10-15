package company

import (
	"github.com/jmoiron/sqlx"
)

// Repository is the datastore inteface for the company service
type Repository interface {
	save(*companyDTO) (*companyDTO, error)
	delete(int64) error
	find(int64) (*companyDTO, error)
	findAll() ([]*companyDTO, error)
}

type repository struct {
	db *sqlx.DB
}

// NewRepository returns an initialized datastore repository
func NewRepository(db *sqlx.DB) Repository {
	return repository{
		db: db,
	}
}

func (r repository) save(company *companyDTO) (*companyDTO, error) {
	var stmt *sqlx.NamedStmt
	{
		var err error
		if company.ID == 0 {
			stmt, err = r.db.PrepareNamed(sqlInsertCompany)
		} else {
			stmt, err = r.db.PrepareNamed(sqlUpdateCompany)
		}
		if err != nil {
			return nil, err
		}
	}

	var saved companyDTO
	if err := stmt.QueryRowx(company).StructScan(&saved); err != nil {
		return nil, err
	}

	return &saved, nil
}

func (r repository) delete(id int64) error {
	return nil
}

func (r repository) find(id int64) (*companyDTO, error) {
	return &companyDTO{}, nil
}

func (r repository) findAll() ([]*companyDTO, error) {
	return []*companyDTO{}, nil
}

const sqlInsertCompany = `
	INSERT INTO companies (name)
	VALUES (:name)
	RETURNING id, name, created_at, updated_at;`

const sqlUpdateCompany = `
	UPDATE companies SET name = :name
	WHERE id = :id
	RETURNING id, name, created_at, updated_at;`
