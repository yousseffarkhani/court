package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/yousseffarkhani/court/middlewares"
	"github.com/yousseffarkhani/court/model"
)

/* Court Handlers */
func (server *BasketServer) APIGetAllCourts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", JsonContentType)
	courts := server.database.GetAllCourts()
	if err := json.NewEncoder(w).Encode(courts); err != nil {
		log.Fatalln(err)
	}
}

func (server *BasketServer) getComments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", JsonContentType)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	comments, err := server.database.GetComments(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(comments) > 0 {
		claims, ok := context.Get(r, "claims").(*middlewares.Claims)
		if ok {
			for i, comment := range comments {
				if comment.Author == claims.Username {
					comments[i].IsAuthor = true
				} else {
					comments[i].IsAuthor = false
				}
			}
		}
	}

	json.NewEncoder(w).Encode(&comments)
}

func (server *BasketServer) addComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	claims, ok := context.Get(r, "claims").(*middlewares.Claims)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	comment := model.Comment{
		Message: extractFieldFromForm("message", r),
		CourtID: id,
		Author:  claims.Username,
	}

	errors := extractEmptyFieldErrors(comment)

	if len(errors) > 0 {
		http.Redirect(w, r, fmt.Sprintf("/court/%d", id), http.StatusFound)
		return
	}

	err = server.database.AddComment(comment)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/court/%d", id), http.StatusFound)
}

func (server *BasketServer) deleteComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	courtID, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	commentID, err := strconv.Atoi(vars["commentID"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	comment, err := server.database.GetComment(courtID, commentID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	claims, ok := context.Get(r, "claims").(*middlewares.Claims)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if comment.Author != claims.Username {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = server.database.DeleteComment(courtID, commentID)
	if err != nil {
		if err.Error() == "Comment doesn't exist." {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (server *BasketServer) updateComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	courtID, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	commentID, err := strconv.Atoi(vars["commentID"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	comment, err := server.database.GetComment(courtID, commentID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	claims, ok := context.Get(r, "claims").(*middlewares.Claims)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if comment.Author != claims.Username {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	message := struct {
		Id      int
		Message string
	}{}
	err = json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = server.database.UpdateComment(courtID, commentID, message.Message)
	if err != nil {
		if err.Error() == "Comment doesn't exist." {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}
