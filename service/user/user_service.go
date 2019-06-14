package user

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/lib/pq"
	"github.com/training_project/core/db"
	"github.com/training_project/util"
)

type (
	User struct {
		UserID     int64  `json:"user_id",db:"user_id"`
		Fullname   string `json:"full_name",db:"full_name"`
		MSISDN     string `json:"msisdn",db:"msisdn"`
		Email      string `json:"user_email",db:"user_email"`
		BirthDate  string `json:"birth_date",db:"birth_date"`
		CreateTime string `json:"create_time",db:"create_time"`
		Age        string `json:"age",db:"age"`
	}
	UserService interface {
		List(string) []User
	}
	userServiceImpl struct {
		db db.DB
	}
)

var (
	lock     = &sync.Mutex{}
	instance UserService
)

func (userService *userServiceImpl) List(userName string) (users []User) {
	var rows *sql.Rows
	var err error

	if userName == "" {
		rows, err = userService.db.Select(LIST_OF_USER_SECOND)
	} else {
		rows, err = userService.db.Select(LIST_OF_USER_FIRST, userName)
	}

	if err == nil {
		for rows.Next() {
			temp := struct {
				UserID     int64       `json:"user_id",db:"user_id"`
				Fullname   string      `json:"full_name",db:"full_name"`
				MSISDN     string      `json:"msisdn",db:"msisdn"`
				Email      string      `json:"user_email",db:"user_email"`
				BirthDate  pq.NullTime `json:"birth_date",db:"birth_date"`
				CreateTime pq.NullTime `json:"create_time",db:"create_time"`
			}{}
			if err = rows.Scan(
				&temp.UserID,
				&temp.Fullname,
				&temp.Email,
				&temp.BirthDate,
				&temp.MSISDN,
				&temp.CreateTime,
			); err != nil {
				log.Println("Err: ", err)
			} else {
				user := User{
					UserID:   temp.UserID,
					Fullname: temp.Fullname,
					MSISDN:   temp.MSISDN,
					Email:    temp.Email,
				}
				if temp.BirthDate.Valid {
					year, _, _, _, _, _ := util.Diff(temp.BirthDate.Time, time.Now())
					user.Age = fmt.Sprintf("%d year(s)", year)
					user.BirthDate = temp.BirthDate.Time.Local().Format("2006/01/02 15:04:05")
				}
				if temp.CreateTime.Valid {
					user.CreateTime = temp.CreateTime.Time.Local().Format("2006/01/02 15:04:05")
				}
				users = append(users, user)
			}
		}
	} else {
		log.Fatal(err)
	}
	return
}

func New() UserService {
	lock.Lock()
	defer lock.Unlock()

	if instance == nil {
		instance = &userServiceImpl{
			db: db.New(),
		}
	}
	return instance
}
