package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"myapp/internal/models"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/uptrace/bun"
)

func NewOperator(Bun *bun.DB) *Operator {
	return &Operator{
		Bun: Bun,
	}
}

type Operator struct {
	Bun *bun.DB
}

func (db *Operator) CreateOperator(ctx context.Context, o models.Operator) error {
	_, err := db.Bun.NewInsert().Model(&o).Exec(ctx)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("Operator - CreateOperator - db.Bun.NewInsert: %w", err)
	}
	return nil
}

func (db *Operator) GetAllOperators(ctx context.Context) ([]models.Operator, error) {

	operators := []models.Operator{}
	operator := &models.Operator{}

	err := db.Bun.NewSelect().Model(operator).Scan(ctx, &operators)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("Operator - GetAllOperators - db.Bun.NewSelect: %w", err)
	}

	return operators, nil
}

func (db *Operator) DeleteOperator(ctx context.Context, id uuid.UUID) error {
	operator := &models.Operator{}
	err := db.Bun.NewDelete().
		Model(operator).
		Where(`uuid = ?`, id.String()).
		Returning("uuid").
		Scan(ctx, &id)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("Operator - DeleteOperator - db.Bun.NewDelete: %s", fmt.Sprintf("оператора с id = %s не существует", id.String()))
		}
		log.Println(err)
		return fmt.Errorf("Operator - DeleteOperator - db.Bun.NewDelete: %w", err)
	}
	return nil
}

func (db *Operator) UpdateOperator(ctx context.Context, id uuid.UUID, p models.Operator) (*models.Operator, error) {

	err := db.Bun.NewUpdate().
		Model(&p).
		Column("first_name", "last_name", "patronymic", "city", "phone", "email").
		Where(`uuid = ?`, id.String()).
		Returning("uuid, password").
		Scan(ctx, &p.Id, &p.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Operator - UpdateOperator - db.Bun.NewUpdate: %s", fmt.Sprintf("оператора с id = %s не существует", id.String()))
		}
		log.Println(err)
		return nil, fmt.Errorf("Operator - UpdateOperator - db.Bun.NewUpdate: %w", err)
	}

	return &p, nil
}

func (db *Operator) GetOneOperator(ctx context.Context, id uuid.UUID) (*models.Operator, error) {
	operators := models.Operator{}

	err := db.Bun.NewSelect().Model(&operators).Where("uuid = ?", id.String()).Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Operator - GetOneOperator - db.Bun.NewSelect: %s", fmt.Sprintf("оператора с id = %s не существует", id.String()))
		}
		log.Println(err)
		return nil, fmt.Errorf("Operator - GetOneOperator - db.Bun.NewSelect: %w", err)
	}

	return &operators, nil
}

func (db *Operator) GetOperatorsById(ctx context.Context, id []string) ([]*models.Operator, error) {
	operator := models.Operator{}
	operators := []*models.Operator{}

	err := db.Bun.NewSelect().Model(&operator).Where("uuid IN (?)", bun.In(id)).Scan(ctx, &operators)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Println(err)
		return nil, fmt.Errorf("Operator - GetOperatorsById - db.Bun.NewSelect: %w", err)
	}
	return operators, nil
}
