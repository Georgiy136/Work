package usecase

import (
	"context"
	"fmt"
	"myapp/internal/models"

	"github.com/google/uuid"
)

type ProjectUseCases struct {
	store ProjectStrore
}

func NewProjectUsecases(st ProjectStrore) *ProjectUseCases {
	return &ProjectUseCases{
		store: st,
	}
}

func (us *ProjectUseCases) AddProject(ctx context.Context, p models.Project) (*models.Project, error) {
	p.Id = uuid.New()
	err := us.store.CreateProject(ctx, p)
	if err != nil {
		return nil, fmt.Errorf("ProjectUseCases - AddProject - us.store.CreateProject: %w", err)
	}
	return &p, nil
}

func (us *ProjectUseCases) GetAllProjects(ctx context.Context) ([]models.Project, error) {
	p, err := us.store.GetAllProjects(ctx)
	if err != nil {
		return nil, fmt.Errorf("ProjectUseCases - GetAllProjects - us.store.GetAllProjects: %w", err)
	}
	return p, nil
}

func (us *ProjectUseCases) DeleteProject(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("ProjectUseCases - DeleteProject - uuid.Parse: %w", err)
	}
	err = us.store.DeleteProject(ctx, uid)
	if err != nil {
		return fmt.Errorf("ProjectUseCases - DeleteProject - us.store.DeleteProject: %w", err)
	}
	return nil
}

func (us *ProjectUseCases) UpdateProject(ctx context.Context, id string, p models.Project) (*models.Project, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("ProjectUseCases - UpdateProject - uuid.Parse: %w", err)
	}
	project, err := us.store.UpdateProject(ctx, uid, p)
	if err != nil {
		return nil, fmt.Errorf("ProjectUseCases - UpdateProject - us.store.UpdateProject: %w", err)
	}
	return project, nil
}

func (us *ProjectUseCases) GetOneProject(ctx context.Context, id string) (*models.Project, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("ProjectUseCases - GetOneProject - uuid.Parse: %w", err)
	}
	p, err := us.store.GetOneProject(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("ProjectUseCases - GetOneProject - us.store.GetOneProject: %w", err)
	}

	return p, nil
}

func (us *ProjectUseCases) AddOperatorToProject(ctx context.Context, project_id string, operator_id string) (*models.Project, error) {

	pr_id, err := uuid.Parse(project_id)
	if err != nil {
		return nil, fmt.Errorf("ProjectUseCases - AddOperatorToProject - uuid.Parse: %w", err)
	}
	op_id, err := uuid.Parse(operator_id)
	if err != nil {
		return nil, fmt.Errorf("ProjectUseCases - AddOperatorToProject - uuid.Parse: %w", err)
	}
	project, err := us.store.AddOperatorToProject(ctx, pr_id, op_id)
	if err != nil {
		return nil, fmt.Errorf("ProjectUseCases - AddOperatorToProject - us.store.AddOperatorToProject: %w", err)
	}
	return project, nil
}

func (us *ProjectUseCases) DeleteOperatorFromProject(ctx context.Context, project_id string, operator_id string) (*models.Project, error) {
	pr_id, err := uuid.Parse(project_id)
	if err != nil {
		return nil, fmt.Errorf("ProjectUseCases - DeleteOperatorFromProject - uuid.Parse: %w", err)
	}
	op_id, err := uuid.Parse(operator_id)
	if err != nil {
		return nil, fmt.Errorf("ProjectUseCases - DeleteOperatorFromProject - uuid.Parse: %w", err)
	}
	p, err := us.store.DeleteOperatorFromProject(ctx, pr_id, op_id)
	if err != nil {
		return nil, fmt.Errorf("ProjectUseCases - DeleteOperatorFromProject - us.store.DeleteOperatorFromProject: %w", err)
	}
	return p, nil
}
