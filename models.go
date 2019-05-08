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

func readyDB() {

	e := godotenv.Load() //Load .env file
	if e != nil {
		fmt.Print(e)
	}

	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")

	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, username, dbName, password) //Build connection string
	fmt.Println(dbURI)

	if os.Getenv("DATABASE_URL") != "" {
		dbURI = os.Getenv("DATABASE_URL")
	}

	conn, err := gorm.Open("postgres", dbURI)
	if err != nil {
		fmt.Print(err)
	}

	db = conn
	db.Debug().AutoMigrate(&userAccount{}) //Database migration
}

//returns a handle to the DB object
func getDB() *gorm.DB {
	return db
}

type token struct {
	UserID uint
	jwt.StandardClaims
}

type userAccount struct {
	gorm.Model
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token" sql:"-"`
}

func (someUser *userAccount) Validate() (map[string]interface{}, bool) {
	if !strings.Contains(someUser.Email, "@") {
		return message(false, "Email is required"), false
	}

	temp := &userAccount{}

	err := getDB().Table("user_accounts").Where("email = ?", someUser.Email).First(temp).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return message(false, "Connection error with DB."), false
	}

	if temp.Email != "" {
		return message(false, "Email already used"), false
	}

	return message(false, "Validated"), true
}

func (someUser *userAccount) Create() map[string]interface{} {
	if resp, ok := someUser.Validate(); !ok {
		return resp
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(someUser.Password), bcrypt.DefaultCost)
	someUser.Password = string(hashedPassword)

	getDB().Create(someUser)

	if someUser.ID <= 0 {
		return message(false, "Failed to create, some error.")
	}

	tk := &token{UserID: someUser.ID}
	tk.ExpiresAt = time.Now().UTC().Add((600 * time.Second)).Unix()
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("TOKEN_PASSWORD")))
	someUser.Token = tokenString

	response := message(true, "Account created.")
	someUser.Password = ""
	response["userAccount"] = someUser
	return response
}

func login(email, password string) map[string]interface{} {
	someUser := &userAccount{}
	err := getDB().Table("user_accounts").Where("email = ?", email).First(someUser).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return message(false, "Email address not found")
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(someUser.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return message(false, "Invalid login creds.")
	}

	someUser.Password = ""

	tk := &token{UserID: someUser.ID}
	tk.ExpiresAt = time.Now().UTC().Add((600 * time.Second)).Unix()
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("TOKEN_PASSWORD")))
	someUser.Token = tokenString

	resp := message(true, "Logged In")
	resp["userAccount"] = someUser
	return resp
}

func getuserAccount(uid uint) *userAccount {
	uAcc := &userAccount{}
	getDB().Table("user_accounts").Where("id = ?", uid).First(uAcc)
	if uAcc.Email == "" {
		return nil
	}
	uAcc.Password = ""
	return uAcc
}

func prepareQuotes() *[]string {
	bytes, err := ioutil.ReadFile("einstein_quotes.txt")
	if err != nil {
		fmt.Printf("Error opening file: %s", err)
		return nil
	}
	s := string(bytes)
	q := strings.Split(s, "\n")
	return &q
}

func getRandomQuote() string {
	if quotes == nil {
		quotes = prepareQuotes()
	}

	rand.Seed(time.Now().UnixNano())
	num := rand.Intn(len(*quotes))
	s := (*quotes)[num] // stupid syntax, more info: https://flaviocopes.com/golang-does-not-support-indexing/
	return fmt.Sprintf("Einstein said: %s", s)
}
