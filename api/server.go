package api

import (
	"fmt"

	db "github.com/EphemSpirit/simple_bank/db/sqlc"
	"github.com/EphemSpirit/simple_bank/token"
	"github.com/EphemSpirit/simple_bank/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	store db.Store
	router *gin.Engine
	tokenMaker token.Maker
	config util.Config
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("could not create token maker %w", err)
	}
	server := &Server{
		config: config, 
		store: store, 
		tokenMaker: tokenMaker,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()

	return server, nil
}

func (s *Server) setupRouter() {
	router := gin.Default();

	router.POST("/accounts", s.CreateAccount)
	router.GET("/accounts/:id", s.GetAccount)
	router.GET("/accounts", s.GetAccounts)
	router.POST("/transfers", s.CreateTransfer)
	router.POST("/users", s.CreateUser)
	router.POST("/users/login", s.loginUser)

	s.router = router
}

// run server on specified server port
func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func errResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}