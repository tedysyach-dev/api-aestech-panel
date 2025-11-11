package service

import (
	"backend/core/repo"
	"backend/web/entity"
	"backend/web/model"
	"backend/web/repository"
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type BranchsService struct {
	DB               *gorm.DB
	Log              *logrus.Logger
	Validate         *validator.Validate
	BranchRepository *repository.BranchsRepository
}

func NewBrandsService(
	db *gorm.DB,
	logger *logrus.Logger,
	validate *validator.Validate,
	branchRepository *repository.BranchsRepository,
) *BranchsService {
	return &BranchsService{
		DB:               db,
		Log:              logger,
		Validate:         validate,
		BranchRepository: branchRepository,
	}
}

func (s *BranchsService) AddNewManagement(ctx context.Context, request *model.CreateManagementRequest) (*model.CreateManagementResponse, error) {

	management := &entity.Branch{}

	err := repo.ExecuteInTransaction(ctx, s.DB, func(tx *gorm.DB) error {
		// validasi request
		if err := s.Validate.Struct(request); err != nil {
			s.Log.Warnf("Invalid request body : %+v", err)
			return fiber.ErrBadRequest
		}

		management.ID = request.ID
		management.NamaCabang = request.Name
		management.Alamat = request.Address
		management.Email = request.Email

		// 🧱 Field default/hardcoded (mengikuti contoh data)
		management.Kota = "Bandung"
		management.Kontak = "-"
		management.Koordinat = ""
		management.Sipa = ""
		management.IsPrivate = false
		management.Pettycash = 0
		management.BpomMode = true
		management.PPN = 0
		timezone := "Asia/Jakarta"
		management.Datetime = &timezone
		management.Upline = "0"
		isManajemen := true
		management.IsManajemen = &isManajemen
		round := "up"
		management.RoundPPN = &round
		expDate, _ := time.Parse("2006-01-02", "2030-01-01")
		management.ExpireDate = &expDate
		isPaid := false
		dev := false
		isDelete := false
		management.IsPaid = &isPaid
		management.Dev = &dev
		management.IsDelete = &isDelete
		management.AvgGuest = new(int)
		management.AvgTransaction = new(int)
		management.GuestCommentRate = new(int)
		management.GuestTotalByMonth = new(int)
		management.TrxTotalByMonth = new(int)
		rate := 0.0
		management.RateReceptionist = &rate
		management.RateDoctor = &rate
		management.RateBeautician = &rate
		management.IDKlien = uuid.New().String()
		management.AccessStatus = new(bool)

		if err := s.BranchRepository.Create(tx, management); err != nil {
			s.Log.Warnf("Failed to create branch : %+v", err)
			return fiber.ErrInternalServerError
		}
		return nil
	})

	if err != nil {
		s.Log.Warnf("Message : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return &model.CreateManagementResponse{
		ID:      management.ID,
		Name:    management.NamaCabang,
		Address: management.Alamat,
	}, nil
}

func (s *BranchsService) AddNewBranch(ctx context.Context, request *model.CreateBranchRequest) (*model.CreateBranchResponse, error) {

	management := &entity.Branch{}

	err := repo.ExecuteInTransaction(ctx, s.DB, func(tx *gorm.DB) error {
		// validasi request
		if err := s.Validate.Struct(request); err != nil {
			s.Log.Warnf("Invalid request body : %+v", err)
			return fiber.ErrBadRequest
		}

		management.ID = request.ID
		management.NamaCabang = request.Name
		management.Alamat = request.Address
		management.Email = request.Email
		management.Upline = request.Upline

		// 🧱 Field default/hardcoded (mengikuti contoh data)
		management.Kota = "Bandung"
		management.Kontak = "-"
		management.Koordinat = ""
		management.Sipa = ""
		management.IsPrivate = false
		management.Pettycash = 0
		management.BpomMode = true
		management.PPN = 0
		timezone := "Asia/Jakarta"
		management.Datetime = &timezone
		isManajemen := false
		management.IsManajemen = &isManajemen
		round := "up"
		management.RoundPPN = &round
		expDate, _ := time.Parse("2006-01-02", "2030-01-01")
		management.ExpireDate = &expDate
		isPaid := true
		dev := false
		isDelete := false
		management.IsPaid = &isPaid
		management.Dev = &dev
		management.IsDelete = &isDelete
		management.AvgGuest = new(int)
		management.AvgTransaction = new(int)
		management.GuestCommentRate = new(int)
		management.GuestTotalByMonth = new(int)
		management.TrxTotalByMonth = new(int)
		rate := 0.0
		management.RateReceptionist = &rate
		management.RateDoctor = &rate
		management.RateBeautician = &rate
		management.IDKlien = uuid.New().String()
		management.AccessStatus = new(bool)

		if err := s.BranchRepository.Create(tx, management); err != nil {
			s.Log.Warnf("Failed to create branch : %+v", err)
			return fiber.ErrInternalServerError
		}
		return nil
	})

	if err != nil {
		s.Log.Warnf("Message : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return &model.CreateBranchResponse{
		ID:      management.ID,
		Name:    management.NamaCabang,
		Address: management.Alamat,
	}, nil
}
