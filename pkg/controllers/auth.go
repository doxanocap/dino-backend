package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/doxanocap/reactNative/dino-back/pkg/database"
	"github.com/doxanocap/reactNative/dino-back/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(ctx *gin.Context) {
	var data map[string]string
	if err := ctx.BindJSON(&data); err != nil {
		panic(err)
	}
	user := &models.User{Username: data["username"], Email: data["email"], Password: nil}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)
	if user.Username == "" || user.Email == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid input"})
		return
	}
	fmt.Println(data)
	res, err := database.DB.Query(fmt.Sprintf("INSERT INTO dinoUsers (email, username, password) VALUES('%s','%s','%s')", user.Email, user.Username, password))
	if err != nil {
		panic(err)
	}
	defer res.Close()
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func SignIn(ctx *gin.Context) {
	var data map[string]string

	if err := ctx.BindJSON(&data); err != nil {
		panic(err)
	}
	res, err := database.DB.Query(fmt.Sprintf("SELECT * FROM dinoUsers WHERE email = '%s'", data["email"]))
	if err != nil {
		panic(err)
	}
	defer res.Close()

	var newUser models.User
	for res.Next() {
		err = res.Scan(&newUser.Id, &newUser.Email, &newUser.Username, &newUser.Password)
		if err != nil {
			panic(err)
		}
		if err := bcrypt.CompareHashAndPassword(newUser.Password, []byte(data["password"])); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "incorrect password"})
			return
		}
		break
	}
	if newUser.Id == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "unsuccessful"})
		return
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
		Issuer:    strconv.Itoa(int(newUser.Id)),
		ExpiresAt: &jwt.NumericDate{time.Now().Add(time.Hour * 24)},
	})

	token, err := claims.SignedString([]byte("SecretKey"))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "error with creation of the token"})
		return
	}

	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:    "jwt4",
		Value:   token,
		Expires: time.Now().Add(time.Hour * 24),
	})

	ctx.JSON(http.StatusOK, gin.H{"jwt": token, "userInfo": newUser})

}

func User(ctx *gin.Context) {
	cookie, _ := ctx.Cookie("jwt4")
	token, err := jwt.ParseWithClaims(cookie, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("SecretKey"), nil
	})

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "unauthenticated"})
		return
	}

	claims := token.Claims.(*jwt.RegisteredClaims)

	res, err := database.DB.Query(fmt.Sprintf("SELECT * FROM dinoUsers WHERE id = '%s'", claims.Issuer))

	if err != nil {
		panic(err)
	}

	defer res.Close()

	var user models.User
	for res.Next() {
		err = res.Scan(&user.Id, &user.Email, &user.Username, &user.Password)
		if err != nil {
			panic(err)
		}
		break
	}
	ctx.JSON(http.StatusOK, gin.H{"jwt": token.Raw, "userInfo": user})
}

func SignOut(ctx *gin.Context) {
	cookie := &http.Cookie{
		Name:     "jwt4",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	}
	http.SetCookie(ctx.Writer, cookie)
	ctx.JSON(http.StatusOK, gin.H{"message": "deleted cookie"})
}
