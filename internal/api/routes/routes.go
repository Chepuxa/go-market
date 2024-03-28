package routes

import (
	"training/proj/internal/api/handlers"
	"training/proj/internal/api/middleware"
	"training/proj/internal/config"
	"training/proj/internal/customerrors"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
)

var tokenAuth *jwtauth.JWTAuth

func SetupRoutes(r *chi.Mux, h *handlers.Handlers, cfg *config.Config) {
	tokenAuth = jwtauth.New("HS256", []byte(cfg.JwtSecret), nil)

	r.NotFound(customerrors.NotFoundResponse)
	r.MethodNotAllowed(customerrors.MethodNotAllowedResponse)

	r.Use(middleware.RecoverPanic)

	r.Route("/api/v1/", func(r chi.Router) {
		r.Mount("/categories", categoryRoutes(h.CategoryHandler))
		r.Mount("/items", itemsRoutes(h.ItemHandler))
		r.Mount("/users", usersRoutes(h.UserHandler))
	})
}

func categoryRoutes(h *handlers.CategoryHandler) *chi.Mux {

	r := chi.NewRouter()

	r.Get("/", h.GetAllCategories)
	r.Get("/{category_id}", h.GetCategory)
	r.Get("/{category_id}/items", h.GetCategoryItems)

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(middleware.Authenticator(tokenAuth))
		r.Put("/{category_id}", h.PutCategory)
		r.Delete("/{category_id}", h.DeleteCategory)
		r.Post("/", h.PostCategory)
		r.Put("/{category_id}/items/{item_id}", h.PutCategoryItem)
	})

	return r
}

func itemsRoutes(h *handlers.ItemHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", h.GetAllItems)
	r.Get("/{item_id}", h.GetItem)
	r.Get("/{item_id}/categories", h.GetItemCategories)

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(middleware.Authenticator(tokenAuth))
		r.Post("/", h.PostItem)
		r.Put("/{item_id}", h.PutItem)
		r.Delete("/{item_id}", h.DeleteItem)
	})
	return r
}

func usersRoutes(h *handlers.UserHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Post("/signup", h.PostUser)
	r.Get("/auth", h.Login)

	return r
}
