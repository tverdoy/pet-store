package domain

// # TODO add image upload

import (
	"context"
	"errors"
	"fmt"
)

type PetStatus string

var ErrPetNotFound = errors.New("pet not found")
var ErrCategoryNotFound = errors.New("category not found")
var ErrTagNotFound = errors.New("tag not found")

const (
	PetStatusAvailable PetStatus = "available"
	PetStatusPending             = "pending"
	PetStatusSold                = "sold"
)

type Pet struct {
	Id        int       `json:"id"`
	Category  *Category `json:"category"`
	Name      string    `json:"name"`
	Tags      []*Tag    `json:"tags"`
	Status    PetStatus `json:"status"`
	PhotoUrls []string  `json:"photoUrls"`
}

type PetDTO struct {
	Id         int       `json:"id"`
	CategoryId int       `json:"category_id"`
	Name       string    `json:"name"`
	Status     PetStatus `json:"status"`
}

type Category struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Tag struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type PhotoDTO struct {
	Id    int `json:"id"`
	PetId int `json:"pet_id"`
}

type PetUsecase interface {
	Get(ctx context.Context, id int) (*Pet, error)
	Create(ctx context.Context, pet *Pet) error
	Update(ctx context.Context, pet *Pet) error
	Delete(ctx context.Context, id int) error

	UploadImage(ctx context.Context, photo *PhotoDTO) error

	GetByStatus(ctx context.Context, status PetStatus) ([]*Pet, error)
}

type PetRepository interface {
	Get(ctx context.Context, id int) (*PetDTO, error)
	Create(ctx context.Context, pet *PetDTO) error
	Update(ctx context.Context, pet *PetDTO) error
	Delete(ctx context.Context, id int) error

	GetByStatus(ctx context.Context, status PetStatus) ([]*PetDTO, error)
}

type CategoryRepository interface {
	Create(ctx context.Context, category *Category) error
	Get(ctx context.Context, id int) (*Category, error)
	GetByName(ctx context.Context, name string) (*Category, error)

	// GetElseCreate get else create category.
	// If category is not exist, create it first.
	// Check if category exist by name.
	GetElseCreate(ctx context.Context, category *Category) (*Category, error)
}

type TagRepository interface {
	Create(ctx context.Context, tag *Tag) error
	GetByName(ctx context.Context, name string) (*Tag, error)

	// GetElseCreate get else create tag.
	// If tag is not exist, create it first.
	// Check if tag exist by name.
	GetElseCreate(ctx context.Context, tag *Tag) (*Tag, error)

	AddTagsToPet(ctx context.Context, petId int, tagIds []int) error
	GetPetTags(ctx context.Context, petId int) ([]*Tag, error)
	RemovePetTags(ctx context.Context, petId int) error
}

type PhotoRepository interface {
	Create(ctx context.Context, photo PhotoDTO) error
	GetByPet(ctx context.Context, petId int) ([]*PhotoDTO, error)
}

func PetDTOToPet(petDTO *PetDTO, category *Category, tags []*Tag) *Pet {
	return &Pet{
		Id:       petDTO.Id,
		Category: category,
		Name:     petDTO.Name,
		Tags:     tags,
		Status:   petDTO.Status,
	}
}

func PetToPetDTO(pet *Pet) *PetDTO {
	return &PetDTO{
		Id:         pet.Id,
		CategoryId: pet.Category.Id,
		Name:       pet.Name,
		Status:     pet.Status,
	}
}

func (p *PhotoDTO) GetPublicURL() string {
	return fmt.Sprintf("/static/pets/%d/photos/%d.jpg", p.PetId, p.Id)
}

func PetStatusFromString(status string) (PetStatus, error) {
	switch status {
	case string(PetStatusAvailable):
		return PetStatusAvailable, nil
	case PetStatusPending:
		return PetStatusPending, nil
	case PetStatusSold:
		return PetStatusSold, nil
	default:
		return PetStatusAvailable, errors.New("invalid status")
	}
}
