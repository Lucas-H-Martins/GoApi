package routes

import (
	"context"
	"database/sql"

	"goapi/controllers"
	"goapi/middleware"
	"goapi/repository"
	"net/http"
	"regexp"
	"strconv"
)

type Router struct {
	db *sql.DB
	*http.ServeMux
	routes []route
}

type route struct {
	pattern *regexp.Regexp
	method  string
	handler http.HandlerFunc
}

func NewRouter(db *sql.DB) *Router {
	r := &Router{
		db:       db,
		ServeMux: http.NewServeMux(),
	}
	r.setupRoutes()
	return r
}

func (r *Router) setupRoutes() {
	// Initialize repositories
	userRepo := repository.NewPostgresUserRepository(r.db)

	// Initialize controllers
	userController := controllers.NewUserController(userRepo)

	// Public routes
	r.Handle("/hello", http.HandlerFunc(controllers.HelloWorld))

	// Protected routes
	r.Handle("/users", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			userController.CreateUser(w, r)
		case http.MethodGet:
			userController.ListUsers(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	// User routes with ID
	r.routes = append(r.routes, route{
		pattern: regexp.MustCompile(`^/users/(\d+)$`),
		method:  "GET",
		handler: userController.GetUser,
	})
	r.routes = append(r.routes, route{
		pattern: regexp.MustCompile(`^/users/(\d+)$`),
		method:  "PUT",
		handler: userController.UpdateUser,
	})
	r.routes = append(r.routes, route{
		pattern: regexp.MustCompile(`^/users/(\d+)$`),
		method:  "DELETE",
		handler: userController.DeleteUser,
	})

	// Handle all requests
	r.Handle("/", http.HandlerFunc(r.routeHandler))
}

func (r *Router) routeHandler(w http.ResponseWriter, req *http.Request) {
	for _, route := range r.routes {
		matches := route.pattern.FindStringSubmatch(req.URL.Path)
		if len(matches) > 0 {
			if req.Method != route.method {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}

			// Update request context with URL parameters
			if len(matches) > 1 {
				id, _ := strconv.ParseInt(matches[1], 10, 64)
				req = req.WithContext(withURLParam(req.Context(), "id", id))
			}

			middleware.AuthMiddleware(http.HandlerFunc(route.handler)).ServeHTTP(w, req)
			return
		}
	}

	http.NotFound(w, req)
}

// Helper functions for URL parameters
type contextKey string

const urlParamsKey contextKey = "urlParams"

func withURLParam(ctx context.Context, key string, value interface{}) context.Context {
	params := ctx.Value(urlParamsKey)
	if params == nil {
		params = make(map[string]interface{})
	}
	paramsMap := params.(map[string]interface{})
	paramsMap[key] = value
	return context.WithValue(ctx, urlParamsKey, paramsMap)
}

func GetURLParam(r *http.Request, key string) interface{} {
	params := r.Context().Value(urlParamsKey)
	if params == nil {
		return nil
	}
	return params.(map[string]interface{})[key]
}
