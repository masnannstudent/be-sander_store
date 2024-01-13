package handler

import (
	"github.com/gofiber/fiber/v2"
	"mime/multipart"
	"ruti-store/module/feature/home/domain"
	"ruti-store/utils/response"
	"ruti-store/utils/upload"
	"ruti-store/utils/validator"
	"strconv"
)

type HomeHandler struct {
	service domain.HomeServiceInterface
}

func NewHomeHandler(service domain.HomeServiceInterface) domain.HomeHandlerInterface {
	return &HomeHandler{
		service: service,
	}
}

func (h *HomeHandler) CreateCarousel(c *fiber.Ctx) error {
	req := new(domain.CreateCarouselRequest)
	file, err := c.FormFile("photo")
	var uploadedURL string
	if err == nil {
		fileToUpload, err := file.Open()
		if err != nil {
			return response.ErrorBuildResponse(c, fiber.StatusInternalServerError, "Error opening file: "+err.Error())
		}
		defer func(fileToUpload multipart.File) {
			_ = fileToUpload.Close()
		}(fileToUpload)

		uploadedURL, err = upload.ImageUploadHelper(fileToUpload)
		if err != nil {
			return response.ErrorBuildResponse(c, fiber.StatusInternalServerError, "Error uploading file: "+err.Error())
		}
	}

	req.Photo = uploadedURL

	if err := c.BodyParser(req); err != nil {
		return response.ErrorBuildResponse(c, fiber.StatusBadRequest, "Failed to parse request body")
	}

	if err := validator.ValidateStruct(req); err != nil {
		return response.ErrorBuildResponse(c, fiber.StatusBadRequest, err.Error())
	}

	result, err := h.service.CreateCarousel(req)
	if err != nil {
		return response.ErrorBuildResponse(c, fiber.StatusInternalServerError, "Internal server error occurred: "+err.Error())
	}

	return response.SuccessBuildResponse(c, fiber.StatusCreated, "Success create carousels", domain.CarouselFormatter(result))

}

func (h *HomeHandler) GetCarouselByID(c *fiber.Ctx) error {
	id := c.Params("id")
	carouselID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return response.ErrorBuildResponse(c, fiber.StatusBadRequest, "Invalid input format.")
	}

	result, err := h.service.GetCarouselById(carouselID)
	if err != nil {
		return response.ErrorBuildResponse(c, fiber.StatusInternalServerError, "Failed to retrieve carousel: "+err.Error())
	}

	return response.SuccessBuildResponse(c, fiber.StatusOK, "Successfully retrieved carousel by ID", domain.CarouselFormatter(result))
}

func (h *HomeHandler) GetAllCarouselItems(c *fiber.Ctx) error {
	currentPage, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		return response.ErrorBuildResponse(c, fiber.StatusBadRequest, "Invalid page number")
	}

	pageSize, err := strconv.Atoi(c.Query("page_size"))
	if err != nil {
		return response.ErrorBuildResponse(c, fiber.StatusBadRequest, "Invalid page size")
	}

	result, totalItems, err := h.service.GetAllCarouselItems(currentPage, pageSize)
	if err != nil {
		return response.ErrorBuildResponse(c, fiber.StatusInternalServerError, "Internal server error occurred: "+err.Error())
	}

	currentPage, totalPages, nextPage, prevPage, err := h.service.GetCarouselPage(currentPage, pageSize)
	if err != nil {
		return response.ErrorBuildResponse(c, fiber.StatusInternalServerError, "Failed to get page info: "+err.Error())
	}

	return response.PaginationBuildResponse(c, fiber.StatusOK, "Success get pagination",
		domain.ResponseArrayCarousel(result), currentPage, int(totalItems), totalPages, nextPage, prevPage)
}