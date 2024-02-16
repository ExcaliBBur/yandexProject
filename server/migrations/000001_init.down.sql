DROP INDEX idempotency_index;

DROP TRIGGER set_time ON expression;
DROP TRIGGER set_time_finish ON expression;

DROP FUNCTION upsert_workers(text);
DROP FUNCTION upsert_task(integer, integer, text, float);
DROP FUNCTION update_workers_and_select(integer);
DROP FUNCTION clean_up_workers(integer);
DROP FUNCTION set_time();
DROP FUNCTION set_time_finish();

DROP TABLE duration;
DROP TABLE idempotency;
DROP TABLE worker;
DROP TABLE task;
DROP TABLE expression;