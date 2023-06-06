package usecase

import (
	"context"
	"crypto/sha256"
	"fmt"
	"myapp/internal/models"
	"time"

	"myapp/config"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const tokenLifetime = time.Hour * 24 * 30 //время жизни токена = 1 месяц

type AuthUseCases struct {
	authStore AuthorizationStore
	roleStore RoleStore
	key       string
}

func NewAuthUsecases(st AuthorizationStore, rt RoleStore, cfg config.Auth) *AuthUseCases {
	return &AuthUseCases{
		authStore: st,
		roleStore: rt,
		key:       cfg.Key,
	}
}

// Регистрация пользователя
func (s *AuthUseCases) RegistClient(ctx context.Context, client models.Client) error {

	if client.ClientRole != "Admin" {
		client.ClientRole = "User"
	}

	client.Id = uuid.New()

	//Генерируем хешированный пароль
	hash := generateSha256Password(client.Password)

	client.Password = string(hash)
	if err := validEmail(client.Login); err != nil {
		return fmt.Errorf("AuthUseCases - CreateClient - validEmail: %w", err)
	}
	if err := s.authStore.CreateClient(ctx, client); err != nil {
		return fmt.Errorf("AuthUseCases - CreateClient - s.authStore.CreateClient: %w", err)
	}
	return nil
}

func (s *AuthUseCases) GetOneClientById(ctx context.Context, id string) (*models.Client, error) {

	client, err := s.authStore.GetOneClientById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("AuthUseCases - GetOneClientById - s.authStore.GetOneClientById: %w", err)
	}
	return client, nil
}

func (s *AuthUseCases) GenerateToken(ctx context.Context, login, password string) (string, error) {

	//Генерируем хешированный пароль
	hash := generateSha256Password(password)

	client, err := s.authStore.GetOneClient(ctx, login, string(hash))
	if err != nil {
		return "", fmt.Errorf("AuthUseCases - GenerateToken - s.authStore.GetOneClient: %w", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenLifetime)),
		Subject:   client.Id.String(),
	})

	tokenString, err := token.SignedString([]byte(s.key))
	if err != nil {
		return "", fmt.Errorf("AuthUseCases - GenerateToken - token.SignedString: %w", err)
	}
	return tokenString, nil
}

func (s *AuthUseCases) ParseToken(tokenString string) (string, int64, error) {

	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.key), nil
	}, jwt.WithLeeway(5*time.Second))

	if err != nil {
		return "", 0, fmt.Errorf("AuthUseCases - ParseToken -  jwt.Parse: %v", err)
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return claims.Subject, claims.ExpiresAt.Time.Unix(), nil
	} else {
		return "", 0, fmt.Errorf("AuthUseCases - ParseToken -  jwt.Parse: %s", "token claims are not of type *tokenClaims")
	}

}

func (s *AuthUseCases) GetRoleRights(ctx context.Context, role string) ([]string, error) {

	var (
		rights []string
		err    error
	)

	//Делаем выборку из Redis
	rights, err = s.roleStore.GetRoleRights(ctx, role)
	if err != nil {
		return nil, fmt.Errorf("AuthUseCases - GetRoleRights -  s.roleStore.GetRoleRights: %v", err)
	}

	//log.Println("Redis rights", rights)

	//Проверяем результат на nil
	if len(rights) != 0 {
		return rights, nil
	}

	//Делаем выборку из PostgreSQL
	rights, err = s.authStore.GetRoleRights(ctx, role)
	if err != nil {
		return nil, fmt.Errorf("AuthUseCases - GetRoleRights - s.authStore.GetRoleRights: %v", err)
	}

	//log.Println("PostgreSQL rights", rights)

	//Добавляем права сроком на 5 минут в Redis
	err = s.roleStore.AddRoleRights(ctx, role, rights, 5*60*time.Second)
	if err != nil {
		return nil, fmt.Errorf("AuthUseCases - GetRoleRights -  s.roleStore.AddRoleRights: %v", err)
	}

	return rights, nil
}

func generateSha256Password(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return fmt.Sprint(hash)
}
