package http

import (
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rendyfutsuy/base-go/helpers/middleware"
	_reqContext "github.com/rendyfutsuy/base-go/helpers/middleware/request"
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/helpers/response"
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/course"
	"github.com/rendyfutsuy/base-go/modules/course/dto"
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

// Create Course
// @Summary      Create course
// @Description  Create a course with optional thumbnail upload
// @Tags         Course
// @Accept       json
// @Accept       multipart/form-data
// @Produce      json
// @Security     BearerAuth
// @Param        payload   body     dto.ReqCreateCourse  true  "Course payload"
// @Param        thumbnail formData file                 false "Thumbnail file"
// @Success      200       {object} response.NonPaginationResponse{data=dto.RespCourse}
// @Failure      400       {object} response.NonPaginationResponse
// @Router       /v1/course [post]
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
	var thumbnailData []byte
	var thumbnailName string
	if thumbnailFile != nil {
		src, err := thumbnailFile.Open()
		if err != nil {
			return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
		}
		defer src.Close()
		thumbnailData, err = io.ReadAll(src)
		if err != nil {
			return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
		}
		thumbnailName = thumbnailFile.Filename
	}

	var userID string
	if user := c.Get("user"); user != nil {
		if u, ok := user.(models.User); ok {
			userID = u.ID.String()
		}
	}
	res, err := h.Usecase.Create(ctx, req, userID, thumbnailData, thumbnailName)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(dto.ToRespCourse(*res))
	return c.JSON(http.StatusOK, resp)
}

// Update Course
// @Summary      Update course
// @Description  Update a course with optional thumbnail upload
// @Tags         Course
// @Accept       json
// @Accept       multipart/form-data
// @Produce      json
// @Security     BearerAuth
// @Param        id        path     string               true  "Course ID"
// @Param        payload   body     dto.ReqUpdateCourse  true  "Course payload"
// @Param        thumbnail formData file                 false "Thumbnail file"
// @Success      200       {object} response.NonPaginationResponse{data=dto.RespCourse}
// @Failure      400       {object} response.NonPaginationResponse
// @Router       /v1/course/{id} [put]
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
	var thumbnailData []byte
	var thumbnailName string
	if thumbnailFile != nil {
		src, err := thumbnailFile.Open()
		if err != nil {
			return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
		}
		defer src.Close()
		thumbnailData, err = io.ReadAll(src)
		if err != nil {
			return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
		}
		thumbnailName = thumbnailFile.Filename
	}

	var userID string
	if user := c.Get("user"); user != nil {
		if u, ok := user.(models.User); ok {
			userID = u.ID.String()
		}
	}
	res, err := h.Usecase.Update(ctx, id, req, userID, thumbnailData, thumbnailName)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(dto.ToRespCourse(*res))
	return c.JSON(http.StatusOK, resp)
}

// Delete Course
// @Summary      Delete course
// @Description  Delete a course by ID
// @Tags         Course
// @Produce      json
// @Security     BearerAuth
// @Param        id   path string true "Course ID"
// @Success      200  {object} response.NonPaginationResponse
// @Failure      400  {object} response.NonPaginationResponse
// @Router       /v1/course/{id} [delete]
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

// Get Courses
// @Summary      Get paginated list of courses
// @Description  Retrieve a paginated list of courses with optional filtering
// @Tags         Course
// @Accept       json
// @Produce      json
// @Param        page        query   int                       false "Page number"     default(1)
// @Param        per_page    query   int                       false "Items per page"  default(10)
// @Param        sort_by     query   string                    false "Sort column"
// @Param        sort_order  query   string                    false "Sort order (asc/desc)"
// @Param        search      query   string                    false "Search query"
// @Param        filter      query   dto.ReqCourseIndexFilter  false "Filter options"
// @Success      200         {object} response.PaginationResponse{data=[]dto.RespCourseIndex}
// @Failure      400         {object} response.NonPaginationResponse
// @Router       /v1/course [get]
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

// Get Course By ID
// @Summary      Get course by ID
// @Description  Retrieve course detail and its parameters
// @Tags         Course
// @Produce      json
// @Param        id   path string true "Course ID"
// @Success      200  {object} response.NonPaginationResponse{data=dto.RespCourse}
// @Failure      400  {object} response.NonPaginationResponse
// @Router       /v1/course/{id} [get]
func (h *CourseHandler) GetByID(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")
	res, err := h.Usecase.GetByID(ctx, id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	level, lang, topics, err := h.Usecase.GetParameterReferences(ctx, id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	out := dto.ToRespCourse(*res)
	out.Level = level
	out.Lang = lang
	out.Topics = topics
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(out)
	return c.JSON(http.StatusOK, resp)
}
