package config

import (
	"encoding/json"
	"fmt"
	"jwt/pkg/models"
	"jwt/pkg/response"
	"net/http"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"

	"golang.org/x/crypto/bcrypt"
	// For mysql framework
	_ "github.com/go-sql-driver/mysql"
)

//Login Controller
func Login(w http.ResponseWriter, r *http.Request) {
	if storedCookie, _ := r.Cookie("token"); storedCookie != nil {
		response.JSON(w, http.StatusBadRequest, "You've logged in !")
		return
	}
	db := connect()
	defer db.Close()
	var user models.User
	var pass string
	var wrong bool = false
	var remspace = func(s string) string {
		return strings.TrimSpace(s)
	}
	json.NewDecoder(r.Body).Decode(&user)
	if remspace(user.Email) == "" {
		wrong = true
	} else if remspace(user.Password) == "" {
		wrong = true
	}
	if wrong {
		response.JSON(w, http.StatusUnauthorized, "Please fill the data")
		return
	}
	if err := db.QueryRow("select password from user where email = ?", user.Email).Scan(&pass); err != nil {
		response.JSON(w, http.StatusUnauthorized, "Wrong email or password")
		return
	}
	hashedFromDB := pass
	if err := bcrypt.CompareHashAndPassword([]byte(hashedFromDB), []byte(user.Password)); err != nil {
		response.JSON(w, http.StatusUnauthorized, "Wrong email or password !")
		return
	}
	if err := db.QueryRow("select * from user where email = ?", user.Email).Scan(&user.ID, &user.Name, &user.Level, &user.Email, &user.Password); err != nil {
		fmt.Println(err)
		return
	}
	token, err := CreateToken(uint32(user.ID), user.Level)
	if err != nil {
		fmt.Println(err)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:  "token",
		Value: token,
	})
	response.Token(w, http.StatusOK, token, "Successfully got your token")
}

//CreateUser Contoller
func CreateUser(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)
	var wrong bool = false
	var remspace = func(s string) string {
		return strings.TrimSpace(s)
	}
	if remspace(user.Name) == "" {
		wrong = true
	} else if remspace(user.Password) == "" {
		wrong = true
	} else if remspace(user.Email) == "" {
		wrong = true
	} else if remspace(user.Level) == "" {
		wrong = true
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
	}
	_, errs := db.Exec("insert into user (nama,level,email,password) values (?,?,?,?)",
		user.Name,
		user.Level,
		user.Email,
		hashed,
	)
	if errs != nil {
		fmt.Println(errs)
		response.JSON(w, http.StatusInternalServerError, "Something went wrong")
		return
	}
	if wrong {
		response.JSON(w, http.StatusBadRequest, "Please fill the data correctly !")
		return
	}
	fmt.Println("Success inserted to database")
}

//Logout Function
func Logout(w http.ResponseWriter, r *http.Request) {
	if _, err := r.Cookie("token"); err == http.ErrNoCookie {
		response.JSON(w, http.StatusBadRequest, "No one has logged in")
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:   "token",
		MaxAge: -1,
	})
	response.JSON(w, http.StatusOK, "Successfully logged out")
}

//Refresh Function for refresh the token
func Refresh(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()
	var status string
	cookie, err := r.Cookie("token")
	if err == http.ErrNoCookie {
		response.JSON(w, http.StatusBadRequest, "Can't refresh the token. Please login first !")
		return
	}
	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		return []byte("MunzirDEV"), nil
	})
	if err != nil {
		fmt.Println(err)
		response.JSON(w, http.StatusInternalServerError, "Something Error")
		return
	}
	claims, _ := token.Claims.(jwt.MapClaims)
	if errs := db.QueryRow("select level from user where id = ?", claims["user_id"]).Scan(&status); errs != nil {
		fmt.Println(errs)
		response.JSON(w, http.StatusInternalServerError, "Something error")
		return
	}
	userID, err := strconv.ParseUint(fmt.Sprint(claims["user_id"]), 10, 32)
	if err != nil {
		fmt.Println(err)
		response.JSON(w, http.StatusInternalServerError, "Something error")
		return
	}
	tokenResponse, err := CreateToken(uint32(userID), status)
	if err != nil {
		fmt.Println(err)
		response.JSON(w, http.StatusInternalServerError, "Something error")
		return
	}
	response.Token(w, http.StatusOK, tokenResponse, "Successfully refreshed your token")
}

//Test Function for check the token
func Test(w http.ResponseWriter, r *http.Request) {
	if err := CookieValid(r); err != nil {
		response.JSON(w, http.StatusUnauthorized, fmt.Sprint(err))
	}
	response.JSON(w, http.StatusOK, "DONE !")
}
