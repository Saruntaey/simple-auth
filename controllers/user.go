package controllers

import (
	"net/http"

	"github.com/saruntaey/simple-auth/config"
)

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
	w.Write([]byte("register"))
}

// @desc    Login user
// @route   < GET | POST > /login
// @access  Public
func (c *Controller) Login(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("login"))
}

// @desc    Get current user
// @route   GET /me
// @access  Private
func (c *Controller) GetMe(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("get me"))
}

// @desc    Update user
// @route   < GET | PUT > /update
// @access  Private
func (c *Controller) Update(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("update details"))
}

// @desc    Logout user
// @route   GET /logout
// @access  Private
func (c *Controller) Logout(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("logout"))
}
