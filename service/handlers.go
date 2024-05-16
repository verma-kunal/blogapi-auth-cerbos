package service

import (
	"log"
	"net/http"

	"github.com/cerbos/cerbos-sdk-go/cerbos"
	"github.com/labstack/echo/v4"
	"github.com/verma-kunal/blogapi-auth-cerbos/db"
)

// view the post with id
func (s *Service) handlePostView(ctx echo.Context) error {

	// retrieve post info from request
	post, err := s.retrievePost(ctx)
	if err != nil {
		log.Printf("ERROR: %v", err)
		return ctx.String(http.StatusBadRequest, "Post not found")
	}

	if !s.isAllowedByCerbos(ctx.Request().Context(), postResource(post), "VIEW") {
		return ctx.String(http.StatusForbidden, "Operation not allowed")
	}

	return ctx.JSON(http.StatusOK, post)

}

func (s *Service) handlePostCreate(ctx echo.Context) error {

	post, err := readPost(ctx.Request().Body)
	if err != nil {
		log.Printf("ERROR: %v", err)
		return ctx.String(http.StatusBadRequest, "Post not found")
	}

	// create cerbos resource
	cerbosResource := cerbos.NewResource("post", "new").
		WithAttr("title", post.Title).
		WithAttr("owner", post.Owner)

	// cerbos auth
	if !s.isAllowedByCerbos(ctx.Request().Context(), cerbosResource, "CREATE") {
		return ctx.String(http.StatusForbidden, "Operation not allowed")
	}

	// get username from current echo context
	username := getCurrentUser(ctx.Request().Context())
	postId := s.posts.CreatePost(username, post)

	return ctx.JSON(http.StatusCreated, struct {
		PostId uint64 `json:"postId"`
	}{PostId: postId})
}

func (s *Service) handlePostUpdate(ctx echo.Context) error {

	// retrieve the post info to update
	post, err := s.retrievePost(ctx)
	if err != nil {
		log.Printf("ERROR: %v", err)
		return ctx.String(http.StatusBadRequest, "Post not found")
	}

	// cerbos auth
	if !s.isAllowedByCerbos(ctx.Request().Context(), postResource(post), "UPDATE") {
		return ctx.String(http.StatusForbidden, "Operation not allowed")
	}

	// read & validate the new post data from the request body
	var newPost db.Post
	if err := ctx.Bind(&newPost); err != nil {
		log.Printf("ERROR: %v", err)
		return ctx.String(http.StatusBadRequest, "Invalid post data")
	}

	// title check
	if newPost.Title == "" {
		return ctx.String(http.StatusBadRequest, "Title cannot be empty")
	}

	// update operation
	if err := s.posts.UpdatePost(post.PostId, newPost); err != nil {
		log.Printf("ERROR: %v", err)
		return ctx.String(http.StatusInternalServerError, "Failed to update post")
	}

	return ctx.String(http.StatusOK, "Post updated")

}

func (s *Service) handlePostDelete(ctx echo.Context) error {

	post, err := s.retrievePost(ctx)
	if err != nil {
		log.Printf("ERROR: %v", err)
		return ctx.String(http.StatusBadRequest, "Post not found")
	}

	// cerbos auth
	if !s.isAllowedByCerbos(ctx.Request().Context(), postResource(post), "DELETE") {
		return ctx.String(http.StatusForbidden, "Operation not allowed")
	}

	// delete operation
	if err := s.posts.DeletePost(post.PostId); err != nil {
		log.Printf("ERROR: %v", err)
		return ctx.String(http.StatusInternalServerError, "Failed to delete post")
	}

	return ctx.JSON(http.StatusOK, "Post Deleted")

}
