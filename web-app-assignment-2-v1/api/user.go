package api

import (
	"a21hc3NpZ25tZW50/model"
	"a21hc3NpZ25tZW50/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type UserAPI interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	GetUserTaskCategory(c *gin.Context)
}

type userAPI struct {
	userService service.UserService
}

func NewUserAPI(userService service.UserService) *userAPI {
	return &userAPI{userService}
}

func (u *userAPI) Register(c *gin.Context) {
	var user model.UserRegister

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("invalid decode json"))
		return
	}

	if user.Email == "" || user.Password == "" || user.Fullname == "" {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("register data is empty"))
		return
	}

	var recordUser = model.User{
		Fullname: user.Fullname,
		Email:    user.Email,
		Password: user.Password,
	}

	recordUser, err := u.userService.Register(&recordUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse("error internal server"))
		return
	}

	c.JSON(http.StatusCreated, model.NewSuccessResponse("register success"))
}

func (u *userAPI) Login(c *gin.Context) {
	var user model.UserLogin

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("invalid decode json"))
		return
	}

	if user.Email == "" || user.Password == "" {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse("email or password is empty"))
		return
	}

	var recordUser = model.User{
		Email:    user.Email,
		Password: user.Password,
	}

	sessionToken, err := u.userService.Login(&recordUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse("error internal server"))
		return
	}

	tokenString := *sessionToken
	claims := &model.Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return model.JwtKey, nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse("error internal server"))
		return
	}

	if !token.Valid {
		c.JSON(http.StatusUnauthorized, model.NewErrorResponse("invalid token"))
		return
	}

	cookie := http.Cookie{
		Name:    "session_token",
		Value:   tokenString,
		Expires: time.Unix(claims.ExpiresAt, 0),
	}

	c.Writer.Header().Add("Set-Cookie", cookie.String())
	c.SetCookie("session_token", tokenString, int(claims.ExpiresAt), "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{
		"user_id": claims.UserID,
		"message": "login success",
	})
}

func (u *userAPI) GetUserTaskCategory(c *gin.Context) {

	userTaskCategory, err := u.userService.GetUserTaskCategory()
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse("error internal server"))
		return
	}

	c.JSON(http.StatusOK, userTaskCategory)
}
