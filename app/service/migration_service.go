package service

import (
	"crud-alumni/app/models"
	"crud-alumni/app/repository"
	"crud-alumni/app/repository/mongo"
	"fmt"
)

func MigratePekerjaanToMongo() error {
	// Get PostgreSQL repository
	pekerjaanPostgres, err := repository.GetAllPekerjaanRepo()
	if err != nil {
		return fmt.Errorf("failed to get data from PostgreSQL: %v", err)
	}

	// Get MongoDB repository
	mongoRepo := mongo.NewPekerjaanRepository()

	// Migrate each record
	for _, p := range pekerjaanPostgres {
		// Copy data ke MongoDB
		err := mongoRepo.Create(&models.Pekerjaan{
			AlumniID:            p.AlumniID,
			NamaPerusahaan:      p.NamaPerusahaan,
			PosisiJabatan:       p.PosisiJabatan,
			BidangIndustri:      p.BidangIndustri,
			LokasiKerja:         p.LokasiKerja,
			GajiRange:           p.GajiRange,
			TanggalMulaiKerja:   p.TanggalMulaiKerja,
			TanggalSelesaiKerja: p.TanggalSelesaiKerja,
			StatusPekerjaan:     p.StatusPekerjaan,
			Deskripsi:           p.Deskripsi,
			CreatedAt:           p.CreatedAt,
			UpdatedAt:           p.UpdatedAt,
			IsDelete:            p.IsDelete,
			DeletedAt:           p.DeletedAt,
			DeletedBy:           p.DeletedBy,
		})
		if err != nil {
			return fmt.Errorf("failed to migrate pekerjaan ID %d: %v", p.ID, err)
		}
	}

	return nil
}
