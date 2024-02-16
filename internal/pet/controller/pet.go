package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"net/http"
	"petstore/internal/domain"
	"petstore/internal/responder"
	"strconv"
)

type petController struct {
	responder  responder.Responder
	petUsecase domain.PetUsecase
}

func NewPetController(r chi.Router, responder responder.Responder, pu domain.PetUsecase) {
	controller := &petController{
		responder:  responder,
		petUsecase: pu,
	}

	r.Route("/pet", func(r chi.Router) {
		r.Get("/{petId}", controller.Get)
		r.Post("/", controller.Create)
		r.Put("/", controller.Update)
		r.Delete("/{petId}", controller.Delete)

		r.Get("/findByStatus", controller.FindByStatus)
	})
}

// Create this function is used to create a new pet in the store.
//
// @Summary		Add a new pet to the store
// @Tags		pet
// @Accept		json
// @Produce		json
// @Security 	ApiKeyAuth
//
// @Param		pet		body		domain.Pet			true	"Pet object that needs to be added to the store"
//
// @Success		200		{object}	domain.Pet			"Pet object that was added"
// @Failure		400		{string}	string				"Invalid input"
// @Router		/pet 	[post]
func (p *petController) Create(w http.ResponseWriter, r *http.Request) {
	var petInput domain.Pet

	err := json.NewDecoder(r.Body).Decode(&petInput)
	if err != nil {
		p.responder.ErrorBadRequest(w, err)
		return
	}

	err = p.petUsecase.Create(r.Context(), &petInput)
	if err != nil {
		p.responder.ErrorInternal(w, err)
		return
	}

	p.responder.OutputJSON(w, responder.Response{
		Success: true,
		Message: "pet created",
		Data:    petInput,
	})
}

// Get this function is used to get a pet from the store.
//
// @Summary		Get a pet by ID
// @Tags		pet
// @Produce		json
// @Security 	ApiKeyAuth
//
//	@Param		petId	path		int				true	"ID of pet to return"
//
// @Success		200		{object}	domain.Pet			"Find pet by ID"
// @Failure		400		{string}	string				"Invalid input"
// @Failure		404		{string}	string				"Pet not found"
// @Router		/pet/{petId} 		[get]
func (p *petController) Get(w http.ResponseWriter, r *http.Request) {
	petId := chi.URLParam(r, "petId")
	if petId == "" {
		p.responder.ErrorBadRequest(w, fmt.Errorf("param petId is required"))
		return
	}

	petIdInt, err := strconv.Atoi(petId)
	if err != nil {
		p.responder.ErrorBadRequest(w, fmt.Errorf("param petId must be integer"))
		return
	}

	pet, err := p.petUsecase.Get(r.Context(), petIdInt)
	if err != nil {
		if errors.Is(err, domain.ErrPetNotFound) {
			p.responder.ErrorNotFound(w, err)
		} else {
			p.responder.ErrorInternal(w, err)
		}

		return
	}

	p.responder.OutputJSON(w, responder.Response{
		Success: true,
		Message: "get pet",
		Data:    pet,
	})
}

// Update this function is used to update a pet in the store.
//
// @Summary		Update a pet in the store with form data
// @Tags		pet
// @Accept		json
// @Produce		json
// @Security 	ApiKeyAuth
//
// @Param		pet		body		domain.Pet			true	"Pet object that needs to update"
//
// @Success		200		{string}	string				"Pet updated"
// @Failure		400		{string}	string				"Invalid input"
// @Failure		404		{string}	string				"Pet not found"
// @Router		/pet 	[put]
func (p *petController) Update(w http.ResponseWriter, r *http.Request) {
	var petInput domain.Pet

	err := json.NewDecoder(r.Body).Decode(&petInput)
	if err != nil {
		p.responder.ErrorBadRequest(w, err)
		return
	}

	err = p.petUsecase.Update(r.Context(), &petInput)
	if err != nil {
		p.responder.ErrorInternal(w, err)
		return
	}

	p.responder.OutputJSON(w, responder.Response{
		Success: true,
		Message: "pet updated",
		Data:    nil,
	})
}

// Delete this function is used to delete a pet from the store.
//
// @Summary		Delete a pet by ID
// @Tags		pet
// @Produce		json
// @Security 	ApiKeyAuth
//
// @Param		petId	path		int				true	"ID of pet to delete"
//
// @Success		200		{string}	string				"Pet deleted"
// @Failure		400		{string}	string				"Invalid input"
// @Failure		404		{string}	string				"Pet not found"
// @Router		/pet/{petId} 		[delete]
func (p *petController) Delete(w http.ResponseWriter, r *http.Request) {
	petId := chi.URLParam(r, "petId")
	if petId == "" {
		p.responder.ErrorBadRequest(w, fmt.Errorf("param petId is required"))
		return
	}

	petIdInt, err := strconv.Atoi(petId)
	if err != nil {
		p.responder.ErrorBadRequest(w, fmt.Errorf("param petId must be integer"))
		return
	}

	err = p.petUsecase.Delete(r.Context(), petIdInt)
	if err != nil {
		if errors.Is(err, domain.ErrPetNotFound) {
			p.responder.ErrorNotFound(w, err)
		} else {
			p.responder.ErrorInternal(w, err)
		}

		return
	}

	p.responder.OutputJSON(w, responder.Response{
		Success: true,
		Message: "get deleted",
		Data:    nil,
	})
}

// FindByStatus this function is used to get a pet from the store by pet status.
//
// @Summary		Delete a pet by ID
// @Tags		pet
// @Produce		json
// @Security 	ApiKeyAuth
//
// @Param		status	query		string				true	"Status values that need to be considered for filter"
//
// @Success		200		{object}	[]domain.Pet		"Pets found by status"
// @Failure		400		{string}	string				"Invalid input"
// @Failure		404		{string}	string				"Pet not found"
// @Router		/pet/findByStatus 		[get]
func (p *petController) FindByStatus(w http.ResponseWriter, r *http.Request) {
	if !r.URL.Query().Has("status") {
		p.responder.ErrorBadRequest(w, fmt.Errorf("param status is required"))
		return
	}

	status, err := domain.PetStatusFromString(r.URL.Query().Get("status"))
	if err != nil {
		p.responder.ErrorBadRequest(w, err)
		return
	}

	pets, err := p.petUsecase.GetByStatus(r.Context(), status)
	if err != nil {
		p.responder.ErrorInternal(w, err)
		return
	}

	p.responder.OutputJSON(w, responder.Response{
		Success: true,
		Message: "find pet by status",
		Data:    pets,
	})
}
