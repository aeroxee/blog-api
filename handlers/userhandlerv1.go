package handlers

import (
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Aeroxee/blog-api/auth"
	"github.com/Aeroxee/blog-api/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

type UserHandlerV1 struct{}

func NewUserHandlerV1() UserHandlerV1 {
	return UserHandlerV1{}
}

type ActivationData struct {
	Email          string
	ExpirationTime time.Time
}

var activationData = make(map[string]ActivationData)

// generateActivationCode is function to get activation code.
func (UserHandlerV1) generateActivationCode() string {
	return uuid.NewString()
}

// send activation code to email target
func (UserHandlerV1) sendActivationEmail(email, activationCode string) bool {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	smtpServer := os.Getenv("SMTP_SERVER")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUsername := os.Getenv("SMTP_USERNAME")
	smtpPassword := os.Getenv("SMTP_PASSWORD")
	HOSTNAME := os.Getenv("HOSTNAME")
	PORT := os.Getenv("PORT")

	to := []string{email}
	subject := "Aktifasi Akun Kamu"
	body := fmt.Sprintf("Klik link ini untuk mengaktifkan akun kamu: http://%s:%s/v1/activate/%s", HOSTNAME, PORT, activationCode)
	message := fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", strings.Join(to, ","), subject, body)

	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpServer)
	err = smtp.SendMail(smtpServer+":"+smtpPort, auth, smtpUsername, to, []byte(message))
	if err == nil {
		return true
	} else {
		return false
	}
}

func (u *UserHandlerV1) Register() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload RegisterPayload
		err := ctx.ShouldBindJSON(&payload)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Payload yang anda gunakan bukan bertipe json. Harap kirim data menggunakan json saja.",
			})
			return
		}

		validate = validator.New(validator.WithRequiredStructEnabled())
		err = validate.Struct(payload)
		if err != nil {
			if _, ok := err.(*validator.InvalidValidationError); ok {
				fmt.Println(err.Error())
				return
			}

			var validations []Validation
			for _, err := range err.(validator.ValidationErrors) {
				validations = append(validations, Validation{
					Field: err.Field(),
					Tag:   err.ActualTag(),
				})
			}

			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Validasi error",
				"errors":  validations,
			})
			return
		}

		user := models.User{
			FirstName: payload.FirstName,
			LastName:  payload.LastName,
			Email:     payload.Email,
			Username:  payload.Username,
			Password:  auth.EncryptionPassword(payload.Password),
		}

		userModel := models.NewUserModel(models.GetDB())
		err = userModel.CreateUser(&user)
		if err != nil {
			if strings.Contains(err.Error(), "users.username") {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"status":  "error",
					"message": fmt.Sprintf("Username `%s` sudah terdaftar, silahkan gunakan username yang lain.", payload.Username),
				})
				return
			} else if strings.Contains(err.Error(), "users.email") {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"status":  "error",
					"message": fmt.Sprintf("Alamat email `%s` sudah terdaftar, silahkan gunakan alamat email yang lain.", payload.Email),
				})
				return
			} else {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"status":  "error",
					"message": err.Error(),
				})
				return
			}
		}

		// generate new activation code
		activationCode := u.generateActivationCode()
		activationData[activationCode] = ActivationData{
			Email:          user.Email,
			ExpirationTime: time.Now().Add(24 * time.Hour),
		}

		// send email verivication
		if !u.sendActivationEmail(activationData[activationCode].Email, activationCode) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Gagal mengirimkan kode aktifasi ke alamat email anda. Harap periksa alamat email anda aktif atau tidak aktif.",
			})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"status":  "success",
			"message": fmt.Sprintf("Berhasil mendaftarkan akun dengan alaman email: %s. Silahkan cek alamat email anda untuk aktifasi akun.", payload.Email),
			"user":    user,
		})
	}
}

// ActivationHandler is function to handle activation code to user.
func (u *UserHandlerV1) ActivationHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		activationCode := ctx.Param("activationCode")

		data, ok := activationData[activationCode]
		if !ok || time.Now().After(data.ExpirationTime) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Kode aktifasi sudah kadaluarsa atau tidak valid.",
			})
			return
		}

		userModel := models.NewUserModel(models.GetDB())
		user, err := userModel.GetUserEmail(data.Email)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Kode aktifasi yang anda kirimkan sudah tidak valid.",
			})
			return
		}

		user.IsActive = true
		models.GetDB().Save(&user)
		delete(activationData, activationCode)

		ctx.Writer.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(ctx.Writer, "<h1>Your account is activated successfully.</h1>")
	}
}

func (UserHandlerV1) GetToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload LoginPayload
		err := ctx.ShouldBindJSON(&payload)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Payload yang anda gunakan bukan bertipe json. Harap kirim data menggunakan json saja.",
			})
			return
		}

		userModel := models.NewUserModel(models.GetDB())
		var user models.User
		if strings.Contains(payload.Username, "@") {
			user, err = userModel.GetUserEmail(payload.Username)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"status":  "error",
					"message": "Email atau katasandi yang anda masukkan salah.",
				})
				return
			}
		} else {
			user, err = userModel.GetUserUsername(payload.Username)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"status":  "error",
					"message": "Email atau katasandi yang anda masukkan salah.",
				})
				return
			}
		}

		// check password
		if !auth.DecryptionPassword(user.Password, payload.Password) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Email atau katasandi yang anda masukkan salah.",
			})
			return
		}

		// check user is not active
		if !user.IsActive {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Akun anda belum aktif, silahkan periksa alamat email anda untuk cek aktifasi kode.",
			})
			return
		}

		credential := auth.Credential{
			UserID: user.ID,
		}

		token, err := auth.GetToken(credential)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Akses token berhasil didapatkan.",
			"token":   token,
		})
	}
}

func (UserHandlerV1) Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, err := getUserContext(ctx.Request)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Autentikasi dibutuhkan.",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Akun anda telah terautentikasi.",
			"user":    user,
		})
	}
}

func (UserHandlerV1) CheckEmail() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		payload := struct {
			Email string `json:"email"`
		}{}
		err := ctx.ShouldBindJSON(&payload)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Payload yang anda gunakan bukan bertipe json. Harap kirim data menggunakan json saja.",
			})
			return
		}

		userModel := models.NewUserModel(models.GetDB())
		_, err = userModel.GetUserEmail(payload.Email)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": fmt.Sprintf("Alamat email `%s` tidak ditemukan.", payload.Email),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": fmt.Sprintf("Alamat email `%s` ditemukan.", payload.Email),
		})
	}
}

func (UserHandlerV1) CheckUsername() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		payload := struct {
			Username string `json:"username"`
		}{}
		err := ctx.ShouldBindJSON(&payload)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Payload yang anda gunakan bukan bertipe json. Harap kirim data menggunakan json saja.",
			})
			return
		}

		userModel := models.NewUserModel(models.GetDB())
		_, err = userModel.GetUserUsername(payload.Username)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": fmt.Sprintf("User dengan username `%s` tidak ditemukan.", payload.Username),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": fmt.Sprintf("User dengan username `%s` ditemukan.", payload.Username),
		})
	}
}

func (UserHandlerV1) GetUserFromUsername() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		username := ctx.Param("username")

		userModel := models.NewUserModel(models.GetDB())
		user, err := userModel.GetUserUsername(username)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Pengguna dengan username: " + username + " tidak ditemukan.",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "",
			"user":    user,
		})
	}
}

func (UserHandlerV1) GetUserFromID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.Param("userId")
		userIdInt, err := strconv.Atoi(userId)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Pengguna dengan id: " + userId + " tidak ditemukan.",
			})
			return
		}

		userModel := models.NewUserModel(models.GetDB())
		user, err := userModel.GetUserByID(userIdInt)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Pengguna dengan id: " + userId + " tidak ditemukan.",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "",
			"user":    user,
		})
	}
}
