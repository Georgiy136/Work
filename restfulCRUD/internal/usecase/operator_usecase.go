package usecase

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"myapp/internal/models"
	"net/mail"
	"strings"
	"time"

	"github.com/google/uuid"
)

type OperatorUseCases struct {
	store OperatorStrore
}

func NewOperatorUsecases(st OperatorStrore) *OperatorUseCases {
	return &OperatorUseCases{
		store: st,
	}
}

func (us *OperatorUseCases) AddOperator(ctx context.Context, p models.Operator) (*models.Operator, error) {
	if err := validPhone(p.Phone); err != nil {
		return nil, fmt.Errorf("OperatorUseCases - AddOperator - validPhone: %w", err)
	}
	if err := validEmail(p.Email); err != nil {
		return nil, fmt.Errorf("OperatorUseCases - AddOperator - validEmail: %w", err)
	}
	p.Id = uuid.New()
	p.Password = generatePassword()
	err := us.store.CreateOperator(ctx, p)
	if err != nil {
		return nil, fmt.Errorf("OperatorUseCases - AddOperator - us.store.CreateOperator: %w", err)
	}
	return &p, nil
}

func (us *OperatorUseCases) GetAllOperators(ctx context.Context) ([]models.Operator, error) {
	p, err := us.store.GetAllOperators(ctx)
	if err != nil {
		return nil, fmt.Errorf("OperatorUseCases - GetAllOperators - us.store.GetAllOperators: %w", err)
	}
	return p, nil
}

func (us *OperatorUseCases) DeleteOperator(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("OperatorUseCases - DeleteOperator - uuid.Parse: %w", err)
	}
	err = us.store.DeleteOperator(ctx, uid)
	if err != nil {
		return fmt.Errorf("OperatorUseCases - DeleteOperator - us.store.DeleteOperator: %w", err)
	}
	return nil
}

func (us *OperatorUseCases) UpdateOperator(ctx context.Context, id string, p models.Operator) (*models.Operator, error) {
	if err := validPhone(p.Phone); err != nil {
		return nil, fmt.Errorf("OperatorUseCases - UpdateOperator - validPhone: %w", err)
	}
	if err := validEmail(p.Email); err != nil {
		return nil, fmt.Errorf("OperatorUseCases - UpdateOperator - validEmail: %w", err)
	}
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("OperatorUseCases - UpdateOperator - uuid.Parse: %w", err)
	}
	operator, err := us.store.UpdateOperator(ctx, uid, p)
	if err != nil {
		return nil, fmt.Errorf("OperatorUseCases - UpdateOperator - us.store.UpdateOperator: %w", err)
	}
	return operator, nil
}

func (us *OperatorUseCases) GetOneOperator(ctx context.Context, id string) (*models.Operator, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("OperatorUseCases - GetOneOperator - uuid.Parse: %w", err)
	}
	p, err := us.store.GetOneOperator(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("OperatorUseCases - GetOneOperator - us.store.GetOneOperator: %w", err)
	}
	return p, nil
}

func generatePassword() string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")
	length := 10
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}

func validEmail(email string) error {
	if _, err := mail.ParseAddress(email); err != nil {
		return fmt.Errorf("%s - данный email не валиден", email)
	}
	return nil
}

func validPhone(phone string) error {

	pattern := "0123456789"

	if len(phone) != 11 || phone[0] != byte('8') {
		return fmt.Errorf("телефон %s не соответсвует требуемому шаблону", phone)
	}
	for _, num := range phone {
		if ok := strings.Contains(pattern, string(num)); !ok {
			log.Println(string(num))
			return fmt.Errorf("телефон %s не соответсвует требуемому шаблону", phone)
		}
	}
	return nil
}
