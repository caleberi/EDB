package engine

import (
	"yc-backend/common"
	"yc-backend/controllers"
)

func (srv *Application) RegisterRoute() *Application {
	r := srv.mux
	r.RemoveExtraSlash = true
	r.RedirectFixedPath = false
	r.RedirectTrailingSlash = false

	r.GET("/ping", controllers.PingDb)
	r.POST("/webhook/yellow-card", controllers.YellowCardWebHook)

	authorizedRouter := r.Group("/auth")
	{
		authorizedRouter.POST("/register", controllers.RegisterUser)
		authorizedRouter.POST("/login", controllers.LoginUser)

		authorizedRouter.Use(common.AuthorizeUser()).
			POST("/logout", controllers.LogoutUser)
	}

	managementRouter := r.Group("/employee")
	managementRouter.Use(common.AuthorizeUser())
	{
		managementRouter.POST("/", (controllers.AddEmployee))
		managementRouter.PUT("/:employeeId", (controllers.UpdateEmployee))
		managementRouter.DELETE("/:employeeId", (controllers.DeleteEmployee))
	}

	// include admin route check here
	disbursementRouter := r.Group("/disbursements")
	disbursementRouter.Use(common.AuthorizeUser())
	{
		disbursementRouter.POST("/:employeeId", (controllers.MakeDisbursmentToEmployee))
	}

	return srv
}
