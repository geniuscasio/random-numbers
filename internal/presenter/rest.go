package presenter

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"random-numbers/internal/common"
	"random-numbers/internal/logic"
	"strings"
)

type RestAPI struct {
	log      *logrus.Entry
	service  logic.Service
	sessions logic.SessionWorker
}

func NewREST(log *logrus.Entry, service logic.Service, sessions logic.SessionWorker) (*RestAPI, error) {
	return &RestAPI{
		log:      log,
		service:  service,
		sessions: sessions,
	}, nil
}

const (
	sessionHeader = "Session-Id"
)

func (r *RestAPI) LoggingMiddleware(c *gin.Context) {
	r.log.WithFields(logrus.Fields{
		"url": c.Request.URL.Path,
	}).Info("new request")

	c.Next()
}

func (r *RestAPI) AuthMiddleware(c *gin.Context) {
	session := r.getSession(c)

	err := r.sessions.IsAuthenticated(session)
	if err != nil {
		r.handleError(c, err)

		return
	}

	c.Next()
}

func (r *RestAPI) StatsMiddleware(c *gin.Context) {
	session := r.getSession(c)

	userID, err := r.sessions.GetUserID(session)
	if err != nil {
		r.handleError(c, err)

		return
	}

	err = r.service.IncrementStatistic(userID)
	if err != nil {
		r.handleError(c, err)

		return
	}

	c.Next()
}

func (r *RestAPI) Generate(c *gin.Context) {
	var req = new(common.GenerateRequest)
	var order = new(common.Order)

	if err := c.ShouldBindJSON(req); err != nil {
		r.handleError(c, err)

		return
	}

	err := c.BindQuery(order)
	if err != nil {
		r.log.WithError(err).Error("cant bind order")
	}

	var orderDesc bool

	switch order.Order {
	case "asc":
		orderDesc = false
	case "desc", "": // stick to desc order as default
		orderDesc = true
	default:
		r.handleError(c, errors.New("wrong order value"))

		return
	}

	res, err := r.service.Generate(req.From, req.To, req.Total, orderDesc)
	if err != nil {
		r.handleError(c, err)

		return
	}

	c.JSON(http.StatusOK, res)
}

func (r *RestAPI) Register(c *gin.Context) {
	var req = new(common.User)

	if err := c.ShouldBindJSON(req); err != nil {
		r.handleError(c, err)

		return
	}

	if err := r.service.CreateUser(*req); err != nil {
		r.handleError(c, err)

		return
	}
}

func (r *RestAPI) Login(c *gin.Context) {
	var req = new(common.UserCredentials)

	if err := c.ShouldBindJSON(req); err != nil {
		r.handleError(c, err)

		return
	}

	userID, err := r.service.LoginUser(req.Email, req.Password)
	if err != nil {
		r.handleError(c, err)

		return
	}

	session, err := r.sessions.CreateSession(userID)
	if err != nil {
		r.handleError(c, err)

		return
	}

	c.Header(sessionHeader, session)
}

func (r *RestAPI) Details(c *gin.Context) {
	session := r.getSession(c)

	userID, err := r.sessions.GetUserID(session)
	if err != nil {
		r.handleError(c, err)

		return
	}

	res, err := r.service.GetStatistic(userID)
	if err != nil {
		r.handleError(c, err)

		return
	}

	c.JSON(http.StatusOK, res)
}

func (r *RestAPI) getSession(c *gin.Context) string {
	session := strings.Join(c.Request.Header[sessionHeader], "")

	return session
}

func (r *RestAPI) handleError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	c.Abort()

	r.log.WithError(err).Error("error occurred")
	c.JSON(http.StatusOK, struct {
		Error string `json:"error"`
	}{err.Error()})
}
