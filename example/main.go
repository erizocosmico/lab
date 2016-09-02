package main

import (
	"log"
	"net/http"
	"sync"

	"gopkg.in/mgo.v2/bson"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/mvader/lab"
)

const (
	IsAdmin            = "is-admin"
	AdminOnly          = "admin-only"
	NewWelcomeMessage  = "new-welcome-message"
	RealNameOnRegister = "real-name-on-register"
)

func CreateLab() lab.Lab {
	l := lab.New()
	l.DefineStrategy(IsAdmin, isAdminStrategy)
	l.Experiment(AdminOnly).
		Aim(lab.AimStrategy(IsAdmin, nil))
	l.Experiment(NewWelcomeMessage).
		Aim(lab.AimRandom())
	l.Experiment(RealNameOnRegister).
		Aim(lab.AimPercent(50))
	return l
}

func isAdminStrategy(v lab.Visitor, p lab.Params) bool {
	if visitor, ok := v.(*user); ok {
		return visitor.isAdmin
	}
	return false
}

type user struct {
	id       string
	password string
	isAdmin  bool
	name     string
}

func (u user) ID() string {
	return u.id
}

type Database struct {
	mut   sync.RWMutex
	users map[string]*user
}

func NewDatabase() *Database {
	return &Database{
		users: make(map[string]*user),
	}
}

func (d *Database) Register(u *user) {
	d.mut.Lock()
	defer d.mut.Unlock()
	d.users[u.id] = u
}

func (d *Database) Get(id string) *user {
	d.mut.Lock()
	defer d.mut.Unlock()
	return d.users[id]
}

var (
	database   = NewDatabase()
	laboratory = CreateLab()
	guestUser  = &user{}
)

func main() {
	r := gin.Default()
	store := sessions.NewCookieStore([]byte("so secret"))
	r.Use(sessions.Sessions("cool-session", store))
	r.Use(sessionMiddleware)
	r.LoadHTMLGlob("templates/*")
	r.Any("/", indexHandler)
	r.Any("/signup", signupHandler)

	log.Fatal(r.Run(":8080"))
}

func sessionMiddleware(c *gin.Context) {
	session := sessions.Default(c)
	var (
		s lab.Session
		u *user
	)

	if v, ok := session.Get("user").(string); ok {
		u = database.Get(v)
		c.Set("loggedIn", u != nil)
	}

	if u == nil {
		u = &user{
			id: "guest-" + bson.NewObjectId().Hex(),
		}
		c.Set("loggedIn", false)
	}

	s = laboratory.Session(u)
	c.Set("session", s)
	c.Next()
}

func session(c *gin.Context) lab.Session {
	return c.MustGet("session").(lab.Session)
}

func indexHandler(c *gin.Context) {
	s := session(c)
	c.HTML(http.StatusOK, "index.html", gin.H{
		"showNewWelcome": s.Launch(NewWelcomeMessage, nil),
		"adminOnly":      s.Launch(AdminOnly, nil),
		"loggedIn":       c.MustGet("loggedIn").(bool),
	})
}

func signupHandler(c *gin.Context) {
	if c.MustGet("loggedIn").(bool) {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	s := session(c)
	if c.Request.Method == http.MethodPost {
		username := c.PostForm("username")
		password := c.PostForm("password")
		admin := c.PostForm("admin") == "yes"

		user := &user{
			id:       username,
			password: password,
			isAdmin:  admin,
		}

		s.Launch(RealNameOnRegister, func() {
			realName := c.PostForm("name")
			user.name = realName
		})

		database.Register(user)
		session := sessions.Default(c)
		session.Set("user", username)
		session.Save()
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	c.HTML(http.StatusOK, "signup.html", gin.H{
		"showRealName": s.Launch(RealNameOnRegister, nil),
	})
}
