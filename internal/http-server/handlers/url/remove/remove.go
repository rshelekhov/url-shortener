package remove

import (
	"errors"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	resp "github.com/rshelekhov/url-shortener/internal/lib/api/response"
	"github.com/rshelekhov/url-shortener/internal/lib/logger/sl"
	"github.com/rshelekhov/url-shortener/internal/storage"
	"log/slog"
	"net/http"
)

type URLDeleter interface {
	DeleteURL(alias string) error
}

func Url(log *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.url.remove.Url"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")

			render.JSON(w, r, resp.Error("invalid reqest"))

			return
		}

		err := urlDeleter.DeleteURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", "alias", alias)

			render.JSON(w, r, resp.Error("not found"))

			return
		}
		if err != nil {
			log.Error("failed to get url", sl.Err(err))

			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		log.Info("deleted url", slog.String("alias", alias))

		render.JSON(w, r, http.StatusOK)
	}
}
