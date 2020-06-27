package models

import (
	"bfimpl/services"
	"errors"
	"strconv"
	"time"
)

var (
	UserList map[string]*User
)

func init() {
	UserList = make(map[string]*User)
	UserList["wetest"] = &User{"demo1", "demo1", Profile{"male", 20, "a@demo.com"}}
	UserList["cloud"] = &User{"demo2", "demo2", Profile{"male", 20, "b@demo.com"}}
}

type User struct {
	Id       string
	Username string
	Profile  Profile
}

type Profile struct {
	Gender string
	Age    int
	Email  string
}

func AddUser(u User) string {
	u.Id = "user_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	UserList[u.Id] = &u
	return u.Id
}

func GetUser(uid string) (u *User, err error) {
	if u, ok := UserList[uid]; ok {
		return u, nil
	}
	return nil, errors.New("User not exists")
}

func GetAllUsers() map[string]*User {
	return UserList
}

func UpdateUser(uid string, uu *User) (a *User, err error) {
	if u, ok := UserList[uid]; ok {
		if uu.Username != "" {
			u.Username = uu.Username
		}
		if uu.Profile.Age != 0 {
			u.Profile.Age = uu.Profile.Age
		}
		if uu.Profile.Gender != "" {
			u.Profile.Gender = uu.Profile.Gender
		}
		if uu.Profile.Email != "" {
			u.Profile.Email = uu.Profile.Email
		}
		return u, nil
	}
	return nil, errors.New("User Not Exist")
}

func DeleteUser(uid string) {
	delete(UserList, uid)
}

type Table struct {
	TableSchema   string
	TableName     string
	Engine        string
	AutoIncrement int64
}

func (t *Table) GetAllTables() *[]Table {
	tables := make([]Table, 0)
	services.Slave().
		Table("information_schema.tables").
		Where("table_schema != ?", "information_schema").
		Where("table_schema != ?", "mysql").
		Select("table_schema, table_name, engine, auto_increment").
		Scan(&tables)
	return &tables
}
