package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type Task struct {
	TaskId      int    `json:"task_id"`
	TaskDesc    string `json:"task_desc"`
	TaskDate    string `json:"task_date"`
	TaskTime    string `json:"task_time"`
	CreatedTime string `json:"created_time"`
	Done        int    `json:"done"`
}

// TaskModel 任务模版
type TaskModel struct {
	TaskId      int    `json:"task_id"`
	TaskDesc    string `json:"task_desc"`
	TaskTime    string `json:"task_time"`
	CreatedTime string `json:"created_time"`
}

type server struct {
	db *sql.DB
}

func formateTime(timeStr string) string {
	// 使用字符串格式化将小时和分钟转换为 2 位数
	if len(timeStr) <= 4 { //9:30
		return "0" + timeStr
	}
	return timeStr

}
func (s *server) checkTaskExists(taskDesc string) (bool, error) {
	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM tasks WHERE task_desc = ? LIMIT 1)", taskDesc).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s *server) getTasks(w http.ResponseWriter, r *http.Request) {
	today := time.Now().Format("2006-01-02")
	query := fmt.Sprintf(`SELECT task_id, task_desc, task_date, created_time, done,task_time FROM tasks Where task_date="%s" ORDER BY task_time`, today)
	rows, err := s.db.Query(query)
	if err != nil {
		log.Println("getTasks:", err)
		log.Println("sql is:", query)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	defer rows.Close()

	tasks := []Task{}
	for rows.Next() {
		task := Task{}
		err := rows.Scan(&task.TaskId, &task.TaskDesc, &task.TaskDate, &task.CreatedTime, &task.Done, &task.TaskTime)
		if err != nil {
			log.Println("get tasks err: ", err)
			continue
		}
		tasks = append(tasks, task)
	}

	json.NewEncoder(w).Encode(tasks)
}

func (s *server) getDoneTasks(w http.ResponseWriter, r *http.Request) {
	rows, err := s.db.Query("SELECT task_id, task_desc, task_date, created_time, done,task_time FROM tasks WHERE done = 1 ORDER BY task_time")
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	defer rows.Close()

	tasks := []Task{}
	for rows.Next() {
		task := Task{}
		err := rows.Scan(&task.TaskId, &task.TaskDesc, &task.TaskDate, &task.CreatedTime, &task.Done)
		if err != nil {
			json.NewEncoder(w).Encode(err.Error())
			return
		}
		tasks = append(tasks, task)
	}

	json.NewEncoder(w).Encode(tasks)
}

func (s *server) createTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	desc := r.FormValue("task_desc")
	date := r.FormValue("task_date")
	timeExe := formateTime(r.FormValue("task_time"))
	task.TaskTime = timeExe
	task.TaskDesc = desc
	task.TaskDate = date

	has, err := s.checkTaskExists(desc)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	if has {
		json.NewEncoder(w).Encode("ok")
		return
	}

	task.CreatedTime = time.Now().Format("2006-01-02 15:04:05")

	result, err := s.db.Exec("INSERT INTO tasks (task_desc, task_date, created_time, done,task_time) VALUES (?, ?, ?, ?,?)",
		task.TaskDesc, task.TaskDate, task.CreatedTime, 0, task.TaskTime)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	task.TaskId = int(id)

	json.NewEncoder(w).Encode(task)
}

func (s *server) updateTask(w http.ResponseWriter, r *http.Request) {
	done := r.FormValue("done")
	vars := mux.Vars(r)
	taskId, err := strconv.Atoi(vars["task_id"])
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	_, err = s.db.Exec("UPDATE tasks SET done = ? WHERE task_id = ?", done, taskId)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	json.NewEncoder(w).Encode("ok")
}

func (s *server) deleteTask(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	taskId, err := strconv.Atoi(vars["task_id"])
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	_, err = s.db.Exec("DELETE FROM tasks WHERE task_id = ?", taskId)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *server) createModel(w http.ResponseWriter, r *http.Request) {

	var model TaskModel
	desc := r.FormValue("task_desc")
	timeExe := formateTime(r.FormValue("task_time"))
	model.TaskTime = timeExe
	model.TaskDesc = desc

	model.CreatedTime = time.Now().Format("2006-01-02 15:04:05")

	_, err := s.db.Exec("INSERT INTO models (task_desc, created_time,task_time) VALUES (?, ?, ?)",
		model.TaskDesc, model.CreatedTime, model.TaskTime)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	json.NewEncoder(w).Encode("ok")
}

func (s *server) updateModel(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	taskId, err := strconv.Atoi(vars["task_id"])
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	desc := r.FormValue("task_desc")
	timeExe := formateTime(r.FormValue("task_time"))
	_, err = s.db.Exec("UPDATE models SET task_desc = ?,task_time=? WHERE task_id = ?", desc, timeExe, taskId)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	json.NewEncoder(w).Encode("ok")
}

func (s *server) deleteModel(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	taskId, err := strconv.Atoi(vars["task_id"])
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	_, err = s.db.Exec("DELETE FROM models WHERE task_id = ?", taskId)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *server) getModels(w http.ResponseWriter, r *http.Request) {

	query := "SELECT task_id, task_desc, created_time,task_time FROM models  ORDER BY task_time"

	rows, err := s.db.Query(query)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	defer rows.Close()

	models := []TaskModel{}
	for rows.Next() {
		model := TaskModel{}
		err := rows.Scan(&model.TaskId, &model.TaskDesc, &model.CreatedTime, &model.TaskTime)
		if err != nil {
			log.Println("get models err: ", err)
			continue
		}
		models = append(models, model)
	}

	json.NewEncoder(w).Encode(models)
}

func (s *server) modelsToTask(w http.ResponseWriter, r *http.Request) {

	rows, err := s.db.Query("SELECT task_id, task_desc, created_time,task_time FROM models  ORDER BY task_time")
	if err != nil {
		log.Println("query model err:", err)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	tasks := []Task{}
	for rows.Next() {
		model := TaskModel{}
		err := rows.Scan(&model.TaskId, &model.TaskDesc, &model.CreatedTime, &model.TaskTime)
		if err != nil {
			log.Println("list model err: ", err)
			continue
		}
		has, err := s.checkTaskExists(model.TaskDesc)
		if err != nil {
			log.Println("check exist err")
			continue
		}
		if has {
			log.Println("already has")
			continue
		}

		log.Println("insert task: ", model.TaskDesc)

		createTime := time.Now().Format("2006-01-02 15:04:05")
		execDate := time.Now().Format("2006-01-02")

		task := Task{}
		task.TaskDesc = model.TaskDesc
		task.TaskTime = model.TaskTime
		task.TaskDate = execDate
		task.CreatedTime = createTime
		tasks = append(tasks, task)

	}
	rows.Close()

	for _, task := range tasks {
		_, err = s.db.Exec("INSERT INTO tasks (task_desc, task_date, created_time, done,task_time) VALUES (?, ?, ?, ?,?)",
			task.TaskDesc, task.TaskDate, task.CreatedTime, 0, task.TaskTime)
		if err != nil {
			log.Println("model to task err: ", err.Error(), " id:", task.TaskId, " desc:", task.TaskDesc)
			continue
		}
	}

	json.NewEncoder(w).Encode("ok")
}

func main() {
	db, err := sql.Open("sqlite3", "./todolist.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	s := &server{
		db: db,
	}
	router := mux.NewRouter()

	// Add a middleware that allows cross-origin resource sharing (CORS)
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Set the Access-Control-Allow-Origin header to allow requests from any origin
			w.Header().Set("Access-Control-Allow-Origin", "*")

			// Call the next handler
			next.ServeHTTP(w, r)
		})
	})

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./index.html")
	})
	router.HandleFunc("/api/tasks", s.getTasks).Methods("GET")
	router.HandleFunc("/api/done", s.getDoneTasks).Methods("GET")
	router.HandleFunc("/api/tasks", s.createTask).Methods("POST")
	router.HandleFunc("/api/tasks/{task_id}", s.updateTask).Methods("PUT")
	router.HandleFunc("/api/tasks/{task_id}", s.deleteTask).Methods("DELETE")
	router.HandleFunc("/api/model", s.createModel).Methods("POST")
	router.HandleFunc("/api/models/{task_id}", s.updateModel).Methods("PUT")
	router.HandleFunc("/api/models/{task_id}", s.deleteModel).Methods("DELETE")
	router.HandleFunc("/api/models", s.getModels).Methods("GET")
	router.HandleFunc("/api/modeltotask", s.modelsToTask).Methods("GET")
	log.Fatal(http.ListenAndServe(":2305", router))
}
