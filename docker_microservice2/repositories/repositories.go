package repositories

import (
	"github.com/chetanmeniyabacncy/docker_microservice2/models"

	"github.com/jmoiron/sqlx"
)

// CompanyRepository ..
type CompanyRepository interface {
	FindByID(ID int64) (*models.Company, error)
	Save(company *models.Company) error
}

// UserRepo implements models.UserRepository
type UserRepo struct {
	db *sqlx.DB
}

// NewUserRepo ..
func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

// FindByID ..
func (r *UserRepo) FindByID(ID int64) (*models.Company, error) {
	var company models.Company
	err := r.db.Get(&company, "SELECT id,name,status FROM companies where id = ?", ID)
	if err != nil {
		return &company, err
	}
	return &company, nil
}

// Save ..
func (r *UserRepo) Save(company *models.Company) error {
	return nil
}
