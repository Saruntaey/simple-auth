package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/saruntaey/simple-auth/config"
	"github.com/saruntaey/simple-auth/models"
	"github.com/zebresel-com/mongodm"
	"gopkg.in/mgo.v2/bson"
)

type data struct {
	Title    string
	FlashMsg *models.FlashMsg
	User     *models.User
}

type Controller struct {
	appConfig *config.Config
}

func New(appConfig *config.Config) *Controller {
	return &Controller{
		appConfig: appConfig,
	}
}

// @desc    Register user
// @route   < GET | POST > /register
// @access  Public
func (c *Controller) Register(w http.ResponseWriter, r *http.Request) {
	session, errGetSession := c.NewSession().GetFromCookie(w, r)
	switch r.Method {

	case http.MethodGet:
		// check if user already login
		if errGetSession == nil && len(session.SessionModel.User) > 0 {
			http.Redirect(w, r, "/me", http.StatusSeeOther)
			return
		}

		data := data{
			Title: "register",
		}
		// add flash message to data
		if errGetSession == nil {
			data.FlashMsg = session.FlashMsg
		}
		c.render(w, "register", data)

	case http.MethodPost:
		r.ParseForm()

		User := c.appConfig.DbConn.Model("User")
		user := &models.User{}

		User.New(user)

		// fill data
		user.Email = r.PostForm.Get("email")
		user.PasswordRaw = r.PostForm.Get("password")
		user.Name = r.PostForm.Get("name")

		if valid, issues := user.ValidateCreate(); !valid {
			msg := ""
			for _, v := range issues {
				msg += v.Error()
				msg += ", "
			}
			msg = strings.TrimRight(msg, ", ")
			session.FlashAndRedirect(w, r, "danger", msg, "/register")
			return

		}

		err := user.HashPassword()
		if err != nil {
			msg := err.Error()
			session.FlashAndRedirect(w, r, "danger", msg, "/register")
			return
		}

		err = user.Save()
		if _, ok := err.(*mongodm.DuplicateError); ok {
			msg := fmt.Sprintf("the email %s is already taken", user.Email)
			session.FlashAndRedirect(w, r, "danger", msg, "/register")
			return
		} else if err != nil {
			fmt.Print(err)
			msg := "server error"
			session.FlashAndRedirect(w, r, "danger", msg, "/register")
			return
		}

		if errGetSession != nil {
			session.InitModel()
		}
		session.AddUser(user.Id)
		session.FlashAndRedirect(w, r, "success", "Welcome to Simple-Auth", "/me")

	default:
		session.FlashAndRedirect(w, r, "danger", "Method not allow", "/register")
	}
}

// @desc    Login user
// @route   < GET | POST > /login
// @access  Public
func (c *Controller) Login(w http.ResponseWriter, r *http.Request) {
	session, errGetSession := c.NewSession().GetFromCookie(w, r)

	switch r.Method {

	case http.MethodGet:
		// check if user already login
		if errGetSession == nil && len(session.SessionModel.User) > 0 {
			http.Redirect(w, r, "/me", http.StatusSeeOther)
			return
		}

		data := data{
			Title: "login",
		}
		// add flash message to data
		if errGetSession == nil {
			data.FlashMsg = session.FlashMsg
		}
		c.render(w, "login", data)

	case http.MethodPost:
		r.ParseForm()

		User := c.appConfig.DbConn.Model("User")
		user := &models.User{}

		email := r.PostForm.Get("email")
		passwordRaw := r.PostForm.Get("password")

		if len(email) == 0 || len(passwordRaw) == 0 {
			session.FlashAndRedirect(w, r, "danger", "Please provide email and password", "/login")
			return
		}

		query := bson.M{
			"email": email,
		}

		err := User.FindOne(query).Exec(user)
		if _, ok := err.(*mongodm.NotFoundError); ok {
			session.FlashAndRedirect(w, r, "danger", "Invalid credential", "/login")
			return
		} else if err != nil {
			session.FlashAndRedirect(w, r, "warning", "Server error", "/login")
			return
		}

		if !user.MatchPassword(passwordRaw) {
			session.FlashAndRedirect(w, r, "danger", "Invalid credential", "/login")
			return
		}

		if errGetSession != nil {
			session.InitModel()
		}
		session.AddUser(user.Id)
		session.FlashAndRedirect(w, r, "success", "Welcome back", "/me")

	default:
		session.FlashAndRedirect(w, r, "danger", "Method not allow", "/login")
	}
}

// @desc    Get current user
// @route   GET /me
// @access  Private
func (c *Controller) GetMe(w http.ResponseWriter, r *http.Request) {
	session, errGetSession := c.NewSession().GetFromCookie(w, r)
	switch r.Method {

	case http.MethodGet:
		if errGetSession != nil || len(session.SessionModel.User) == 0 {
			session.FlashAndRedirect(w, r, "danger", "Please login first", "/login")
			return
		}

		User := c.appConfig.DbConn.Model("User")
		user := &models.User{}
		User.FindId(session.SessionModel.User).Exec(user)
		data := data{
			Title:    "Profile",
			FlashMsg: session.FlashMsg,
			User:     user,
		}
		c.render(w, "me", data)

	default:
		if errGetSession == nil && len(session.SessionModel.User) != 0 {
			session.FlashAndRedirect(w, r, "danger", "Method not allow", "/me")
			return
		}
		session.FlashAndRedirect(w, r, "danger", "Method not allow", "/login")
	}
}

// @desc    Update user
// @route   < GET | POST > /update
// @access  Private
func (c *Controller) Update(w http.ResponseWriter, r *http.Request) {
	session, errGetSession := c.NewSession().GetFromCookie(w, r)

	switch r.Method {

	case http.MethodGet:
		// check if user login
		if errGetSession != nil || len(session.SessionModel.User) == 0 {
			session.FlashAndRedirect(w, r, "danger", "Please login first", "/login")
			return
		}

		// find user in DB
		User := c.appConfig.DbConn.Model("User")
		user := &models.User{}
		User.FindId(session.SessionModel.User).Exec(user)

		data := data{
			Title: "Update profile",
			User:  user,
		}
		// add flash message to data
		if errGetSession == nil {
			data.FlashMsg = session.FlashMsg
		}
		c.render(w, "update", data)

	case http.MethodPost:
		// check if user login
		if errGetSession != nil || len(session.SessionModel.User) == 0 {
			session.FlashAndRedirect(w, r, "danger", "Please login first", "/login")
			return
		}

		r.ParseForm()

		name := r.PostForm.Get("name")
		if len(name) == 0 {
			session.FlashAndRedirect(w, r, "danger", "Please provide name", "/update")
			return
		}

		User := c.appConfig.DbConn.Model("User")
		user := &models.User{}

		err := User.FindId(session.SessionModel.User).Exec(user)
		if err != nil {
			session.FlashAndRedirect(w, r, "warning", "Server error", "/update")
			return
		}

		d := map[string]interface{}{
			"name": name,
		}

		// The Update method is incompleted so the error is not handled
		// see https://github.com/zebresel-com/mongodm/issues/20
		user.Update(d)
		err = user.Save()
		if err != nil {
			session.FlashAndRedirect(w, r, "warning", "Server error", "/update")
			return
		}
		session.FlashAndRedirect(w, r, "success", "Your name was updated", "/me")

	default:
		if errGetSession == nil && len(session.SessionModel.User) != 0 {
			session.FlashAndRedirect(w, r, "danger", "Method not allow", "/update")
			return
		}
		session.FlashAndRedirect(w, r, "danger", "Method not allow", "/login")
	}
}

// @desc    Logout user
// @route   GET /logout
// @access  Private
func (c *Controller) Logout(w http.ResponseWriter, r *http.Request) {
	session, errGetSession := c.NewSession().GetFromCookie(w, r)

	switch r.Method {

	case http.MethodGet:
		if errGetSession != nil || len(session.SessionModel.User) == 0 {
			session.FlashAndRedirect(w, r, "danger", "You did not login", "/login")
			return
		}
		session.SessionModel.User = ""
		session.FlashAndRedirect(w, r, "success", "You successfully logout", "/login")

	default:
		if errGetSession == nil && len(session.SessionModel.User) != 0 {
			session.FlashAndRedirect(w, r, "danger", "Method not allow", "/me")
			return
		}
		session.FlashAndRedirect(w, r, "danger", "Method not allow", "/login")
	}
}
