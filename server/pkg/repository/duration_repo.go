package repository

import (
	"database/sql"
	entity "server/model/entity"
)

type DurationRepo struct {
	db *sql.DB
}

func NewDurationRepo(db *sql.DB) *DurationRepo {
	return &DurationRepo{db: db}
}

func (r *DurationRepo) SendMessage(bytes []byte) error {
	return SendMessage(bytes, "SendDuration")
}

func (r *DurationRepo) UpdateDuration(duration entity.Duration) error {
	query, err := r.db.Prepare(`UPDATE duration SET plus_duration_ms=$1, minus_duration_ms=$2, ` +
		`mul_duration_ms=$3, div_duration_ms=$4, heartbeat_duration_s=$5 WHERE id=$6`)
	if err != nil {
		return err
	}
	defer query.Close()
	query.Exec(duration.PlusDuration, duration.MinusDuration,
		duration.MulDuration, duration.DivDuration, duration.HeartBeatDuration, 1)
	return nil
}

func (r *DurationRepo) GetDuration() (entity.Duration, error) {
	query, err := r.db.Prepare("SELECT plus_duration_ms, minus_duration_ms, mul_duration_ms, " +
		"div_duration_ms, heartbeat_duration_s FROM duration WHERE id=$1")

	if err != nil {
		return entity.Duration{}, err
	}
	defer query.Close()

	duration := entity.Duration{}
	row, err := query.Query(1)
	if err != nil {
		return entity.Duration{}, err
	}

	for row.Next() {
		err = row.Scan(&duration.PlusDuration, &duration.MinusDuration, &duration.MulDuration,
			&duration.DivDuration, &duration.HeartBeatDuration)
		if err != nil {
			return entity.Duration{}, err
		}
	}
	return duration, nil
}
