package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/RITesh-namifyx/RssAgg/internal/database"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerCreateFeedFollow(w http.ResponseWriter , r *http.Request , user database.User) {
	type parameters struct {
		FeedId uuid.UUID `json:"feed_id"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400 , fmt.Sprintf("Error parsing JSON:%v", err))
		return
	}

	FeedFollow, err := apiCfg.DB.CreateFeedFollow(r.Context(),database.CreateFeedFollowParams{

		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID: user.ID,
		FeedID: params.FeedId,

	})
	if  err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couidn't create user:%v", err))
		return
	}
	respondWithJSON(w,200, databaseFeedFollowToFeedFollow(FeedFollow))

}

func (apiCfg *apiConfig) GetAllFeedFollows(w http.ResponseWriter , r *http.Request , user database.User) {
	
	FeedFollow, err := apiCfg.DB.GetAllFeedFollows(r.Context(), user.ID)
	if  err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couidn't create user:%v", err))
		return
	}
	respondWithJSON(w,200, databaseFeedFollowsToFeedFollows(FeedFollow))

}

func (apiCfg *apiConfig) handlerDeleteFeedFollow(w http.ResponseWriter , r *http.Request , user database.User) {
	feedFollowIDStr := chi.URLParam(r ,"feedFollowID")
	feedFollowId, err := uuid.Parse(feedFollowIDStr)

	if  err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couidn't Parse feed follow ID:%v", err))
		return
	}

	err = apiCfg.DB.DeleteFeedFollow(r.Context(),database.DeleteFeedFollowParams{
		ID: feedFollowId,
		UserID: user.ID,
	})
	if  err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couidn't unfollow:%v", err))
		return
	}
	respondWithJSON(w,200, struct{
		Message string `json:"message"`
		Success bool `json:"success"`
		 }{
		"feed unfolowed successfully",
		true,
	} )
	// or to return just empty struct respondWithJSON(w,200, struct{}{})

}


