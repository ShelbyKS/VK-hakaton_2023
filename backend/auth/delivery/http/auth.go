package http

import (
	"encoding/json"
	"io"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"github.com/gorilla/mux"

	"hakaton/domain"
	logs "hakaton/logger"
)

type AuthHandler struct {
	AuthUsecase domain.AuthUsecase
}

func NewAuthHandler(authMwRouter *mux.Router, mainRouter *mux.Router, u domain.AuthUsecase) {
	handler := &AuthHandler{
		AuthUsecase: u,
	}

	mainRouter.HandleFunc("/api/v1/auth/login", handler.Login).Methods(http.MethodPost, http.MethodOptions)
	mainRouter.HandleFunc("/api/v1/auth/register", handler.Register).Methods(http.MethodPost, http.MethodOptions)

	authMwRouter.HandleFunc("/v1/auth/check", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNoContent) }).Methods(http.MethodPost, http.MethodOptions)
	authMwRouter.HandleFunc("/v1/auth/logout", handler.Logout).Methods(http.MethodPost, http.MethodOptions)
}

// Login godoc
// @Summary      login user
// @Description  create user session and put it into cookie
// @Tags         auth
// @Accept       json
// @Param 		 body body domain.Credentials true "user credentials"
// @Success      200  {object} object{body=object{id=int}}
// @Failure      400  {object} object{err=string}
// @Failure      403  {object} object{err=string}
// @Failure      404  {object} object{err=string}
// @Failure      500  {object} object{err=string}
// @Router       /api/v1/auth/login [post]
func (a *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	auth, err := a.auth(r)
	if auth == true {
		domain.WriteError(w, "you must be unauthorised", http.StatusForbidden)
		logs.LogError(logs.Logger, "auth_http", "Login", err, "User is already logged in")
		return
	}

	var credentials domain.Credentials

	err = json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		domain.WriteError(w, err.Error(), http.StatusBadRequest)
		logs.LogError(logs.Logger, "auth_http", "Login", err, "Failed to decode json from body")
		return
	}
	logs.Logger.Debug("Login credentials:", credentials)
	defer a.CloseAndAlert(r.Body)

	if err = checkCredentials(credentials); err != nil {
		domain.WriteError(w, err.Error(), domain.GetStatusCode(err))
		logs.LogError(logs.Logger, "auth_http", "Login", err, "Credentials are incorrect")
		return
	}
	credentials.Email = strings.TrimSpace(credentials.Email)

	session, userID, err := a.AuthUsecase.Login(credentials)
	if err != nil {
		domain.WriteError(w, err.Error(), domain.GetStatusCode(err))
		logs.LogError(logs.Logger, "auth_http", "Login", err, "Failed to login")
		return
	}
	logs.Logger.Debug("Login: session:", session)

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    session.Token,
		Expires:  session.ExpiresAt,
		Path:     "/",
		HttpOnly: true,
	})

	domain.WriteResponse(
		w,
		map[string]interface{}{
			"id": userID,
		},
		http.StatusOK,
	)
}

// Logout godoc
// @Summary      logout user
// @Description  delete current session and nullify cookie
// @Tags         auth
// @Success      204
// @Failure      400  {object} object{err=string}
// @Failure      403  {object} object{err=string}
// @Failure      404  {object} object{err=string}
// @Failure      500  {object} object{err=string}
// @Router       /api/v1/auth/logout [post]
func (a *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	sessionToken := c.Value
	logs.Logger.Debug("Logout: session token:", c)

	if err = a.AuthUsecase.Logout(sessionToken); err != nil {
		domain.WriteError(w, err.Error(), http.StatusInternalServerError)
		logs.LogError(logs.Logger, "auth_http", "Logout", err, "Failed to logout")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now(),
		Path:     "/",
		HttpOnly: true,
	})

	w.WriteHeader(http.StatusNoContent)
}

// Register godoc
// @Summary      register user
// @Description  add new user to db and return it id
// @Tags         auth
// @Produce      json
// @Accept       json
// @Param 		 body body domain.Credentials true "user credentials"
// @Success      200  {object} object{body=object{id=int}}
// @Failure      400  {object} object{err=string}
// @Failure      403  {object} object{err=string}
// @Failure      500  {object} object{err=string}
// @Router       /api/v1/auth/register [post]
func (a *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	auth, err := a.auth(r)
	if auth == true {
		domain.WriteError(w, "you must be unauthorised", http.StatusForbidden)
		logs.LogError(logs.Logger, "auth_http", "Register.auth", err, "user is authorised")
		return
	}

	var user domain.User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		domain.WriteError(w, "you must be unauthorised", http.StatusBadRequest)
		logs.LogError(logs.Logger, "auth_http", "Register.decode", err, "Failed to decode json from body")
		return
	}
	//logs.Logger.Debug("Register user:", user)
	defer a.CloseAndAlert(r.Body)

	user.Email = strings.TrimSpace(user.Email)
	if err = checkCredentials(domain.Credentials{Email: user.Email, Password: user.Password}); err != nil {
		domain.WriteError(w, err.Error(), domain.GetStatusCode(err))
		logs.LogError(logs.Logger, "auth_http", "Register.credentials", err, "creds are invalid")
		return
	}

	var id int
	if id, err = a.AuthUsecase.Register(user); err != nil {
		domain.WriteError(w, err.Error(), domain.GetStatusCode(err))
		logs.LogError(logs.Logger, "auth_http", "Register.register", err, "Failed to register")
		return
	}

	session, _, err := a.AuthUsecase.Login(domain.Credentials{Email: user.Email, Password: user.Password})
	if err != nil {
		domain.WriteError(w, err.Error(), domain.GetStatusCode(err))
		logs.LogError(logs.Logger, "auth_http", "Register.login", err, "Failed to login")
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    session.Token,
		Expires:  session.ExpiresAt,
		Path:     "/",
		HttpOnly: true,
	})

	domain.WriteResponse(
		w,
		map[string]interface{}{
			"id": id,
		},
		http.StatusOK,
	)
}

func (a *AuthHandler) auth(r *http.Request) (bool, error) {
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			return false, domain.ErrUnauthorized
		}

		return false, domain.ErrBadRequest
	}
	if c.Expires.After(time.Now()) {
		return false, domain.ErrUnauthorized
	}
	sessionToken := c.Value
	exists, err := a.AuthUsecase.IsAuth(sessionToken)
	if err != nil {
		return false, domain.ErrInternalServerError
	}
	if !exists {
		return false, domain.ErrUnauthorized
	}

	return true, nil
}

func (a *AuthHandler) CloseAndAlert(body io.ReadCloser) {
	err := body.Close()
	if err != nil {
		logs.LogError(logs.Logger, "auth_http", "CloseAndAlert", err, "Failed to close body")
	}
}

func valid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func checkCredentials(cred domain.Credentials) error {
	if cred.Email == "" || len(cred.Password) == 0 {
		return domain.ErrWrongCredentials
	}

	if !valid(cred.Email) {
		return domain.ErrWrongCredentials
	}

	return nil
}
