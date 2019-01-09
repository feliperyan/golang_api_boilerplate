package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

var db *gorm.DB //database
var quotes *[]string

func ReadyDB() {

	e := godotenv.Load() //Load .env file
	if e != nil {
		fmt.Print(e)
	}

	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")

	dbUri := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, username, dbName, password) //Build connection string
	fmt.Println(dbUri)

	if os.Getenv("DATABASE_URL") != "" {
		dbUri = os.Getenv("DATABASE_URL")
	}

	conn, err := gorm.Open("postgres", dbUri)
	if err != nil {
		fmt.Print(err)
	}

	db = conn
	db.Debug().AutoMigrate(&UserAccount{}) //Database migration
}

//returns a handle to the DB object
func GetDB() *gorm.DB {
	return db
}

type Token struct {
	UserId uint
	jwt.StandardClaims
}

type UserAccount struct {
	gorm.Model
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token;sql:"-"`
}

func (userAccount *UserAccount) Validate() (map[string]interface{}, bool) {
	if !strings.Contains(userAccount.Email, "@") {
		return Message(false, "Email is required"), false
	}

	temp := &UserAccount{}

	err := GetDB().Table("user_accounts").Where("email = ?", userAccount.Email).First(temp).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return Message(false, "Connection error with DB."), false
	}

	if temp.Email != "" {
		return Message(false, "Email already used"), false
	}

	return Message(false, "Validated"), true
}

func (userAccount *UserAccount) Create() map[string]interface{} {
	if resp, ok := userAccount.Validate(); !ok {
		return resp
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(userAccount.Password), bcrypt.DefaultCost)
	userAccount.Password = string(hashedPassword)

	GetDB().Create(userAccount)

	if userAccount.ID <= 0 {
		return Message(false, "Failed to create, some error.")
	}

	tk := &Token{UserId: userAccount.ID}
	tk.ExpiresAt = time.Now().UTC().Add((600 * time.Second)).Unix()
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("TOKEN_PASSWORD")))
	userAccount.Token = tokenString

	response := Message(true, "Account created.")
	userAccount.Password = ""
	response["userAccount"] = userAccount
	return response
}

func Login(email, password string) map[string]interface{} {
	userAccount := &UserAccount{}
	err := GetDB().Table("user_accounts").Where("email = ?", email).First(userAccount).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return Message(false, "Email address not found")
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(userAccount.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return Message(false, "Invalid login creds.")
	}

	userAccount.Password = ""

	tk := &Token{UserId: userAccount.ID}
	tk.ExpiresAt = time.Now().UTC().Add((600 * time.Second)).Unix()
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	userAccount.Token = tokenString

	resp := Message(true, "Logged In")
	resp["userAccount"] = userAccount
	return resp
}

func GetUserAccount(uid uint) *UserAccount {
	uAcc := &UserAccount{}
	GetDB().Table("user_accounts").Where("id = ?", uid).First(uAcc)
	if uAcc.Email == "" {
		return nil
	}
	uAcc.Password = ""
	return uAcc
}

func PrepareQuotes() {
	bytes, err := ioutil.ReadFile("einstein_quotes.txt")
	if err != nil {
		fmt.Printf("Error opening file: %s", err)
		return
	} else {
		s := string(bytes)
		q := strings.Split(s, "\n")
		quotes = &q
	}
}

func GetRandomQuote() string {
	if quotes == nil {
		PrepareQuotes()
	}

	rand.Seed(time.Now().UnixNano())
	num := rand.Intn(len(*quotes))
	return (*quotes)[num] // stupid syntax, more info: https://flaviocopes.com/golang-does-not-support-indexing/
}
