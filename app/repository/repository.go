package repository

import "crud-alumni/app/models"

type AlumniRepository interface {
	Create(alumni *models.Alumni) error
	FindAll() ([]models.Alumni, error)
	FindByID(id string) (*models.Alumni, error)
	Update(id string, alumni *models.Alumni) error
	Delete(id string) error
}
