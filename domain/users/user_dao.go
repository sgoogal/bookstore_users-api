package users

import (
	"fmt"

	"github.com/sgoogal/bookstore_users-api/datasources/mysql/users_db"
	"github.com/sgoogal/bookstore_users-api/logger"
	"github.com/sgoogal/bookstore_users-api/utils/errors"
)

const (
	queryInsertUser       = "INSERT INTO users(first_name, last_name, email, date_created, status, password) VALUES(?, ?, ?, ?, ?, ?);"
	queryGetUser          = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE id=?;"
	queryUpdateUser       = "UPDATE users SET first_name=?, last_name=?, email=? WHERE id=?;"
	queryDeleteUser       = "DELETE FROM users WHERE id=?;"
	queryFindUserByStatus = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE status=?;"
	queryGetLoginUser     = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE email=? AND password=?;"
)

func (user *User) Get() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryGetUser)
	if err != nil {
		logger.Error("error while trying to prepare get user statement", err)
		return errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.Id)
	if getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); getErr != nil {
		logger.Error("error while trying to get user by id", getErr)
		return errors.NewInternalServerError("database error")
	}
	logger.Info("get user executed")
	return nil
}

func (user *User) Save() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		logger.Error("error while trying to create user statement", err)
		return errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	//user.DateCreated = date_utils.GetNowString()

	insertResult, saveErr := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated, user.Status, user.Password)
	if saveErr != nil {
		logger.Error("error while trying to create user", saveErr)
		return errors.NewInternalServerError("database error")
	}

	userId, err := insertResult.LastInsertId()
	if err != nil {
		logger.Error("error while trying to get last insert id after creating new user", err)
		return errors.NewInternalServerError("database error")
	}
	user.Id = userId
	return nil
}

func (user *User) Update() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryUpdateUser)
	if err != nil {
		logger.Error("error while trying to prepare update statement", err)
		return errors.NewInternalServerError("database error")
	}

	_, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.Id)
	if err != nil {
		logger.Error("error while trying to execute update statement", err)
		return errors.NewInternalServerError("database error")
	}
	return nil
}

func (user *User) Delete() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryDeleteUser)
	if err != nil {
		logger.Error("error while trying to prepare delete statement", err)
		return errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	if _, err = stmt.Exec(user.Id); err != nil {
		logger.Error("error while trying to execute delete statement", err)
		return errors.NewInternalServerError("database error")
	}
	return nil
}

func (user *User) FindByStatus(status string) ([]User, *errors.RestErr) {
	stmt, err := users_db.Client.Prepare(queryFindUserByStatus)
	if err != nil {
		logger.Error("error while trying to prepare find users by status statement", err)
		return nil, errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	rows, err := stmt.Query(status)
	if err != nil {
		logger.Error("error while trying to execute find users by status statement", err)
		return nil, errors.NewInternalServerError("database error")
	}
	defer rows.Close()

	results := make([]User, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); err != nil {
			logger.Error("error while trying scan users while finding them", err)
			return nil, errors.NewInternalServerError("database error")
		}
		results = append(results, user)
	}

	if len(results) == 0 {
		return nil, errors.NewNotFoundError(fmt.Sprintf("no users found matching status %s", status))
	}
	return results, nil
}

func (user *User) GetLoginUser() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryGetLoginUser)
	if err != nil {
		logger.Error("error while trying to prepare get user by email and password statement", err)
		return errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.Email, user.Password)
	if getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); getErr != nil {
		logger.Error("error while trying to get user by email and password", getErr)
		return errors.NewInternalServerError("database error")
	}
	logger.Info("get user executed")
	return nil
}
