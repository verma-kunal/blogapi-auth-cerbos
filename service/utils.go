package service

import (
	"context"
	"encoding/json"
	"io"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/verma-kunal/blogapi-auth-cerbos/db"
)

func (s *Service) retrievePost(ctx echo.Context) (db.Post, error) {

	// get id from request
	postId := ctx.Param("postId")
	id, err := strconv.ParseUint(postId, 10, 64)
	if err != nil {
		return db.Post{}, err
	}

	return s.posts.GetPost(id)
}

func readPost(r io.Reader) (db.Post, error) {

	decoder := json.NewDecoder(r)

	var post db.Post
	err := decoder.Decode(&post)

	return post, err
}

func getCurrentUser(ctx context.Context) string {
	actx := getAuthContext(ctx)
	if actx == nil {
		return ""
	}

	return actx.username
}
