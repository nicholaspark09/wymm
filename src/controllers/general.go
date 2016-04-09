package controllers

import (
	"net/http"
	"appengine"
)

type GeneralController struct{
	Controller
}

func NewGeneralController() *GeneralController{
	return &GeneralController{
	}
}

func (controller *GeneralController) Init(r *http.Request, w http.ResponseWriter){
	controller.R = r
	controller.W = w
	controller.Data = make(map[interface{}]interface{})
	controller.C = appengine.NewContext(r)
	controller.Layout = make([]string,0)
	controller.TemplateNames = "body"
	controller.Name = "users"
	controller.DetermineRoutes()


	controller.Data["Language"] = controller.Locale
	controller.Layout = append(controller.Layout,"layouts/header.tmpl")
	controller.Layout = append(controller.Layout,"layouts/footer.tmpl")
	if controller.CheckSession(){
		controller.Layout = append(controller.Layout,"layouts/logged.tmpl")
	}else{
				controller.Layout = append(controller.Layout,"layouts/nav.tmpl")
	}
	cookie, err := controller.R.Cookie("QuicklishLanguage")
	if err!= nil{
		controller.Locale = "en"
	}else if cookie.Value != ""{
		controller.Locale = cookie.Value
	}
	controller.Data["Language"] = controller.Locale
}

func (controller *GeneralController) Serve(action string){
			if action == "home"{
				controller.Home()
			}
}

func (controller *GeneralController) MobileServe(action string){
	
}


func (controller *GeneralController) Index(){
	controller.Render(true)
}

func (controller *GeneralController) Edit(){
	controller.Render(true)
}

func (controller *GeneralController) Delete(){

}

func (controller *GeneralController) View(){
	//http.ServeFile(controller.W,controller.R,"views/layout.html")
	

	//
}

func (controller *GeneralController) Home(){
	//http.ServeFile(controller.W,controller.R,"views/layout.html")
	
	controller.Layout = append(controller.Layout,"general/online.tmpl")
	controller.Render(true)
	//
}

/*
		!!! This is for apps only !!!
			All web logins must go to the regular Login() method 

*/
func (controller *GeneralController) MobileAdd(){
	
}

//This is the final step to login an Android App User
func (controller *GeneralController) CompleteApp(){
	
}



func (controller *GeneralController) Render(show bool) error{
	//controller.Controller.Render(show)
	controller.Controller.Render(show)
	return nil
}