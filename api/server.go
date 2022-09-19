package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	db "github.com/ozan1338/db/sqlc"
	"github.com/ozan1338/token"
	"github.com/ozan1338/util"
)

// Server serves HTTP Request for our banking services
type Server struct {
	config util.Config
	store db.Store
	tokenMaker token.Maker
	router *gin.Engine
}

// NewServer creates a new HTTP server and setup routing 
func NewServer(config util.Config,store db.Store) (*Server,error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create tken maker : %w", err)
	}
	server := &Server{
		config: config,
		store: store,
		tokenMaker: tokenMaker,
	}

	server.setupRouter()

	return server,nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	// if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
	// 	v.RegisterValidator("currency", validCurrency)
	// }

	//add routes to routes
	router.POST("/register", server.createUser)
	router.POST("/login", server.loginUser)
	router.POST("/token/renew-token", server.renewAccessToken)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccountById)
	authRoutes.GET("/accounts", server.listAccount)

	
	authRoutes.POST("/transfer", server.transferApi)

	server.router = router

}

// Start runs the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}