package router

import (
	"medcard/beck/controllers"
	"os"

	"github.com/gin-gonic/gin"
)

func Routers() {
	//"""""""""""" controllers""""""""""""
	cors := controllers.Cors
	login := controllers.Login
	signupDocktors := controllers.SignupDr
	signupClient := controllers.SignupCl
	addQuestion := controllers.AddQuestion
	adminProfileChange := controllers.AdminProfileChange
	removeQuestion := controllers.RemoveQuestion
	statistics := controllers.Statistics
	users_clients := controllers.Users_clients
	questions_get := controllers.Questions_get
	logout := controllers.Logout
	doctors_get := controllers.Doctors_get
	clients_get := controllers.Clients_get
	doctorProfileChange := controllers.DoctorProfileChange
	client_prof_change := controllers.Client_prof_change
	accept_decline := controllers.Accept_decline
	signup_cl_view := controllers.Signup_cl_view
	views_get_dr := controllers.Views_get_dr
	views_get_cl := controllers.Views_get_cl
	emc_get := controllers.Emc_get
	// """"""""""""""routers""""""""""""""
	r := gin.Default()
	r.StaticFS("/static", gin.Dir("./static", true))
	r.Use(cors)
	// """"""""""""""""""""""""""""POST requests""""""""""""""""""""""""""""
	// Login && password
	r.POST("/login", login)                                   //Done
	//
	r.GET("/logout", logout)                                 //Done
	// DoctorLog "struct is used for"
	r.POST("/signup_dr", signupDocktors)                      //Done
	// ClientLog "struct is used for"
	r.POST("/signup_cl", signupClient)                        //Done
	// ViewReq "struct is used for"
	r.POST("/signup_cl_view", signup_cl_view)                 //Done
	// FrequentlyAskedQuestion "struct is used for"
	r.POST("/question_add", addQuestion)                      //Done
	// AdminLog "struct is used for"
	r.POST("/admin_prof_change", adminProfileChange)          //Done
	// DoctorLog "struct is used for"
	r.POST("/doctor_prof_change", doctorProfileChange)        //Done
	// ClientLog "struct is used for"
	r.POST("/client_prof_change", client_prof_change)         //Done
	// Accept_Decline "struct is used for"
	r.POST("/accept_decline", accept_decline)
											// """""""""""""""""""GET requests"""""""""""""""""""
	r.GET("/statistics", statistics)                          //Done
	r.GET("/users_clients", users_clients)                    //Done
	r.GET("/clients_get", clients_get)                        //Done
	r.GET("/views_get_dr", views_get_dr)                      //Done
	r.GET("/views_get_cl", views_get_cl)                      //Done
	r.GET("/emc_get", emc_get)                                //Done
	r.GET("/questions_get", questions_get)                    //Done
	r.GET("/doctors_get", doctors_get)                        //Done
	r.DELETE("/question_rm", removeQuestion)                  

	r.Run(":"+os.Getenv("PORT"))

}
