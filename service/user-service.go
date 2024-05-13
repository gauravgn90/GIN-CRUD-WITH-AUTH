package service

import (
	"database/sql"
	"gauravgn90/gin-crud-with-auth/v2/connection"
	"gauravgn90/gin-crud-with-auth/v2/model"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	SaveUser(user model.User) (model.User, error)
	FindAll() ([]model.User, error)
	Delete(id int) error
	Update(id int, user model.User) error
}

type UserServiceImpl struct{}

func New() UserService {
	return &UserServiceImpl{}
}

func (service *UserServiceImpl) Update(id int, user model.User) error {

	db := connection.GetDB()

	resultChan := make(chan sql.Result)
	errorChan := make(chan error)

	worker := NewWorker(db)
	worker.Start()

	worker.JobQueue <- Job{
		Id:         id,
		Type:       "update",
		User:       user,
		ResultChan: resultChan,
		ErrorChan:  errorChan,
	}

	select {
	case err := <-errorChan:
		return err
	case record := <-resultChan:
		_, err := record.RowsAffected()
		if err != nil {
			return err
		}
		return nil
	}
}

func UpdateUser(user model.User) error {
	db := connection.GetDB()

	resultChan := make(chan sql.Result)
	errorChan := make(chan error)

	worker := NewWorker(db)
	worker.Start()

	worker.JobQueue <- Job{
		Type:       "update",
		User:       user,
		ResultChan: resultChan,
		ErrorChan:  errorChan,
	}

	select {
	case err := <-errorChan:
		return err
	case record := <-resultChan:
		_, err := record.RowsAffected()
		if err != nil {
			return err
		}
		return nil
	}
}

func (service *UserServiceImpl) Delete(id int) error {

	db := connection.GetDB()

	resultChan := make(chan sql.Result)
	errorChan := make(chan error)

	worker := NewWorker(db)
	worker.Start()

	worker.JobQueue <- Job{
		Type:       "delete",
		ResultChan: resultChan,
		ErrorChan:  errorChan,
	}

	select {
	case err := <-errorChan:
		return err
	case record := <-resultChan:
		_, err := record.RowsAffected()
		if err != nil {
			return err
		}
		return nil
	}
}

func (service *UserServiceImpl) SaveUser(user model.User) (model.User, error) {

	db := connection.GetDB()

	resultChan := make(chan sql.Result)
	errorChan := make(chan error)

	worker := NewWorker(db)
	worker.Start()

	worker.JobQueue <- Job{
		Type:       "save",
		User:       user,
		ResultChan: resultChan,
		ErrorChan:  errorChan,
	}

	select {
	case err := <-errorChan:
		return model.User{}, err
	case record := <-resultChan:
		lastInsertId, err := record.LastInsertId()
		if err != nil {
			return model.User{}, err
		}
		user.Id = int(lastInsertId)
		return user, nil
	}
}

func (service *UserServiceImpl) FindAll() ([]model.User, error) {
	db := connection.GetDB()

	resultChan := make(chan []model.User)
	errorChan := make(chan error)

	worker := NewWorker(db)
	worker.Start()

	worker.JobQueue <- Job{
		Type:       "find",
		ResultChan: resultChan,
		ErrorChan:  errorChan,
	}

	select {
	case err := <-errorChan:
		return nil, err
	case users := <-resultChan:
		return users, nil
	}
}

type Job struct {
	Type       string
	Id         int
	User       model.User
	ResultChan interface{}
	ErrorChan  chan<- error
}

type Worker struct {
	JobQueue chan Job
}

func NewWorker(db *sql.DB) *Worker {
	return &Worker{JobQueue: make(chan Job)}
}

func (w *Worker) Start() {
	for i := 0; i < 10; i++ { // Limit the number of worker goroutines
		go w.process()
	}
}

func (w *Worker) process() {

	db := connection.GetDB()

	for job := range w.JobQueue {
		switch job.Type {
		case "save":
			w.handleSaveUserJob(db, job, job.ResultChan.(chan sql.Result))
		case "find":
			w.handleFindAllUsersJob(db, job, job.ResultChan.(chan []model.User))
		case "delete":
			w.handleDeleteUserJob(db, job, job.ResultChan.(chan sql.Result))
		case "update":
			w.handleUpdateUserJob(db, job, job.ResultChan.(chan sql.Result))

		}
		// switch resultChan := job.ResultChan.(type) {
		// case chan sql.Result:
		// 	w.handleSaveUserJob(db, job, resultChan)

		// case chan []model.User:
		// 	w.handleFindAllUsersJob(db, job, resultChan)
		// }
	}
}

func (w *Worker) handleSaveUserJob(db *sql.DB, job Job, resultChan chan<- sql.Result) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(job.User.Password), bcrypt.DefaultCost)
	if err != nil {
		job.ErrorChan <- err
		close(resultChan)
		close(job.ErrorChan)
		return
	}
	job.User.Password = string(hashedPassword)
	stmt, err := db.Prepare("INSERT INTO users (name, username, email, password) VALUES (?, ?, ?, ?)")
	if err != nil {
		job.ErrorChan <- err
		close(resultChan)
		close(job.ErrorChan)
		return
	}
	defer stmt.Close()
	result, err := stmt.Exec(job.User.Name, job.User.Username, job.User.Email, job.User.Password)
	if err != nil {
		job.ErrorChan <- err
	} else {
		resultChan <- result
	}
	close(resultChan)
	close(job.ErrorChan)
}

func (w *Worker) handleUpdateUserJob(db *sql.DB, job Job, resultChan chan<- sql.Result) {
	stmt, err := db.Prepare("UPDATE users SET name = ?, username = ?, email = ?, password = ? WHERE id = ?")
	if err != nil {
		job.ErrorChan <- err
		close(resultChan)
		close(job.ErrorChan)
		return
	}
	defer stmt.Close()
	result, err := stmt.Exec(job.User.Name, job.User.Username, job.User.Email, job.User.Password, job.Id)
	if err != nil {
		job.ErrorChan <- err
	} else {
		resultChan <- result
	}
	close(resultChan)
	close(job.ErrorChan)
}

func (w *Worker) handleDeleteUserJob(db *sql.DB, job Job, resultChan chan<- sql.Result) {
	stmt, err := db.Prepare("DELETE FROM users WHERE id = ?")
	if err != nil {
		job.ErrorChan <- err
		close(resultChan)
		close(job.ErrorChan)
		return
	}
	defer stmt.Close()
	result, err := stmt.Exec(job.Id)
	if err != nil {
		job.ErrorChan <- err
	} else {
		resultChan <- result
	}
	close(resultChan)
	close(job.ErrorChan)
}

func (w *Worker) handleFindAllUsersJob(db *sql.DB, job Job, resultChan chan<- []model.User) {
	rows, err := db.Query("SELECT * FROM users LIMIT 10")
	if err != nil {
		job.ErrorChan <- err
		return
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.Id, &user.Name, &user.Username, &user.Email, &user.Password)
		if err != nil {
			job.ErrorChan <- err
			return
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		job.ErrorChan <- err
		return
	}

	resultChan <- users
	close(resultChan)
	close(job.ErrorChan)
}
