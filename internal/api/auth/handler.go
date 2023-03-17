package auth

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/jmoiron/sqlx"
	"github.com/wisnuragaprawida/project/generated/users"
	"github.com/wisnuragaprawida/project/internal/api/response"
	"github.com/wisnuragaprawida/project/pkg/crashy"
	"github.com/wisnuragaprawida/project/pkg/one-go/utils"
)

type AuthHandler struct {
	db       *sqlx.DB
	userRepo *users.Queries
}

func NewAuthHandler(db *sqlx.DB) *AuthHandler {
	return &AuthHandler{
		db:       db,
		userRepo: users.New(db),
	}
}

func (ah *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	//register user with bcrypt as hash pasword

	var (
		req RegisterRequest
	)

	if err := render.Bind(r, &req); err != nil {

		response.Nay(w, r, crashy.Wrapc(err, crashy.ErrCode(err.Error())), http.StatusBadRequest)
		// response.Nay(w, r, crashy.Wrapc(err, crashy.ErrParsingData), http.StatusBadRequest)
		return
	}

	//check if email already exist
	_, err := ah.userRepo.FindUserByEmail(r.Context(), req.Email)
	if err == nil {
		response.Nay(w, r, crashy.Wrapc(err, crashy.ErrCode("Email Alredy Exist!")), http.StatusBadRequest)
		return
	}

	//hash pasword
	PasswordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		response.Nay(w, r, crashy.Wrapc(err, crashy.ErrCode("Hash Password Failed")), http.StatusBadRequest)
		return
	}

	//register user use bycrypt to encrypt password
	_, err = ah.userRepo.RegisterUser(r.Context(), users.RegisterUserParams{
		Name:     req.Name,
		Email:    req.Email,
		Password: PasswordHash,
	})

	if err != nil {
		response.Nay(w, r, crashy.Wrapc(err, crashy.ErrCode("Register Failed")), http.StatusBadRequest)
		return
	}

	response.Yay(w, r, "Register Success", http.StatusOK)
}

func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	//login user with bcrypt as hash pasword

	var (
		req LoginRequest
	)

	if err := render.Bind(r, &req); err != nil {
		response.Nay(w, r, crashy.Wrapc(err, crashy.ErrCode(err.Error())), http.StatusBadRequest)
		return
	}

	//check if email already exist
	user, err := ah.userRepo.FindUserByEmail(r.Context(), req.Email)
	if err != nil {
		response.Nay(w, r, crashy.Wrapc(err, crashy.ErrCode("Invalid Credential")), http.StatusBadRequest)
		return
	}

	//check if password match
	ok := utils.CheckPasswordHash(req.Password, user.Password)
	if !ok {
		response.Nay(w, r, crashy.Wrapc(err, crashy.ErrCode("Invalid Credential")), http.StatusBadRequest)
		return
	}

	response.Yay(w, r, "Login Success", http.StatusOK)
}
