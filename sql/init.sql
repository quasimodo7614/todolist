-- sqlite3 todolist.db


CREATE TABLE IF NOT EXISTS tasks (
    task_id int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
    task_desc TEXT NOT NULL,
    task_date TEXT NOT NULL,
    task_time TEXT NOT NULL,
    created_time TEXT NOT NULL,
    done INTEGER NOT NULL,
    PRIMARY KEY (task_id)
    );

CREATE TABLE IF NOT EXISTS models (
    task_id int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
    task_desc TEXT NOT NULL,
    task_time TEXT NOT NULL,
    created_time TEXT NOT NULL,
    PRIMARY KEY (task_id)
    );