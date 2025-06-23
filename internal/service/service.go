package service

import (
	"Templatest/internal/dto"
	repo2 "Templatest/internal/repo"
	"Templatest/pkg/validator"
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type Service interface {
	Create(ctx *fiber.Ctx) error
	Read(ctx *fiber.Ctx) error
	ReadAll(ctx *fiber.Ctx) error
	Update(ctx *fiber.Ctx) error
	Delete(ctx *fiber.Ctx) error
}

type service struct {
	repo repo2.Repository
	log  *zap.SugaredLogger
}

func NewService(repo repo2.Repository, logger *zap.SugaredLogger) Service {
	return &service{
		repo: repo,
		log:  logger,
	}
}

func (s *service) Create(ctx *fiber.Ctx) error {
	var obj PostRequest

	// deserialize  JSON-request
	if err := json.Unmarshal(ctx.Body(), &obj); err != nil {
		s.log.Error("Invalid request body", zap.Error(err))
		return dto.BadResponseError(ctx, dto.FieldBadFormat, "Invalid request body")
	}

	// validation
	if vErr := validator.Validate(ctx.Context(), obj); vErr != nil {
		s.log.Error("Invalid request data", zap.Error(vErr))
		return dto.BadResponseError(ctx, dto.FieldIncorrect, vErr.Error())
	}

	// adds to memory
	dataObj := repo2.DataObject{
		Title: obj.Title,
		Data:  obj.Data,
	}
	id, err := s.repo.Post(ctx.Context(), dataObj)
	if err != nil {
		s.log.Error("Failed to insert object", zap.Error(err))
		return dto.InternalServerError(ctx)
	}
	s.log.Infof("object was appended %s", dataObj.Title)

	// forms the answer
	response := dto.Response{
		Status: "success",
		Data:   map[string]int{"task_id": id},
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (s *service) Read(ctx *fiber.Ctx) error {
	req := RequestWithId{ID: ctx.Params("id")}

	// Validation
	if vErr := validator.Validate(ctx.Context(), req); vErr != nil {
		s.log.Error("Invalid request id", zap.Error(vErr))
		return dto.BadResponseError(ctx, dto.FieldIncorrect, vErr.Error())
	}

	// Gets from memory
	objPtr, err := s.repo.Get(ctx.Context(), req.ID)
	if errors.Is(err, dto.ErrInvalidID) {
		s.log.Error("Failed to get task", zap.Error(err))
		return dto.BadResponseError(ctx, dto.NotFound, err.Error())
	}
	if err != nil {
		s.log.Error("Failed to parse id", zap.Error(err))
		return dto.InternalServerError(ctx)
	}

	// Forms answer
	response := dto.Response{
		Status: "success",
		Data:   objPtr,
	}
	return ctx.Status(fiber.StatusOK).JSON(response)

}

func (s *service) ReadAll(ctx *fiber.Ctx) error {
	// Gets from memory
	objsPtr, err := s.repo.GetAll(ctx.Context())
	if err != nil {
		s.log.Error("Failed to get task", zap.Error(err))
		return dto.InternalServerError(ctx)
	}

	// Forms answer
	response := dto.Response{
		Status: "success",
		Data:   objsPtr,
	}
	s.log.Info("whole memory was read and sent")
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (s *service) Update(ctx *fiber.Ctx) error {

	req := RequestWithId{ID: ctx.Params("id")}

	// Validation
	if vErr := validator.Validate(ctx.Context(), req); vErr != nil {
		s.log.Error("Invalid request id", zap.Error(vErr))
		return dto.BadResponseError(ctx, dto.FieldIncorrect, vErr.Error())
	}

	// Gets from memory
	objPtr, err := s.repo.Put(ctx.Context(), req.ID)
	if err != nil {
		s.log.Error("Failed to get task", zap.Error(err))
		return dto.InternalServerError(ctx)
	}

	// Forms answer
	response := dto.Response{
		Status: "success",
		Data:   objPtr,
	}
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (s *service) Delete(ctx *fiber.Ctx) error {
	req := RequestWithId{ID: ctx.Params("id")}

	// Validation
	if vErr := validator.Validate(ctx.Context(), req); vErr != nil {
		s.log.Error("Invalid request id", zap.Error(vErr))
		return dto.BadResponseError(ctx, dto.FieldIncorrect, vErr.Error())
	}

	// Gets from memory
	objPtr, err := s.repo.Delete(ctx.Context(), req.ID)
	if err != nil {
		s.log.Error("Failed to get task", zap.Error(err))
		return dto.InternalServerError(ctx)
	}

	// Forms answer
	response := dto.Response{
		Status: "success",
		Data:   objPtr,
	}

	s.log.Infof("object with id %s was deleted", objPtr.ID)
	return ctx.Status(fiber.StatusOK).JSON(response)
}
