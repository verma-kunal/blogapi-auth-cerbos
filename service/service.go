package service

import (
	"context"
	"log"
	"strconv"

	"github.com/cerbos/cerbos-sdk-go/cerbos"
	"github.com/labstack/echo/v4"
	"github.com/verma-kunal/blogapi-auth-cerbos/db"
	"golang.org/x/crypto/bcrypt"
)

type authCtxKeyType struct{}

var authCtxKey = authCtxKeyType{}

type authContext struct {
	username  string
	principal *cerbos.Principal
}

// Service implements the post API.
type Service struct {
	cerbos *cerbos.GRPCClient
	posts  *db.PostDB
}

// create cerbos resource for the given post
func postResource(post db.Post) *cerbos.Resource {

	return cerbos.NewResource("post", strconv.FormatUint(post.PostId, 10)).
		WithAttr("title", post.Title).
		WithAttr("owner", post.Owner)
}

// new cerbos instance
func New(cerbosAddr string) (*Service, error) {
	cerbosInstance, err := cerbos.New(cerbosAddr, cerbos.WithPlaintext())
	if err != nil {
		return nil, err
	}

	return &Service{cerbos: cerbosInstance, posts: db.NewPostDB()}, nil
}

func (s *Service) Handler() *echo.Echo {

	e := echo.New()
	e.Use(AuthMiddleware) // implementing the auth middleware

	// api routes
	e.PUT("/posts", s.handlePostCreate)
	e.GET("/posts/:postId", s.handlePostView)
	e.POST("/posts/:postId", s.handlePostUpdate)
	e.DELETE("/posts/:postId", s.handlePostDelete)

	return e
}

// auth middleware
func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {

		// get basic auth creds from the request
		user, pass, ok := ctx.Request().BasicAuth()

		if ok {
			// check the password and retrieve the auth context.
			authCtx, err := buildAuthContext(user, pass, ctx)
			if err != nil {
				log.Printf("Failed to authenticate user [%s]: %v", user, err)
			} else {
				// Add the retrieved principal to the context.
				newCtx := context.WithValue(ctx.Request().Context(), authCtxKey, authCtx)
				ctx.SetRequest(ctx.Request().WithContext(newCtx)) // setting the new request context

				return next(ctx)
			}
		}

		return next(ctx)
	}
}

// verify username & passwords + create cerbos principal
func buildAuthContext(username, password string, ctx echo.Context) (*authContext, error) {

	// get user from db
	user, err := db.FindUser(ctx.Request().Context(), username)
	if err != nil {
		return nil, err
	}

	// compare request password with the db password
	if err := (bcrypt.CompareHashAndPassword(user.Password, []byte(password))); err != nil {
		return nil, err
	}

	// create cerbos principal object using info from db & request
	newPrincipal := cerbos.NewPrincipal(username).
		WithRoles(user.Roles...).
		WithAttr("blocked", user.Blocked).
		WithAttr("ipAddress", ctx.Request().RemoteAddr)

	return &authContext{username: username, principal: newPrincipal}, nil
}

func (s *Service) isAllowedByCerbos(ctx context.Context, resource *cerbos.Resource, action string) bool {

	// get current cerbos principal
	principalCtx := s.principalContext(ctx)
	if principalCtx == nil {
		return false
	}

	// using the IsAllowed() utility function from "cerbos Principal Context"
	allowed, err := principalCtx.IsAllowed(ctx, resource, action)
	if err != nil {
		return false
	}

	return allowed
}

// retreive cerbos principal from the current context
func (s *Service) principalContext(ctx context.Context) cerbos.PrincipalContext {
	actx := getAuthContext(ctx)
	if actx == nil {
		log.Fatal("ERROR: auth context is nil")
	}

	return s.cerbos.WithPrincipal(actx.principal) // attaching the current principal
}

// get current auth context
func getAuthContext(ctx context.Context) *authContext {
	ac := ctx.Value(authCtxKey)
	if ac == nil {
		return nil
	}
	return ac.(*authContext)
}
