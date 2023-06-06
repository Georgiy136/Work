package usecase

import (
	"context"
	"myapp/internal/models"
	"time"

	"github.com/google/uuid"
)

type OperatorStrore interface {
	CreateOperator(ctx context.Context, p models.Operator) error
	GetAllOperators(ctx context.Context) ([]models.Operator, error)
	DeleteOperator(ctx context.Context, id uuid.UUID) error
	UpdateOperator(ctx context.Context, id uuid.UUID, p models.Operator) (*models.Operator, error)
	GetOneOperator(ctx context.Context, id uuid.UUID) (*models.Operator, error)
}

type ProjectStrore interface {
	CreateProject(ctx context.Context, p models.Project) error
	GetAllProjects(ctx context.Context) ([]models.Project, error)
	DeleteProject(ctx context.Context, id uuid.UUID) error
	UpdateProject(ctx context.Context, id uuid.UUID, p models.Project) (*models.Project, error)
	GetOneProject(ctx context.Context, id uuid.UUID) (*models.Project, error)

	AddOperatorToProject(ctx context.Context, project_id uuid.UUID, operator_id uuid.UUID) (*models.Project, error)
	DeleteOperatorFromProject(ctx context.Context, project_id uuid.UUID, operator_id uuid.UUID) (*models.Project, error)
}

type AuthorizationStore interface {
	CreateClient(ctx context.Context, c models.Client) error
	GetOneClient(ctx context.Context, login, password string) (*models.Client, error)
	GetOneClientById(ctx context.Context, id string) (*models.Client, error)
	
	GetRoleRights(ctx context.Context, role string) ([]string, error)
}

type RoleStore interface {
	GetRoleRights(ctx context.Context, role string) ([]string, error)
	AddRoleRights (ctx context.Context, role string, rights []string, period time.Duration) error 
}
