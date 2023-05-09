-- sqlite3 todolist.db


-- 创建 tasks 表
CREATE TABLE IF NOT EXISTS tasks (
                                     task_id INTEGER PRIMARY KEY AUTOINCREMENT,
                                     task_desc TEXT NOT NULL,
                                     task_date TEXT NOT NULL,
                                     task_time TEXT NOT NULL,
                                     created_time TEXT NOT NULL,
                                     done INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS models (
                                     task_id INTEGER PRIMARY KEY AUTOINCREMENT,
                                     task_desc TEXT NOT NULL,
                                     task_time TEXT NOT NULL,
                                     created_time TEXT NOT NULL
);