package service

import (
	"context"
	"errors"
	"fmt"
	"sangpalito/backend/internal/apperror"
	"sangpalito/backend/internal/domain/entity"
	"sangpalito/backend/internal/domain/repository"
	"sangpalito/backend/internal/dto"

	"github.com/google/uuid"
)

type VisitService interface {
	CreateVisit(ctx context.Context, req *dto.CreateVisitRequest, createdBy *uuid.UUID) (*dto.VisitResponse, error)
	GetVisitByID(ctx context.Context, id uuid.UUID) (*dto.VisitResponse, error)
	GetVisitsByPatientID(ctx context.Context, patientID uuid.UUID) ([]dto.VisitResponse, error)
	UpdateVisit(ctx context.Context, id uuid.UUID, req *dto.UpdateVisitRequest) error
	SearchICD10(ctx context.Context, query string, limit int) ([]dto.ICD10Response, error)
	SeedICD10(ctx context.Context) error
}

type visitService struct {
	visitRepo repository.VisitRepository
	icd10Repo repository.ICD10Repository
}

func NewVisitService(visitRepo repository.VisitRepository, icd10Repo repository.ICD10Repository) VisitService {
	return &visitService{
		visitRepo: visitRepo,
		icd10Repo: icd10Repo,
	}
}

func (s *visitService) CreateVisit(ctx context.Context, req *dto.CreateVisitRequest, createdBy *uuid.UUID) (*dto.VisitResponse, error) {
	// Pastikan minimal ada 1 diagnosa yang di-set sebagai diagnosa primer
	hasPrimary := false
	for _, d := range req.Diagnoses {
		if d.IsPrimary {
			hasPrimary = true
			break
		}
	}
	if !hasPrimary && len(req.Diagnoses) > 0 {
		req.Diagnoses[0].IsPrimary = true
	}

	doctorID := req.DoctorID
	vitals := req.Vitals

	visit := &entity.Visit{
		ID:            uuid.New(),
		AppointmentID: req.AppointmentID,
		PatientID:     req.PatientID,
		DoctorID:      &doctorID,
		Vitals:        &vitals,
		Subjective:    req.Subjective,
		Objective:     req.Objective,
		Assessment:    req.Assessment,
		Plan:          req.Plan,
		CreatedBy:     createdBy,
	}

	var diagnoses []entity.Diagnosis
	for _, d := range req.Diagnoses {
		var notesPtr *string
		if d.Notes != "" {
			notesVal := d.Notes
			notesPtr = &notesVal
		}

		diagnoses = append(diagnoses, entity.Diagnosis{
			ID:        uuid.New(),
			VisitID:   visit.ID,
			ICD10Code: d.ICD10Code,
			IsPrimary: d.IsPrimary,
			Notes:     notesPtr,
		})
	}

	if err := s.visitRepo.Create(ctx, visit, diagnoses); err != nil {
		return nil, apperror.Internal(fmt.Errorf("failed to record clinical visit: %w", err))
	}

	return s.GetVisitByID(ctx, visit.ID)
}

func (s *visitService) GetVisitByID(ctx context.Context, id uuid.UUID) (*dto.VisitResponse, error) {
	visit, diagnoses, err := s.visitRepo.FindByID(ctx, id)
	if err != nil {
		return nil, apperror.Internal(fmt.Errorf("failed to get visit: %w", err))
	}
	if visit == nil {
		return nil, apperror.NotFound("Visit", nil)
	}

	diagResponses := make([]dto.DiagnosisResponse, 0, len(diagnoses))
	for _, d := range diagnoses {
		desc := d.ICD10Description
		if d.ICD10 != nil && d.ICD10.Description != "" {
			desc = d.ICD10.Description
		}
		notes := ""
		if d.Notes != nil {
			notes = *d.Notes
		}
		diagResponses = append(diagResponses, dto.DiagnosisResponse{
			ID:          d.ID,
			ICD10Code:   d.ICD10Code,
			Description: desc,
			IsPrimary:   d.IsPrimary,
			Notes:       notes,
		})
	}

	var vitals entity.Vitals
	if visit.Vitals != nil {
		vitals = *visit.Vitals
	}
	doctorID := uuid.Nil
	if visit.DoctorID != nil {
		doctorID = *visit.DoctorID
	}

	return &dto.VisitResponse{
		ID:            visit.ID,
		AppointmentID: visit.AppointmentID,
		PatientID:     visit.PatientID,
		DoctorID:      doctorID,
		Vitals:        vitals,
		Subjective:    visit.Subjective,
		Objective:     visit.Objective,
		Assessment:    visit.Assessment,
		Plan:          visit.Plan,
		Diagnoses:     diagResponses,
		CreatedAt:     visit.CreatedAt,
		UpdatedAt:     visit.UpdatedAt,
	}, nil
}

func (s *visitService) GetVisitsByPatientID(ctx context.Context, patientID uuid.UUID) ([]dto.VisitResponse, error) {
	visits, err := s.visitRepo.FindByPatientID(ctx, patientID)
	if err != nil {
		return nil, apperror.Internal(fmt.Errorf("failed to get patient visits: %w", err))
	}

	var results []dto.VisitResponse
	for _, v := range visits {
		var vitals entity.Vitals
		if v.Vitals != nil {
			vitals = *v.Vitals
		}
		docID := uuid.Nil
		if v.DoctorID != nil {
			docID = *v.DoctorID
		}

		results = append(results, dto.VisitResponse{
			ID:            v.ID,
			AppointmentID: v.AppointmentID,
			PatientID:     v.PatientID,
			DoctorID:      docID,
			Vitals:        vitals,
			Subjective:    v.Subjective,
			Objective:     v.Objective,
			Assessment:    v.Assessment,
			Plan:          v.Plan,
			CreatedAt:     v.CreatedAt,
			UpdatedAt:     v.UpdatedAt,
		})
	}
	return results, nil
}

func (s *visitService) UpdateVisit(ctx context.Context, id uuid.UUID, req *dto.UpdateVisitRequest) error {
	visit, _, err := s.visitRepo.FindByID(ctx, id)
	if err != nil {
		return apperror.Internal(fmt.Errorf("failed to find visit: %w", err))
	}
	if visit == nil {
		return apperror.NotFound("Visit", nil)
	}

	if req.Vitals != nil {
		visit.Vitals = req.Vitals
	}
	if req.Subjective != "" {
		visit.Subjective = req.Subjective
	}
	if req.Objective != "" {
		visit.Objective = req.Objective
	}
	if req.Assessment != "" {
		visit.Assessment = req.Assessment
	}
	if req.Plan != "" {
		visit.Plan = req.Plan
	}

	if err := s.visitRepo.Update(ctx, visit); err != nil {
		return apperror.Internal(fmt.Errorf("failed to update visit: %w", err))
	}
	return nil
}

func (s *visitService) SearchICD10(ctx context.Context, query string, limit int) ([]dto.ICD10Response, error) {
	codes, err := s.icd10Repo.Search(ctx, query, limit)
	if err != nil {
		return nil, apperror.Internal(fmt.Errorf("failed to search ICD-10: %w", err))
	}

	var res []dto.ICD10Response
	for _, c := range codes {
		res = append(res, dto.ICD10Response{
			Code:        c.Code,
			Description: c.Description,
		})
	}
	return res, nil
}

func (s *visitService) SeedICD10(ctx context.Context) error {
	if err := s.icd10Repo.SeedCommonCodes(ctx); err != nil {
		return errors.New("failed to seed ICD-10 codes: " + err.Error())
	}
	return nil
}
