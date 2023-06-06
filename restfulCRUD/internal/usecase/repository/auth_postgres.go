package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"myapp/internal/models"

	_ "github.com/lib/pq"
	"github.com/uptrace/bun"
)

func NewAuthPostgres(Bun *bun.DB) *AuthPostgres {
	return &AuthPostgres{
		Bun: Bun,
	}
}

type AuthPostgres struct {
	Bun *bun.DB
}

func (db *AuthPostgres) CreateClient(ctx context.Context, c models.Client) error {
	_, err := db.Bun.NewInsert().Model(&c).Exec(ctx)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("AuthPostgres - CreateClient - db.Bun.NewInsert: %w", err)
	}
	return nil
}

func (db *AuthPostgres) GetOneClient(ctx context.Context, login, password string) (*models.Client, error) {
	client := models.Client{}

	err := db.Bun.NewSelect().Model(&client).Where("login = ? AND password = ?", login, password).Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("AuthPostgres - GetOneClient - db.Bun.NewSelect: %s", "неверный логин или пароль")
		}
		log.Println(err)
		return nil, fmt.Errorf("AuthPostgres - GetOneClient - db.Bun.NewSelect: %w", err)
	}
	return &client, nil
}

func (db *AuthPostgres) GetOneClientById(ctx context.Context, id string) (*models.Client, error) {
	clients := models.Client{}

	err := db.Bun.NewSelect().Model(&clients).Where("uuid = ?", id).Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("AuthPostgres - GetOneClientById - db.Bun.NewSelect: %s", "неверный логин или пароль")
		}
		log.Println(err)
		return nil, fmt.Errorf("AuthPostgres - GetOneClientById - db.Bun.NewSelect: %w", err)
	}

	return &clients, nil
}

func (db *AuthPostgres) GetRoleRights(ctx context.Context, role string) ([]string, error) {

	var rights []string

	err := db.Bun.NewSelect().Table("roles").
		Column("rights.rights").
		Join("join roles_rights").
		JoinOn("roles.id = roles_rights.role_id").
		Join("join rights").
		JoinOn("rights.id = roles_rights.right_id").
		Where("roles.role = ?", role).
		Scan(ctx, &rights)

	log.Println(rights)

	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("AuthPostgres - GetRoleRights - db.Bun.NewSelect: %w", err)
	}

	return rights, nil
}
