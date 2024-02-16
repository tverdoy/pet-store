package controller

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi"
	"net/http"
	"petstore/internal/domain"
	"petstore/internal/responder"
	"strconv"
)

type orderController struct {
	orderUsecase domain.OrderUsecase
	responder    responder.Responder
}

// Create this function is used to create a new order in the store.
//
// @Summary		Add a new order to the store
// @Tags 		store
// @Accept		json
// @Produce		json
// @Security 	ApiKeyAuth
//
// @Param		order		body	domain.Order		true	"Order object that needs to be added to the store"
//
// @Success		200		{object}	domain.Order		"Order object that was added"
// @Failure		400		{string}	string				"Invalid input"
// @Router		/store/order 	[post]
func (o *orderController) Create(w http.ResponseWriter, r *http.Request) {
	var orderInput domain.Order
	if err := json.NewDecoder(r.Body).Decode(&orderInput); err != nil {
		o.responder.ErrorBadRequest(w, err)
		return
	}

	if err := o.orderUsecase.Create(r.Context(), &orderInput); err != nil {
		o.responder.ErrorInternal(w, err)
		return
	}

	o.responder.OutputJSON(w, responder.Response{
		Success: true,
		Message: "order created",
		Data:    orderInput,
	})
}

// Get this function is used to get an order from the store.
//
// @Summary		Order an order by ID
// @Tags 		store
// @Produce		json
// @Security 	ApiKeyAuth
//
// @Param		orderId	path		int					true	"ID of order to return"
//
// @Success		200		{object}	domain.Order		"Find order by ID"
// @Failure		400		{string}	string				"Invalid input"
// @Failure		404		{string}	string				"Order not found"
// @Router		/store/order/{orderId} 		[get]
func (o *orderController) Get(w http.ResponseWriter, r *http.Request) {
	orderIdParam := chi.URLParam(r, "orderId")
	if orderIdParam == "" {
		o.responder.ErrorBadRequest(w, errors.New("orderId is required"))
		return
	}

	orderId, err := strconv.Atoi(orderIdParam)
	if err != nil {
		o.responder.ErrorBadRequest(w, err)
		return
	}

	order, err := o.orderUsecase.Get(r.Context(), orderId)
	if err != nil {
		o.responder.ErrorInternal(w, err)
		return
	}

	o.responder.OutputJSON(w, responder.Response{
		Success: true,
		Message: "get order",
		Data:    order,
	})
}

// Delete this function is used to delete an order from the store.
//
// @Summary		Delete order by ID
// @Tags 		store
// @Produce		json
// @Security 	ApiKeyAuth
//
// @Param		orderId	path		int					true	"ID of order to delete"
//
// @Success		200		{string}	string				"Order deleted"
// @Failure		400		{string}	string				"Invalid input"
// @Failure		404		{string}	string				"Order not found"
// @Router		/store/order/{orderId} 	[delete]
func (o *orderController) Delete(w http.ResponseWriter, r *http.Request) {
	orderIdParam := chi.URLParam(r, "orderId")
	if orderIdParam == "" {
		o.responder.ErrorBadRequest(w, errors.New("orderId is required"))
		return
	}

	orderId, err := strconv.Atoi(orderIdParam)
	if err != nil {
		o.responder.ErrorBadRequest(w, err)
		return
	}

	err = o.orderUsecase.Delete(r.Context(), orderId)
	if err != nil {
		o.responder.ErrorInternal(w, err)
		return
	}

	o.responder.OutputJSON(w, responder.Response{
		Success: true,
		Message: "order deleted",
		Data:    nil,
	})
}

func NewOrderController(r chi.Router, responder responder.Responder, orderUsecase domain.OrderUsecase) {
	controller := &orderController{orderUsecase: orderUsecase, responder: responder}

	r.Route("/store/order", func(r chi.Router) {
		r.Get("/{orderId}", controller.Get)
		r.Delete("/{orderId}", controller.Delete)
		r.Post("/", controller.Create)
	})
}
