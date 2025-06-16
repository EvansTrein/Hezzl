package repository

import (
	"context"
	"fmt"
	"hezzl/internal/model"
	"hezzl/pkg/db/postgres"
	"log/slog"
	"strings"
	"time"
)

const (
	tableName     = "goods"
	rollbackTimer = time.Second * 10
)

type goodsRepo struct {
	log *slog.Logger
	*postgres.PostgresDB
}

type GoodsRepoDeps struct {
	*slog.Logger
	*postgres.PostgresDB
}

func NewGoodsRepo(deps *GoodsRepoDeps) *goodsRepo {
	return &goodsRepo{
		log:        deps.Logger,
		PostgresDB: deps.PostgresDB,
	}
}

func (r *goodsRepo) Create(ctx context.Context, data model.ProductCreateRequest) (*model.Product, error) {
	op := "goods repository: creating"
	log := r.log.With(slog.String("operation", op))
	log.Debug("Call func Create", "data", data)

	var product model.Product

	query := fmt.Sprintf(`
		INSERT INTO %s (
			project_id, 
			name, 
			priority
		)
		SELECT 
			$1, 
			$2, 
			COALESCE(MAX(priority), 0) + 1
		FROM goods
		WHERE project_id = $1
		RETURNING id, project_id, name, description, priority, removed, created_at
	`, tableName)

	if err := r.DB.QueryRow(ctx, query, data.ProjectID, data.Name).
		Scan(
			&product.ID,
			&product.ProjectID,
			&product.Name,
			&product.Description,
			&product.Priority,
			&product.Removed,
			&product.CreatedAt,
		); err != nil {
		log.Error("failed to create record", "error", err)
		return nil, err
	}

	log.Info("successfully created")
	return &product, nil
}

func (r *goodsRepo) Update(ctx context.Context, data model.ProductUpdateRequest) (*model.Product, error) {
	op := "goods repository: updating"
	log := r.log.With(slog.String("operation", op))
	log.Debug("Call func Update", "data", data)

	ctxRollback, cancel := context.WithTimeout(context.Background(), rollbackTimer)
	defer cancel()

	var product model.Product

	query := fmt.Sprintf(`
		WITH locked_row AS (
			SELECT * FROM %s
			WHERE id = $1 AND project_id = $2
			FOR UPDATE
		),
		updated_row AS (
			UPDATE %s
			SET
				name = $3,
				description = CASE
					WHEN $4 <> '' THEN $4
					ELSE description
				END
			WHERE id = $1 AND project_id = $2
			RETURNING *
		)
		SELECT * FROM updated_row;
	`, tableName, tableName)

	tx, err := r.DB.Begin(ctx)
	if err != nil {
		log.Error("failed to begin transaction", "error", err)
		return nil, err
	}

	err = tx.QueryRow(ctx, query, data.ID, data.ProjectID, data.Name, data.Description).
		Scan(
			&product.ID,
			&product.ProjectID,
			&product.Name,
			&product.Description,
			&product.Priority,
			&product.Removed,
			&product.CreatedAt,
		)

	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			log.Warn("record not found", "error", err)
			tx.Rollback(ctxRollback)
			return nil, model.ErrNotFound
		}
		log.Error("failed to update record", "error", err)
		tx.Rollback(ctxRollback)
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Error("failed to commit transaction", "error", err)
		return nil, err
	}

	log.Info("successfully updated")
	return &product, nil
}

func (r *goodsRepo) Remove(ctx context.Context, id, projectId int) (*model.ProductRemoveResponce, error) {
	op := "goods repository: removing"
	log := r.log.With(slog.String("operation", op))
	log.Debug("Call func Remove", "id", id, "projectId", projectId)

	var result model.ProductRemoveResponce

	query := fmt.Sprintf(`
        UPDATE %s
        SET
            removed = true
        WHERE id = $1 AND project_id = $2
        RETURNING id, project_id, removed
    `, tableName)

	err := r.DB.QueryRow(ctx, query, id, projectId).Scan(
		&result.ID,
		&result.ProjectID,
		&result.Removed,
	)

	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			log.Warn("record not found", "error", err)
			return nil, model.ErrNotFound
		}
		log.Error("failed to remove record", "error", err)
		return nil, err
	}

	log.Info("successfully removed")
	return &result, nil
}

func (r *goodsRepo) List(ctx context.Context, offset, limit int) (*model.ProductListResponce, error) {
	op := "goods repository: goods list retrieval"
	log := r.log.With(slog.String("operation", op))
	log.Debug("Call func List", "offset", offset, "limit", limit)

	var result model.ProductListResponce

	metaQuery := fmt.Sprintf(`
		WITH filtered_goods AS (
			SELECT *
			FROM %s
			ORDER BY priority DESC
			LIMIT $1 OFFSET $2
		)
		SELECT
			(SELECT COUNT(*) FROM %s) AS total,
			(SELECT COUNT(*) FROM filtered_goods WHERE removed = true) AS removed,
			$1 AS "limit",
			$2 AS "offset";
    `, tableName, tableName)

	if err := r.DB.QueryRow(ctx, metaQuery, limit, offset).
		Scan(
			&result.Meta.Total,
			&result.Meta.Removed,
			&result.Meta.Limit,
			&result.Meta.Offset,
		); err != nil {
		log.Error("failed metaQuery", "error", err)
		return nil, err
	}

	listQuery := fmt.Sprintf(`
			SELECT id, project_id, name, description, priority, removed, created_at
			FROM %s
			ORDER BY priority DESC
			LIMIT $1 OFFSET $2;
		`, tableName)

	rows, err := r.DB.Query(ctx, listQuery, limit, offset)
	if err != nil {
		log.Error("failed to get goods list", "error", err)
		return nil, err
	}
	defer rows.Close()

	list := make([]model.Product, 0, 10)
	for rows.Next() {
		var product model.Product
		if err := rows.Scan(
			&product.ID,
			&product.ProjectID,
			&product.Name,
			&product.Description,
			&product.Priority,
			&product.Removed,
			&product.CreatedAt,
		); err != nil {
			log.Error("failed to scan row", "error", err)
			return nil, err
		}
		list = append(list, product)
	}

	result.Goods = list

	if err := rows.Err(); err != nil {
		log.Error("error while iterating over rows", "error", err)
		return nil, err
	}

	log.Info("successful search")
	return &result, nil
}

func (r *goodsRepo) Reprioritizy(ctx context.Context, data model.ProductReprioritizyRequest) (*model.ProductReprioritizyResponce, error) {
	op := "goods service: reprioritizing"
	log := r.log.With(slog.String("operation", op))
	log.Debug("Call func Reprioritizy", "data", data)

	var result model.ProductReprioritizyResponce

	queryMaxPriority := fmt.Sprintf(`
		SELECT MAX(priority)
		FROM %s
		WHERE project_id = $1
    `, tableName)

	var maxPriority int
	err := r.DB.QueryRow(ctx, queryMaxPriority, data.ProjectID).Scan(&maxPriority)
	if err != nil {
		log.Error("failed to get max priority", "error", err)
		return nil, err
	}

	if data.NewPriority > maxPriority {
		log.Warn("new priority is higher than current max")
		return nil, model.ErrMaxPriority
	}

	queryCurrentPriority := fmt.Sprintf(`
        SELECT priority
        FROM %s
        WHERE id = $1 AND project_id = $2
    `, tableName)

    var currentPriority int
    err = r.DB.QueryRow(ctx, queryCurrentPriority, data.ID, data.ProjectID).Scan(&currentPriority)
    if err != nil {
        log.Error("failed to get current priority", "error", err)
        return nil, err
    }

	if currentPriority == data.NewPriority {
		log.Warn("new priority, equal to current priority")
		return nil, model.ErrCurrentPriority
	}

	query := fmt.Sprintf(`
		WITH target AS (
			SELECT id, priority
			FROM %s
			WHERE id = $1 AND project_id = $2
		),
		new_priority_item AS (
			SELECT id, priority
			FROM %s
			WHERE project_id = $2 AND priority = $3 AND id != (SELECT id FROM target)
		)
		UPDATE %s g
		SET priority = CASE
			WHEN g.id = (SELECT id FROM target) THEN (SELECT priority FROM new_priority_item)
			WHEN g.id = (SELECT id FROM new_priority_item) THEN (SELECT priority FROM target)
			ELSE g.priority
		END
		FROM target, new_priority_item
		WHERE g.id IN (target.id, new_priority_item.id)
		RETURNING g.id, g.project_id, g.priority
    `, tableName, tableName, tableName)

	rows, err := r.DB.Query(ctx, query, data.ID, data.ProjectID, data.NewPriority)
	if err != nil {
		log.Error("failed to execute query", "error", err)
		return nil, err
	}
	defer rows.Close()

	found := false
	for rows.Next() {
		var item struct {
			ID        int `json:"id"`
			ProjectID int `json:"project_id"`
			Priority  int `json:"priority"`
		}
		if err := rows.Scan(&item.ID, &item.ProjectID, &item.Priority); err != nil {
			log.Error("failed to scan row", "error", err)
			return nil, err
		}
		result.Priorities = append(result.Priorities, item)
		found = true
	}

	if err := rows.Err(); err != nil {
		log.Error("error while iterating over rows", "error", err)
		return nil, err
	}

	if !found {
		log.Warn("records not found")
		return nil, model.ErrNotFound
	}

	log.Info("successfully reprioritized")
	return &result, nil
}
