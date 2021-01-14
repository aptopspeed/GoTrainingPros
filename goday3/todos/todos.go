package todos

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/pallat/todos/logger"
)

//CREATE DB todos function
func NewNewTaskHandler(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		db.AutoMigrate(Task{})

		var todo struct {
			Task string `json:"task"`
		}

		logger := logger.Extract(c)
		logger.Info("new task todo........")

		if err := c.Bind(&todo); err != nil {
			log.Error(err.Error())
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": errors.Wrap(err, "new task").Error(),
			})
		}

		if err := db.Create(&Task{
			Task: todo.Task,
		}).Error; err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"err": errors.Wrap(err, "create task").Error(),
			})
		}

		return c.JSON(http.StatusOK, map[string]string{})
	}
}

// Get function Todos
func NewTaskOpenHadler(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {

		db.AutoMigrate(Task{})

		logger := logger.Extract(c)
		logger.Info("get task todo........")

		var task []Task

		if err := db.Find(&task).Error; err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"err": errors.Wrap(err, "get task").Error(),
			})
		}

		return c.JSON(http.StatusOK, task)
	}
}

//PUT function todos
func NewTaskEditHandler(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {

		db.AutoMigrate(Task{})

		logger := logger.Extract(c)
		logger.Info("edit task todo......")

		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"err": errors.Wrap(err, "get task").Error(),
			})
		}

		var task Task

		if err := db.Model(&task).Where("id = ?", id).Update("processed", true).Error; err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"err": errors.Wrap(err, "edit task").Error(),
			})
		}
		return c.JSON(http.StatusOK, map[string]string{})

	}
}

//DELETE function todos
func NewTaskDeleteHandler(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		db.AutoMigrate(Task{})

		logger := logger.Extract(c)
		logger.Info("delete task todo........ %")

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"err": errors.Wrap(err, "delete task").Error(),
			})
		}

		var task Task

		if err := db.Where("id = ?", id).Delete(&task).Error; err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"err": errors.Wrap(err, "delete task").Error(),
			})
		}

		return c.JSON(http.StatusOK, map[string]string{})
	}
}

type Task struct {
	gorm.Model
	Task      string
	Processed bool
}

func (Task) TableName() string {
	return "todos"
}
