package repository

import (
	"context"
	"example-event-sourcing/internal/aggregates"
	"example-event-sourcing/internal/domains/users/models"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Repository struct {
	db *bun.DB
}

func New(db *bun.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Find(ctx context.Context, id string) (aggregates.Aggregate, error) {
	user := &models.User{}
	err := r.db.NewSelect().Model(user).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *Repository) Save(ctx context.Context, aggregate aggregates.Aggregate, eventID string) error {
	user := aggregate.(*models.User)
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = r.Find(ctx, user.GetID())
	if err != nil {
		_, err := tx.NewInsert().Model(user).Exec(ctx)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	_, err = tx.NewUpdate().Model(user).Where("id = ?", user.ID).Exec(ctx)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.NewInsert().Model(&aggregates.Locker{
		EventID:     uuid.MustParse(eventID),
		AggregateID: uuid.MustParse(aggregate.GetID()),
		LockDomain:  "users",
	}).Exec(ctx)

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *Repository) Delete(ctx context.Context, id, eventID string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.NewDelete().Model(&models.User{}).Where("id = ?", id).Exec(ctx)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.NewInsert().Model(&aggregates.Locker{
		EventID:     uuid.MustParse(eventID),
		AggregateID: uuid.MustParse(id),
		LockDomain:  "users",
	}).Exec(ctx)

	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (r *Repository) All(ctx context.Context) ([]aggregates.Aggregate, error) {
	data := make([]aggregates.Aggregate, 0)
	users := make([]*models.User, 0)
	err := r.db.NewSelect().Model(&users).Scan(ctx)
	if err != nil {
		return nil, err
	}
	for _, user := range users {
		data = append(data, user)
	}

	return data, nil
}
