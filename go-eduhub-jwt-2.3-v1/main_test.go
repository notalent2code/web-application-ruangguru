package main_test

import (
	main "a21hc3NpZ25tZW50"
	"a21hc3NpZ25tZW50/db"
	"a21hc3NpZ25tZW50/middleware"
	"a21hc3NpZ25tZW50/model"
	repo "a21hc3NpZ25tZW50/repository"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func SetCookie(mux *gin.Engine) *http.Cookie {
	login := model.UserLogin{
		Email:    "test@mail.com",
		Password: "testing123",
	}

	body, _ := json.Marshal(login)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/user/login", bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	mux.ServeHTTP(w, r)

	var cookie *http.Cookie
	for _, c := range w.Result().Cookies() {
		if c.Name == "session_token" {
			cookie = c
		}
	}

	return cookie
}

var _ = Describe("Go EduHub JWT", Ordered, func() {
	var apiServer *gin.Engine
	var userRepo repo.UserRepository
	var courseRepo repo.CourseRepository

	db := db.NewDB()
	dbCredential := model.Credential{
		Host:         "localhost",
		Username:     "postgres",
		Password:     "kambing",
		DatabaseName: "be4_ruangguru",
		Port:         5432,
		Schema:       "public",
	}

	conn, err := db.Connect(&dbCredential)
	Expect(err).ShouldNot(HaveOccurred())

	userRepo = repo.NewUserRepo(conn)
	courseRepo = repo.NewCourseRepo(conn)

	BeforeAll(func() {
		gin.SetMode(gin.ReleaseMode) //release

		conn.AutoMigrate(&model.User{}, &model.Course{})

		apiServer = gin.New()
		apiServer = main.RunServer(conn, apiServer)

		reqRegister := model.UserRegister{
			Fullname: "test",
			Email:    "test@mail.com",
			Password: "testing123",
		}

		reqBody, _ := json.Marshal(reqRegister)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/user/register", bytes.NewReader(reqBody))
		r.Header.Set("Content-Type", "application/json")
		apiServer.ServeHTTP(w, r)

		Expect(w.Result().StatusCode).To(Equal(http.StatusCreated))
	})

	AfterAll(func() {
		err = db.Reset(conn, "users")
		err = db.Reset(conn, "courses")
		Expect(err).ShouldNot(HaveOccurred())
	})

	Describe("Auth Middleware", func() {
		var (
			router *gin.Engine
			w      *httptest.ResponseRecorder
		)

		BeforeEach(func() {
			router = gin.Default()
			w = httptest.NewRecorder()
		})

		When("valid token is provided", func() {
			It("should set user ID in context and call next middleware", func() {
				claims := &model.Claims{UserID: 123}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				signedToken, _ := token.SignedString(model.JwtKey)
				req, _ := http.NewRequest(http.MethodGet, "/", nil)
				req.AddCookie(&http.Cookie{Name: "session_token", Value: signedToken})

				router.Use(middleware.Auth())
				router.GET("/", func(ctx *gin.Context) {
					userID := ctx.MustGet("id").(int)
					Expect(userID).To(Equal(123))
				})

				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusOK))
			})
		})

		When("session token is missing", func() {
			It("should return unauthorized error response", func() {
				req, _ := http.NewRequest(http.MethodGet, "/", nil)

				router.Use(middleware.Auth())

				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusSeeOther))
			})
		})

		When("session token is invalid", func() {
			It("should return unauthorized error response", func() {
				req, _ := http.NewRequest(http.MethodGet, "/", nil)
				req.AddCookie(&http.Cookie{Name: "session_token", Value: "invalid_token"})

				router.Use(middleware.Auth())

				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
		})
	})

	Describe("Repository", func() {
		Describe("User repository", func() {
			When("fetching a single user data by email from users table in the database", func() {
				It("should return a single student data", func() {
					expectUser := model.User{
						Fullname: "test",
						Email:    "test@mail.com",
						Password: "testing123",
					}

					resUser, err := userRepo.GetUserByEmail("test@mail.com")
					Expect(err).ShouldNot(HaveOccurred())
					Expect(resUser.Fullname).To(Equal(expectUser.Fullname))
					Expect(resUser.Email).To(Equal(expectUser.Email))
					Expect(resUser.Password).To(Equal(expectUser.Password))
				})
			})
		})

		Describe("Course repository", func() {
			When("add course data to course table database postgres", func() {
				It("should save course data to courses table database postgres", func() {
					course := model.Course{
						Name:       "Introduction to Computer Science",
						Schedule:   "Monday and Wednesday, 10am - 12pm",
						Grade:      3.5,
						Attendance: 85,
					}
					err := courseRepo.Store(&course)
					Expect(err).ShouldNot(HaveOccurred())

					result, err := courseRepo.FetchByID(1)
					Expect(err).ShouldNot(HaveOccurred())

					Expect(result.Name).To(Equal(course.Name))
					Expect(result.Schedule).To(Equal(course.Schedule))
					Expect(result.Grade).To(Equal(course.Grade))
					Expect(result.Attendance).To(Equal(course.Attendance))

					err = db.Reset(conn, "courses")
					Expect(err).ShouldNot(HaveOccurred())
				})
			})
		})
	})

	Describe("API", func() {
		Describe("/users/login", func() {
			When("send empty email and password with POST method", func() {
				It("should return a bad request", func() {
					loginData := model.UserLogin{
						Email:    "",
						Password: "",
					}

					body, _ := json.Marshal(loginData)
					w := httptest.NewRecorder()
					r := httptest.NewRequest("POST", "/user/login", bytes.NewReader(body))
					r.Header.Set("Content-Type", "application/json")

					apiServer.ServeHTTP(w, r)

					errResp := model.ErrorResponse{}
					err := json.Unmarshal(w.Body.Bytes(), &errResp)
					Expect(err).To(BeNil())
					Expect(w.Result().StatusCode).To(Equal(http.StatusBadRequest))
					Expect(errResp.Error).To(Equal("invalid decode json"))
				})
			})

			When("send email and password with POST method", func() {
				It("should return a success", func() {
					loginData := model.UserLogin{
						Email:    "test@mail.com",
						Password: "testing123",
					}
					body, _ := json.Marshal(loginData)
					w := httptest.NewRecorder()
					r := httptest.NewRequest("POST", "/user/login", bytes.NewReader(body))
					r.Header.Set("Content-Type", "application/json")
					apiServer.ServeHTTP(w, r)

					var resp = map[string]interface{}{}
					err = json.Unmarshal(w.Body.Bytes(), &resp)
					Expect(err).To(BeNil())
					Expect(w.Result().StatusCode).To(Equal(http.StatusOK))
					Expect(resp["message"]).To(Equal("login success"))
				})
			})
		})

		Describe("/course/delete/:id", func() {
			When("sending without cookie", func() {
				It("should return status code 200", func() {
					course := model.Course{
						Name:       "Data Structures and Algorithms",
						Schedule:   "Monday, Wednesday, and Friday, 9am - 11am",
						Grade:      3.2,
						Attendance: 80,
					}
					err := courseRepo.Store(&course)
					Expect(err).ShouldNot(HaveOccurred())

					w := httptest.NewRecorder()
					r := httptest.NewRequest("DELETE", "/course/delete/1", nil)
					apiServer.ServeHTTP(w, r)

					var response model.SuccessResponse
					json.Unmarshal(w.Body.Bytes(), &response)

					Expect(w.Code).To(Equal(http.StatusOK))
					Expect(response.Message).To(Equal("course delete success"))
				})
			})

			When("deleting existing course", func() {
				It("should return status code 200", func() {
					course := model.Course{
						Name:       "Data Structures and Algorithms",
						Schedule:   "Monday, Wednesday, and Friday, 9am - 11am",
						Grade:      3.2,
						Attendance: 80,
					}
					err := courseRepo.Store(&course)
					Expect(err).ShouldNot(HaveOccurred())

					w := httptest.NewRecorder()
					r := httptest.NewRequest("DELETE", "/course/delete/1", nil)
					r.AddCookie(SetCookie(apiServer))
					apiServer.ServeHTTP(w, r)

					var response model.SuccessResponse
					json.Unmarshal(w.Body.Bytes(), &response)

					Expect(w.Code).To(Equal(http.StatusOK))
					Expect(response.Message).To(Equal("course delete success"))
				})
			})

			When("deleting course with invalid ID", func() {
				It("should return status code 400", func() {
					w := httptest.NewRecorder()
					r := httptest.NewRequest("DELETE", "/course/delete/invalid", nil)
					r.AddCookie(SetCookie(apiServer))
					apiServer.ServeHTTP(w, r)

					var response model.ErrorResponse
					json.Unmarshal(w.Body.Bytes(), &response)

					Expect(w.Code).To(Equal(http.StatusBadRequest))
					Expect(response.Error).NotTo(BeNil())
				})
			})
		})
	})
})
