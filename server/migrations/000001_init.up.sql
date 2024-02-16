CREATE TABLE
    expression (
        id serial PRIMARY KEY,
        expression text NOT NULL,
        result float DEFAULT 0,
        date_start TIMESTAMP NOT NULL,
        date_finish TIMESTAMP NOT NULL,
        is_finished boolean DEFAULT false,
        is_error boolean DEFAULT false
    );

CREATE TABLE 
    task (
        expression_id int REFERENCES expression(id) ON UPDATE CASCADE ON DELETE CASCADE,
        task_id int NOT NULL,
        task text NOT NULL,
        result float DEFAULT 0.0,
        is_completed boolean DEFAULT false,
        PRIMARY KEY (expression_id, task_id)
    );

CREATE TABLE
    worker (
        hostname text PRIMARY KEY,
        id serial,
        last_heartbeat TIMESTAMP,
        is_alive boolean DEFAULT false
    );

CREATE TABLE 
    idempotency (
        idempotency_key text PRIMARY KEY,
        expression_id int NOT NULL REFERENCES expression(id) ON UPDATE CASCADE ON DELETE CASCADE
    );

CREATE TABLE
    duration (
        id int PRIMARY KEY,
        plus_duration_ms int DEFAULT 200,
        minus_duration_ms int DEFAULT 200,
        mul_duration_ms int DEFAULT 200,
        div_duration_ms int DEFAULT 200,
        heartbeat_duration_s int DEFAULT 5
    );

INSERT INTO duration (id) VALUES (1);

CREATE INDEX idempotency_index ON idempotency USING HASH("idempotency_key");
CREATE INDEX worker_index ON worker USING HASH("hostname");

CREATE OR REPLACE FUNCTION upsert_workers(name text) RETURNS VOID AS $$ 
    BEGIN 
        UPDATE worker SET last_heartbeat = NOW()::TIMESTAMP, is_alive=true WHERE worker.hostname = name; 
        IF NOT FOUND THEN 
            INSERT INTO worker (hostname, last_heartbeat, is_alive) values (name, NOW()::TIMESTAMP, true); 
        END IF;
    END; 
    $$ LANGUAGE 'plpgsql';

CREATE OR REPLACE FUNCTION upsert_task(id integer, task_id integer, task text, result float) RETURNS VOID AS $$ 
    BEGIN 
        UPDATE task SET task.id = id, task.task_id = task_id, task.task = task, task.result = result 
            WHERE task.id = id; 
        IF NOT FOUND THEN 
            INSERT INTO task (expression_id, task_id, task, result) VALUES (id, task_id, task, result);
        END IF;
    END; 
    $$ LANGUAGE 'plpgsql';

CREATE OR REPLACE FUNCTION update_workers_and_select(time_delay integer) 
RETURNS TABLE(hostname text, id int, last_heartbeat TIMESTAMP, is_alive boolean) AS $$
    BEGIN
        UPDATE worker SET is_alive = false 
            WHERE EXTRACT(EPOCH FROM (NOW()::TIMESTAMP - worker.last_heartbeat)) > time_delay;
        RETURN query SELECT * from worker;
    END;
    $$ LANGUAGE 'plpgsql';

CREATE OR REPLACE FUNCTION clean_up_workers(time_delay integer) RETURNS VOID AS $$
    BEGIN
        DELETE FROM worker
            WHERE EXTRACT(EPOCH FROM (NOW()::TIMESTAMP - worker.last_heartbeat)) > time_delay;
    END;
    $$ LANGUAGE 'plpgsql';

CREATE OR REPLACE FUNCTION set_time() RETURNS trigger AS $$
            BEGIN
                NEW.date_start := NOW()::TIMESTAMP;
                NEW.date_finish := NOW()::TIMESTAMP;
                RETURN NEW;
            END;
            $$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION set_time_finish() RETURNS trigger AS $$
            BEGIN
                NEW.date_finish := NOW()::TIMESTAMP;
                RETURN NEW;
            END;
            $$ LANGUAGE plpgsql;


CREATE OR REPLACE TRIGGER set_time_finish
            BEFORE UPDATE ON expression
            FOR EACH ROW 
            WHEN (OLD.is_finished is DISTINCT FROM NEW.is_finished)
            EXECUTE PROCEDURE set_time_finish();
            
CREATE OR REPLACE TRIGGER set_time
            BEFORE INSERT ON expression
            FOR EACH ROW EXECUTE PROCEDURE set_time();