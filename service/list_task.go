package service

import (
	"context"
	"fmt"

	"github.com/mongolmongol2022/go_todo_app_HandsOnTraining/entity"
	"github.com/mongolmongol2022/go_todo_app_HandsOnTraining/store"
)

type ListTask struct {
	DB   store.Queryer
	Repo TaskLister
}

func (l *ListTask) ListTasks(ctx context.Context) (entity.Tasks, error) {
	ts, err := l.Repo.ListTasks(ctx, l.DB)
	if err != nil {
		return nil, fmt.Errorf("failed to list: %w", err)
	}
	return ts, nil
}