/*
Copyright 2017 Cms Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package model

/*
CREATE DATABASE IF NOT EXISTS my_db default charset utf8
*/

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	. "idlehanker.com/x-mixplatform/golan/util"
)

// AppConfigContent config content
var ConfigContent map[string]interface{} //config

// Db is gorm DB struct instance
var Db *gorm.DB

//App struct is a model for each app
type App struct {
	gorm.Model
	Identifier string `gorm:"not null;unique"`
	AppName    string `gorm:"not null;unique"`
	Platform   string
	ProjectID  uint32
	Versions   []Version `gorm:"ForeignKey:AppID"`
}

//Project is a project model
type Project struct {
	gorm.Model
	ProjectName   string
	Description   string
	StartDatetime time.Time
	EndDatetime   time.Time
	DepartmentID  uint32

	Apps []App `gorm:"ForeignKey:ProjectID"`
	//Users       []User `gorm:"many2many:project_users;"`
	ProjUserRole []ProjectUserRole `gorm:"ForeignKey:ProjectID"`
}

// Department model  contains current IT-Depts
type Department struct {
	gorm.Model
	DepartmentName string
	Description    string
	Project        []Project `gorm:"ForeignKey:DepartmentID"`
}

// Module Mode contains some helper methods to handle relationship things easily.
type Module struct {
	gorm.Model
	ModuleName string
	Versions   []Version `gorm:"ForeignKey:ModuleID"`
}

// Version Mode contains some versions for Module or/and App
type Version struct {
	gorm.Model
	// VersionName string
	// BuildName   string
	IssueCode string
	BuildCode string
	AppID     uint32
	ModuleID  uint32
}

//User model contains platform users
type User struct {
	gorm.Model
	// GroupID      uint32
	Name     string `gorm:"not null;unique"`
	Password string
	Email    string `gorm:"not null;unique"`
	Mobile   string
	Telphone string
	Addr     string
	IDCard   string
	Sex      string
	Age      uint
	Staff    string //inner or outer
	// Birthday     time.Time         `gorm:"null"`
	Birthday     string            `gorm:"null"`
	ProjUserRole []ProjectUserRole `gorm:"ForeignKey:UserID"`
}

//ProjectUserRole model manage relationship with project user and role
type ProjectUserRole struct {
	gorm.Model
	ProjectID uint32
	UserID    uint32
	RoleID    uint
}

//Role model manage all role for users
type Role struct {
	gorm.Model
	// RoleID       uint32            `gorm:"AUTO_INCREMENT"` //`gorm:"primary_key"`
	RoleName     string `gorm:"not null;unique"`
	Description  string
	ProjUserRole []ProjectUserRole `gorm:"ForeignKey:RoleID"`
}

//Group model manage super users
type Group struct {
	gorm.Model
	// GroupName string `gorm:"not null;unique"`
	Role  []Role `gorm:"ForeignKey:GroupID"`
	Users []User `gorm:"ForeignKey:GroupID"`
}

func loadConfig() {

	var fileName = "util/config.json"
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalln("Can not open config file", err)
	}

	defer file.Close()

	fileInfo, err := os.Lstat(fileName)
	if err != nil {
		P("Getting file-info encounter a err", err)
		return
	}
	var obj interface{} // var obj map[string]interface{}
	// json.Unmarshal([]byte(str), &obj)
	var content = make([]byte, fileInfo.Size())
	if _, err := file.Read(content); err != nil {
		P("read config content error", err)
		return
	}
	if err := json.Unmarshal([]byte(content), &obj); err != nil {

		P("Unmarshal json content err,", err)
		return
	}

	ConfigContent = obj.(map[string]interface{})
}

//UserByEmail to find user by email
func UserByEmail(email string) (user User, err error) {

	// Db.Where("email = ?", email).First(&usr)
	Db.First(&user, "Email = ?", email)

	if &user != nil {
		return user, nil
	}
	return user, errors.New("user not found")
}

func init() {

	fmt.Println("loading & init....")
	loadConfig()

	if ok := connectDatabase(); ok {
		autoMigrateTables()
	}
}

func connectDatabase() bool {

	dbCfg := ConfigContent["db_init"].(map[string]interface{})["config"].(map[string]interface{})
	dbDialect := dbCfg["dialect"].(string)
	dbPasswd := dbCfg["passwd"].(string)
	dbUser := dbCfg["user"].(string)
	dbName := dbCfg["dbname"].(string)
	connectString := fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True&loc=Local", dbUser, dbPasswd, dbName)

	var err error
	Db, err = gorm.Open(dbDialect, connectString)
	Db.Debug()

	if err != nil {

		println(err.Error())
		return false
	}
	println("connect db success!")
	return true
}

func autoMigrateTables() {

	roleCfg := ConfigContent["db_init"].(map[string]interface{})["role"].(map[string]interface{})
	dptCfg := ConfigContent["db_init"].(map[string]interface{})["department"].(map[string]interface{})
	usrCfg := ConfigContent["db_init"].(map[string]interface{})["user"].(map[string]interface{})

	for k, v := range usrCfg {

		P(k, v)
	}
	// dbHost := dbCfg["host"].(string)
	// db, err := gorm.Open("mysql", "scott:753951@/dev_db?charset=utf8&parseTime=True&loc=Local")
	// Migrate the schema

	Db.AutoMigrate(Role{})
	var role Role

	if err := Db.Where(Role{RoleName: "admin"}).First(&role).Error; err != nil {
		fmt.Println("not found, and create", Role{})
		// db.Create(&model.Role{RoleName: "admin"})
		for k, v := range roleCfg {
			Db.Create(&Role{RoleName: k, Description: v.(string)})
		}
	}

	Db.AutoMigrate(&Department{})
	var dpt Department

	if err := Db.Where(&Department{DepartmentName: "programme"}).First(&dpt).Error; err != nil {
		fmt.Println("not found, and create", Department{})

		for k, v := range dptCfg {
			fmt.Println("not found, and create", Department{})
			Db.Create(&Department{DepartmentName: k, Description: v.(string)})
		}
	}

	Db.AutoMigrate(&User{})
	var user User

	name := usrCfg["name"].(string)
	password := usrCfg["password"].(string)
	email := usrCfg["email"].(string)
	mobile := usrCfg["mobile"].(string)
	sex := usrCfg["sex"].(string)
	age := uint(usrCfg["age"].(float64))
	staff := usrCfg["staff"].(string)

	if err := Db.Where(&User{Name: "admin"}).First(&user).Error; err != nil {

		P("not found , and create ", User{})
		Db.Create(&User{Name: name, Password: Encrypt(password), Email: email, Mobile: mobile, Sex: sex, Age: age, Staff: staff})
	}

}
