package usecase

import (
	"context"
	"petstore/internal/domain"
)

type petUsecase struct {
	petRepo      domain.PetRepository
	categoryRepo domain.CategoryRepository
	tagRepo      domain.TagRepository
}

func (p *petUsecase) UploadImage(ctx context.Context, photo *domain.PhotoDTO) error {
	//TODO implement me
	panic("implement me")
}

func (p *petUsecase) GetByStatus(ctx context.Context, status domain.PetStatus) ([]*domain.Pet, error) {
	petsDTO, err := p.petRepo.GetByStatus(ctx, status)
	if err != nil {
		return nil, err
	}

	pets := make([]*domain.Pet, 0, len(petsDTO))

	for _, petDTO := range petsDTO {
		category, err := p.categoryRepo.Get(ctx, petDTO.CategoryId)
		if err != nil {
			return nil, err
		}

		tags, err := p.tagRepo.GetPetTags(ctx, petDTO.Id)
		if err != nil {
			return nil, err
		}

		pet := domain.PetDTOToPet(petDTO, category, tags)
		pets = append(pets, pet)
	}

	return pets, nil
}

func (p *petUsecase) Get(ctx context.Context, id int) (*domain.Pet, error) {
	petDTO, err := p.petRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	category, err := p.categoryRepo.Get(ctx, petDTO.CategoryId)
	if err != nil {
		return nil, err
	}

	tags, err := p.tagRepo.GetPetTags(ctx, id)

	pet := domain.PetDTOToPet(petDTO, category, tags)
	return pet, nil
}

func (p *petUsecase) Update(ctx context.Context, pet *domain.Pet) error {
	petDTO := domain.PetToPetDTO(pet)

	category, err := p.categoryRepo.GetElseCreate(ctx, pet.Category)
	if err != nil {
		return err
	}

	petDTO.CategoryId = category.Id
	if err := p.tagRepo.RemovePetTags(ctx, pet.Id); err != nil {
		return err
	}

	tagsIds := make([]int, 0, len(pet.Tags))
	for _, tag := range pet.Tags {
		tagDb, err := p.tagRepo.GetElseCreate(ctx, tag)
		if err != nil {
			return err
		}

		tagsIds = append(tagsIds, tagDb.Id)
	}

	if err := p.tagRepo.AddTagsToPet(ctx, petDTO.Id, tagsIds); err != nil {
		return err
	}

	return p.petRepo.Update(ctx, petDTO)
}

func (p *petUsecase) Delete(ctx context.Context, id int) error {
	if err := p.tagRepo.RemovePetTags(ctx, id); err != nil {
		return err
	}

	return p.petRepo.Delete(ctx, id)
}

func (p *petUsecase) Create(ctx context.Context, pet *domain.Pet) error {
	category, err := p.categoryRepo.GetElseCreate(ctx, pet.Category)
	if err != nil {
		return err
	}

	petDTO := &domain.PetDTO{
		CategoryId: category.Id,
		Name:       pet.Name,
		Status:     pet.Status,
	}

	err = p.petRepo.Create(ctx, petDTO)
	if err != nil {
		return err
	}

	pet.Id = petDTO.Id

	tagsIds := make([]int, 0, len(pet.Tags))
	for _, tag := range pet.Tags {
		tagDb, err := p.tagRepo.GetElseCreate(ctx, tag)
		if err != nil {
			return err
		}

		tagsIds = append(tagsIds, tagDb.Id)
	}

	if err := p.tagRepo.AddTagsToPet(ctx, petDTO.Id, tagsIds); err != nil {
		return err
	}

	return nil
}

func NewPetUsecase(pr domain.PetRepository, cr domain.CategoryRepository, tr domain.TagRepository) domain.PetUsecase {
	return &petUsecase{
		petRepo:      pr,
		categoryRepo: cr,
		tagRepo:      tr,
	}
}
