package news

import (
	"arctfrex-customers/internal/base"
	"arctfrex-customers/internal/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

type newsHandler struct {
	jwtMiddleware *middleware.JWTMiddleware
	newsUsecase   NewsUsecase
}

func NewNewsHandler(
	engine *gin.Engine,
	jmw *middleware.JWTMiddleware,
	nu NewsUsecase,
) *newsHandler {
	handler := &newsHandler{
		jwtMiddleware: jmw,
		newsUsecase:   nu,
	}

	unprotectedRoutesBackOffice := engine.Group("/backoffice/news")
	unprotectedRoutes, protectedRoutes := engine.Group("/news"), engine.Group("/news")

	unprotectedRoutesBackOffice.POST("", handler.BackOfficeCreateNews)
	unprotectedRoutesBackOffice.GET("", handler.BackOfficeGetNews)
	unprotectedRoutesBackOffice.PATCH("", handler.BackOfficeUpdateNews)
	unprotectedRoutesBackOffice.DELETE("", handler.BackOfficeDeactiveNews)

	// unprotectedRoutesBackOffice.POST("/", handler.BackOfficeCreateNews)
	// unprotectedRoutesBackOffice.GET("/", handler.BackOfficeGetNews)
	// unprotectedRoutesBackOffice.PATCH("/", handler.BackOfficeUpdateNews)
	// unprotectedRoutesBackOffice.DELETE("/", handler.BackOfficeDeactiveNews)

	unprotectedRoutes.GET("/latest", handler.GetNewsLatest)
	unprotectedRoutes.GET("/latest/bulletin", handler.GetNewsLatestBulletin)
	protectedRoutes.Use(jmw.ValidateToken())
	{
	}

	return handler
}

func (nh *newsHandler) BackOfficeCreateNews(c *gin.Context) {
	var news *News
	if err := c.ShouldBindJSON(&news); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := nh.newsUsecase.CreateNews(news)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, base.ApiResponse{
		Message: "Success",
		Data:    news,
	})
}

func (nh *newsHandler) BackOfficeGetNews(c *gin.Context) {
	title, hasTitleQuery := c.GetQuery("title")

	// If title query is not provided, return all news
	if !hasTitleQuery {
		news, err := nh.newsUsecase.GetNewsList()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, base.ApiResponse{
			Message: "Success",
			Data:    news,
		})
		return
	}

	// If title query is provided, return news by title
	news, err := nh.newsUsecase.GetNewsByTitle(title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, base.ApiResponse{
		Message: "Success",
		Data:    news,
	})
}

func (nh *newsHandler) BackOfficeUpdateNews(c *gin.Context) {
	title := c.Query("title")

	var news *News
	if err := c.ShouldBindJSON(&news); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := nh.newsUsecase.UpdateNews(title, news)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	newNews, err := nh.newsUsecase.GetNewsByTitle(news.NewsTitle)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, base.ApiResponse{
		Message: "Success",
		Data:    newNews,
	})
}

func (nh *newsHandler) BackOfficeDeactiveNews(c *gin.Context) {
	title := c.Query("title")

	err := nh.newsUsecase.DeactiveNews(title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, base.ApiResponse{
		Message: "Success",
	})
}

func (nh *newsHandler) GetNewsLatest(c *gin.Context) {
	news, err := nh.newsUsecase.GetNewsLatest()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, base.ApiResponse{
		Message: "Success",
		Data:    news,
	})
}

func (nh *newsHandler) GetNewsLatestBulletin(c *gin.Context) {
	newsBulletin, err := nh.newsUsecase.GetNewsLatestBulletin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, base.ApiResponse{
		Message: "Success",
		Data:    newsBulletin,
	})
}
