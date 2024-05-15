package service

import (
	"database/sql"
	"gauravgn90/gin-crud-with-auth/v2/connection"
	"gauravgn90/gin-crud-with-auth/v2/model"
	"gauravgn90/gin-crud-with-auth/v2/utility"

	"github.com/jinzhu/gorm"
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
	resultChan := make(chan interface{})
	errorChan := make(chan error)

	worker := NewWorker()
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
	case <-resultChan:
		return nil
	}
}

func UpdateUser(user model.UserUpdate) error {

	resultChan := make(chan sql.Result)
	errorChan := make(chan error)

	worker := NewWorker()
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

	resultChan := make(chan sql.Result)
	errorChan := make(chan error)

	worker := NewWorker()
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

	resultChan := make(chan interface{})
	errorChan := make(chan error)

	worker := NewWorker()
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
		lastInsertId := record.(int)
		user.Id = int(lastInsertId)
		return user, nil
	}
}

func (service *UserServiceImpl) FindAll() ([]model.User, error) {
	resultChan := make(chan []model.User)
	errorChan := make(chan error)

	worker := NewWorker()
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
	User       model.UserType
	ResultChan interface{}
	ErrorChan  chan<- error
}

type Worker struct {
	JobQueue chan Job
}

func NewWorker() *Worker {
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
			w.handleSaveUserJob(db, job, job.ResultChan.(chan interface{}))
		case "find":
			w.handleFindAllUsersJob(db, job, job.ResultChan.(chan []model.User))
		case "update":
			w.handleUpdateUserJob(db, job, job.ResultChan.(chan interface{}))
			/*case "delete":
			w.handleDeleteUserJob(db, job, job.ResultChan.(chan sql.Result))
			*/

		}
		// switch resultChan := job.ResultChan.(type) {
		// case chan sql.Result:
		// 	w.handleSaveUserJob(db, job, resultChan)

		// case chan []model.User:
		// 	w.handleFindAllUsersJob(db, job, resultChan)
		// }
	}
}

func (w *Worker) handleSaveUserJob(db *gorm.DB, job Job, resultChan chan<- interface{}) {
	user := job.User.(model.User)
	if err := db.Where("username = ?", user.Username).First(&user).Error; err == nil {
		job.ErrorChan <- utility.NewCustomError(400, "username already registered")
		close(resultChan)
		close(job.ErrorChan)
		return
	}
	if err := db.Where("email = ?", user.Email).First(&user).Error; err == nil {
		job.ErrorChan <- utility.NewCustomError(400, "e-mail already registered")
		close(resultChan)
		close(job.ErrorChan)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		job.ErrorChan <- err
		close(resultChan)
		close(job.ErrorChan)
		return
	}
	user.Password = string(hashedPassword)
	// prepare statement with gorm
	db = db.Create(&user)
	if db.Error != nil {
		job.ErrorChan <- db.Error
	} else {
		resultChan <- user.Id
	}
	close(resultChan)
	close(job.ErrorChan)
}

func (w *Worker) handleUpdateUserJob(db *gorm.DB, job Job, resultChan chan<- interface{}) {

	db = db.Model(&model.UserUpdate{}).Where("id = ?", job.Id).Updates(&job.User)
	if db.Error != nil {
		job.ErrorChan <- db.Error
	} else {
		resultChan <- db.RowsAffected
	}
	close(resultChan)
	close(job.ErrorChan)
}

/*
	func (w *Worker) handleDeleteUserJob(db *gorm.DB, job Job, resultChan chan<- sql.Result) {
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
*/
func (w *Worker) handleFindAllUsersJob(db *gorm.DB, job Job, resultChan chan<- []model.User) {
	var users []model.User
	err := db.Limit(10).Find(&users).Error
	if err != nil {
		job.ErrorChan <- err
	} else {
		resultChan <- users
	}
	close(resultChan)
	close(job.ErrorChan)
}
