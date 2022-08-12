package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/ozan1338/db/sqlc"
)

// Server serves HTTP Request for our banking services
type Server struct {
	store db.Store
	router *gin.Engine
}

// NewServer creates a new HTTP server and setup routing 
func NewServer(store db.Store) *Server {
	server := &Server{
		store: store,
	}
	router := gin.Default()

	// if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
	// 	v.RegisterValidator("currency", validCurrency)
	// }

	//add routes to routes
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccountById)
	router.GET("/accounts", server.listAccount)

	
	router.POST("/transfer", server.transferApi)

	server.router = router

	return server
}

// Start runs the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}