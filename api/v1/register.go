package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/practice2311/api/models"
	"github.com/practice2311/pkg/email"
	"github.com/practice2311/pkg/utils"
	"github.com/practice2311/storage/repo"
)

var (
	ErrWrongEmailOrPass = errors.New("wrong email or password")
	ErrUserNotVerified  = errors.New("user not verified")
)

func checkPassword(password string) bool {
	password = strings.Replace(password, " ", "", -1)
	password = strings.Replace(password, "\n", "", -1)
	password = strings.Replace(password, "\t", "", -1)
	fmt.Println(password)
	if utf8.RuneCountInString(password) < 8 {
		return false
	}
	number, sign, U_letter, L_letter := false, false, false, false
	for _, i := range password {
		if unicode.IsUpper(i) && !U_letter {
			U_letter = true
		} else if unicode.IsLower(i) && !L_letter {
			L_letter = true
		} else if unicode.IsNumber(i) && !number {
			number = true
		} else if (0 < int(i) && int(i) < 48) || (57 < int(i) && int(i) < 65) || (90 < int(i) && int(i) < 97) || (122 < int(i)) {
			sign = true
		}
	}
	return number == sign && sign == U_letter && U_letter == L_letter
}

// @Router /users/auth [post]
// @Summary Create a user
// @Description Create a user
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.CreateUser true "User"
// @Success 201 {object} models.User
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) Register(c *gin.Context) {
	var (
		req models.CreateUser
	)

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if !checkPassword(req.Password) {
		c.JSON(http.StatusInternalServerError, errorResponse(errors.New("parol kamida 8 ta belgi,katta-kichik harf,belgi va raqamlardan iborat bo'lishi kerak")))
		return
	}

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	verificationCode := r1.Intn(9000) + 1000
	go func() {
		err = email.SendEmail(h.cfg, &email.SendEmailRequest{
			To:      []string{req.Email},
			Subject: "Verification email",
			Body: map[string]string{
				"code": strconv.Itoa(verificationCode),
			},
			Type: email.VerificationEmail,
		})
		if err != nil {
			fmt.Println("Failed to send email")
		}
	}()

	c.JSON(http.StatusCreated, models.ResponseOK{
		Message: "Verification code has been sent!",
	})

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user := repo.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
	}

	JsonUsr, err := json.Marshal(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	reqCode := models.VerifyUser{
		Email:            req.Email,
		VerificationCode: verificationCode,
	}

	verJson, err := json.Marshal(reqCode)

	err = h.inMemory.SetWithTTL("user:"+req.Email, string(JsonUsr), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err = h.inMemory.SetWithTTL("code:"+req.Email, string(verJson), 1)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

}

// @Router /auth/verify [post]
// @Summary Verify user
// @Description Verify user
// @Tags auth
// @Accept json
// @Produce json
// @Param verify body models.VerifyUser true "Verify"
// @Success 200 {object} models.User
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) Verify(c *gin.Context) {
	var (
		req models.VerifyUser
	)

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userData, err := h.inMemory.Get("user:" + req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	verifyData, err := h.inMemory.Get("code:" + req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var usr models.User
	err = json.Unmarshal([]byte(userData), &usr)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var verf models.VerifyUser
	err = json.Unmarshal([]byte(verifyData), &verf)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if verf.VerificationCode == req.VerificationCode {
		result, err := h.storage.User().Create(&repo.User{
			Name:     usr.Name,
			Email:    usr.Email,
			Password: usr.Password,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		token, _, err := utils.CreateToken(utils.TokenParam{
			Id:       result.Id,
			Name:     result.Name,
			Email:    result.Email,
			Duration: time.Hour * 24,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		c.JSON(http.StatusCreated, models.User{
			Id:          result.Id,
			Name:        result.Name,
			Email:       result.Email,
			Password:    result.Password,
			AccessToken: token,
		})
	} else {
		c.JSON(http.StatusBadRequest, errorResponse(errors.New("ey botanik ko'zingni katta och bunday kod yubormagandim")))
		return
	}

}
