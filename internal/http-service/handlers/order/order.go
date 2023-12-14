package order

import (
	"L0_azat/internal/domain"
	resp "L0_azat/internal/lib/api/response"
	"L0_azat/internal/lib/logger/sl"
	"L0_azat/internal/storage"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type OrdGetter interface {
	GetOrder(orderUid string) (*domain.Message, error)
}

func New(log *slog.Logger, orderGetter OrdGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "http-service.handlers.order.New"

		log := log.With(
			slog.String("fn", fn),
			slog.String("reques_id", middleware.GetReqID(r.Context())),
		)

		orderUid := chi.URLParam(r, "orderUid")
		if orderUid == "" {
			log.Info("no orderUid is provided")
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}

		resOrder, err := orderGetter.GetOrder(orderUid)
		if errors.Is(err, storage.ErrOrderNotFound) {
			log.Info("order not found", "orderUid", orderUid)
			render.JSON(w, r, resp.Error("not found"))
			return
		}

		if err != nil {
			log.Error("failed to get order", sl.Err(err))
			render.JSON(w, r, resp.Error("internal server error"))
			return
		}

		log.Info("got order", slog.String("orderUid", orderUid))

		render.JSON(w, r, resOrder)
	}
}
