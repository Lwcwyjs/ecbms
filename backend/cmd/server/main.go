package main

import (
	"log"

	"ecbms/config"
	"ecbms/internal/handler"
	"ecbms/internal/middleware"
	"ecbms/internal/model"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	if err := model.InitDB(cfg.Database.Path); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer model.CloseDB()

	gin.SetMode(cfg.Server.Mode)
	r := gin.Default()

	r.Use(middleware.CORSMiddleware())

	api := r.Group("/api")
	{
		api.POST("/auth/login", handler.Login)

		auth := api.Group("")
		auth.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
		{
			auth.GET("/auth/me", handler.GetCurrentUser)
			auth.POST("/auth/change-password", handler.ChangePassword)

			auth.GET("/system/info", handler.GetSystemInfo)
			auth.GET("/system/stats", handler.GetSystemStats)
			auth.GET("/system/logs", handler.GetSystemLogs)
			auth.POST("/system/reboot", handler.RebootSystem)
			auth.POST("/system/shutdown", handler.ShutdownSystem)

			auth.GET("/network/interfaces", handler.GetNetworkInterfaces)
			auth.GET("/network/configs", handler.GetNetworkConfigs)
			auth.POST("/network/configure", handler.ConfigureNetwork)
			auth.POST("/network/apply", handler.ApplyNetworkConfig)
			auth.POST("/network/ping", handler.PingTest)
			auth.GET("/network/dns", handler.GetDNSConfig)

			admin := auth.Group("")
			admin.Use(middleware.AdminMiddleware())
			{
				admin.GET("/users", handler.GetUsers)
				admin.POST("/users", handler.CreateUser)
				admin.DELETE("/users/:id", handler.DeleteUser)
			}
		}
	}

	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
