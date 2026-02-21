package http

import (
	"net/http"
	"bytes"
	"io"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rendyfutsuy/base-go/helpers/middleware"
	_reqContext "github.com/rendyfutsuy/base-go/helpers/middleware/request"
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/helpers/response"
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/course"
	"github.com/rendyfutsuy/base-go/modules/course/dto"
	utilsServices "github.com/rendyfutsuy/base-go/utils/services"
)

type CourseHandler struct {
	Usecase              course.Usecase
	validator            *validator.Validate
	mwPageRequest        _reqContext.IMiddlewarePageRequest
	middlewareAuth       middleware.IMiddlewareAuth
	middlewarePermission middleware.IMiddlewarePermission
}

func NewCourseHandler(e *echo.Echo, uc course.Usecase, mwP _reqContext.IMiddlewarePageRequest, auth middleware.IMiddlewareAuth, middlewarePermission middleware.IMiddlewarePermission) {
	h := &CourseHandler{Usecase: uc, validator: validator.New(), mwPageRequest: mwP, middlewareAuth: auth, middlewarePermission: middlewarePermission}

	// Public routes
	e.GET("/v1/course", h.GetIndex, h.mwPageRequest.PageRequestCtx)
	e.GET("/v1/course/:id", h.GetByID)

	// Protected routes
	r := e.Group("/v1/course")
	r.Use(h.middlewareAuth.AuthorizationCheck)

	permissionToCreate := []string{"course.create"}
	permissionToUpdate := []string{"course.update"}
	permissionToDelete := []string{"course.delete"}

	r.POST("", h.Create, middleware.RequireActivatedUser, h.middlewarePermission.PermissionValidation(permissionToCreate))
	r.PUT("/:id", h.Update, middleware.RequireActivatedUser, h.middlewarePermission.PermissionValidation(permissionToUpdate))
	r.DELETE("/:id", h.Delete, middleware.RequireActivatedUser, h.middlewarePermission.PermissionValidation(permissionToDelete))
}

func (h *CourseHandler) Create(c echo.Context) error {
	ctx := c.Request().Context()

	req := new(dto.ReqCreateCourse)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	thumbnailFile, _ := c.FormFile("thumbnail")

	var userID string
	if user := c.Get("user"); user != nil {
		if u, ok := user.(models.User); ok {
			userID = u.ID.String()
		}
	}
	res, err := h.Usecase.Create(ctx, req, userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	// If thumbnail file is provided, upload and update course with URL
	if thumbnailFile != nil {
		src, err := thumbnailFile.Open()
		if err != nil {
			return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
		}
		defer src.Close()
		fileData, err := io.ReadAll(src)
		if err != nil {
			return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
		}
		var buf bytes.Buffer
		buf.Write(fileData)
		destinatedPath := "courses/thumbnails/" + res.ID.String()
		url, err := utilsServices.UploadFile(buf, thumbnailFile.Filename, destinatedPath)
		if err != nil {
			return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, "Failed to upload thumbnail file"))
		}
		// Update course with uploaded thumbnail URL
		updateReq := &dto.ReqUpdateCourse{
			Title:            res.Title,
			Description:      res.Description,
			ShortDescription: res.ShortDescription,
			Price:            res.Price,
			DiscountRate:     res.DiscountRate,
			ThumbnailURL:     &url,
			LevelID:          req.LevelID,
			LangID:           req.LangID,
			TopicIDs:         req.TopicIDs,
		}
		res, err = h.Usecase.Update(ctx, res.ID.String(), updateReq, userID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
		}
	}

	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(dto.ToRespCourse(*res))
	return c.JSON(http.StatusOK, resp)
}

func (h *CourseHandler) Update(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	req := new(dto.ReqUpdateCourse)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	thumbnailFile, _ := c.FormFile("thumbnail")
	// If file provided, upload and set req.ThumbnailURL using same flow as avatar upload
	if thumbnailFile != nil {
		src, err := thumbnailFile.Open()
		if err != nil {
			return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
		}
		defer src.Close()
		fileData, err := io.ReadAll(src)
		if err != nil {
			return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
		}
		var buf bytes.Buffer
		buf.Write(fileData)
		destinatedPath := "courses/thumbnails/" + id
		url, err := utilsServices.UploadFile(buf, thumbnailFile.Filename, destinatedPath)
		if err != nil {
			return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, "Failed to upload thumbnail file"))
		}
		req.ThumbnailURL = &url
	}

	var userID string
	if user := c.Get("user"); user != nil {
		if u, ok := user.(models.User); ok {
			userID = u.ID.String()
		}
	}
	res, err := h.Usecase.Update(ctx, id, req, userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(dto.ToRespCourse(*res))
	return c.JSON(http.StatusOK, resp)
}

func (h *CourseHandler) Delete(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")
	if err := h.Usecase.Delete(ctx, id, ""); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(struct {
		Message string `json:"message"`
	}{Message: "Successfully deleted Course"})
	return c.JSON(http.StatusOK, resp)
}

func (h *CourseHandler) GetIndex(c echo.Context) error {
	ctx := c.Request().Context()
	pageRequest := c.Get("page_request").(*request.PageRequest)

	filter := new(dto.ReqCourseIndexFilter)
	if err := c.Bind(filter); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	if err := c.Validate(filter); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	res, total, err := h.Usecase.GetIndex(ctx, *pageRequest, *filter)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	respCourses := make([]dto.RespCourseIndex, 0, len(res))
	for _, v := range res {
		respCourses = append(respCourses, dto.ToRespCourseIndex(v))
	}
	respPag := response.PaginationResponse{}
	respPag, err = respPag.SetResponse(respCourses, total, pageRequest.PerPage, pageRequest.Page)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	return c.JSON(http.StatusOK, respPag)
}

func (h *CourseHandler) GetByID(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")
	res, err := h.Usecase.GetByID(ctx, id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(dto.ToRespCourse(*res))
	return c.JSON(http.StatusOK, resp)
}
