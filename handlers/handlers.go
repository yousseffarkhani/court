package handlers

import (
	"net/http"

	"github.com/yousseffarkhani/court/courtdb"
)

type BasketServer struct {
	store *courtdb.CourtStore
	http.Handler
}
