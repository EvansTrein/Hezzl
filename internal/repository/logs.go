package repository

import (
	"context"
	"fmt"
	"hezzl/internal/model"
	"hezzl/pkg/db/clickhouse"
	"log/slog"
	"time"
)

const (
	logTableName = "goods"
	logTimer     = time.Second * 5
)

type logsRepo struct {
	log *slog.Logger
	*clickhouse.ClickhouseDB
}

type LogsRepoDeps struct {
	*slog.Logger
	*clickhouse.ClickhouseDB
}

func NewLogsRepo(deps *LogsRepoDeps) *logsRepo {
	return &logsRepo{
		log:          deps.Logger,
		ClickhouseDB: deps.ClickhouseDB,
	}
}

func (r *logsRepo) Create(data *model.Product) {
	op := "logs repository: creating"
	log := r.log.With(slog.String("operation", op))
	log.Debug("Call func Create", "data", data)

	ctx, cancel := context.WithTimeout(context.Background(), logTimer)
	defer cancel()

	query := fmt.Sprintf(`
		INSERT INTO %s (
			Id,
			ProjectId,
			Name,
			Description,
			Priority,
			Removed
		) VALUES (?, ?, ?, ?, ?, ?)
	`, logTableName)

	if _, err := r.ClickhouseDB.DB.ExecContext(
		ctx,
		query,
		data.ID,
		data.ProjectID,
		data.Name,
		data.Description,
		data.Priority,
		data.Removed,
	); err != nil {
		log.Error("failed to insert into clickhouse", "error", err)
		return
	}

	log.Info("successfully created")
}
