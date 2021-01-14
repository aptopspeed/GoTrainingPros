package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"

	"github.com/pallat/todos/auth"
	"github.com/pallat/todos/captcha"
	"github.com/pallat/todos/logger"
	"github.com/pallat/todos/todos"
)

func main() {
	viper.SetDefault("app.addr", "0.0.0.0:8888")
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		log.Println("Warning:", err)
	}

	l, _ := zap.NewProduction()
	defer l.Sync()

	dsn := "sqlserver://beer:1234@localhost:1433?database=DTADB"
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	//db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	router := echo.New()
	router.Use(middleware.Logger())
	router.Use(middleware.Recover())
	router.Use(logger.Middleware(l))

	router.GET("/ping", func(c echo.Context) error {
		return c.String(200, "pong")
	})

	router.GET("/captcha", captchaHandler)
	router.POST("/exchange", exchangeHandler)

	router.POST("/todos", todos.NewNewTaskHandler(db))
	router.GET("/todos", todos.NewTaskOpenHadler(db))
	router.PUT("/todos/:id", todos.NewTaskEditHandler(db))
	router.DELETE("/todos/:id", todos.NewTaskDeleteHandler(db))

	srv := &http.Server{
		Addr:         viper.GetString("app.addr"),
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("serve on: %s\n", viper.GetString("app.addr"))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}

func captchaHandler(c echo.Context) error {
	key, captcha := captcha.KeyQuestion()
	return c.JSON(http.StatusOK, map[string]string{
		"key":     key,
		"captcha": captcha,
	})
}

func exchangeHandler(c echo.Context) error {
	var ans struct {
		Key    string `json:"key"`
		Answer int    `json:"answer"`
	}

	if err := c.Bind(&ans); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	if !captcha.Answer(ans.Key, ans.Answer) {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "wrong answer",
		})
	}

	t, err := auth.Token()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}
