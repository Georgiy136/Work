package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"myapp/internal/models"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/uptrace/bun"
)

func NewProject(Bun *bun.DB) *Project {
	return &Project{
		Bun: Bun,
	}
}

type Project struct {
	Bun *bun.DB
}

// Проекты
func (db *Project) CreateProject(ctx context.Context, p models.Project) error {
	_, err := db.Bun.NewInsert().Model(&p).Exec(ctx)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("Project - CreateProject - db.Bun.NewInsert: %w", err)
	}
	return nil
}

func (db *Project) GetAllProjects(ctx context.Context) ([]models.Project, error) {

	projects := []models.Project{}
	project := models.Project{}
	var operatorsId []string

	rows, err := db.Bun.NewSelect().
		Table("projects").
		Column("uuid", "project_name", "project_types.project_type", "operators").
		Join("join project_types").
		JoinOn("projects.project_type = project_types.id").
		Rows(ctx)

	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("Project - GetAllProjects - db.Bun.NewSelect: %w", err)
	}
	for rows.Next() {
		err = rows.Scan(&project.Id, &project.ProjectName, &project.ProjectType, (pq.Array)(&operatorsId))
		if err != nil {
			return nil, fmt.Errorf("Project - GetAllProjects - rows.Scan: %w", err)
		}

		newOperator := NewOperator(db.Bun)
		project.Operators = nil

		if len(operatorsId) != 0 {
			if operators, err := newOperator.GetOperatorsById(ctx, operatorsId); err == nil {
				project.Operators = append(project.Operators, operators...)
			} else {
				log.Println(err)
			}
		}

		projects = append(projects, project)
	}
	defer rows.Close()

	return projects, nil
}

func (db *Project) GetOneProject(ctx context.Context, id uuid.UUID) (*models.Project, error) {

	project := models.Project{}
	var operatorsId []string

	err := db.Bun.NewSelect().
		Table("projects").
		Column("uuid", "project_name", "project_types.project_type", "operators").
		Join("join project_types").
		JoinOn("projects.project_type = project_types.id").
		Where("uuid = ?", id.String()).
		Scan(ctx, &project.Id, &project.ProjectName, &project.ProjectType, (pq.Array)(&operatorsId))

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Project - GetOneProject - db.Bun.NewSelect: %s", fmt.Sprintf("проекта с id = %s не существует", id.String()))
		}
		log.Println(err)
		return nil, fmt.Errorf("Project - GetOneProject - db.Bun.NewSelect: %w", err)
	}

	newOperator := NewOperator(db.Bun)
	project.Operators = nil

	if len(operatorsId) != 0 {
		if operators, err := newOperator.GetOperatorsById(ctx, operatorsId); err == nil {
			project.Operators = append(project.Operators, operators...)
		} else {
			log.Println(err)
		}
	}

	return &project, nil
}

func (db *Project) UpdateProject(ctx context.Context, id uuid.UUID, project models.Project) (*models.Project, error) {

	var operatorsId []string

	err := db.Bun.NewUpdate().
		Model(&project).
		Column("project_name", "project_type").
		Where(`uuid = ?`, id.String()).
		Returning("uuid, operators").
		Scan(ctx, &project.Id, pq.Array(&operatorsId))

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Project - UpdateProject - db.Bun.NewUpdate: %s", fmt.Sprintf("проекта с id = %s не существует", id.String()))
		}
		log.Println(err)
		return nil, fmt.Errorf("Project - UpdateProject - db.Bun.NewUpdate: %w", err)
	}

	newOperator := NewOperator(db.Bun)

	project.Operators = nil

	if len(operatorsId) != 0 {
		if operators, err := newOperator.GetOperatorsById(ctx, operatorsId); err == nil {
			project.Operators = append(project.Operators, operators...)
		} else {
			log.Println(err)
		}
	}

	return &project, nil
}

func (db *Project) DeleteProject(ctx context.Context, id uuid.UUID) error {

	project := &models.Project{}
	err := db.Bun.NewDelete().
		Model(project).
		Where(`uuid = ?`, id.String()).
		Returning("uuid").
		Scan(ctx, &id)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("Project - DeleteProject - db.Bun.NewDelete: %s", fmt.Sprintf("проекта с id = %s не существует", id.String()))
		}
		log.Println(err)
		return fmt.Errorf("Project - DeleteProject - db.Bun.NewDelete: %w", err)
	}
	return nil
}

func (db *Project) AddOperatorToProject(ctx context.Context, projectId uuid.UUID, operatorId uuid.UUID) (*models.Project, error) {

	var operatorsId []string

	project, err := db.GetOneProject(ctx, projectId)
	if err != nil {
		return nil, fmt.Errorf("Project - AddOperatorToProject - db.Bun.GetOneProject: %w", err)
	}

	for _, operator := range project.Operators {
		if operator.Id.String() == operatorId.String() {
			return nil, fmt.Errorf("Project - AddOperatorToProject: %s", fmt.Sprintf("Данный оператор %s уже участвует в проекте %s", operatorId.String(), projectId))
		}
		operatorsId = append(operatorsId, operator.Id.String())
	}

	newOperator := NewOperator(db.Bun)

	_, err = newOperator.GetOneOperator(ctx, operatorId)
	if err != nil {
		return nil, fmt.Errorf("Project - AddOperatorToProject - db.Bun.GetOneProject: %w", err)
	}

	operatorsId = append(operatorsId, operatorId.String())

	_, err = db.Bun.NewUpdate().
		Model(project).
		Set("operators = ?", pq.Array(operatorsId)).
		Where(`uuid = ?`, projectId.String()).
		Exec(ctx)

	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("Project - AddOperatorToProject - db.Bun.NewUpdate: %w", err)
	}

	log.Println(operatorsId)
	project.Operators = nil

	if len(operatorsId) != 0 {
		if operators, err := newOperator.GetOperatorsById(ctx, operatorsId); err == nil {
			project.Operators = append(project.Operators, operators...)
		}
	}
	return project, nil
}

func (db *Project) DeleteOperatorFromProject(ctx context.Context, projectId uuid.UUID, operatorId uuid.UUID) (*models.Project, error) {

	var operatorsId []string

	project, err := db.GetOneProject(ctx, projectId)
	if err != nil {
		return nil, fmt.Errorf("Project - DeleteOperatorFromProject - db.Bun.GetOneProject: %w", err)
	}

	exist := false
	for _, operator := range project.Operators {
		if operator.Id.String() == operatorId.String() {
			exist = true
			continue
		}
		operatorsId = append(operatorsId, operator.Id.String())
	}

	if !exist {
		return nil, fmt.Errorf("Project - DeleteOperatorFromProject - db.Bun.NewUpdate: %s", fmt.Sprintf("Данный оператор %s не участвует в проекте %s", operatorId.String(), projectId.String()))
	}

	_, err = db.Bun.NewUpdate().
		Model(project).
		Set("operators = ?", pq.Array(operatorsId)).
		Where(`uuid = ?`, projectId.String()).
		Exec(ctx)

	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("Project - DeleteOperatorFromProject - db.Bun.NewUpdate: %w", err)
	}

	log.Println(operatorsId)
	project.Operators = nil

	newOperator := NewOperator(db.Bun)

	if len(operatorsId) != 0 {
		if operator, err := newOperator.GetOperatorsById(ctx, operatorsId); err == nil {
			project.Operators = append(project.Operators, operator...)
		} else {
			log.Println(err)
		}
	}
	return project, nil
}
