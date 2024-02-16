package repository

import (
	"database/sql"
	"log"
	entity "server/model/entity"
)

type WorkerRepo struct {
	db               *sql.DB
	heartBeatChannel chan entity.Worker
}

func NewWorkerRepo(db *sql.DB, heartBeatChannel chan entity.Worker) *WorkerRepo {
	return &WorkerRepo{db: db, heartBeatChannel: heartBeatChannel}
}

func (r *WorkerRepo) GetHeartBeatChannel() chan entity.Worker {
	return r.heartBeatChannel
}

func (r *WorkerRepo) UpsertWorker(worker entity.Worker) error {
	query, err := r.db.Prepare("SELECT * FROM upsert_workers($1)")
	if err != nil {
		return err
	}
	defer query.Close()

	query.Exec(worker.Hostname)
	return nil
}

func (r *WorkerRepo) UpdateWorkers(delay int) ([]entity.Worker, error) {
	var workers []entity.Worker
	worker := entity.Worker{}

	query, err := r.db.Prepare("SELECT * FROM update_workers_and_select($1)")
	if err != nil {
		return nil, err
	}
	defer query.Close()

	rows, err := query.Query(delay)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err = rows.Scan(&worker.Hostname, &worker.Id, &worker.LastHeartBeat, &worker.IsAlive)
		if err != nil {
			return nil, err
		}
		workers = append(workers, worker)
	}

	return workers, nil
}

func (r *WorkerRepo) CleanUpWorkers(delay_seconds int) error {
	query, err := r.db.Prepare("SELECT * FROM clean_up_workers($1)")
	if err != nil {
		return err
	}
	defer query.Close()

	query.Exec(delay_seconds)
	log.Printf("Workers cleaned up succesfully")
	return nil
}
