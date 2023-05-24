-- sqlite3 todolist.db
-- psql -h localhost -U postgres -d postgres
-- \c zze
CREATE TABLE IF NOT EXISTS tasks (
                                     task_id SERIAL PRIMARY KEY,
                                     task_desc TEXT NOT NULL,
                                     task_date TEXT NOT NULL,
                                     task_time TEXT NOT NULL,
                                     created_time TEXT NOT NULL,
                                     done INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS models (
                                      task_id SERIAL PRIMARY KEY,
                                      task_desc TEXT NOT NULL,
                                      task_time TEXT NOT NULL,
                                      created_time TEXT NOT NULL
);