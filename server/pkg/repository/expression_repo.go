package repository

import (
	"database/sql"
	"server/dto"
	entity "server/model/entity"
)

type ExpressionRepo struct {
	db      *sql.DB
	channel chan entity.Task
}

func NewExpressionRepo(db *sql.DB, channel chan entity.Task) *ExpressionRepo {
	return &ExpressionRepo{db: db, channel: channel}
}

func (r *ExpressionRepo) CreateExpression(expression entity.Expression) (int, error) {
	var id int
	query, err := r.db.Prepare("INSERT INTO expression (expression) VALUES ($1) RETURNING id")
	if err != nil {
		return 0, err
	}
	defer query.Close()
	row, err := query.Query(expression.Expression)

	if err != nil {
		return 0, err
	}

	for row.Next() {
		if err := row.Scan(&id); err != nil {
			return 0, err
		}
	}
	return id, nil
}

func (r *ExpressionRepo) UpdateExpression(expression entity.Expression) error {
	query, err := r.db.Prepare("UPDATE expression SET result=$1, is_finished=$2 , is_error=$3 WHERE id=$4")
	if err != nil {
		return err
	}
	defer query.Close()
	query.Exec(expression.Result, expression.IsFinished, expression.IsError, expression.Id)
	return nil
}

func (r *ExpressionRepo) CreateTask(task entity.Task) error {
	query, err := r.db.Prepare("INSERT INTO task (expression_id, task_id, task, result) VALUES ($1, $2, $3, $4)")
	if err != nil {
		return err
	}
	defer query.Close()
	query.Exec(task.ExpressionId, task.TaskId, task.Expression, task.Result)
	return nil
}

func (r *ExpressionRepo) UpdateTask(task entity.Task) error {
	query, err := r.db.Prepare(
		`UPDATE task SET task=$1, result=$2, is_completed=$3 ` +
			`WHERE expression_id=$4 AND task_id=$5`)
	if err != nil {
		return err
	}
	defer query.Close()

	query.Exec(task.Expression, task.Result, true, task.ExpressionId, task.TaskId)
	return nil
}

func (r *ExpressionRepo) SendMessage(bytes []byte) error {
	return SendMessage(bytes, "SendTask")
}

func (r *ExpressionRepo) GetChannel() chan entity.Task {
	return r.channel
}

func (r *ExpressionRepo) GetExpressions(pageNumber, pageSize int) ([]dto.ExpressionResponse, error) {
	var expressions []dto.ExpressionResponse
	expr := dto.ExpressionResponse{}

	query, err := r.db.Prepare("SELECT * FROM expression limit $1 offset $2")
	if err != nil {
		return nil, err
	}
	defer query.Close()
	rows, err := query.Query(pageSize, pageSize*pageNumber)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err = rows.Scan(&expr.Id, &expr.Expression, &expr.Result, &expr.DateStart, &expr.DateFinish, &expr.IsFinished, &expr.IsError)
		if err != nil {
			return nil, err
		}
		expressions = append(expressions, expr)
	}

	return expressions, nil
}

func (r *ExpressionRepo) GetExpression(id int) (dto.ExpressionResponse, error) {
	expr := dto.ExpressionResponse{}

	query, err := r.db.Prepare("SELECT * FROM expression WHERE id = $1")
	if err != nil {
		return dto.ExpressionResponse{}, err
	}
	defer query.Close()
	row, err := query.Query(id)
	if err != nil {
		return dto.ExpressionResponse{}, err
	}
	for row.Next() {
		err = row.Scan(&expr.Id, &expr.Expression, &expr.Result, &expr.DateStart, &expr.DateFinish, &expr.IsFinished, &expr.IsError)
		if err != nil {
			return dto.ExpressionResponse{}, err
		}
	}
	return expr, nil
}

func (r *ExpressionRepo) GetUnfinishedExpressions() ([]entity.Expression, error) {
	var expressions []entity.Expression
	expr := entity.Expression{}

	query, err := r.db.Prepare(`SELECT id, expression, result, is_finished, is_error ` +
		`FROM expression WHERE is_finished = false`)
	if err != nil {
		return nil, err
	}
	rows, err := query.Query()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err = rows.Scan(&expr.Id, &expr.Expression, &expr.Result, &expr.IsFinished, &expr.IsError)
		if err != nil {
			return nil, err
		}
		expressions = append(expressions, expr)
	}

	return expressions, nil
}
