package controllers

import (
	"net/http"
	"appengine"
	"appengine/datastore"
	"time"
	"github.com/parkn09/wymm/src/models"
	"fmt"
	"strconv"
	"encoding/json"
)

type PermissionsController struct{
	Controller
}

func NewPermissionsController() *PermissionsController{
	return &PermissionsController{
	}
}

func (controller *PermissionsController) Init(r *http.Request, w http.ResponseWriter){
	controller.R = r
	controller.W = w
	controller.Data = make(map[interface{}]interface{})
	controller.C = appengine.NewContext(r)
	controller.Layout = make([]string,0)
	controller.TemplateNames = "body"
	controller.Name = "videos"
	controller.DetermineRoutes()


	controller.Data["Language"] = controller.Locale
	controller.Layout = append(controller.Layout,"layouts/header.tmpl")
	controller.Layout = append(controller.Layout,"layouts/footer.tmpl")
	controller.Layout = append(controller.Layout,"layouts/layout.tmpl")
	if controller.CheckSession(){
		controller.Layout = append(controller.Layout,"layouts/logged.tmpl")
	}else{
				controller.Layout = append(controller.Layout,"layouts/nav.tmpl")
	}
	cookie, err := controller.R.Cookie("GolferLanguage")
	if err!= nil{
		controller.Locale = "en"
	}else if cookie.Value != ""{
		controller.Locale = cookie.Value
	}
	controller.Data["Language"] = controller.Locale
}

func (controller *PermissionsController) Serve(action string){
			if action == "home"{
				controller.Home()
			}else if action =="mobileadd"{
				controller.MobileAdd()
			}else if action =="mobileindex"{
				controller.MobileIndex()
			}else if action =="mobileedit"{
				controller.MobileEdit()
			}else if action =="index"{
				controller.Index()
			}else if action =="deleteall"{
				controller.DeleteAll()
			}else if action =="view"{
				controller.View()
			}else if action =="delete"{
				controller.Delete()
			}else if action =="home"{
				controller.Home()
			}else if action =="updateall"{
				controller.UpdateAll()
			}else if action =="pros"{
				controller.Pros()
			}else if action =="count"{
				controller.Count()
			}
}

func (controller *PermissionsController) MobileServe(action string){
	
}


func (controller *PermissionsController) Index(){

	c := appengine.NewContext(controller.R)
	if _,userKey, err := controller.GetUser(c); err != nil{
		controller.RenderError(err.Error())
	}else if userKey != nil{

		controller.Layout = append(controller.Layout,"videos/index.tmpl")
	}else{
	
	}
	controller.Render(true)
}

func (controller *PermissionsController) Edit(){
	controller.Render(true)
}

func (controller *PermissionsController) Delete(){
	//Adding a video from a phone / Get the video + images
	c := appengine.NewContext(controller.R)
	m := make(map[string]interface{})
	m["Result"] = "Failure"
	if !controller.CheckSession(){
		//Checks to see if this is from the web
		controller.UserToken = controller.R.FormValue("token")
	}
	if userKey, err := controller.GetUserKey(c); err != nil{
		m["Error"] = err.Error()
	}else{
		if videoKey,err := datastore.DecodeKey(controller.R.FormValue("video")); err != nil{
			m["Error"] = err.Error()
		}else if videoKey.Parent().Encode() == userKey.Encode(){
			video := models.Permission{}
			if err := datastore.Get(c,videoKey,&video); err != nil{
				m["Error"] = err.Error()
			}else{
				video.Active = false
				video.Modified = time.Now()
				if _, err := datastore.Put(c,videoKey,&video); err != nil{
					m["Error"] = err.Error()
				}else{
					m["Result"] = "Success"
					
				}
			}
		}
	}
	data, _ := json.Marshal(m)
	fmt.Fprintf(controller.W,"%s",data)
}

func (controller *PermissionsController) DeleteAll(){
	c := appengine.NewContext(controller.R)
	q := datastore.NewQuery("Permission").KeysOnly()
	keys, err := q.GetAll(c, nil)
	if err !=nil{
		println(err.Error())
	}
	err = datastore.DeleteMulti(c, keys)
	if err != nil{
		println(err.Error())
	}
	/*q = datastore.NewQuery("Subject").KeysOnly()
	keys, err = q.GetAll(c, nil)
	if err !=nil{
		println(err.Error())
	}
	err = datastore.DeleteMulti(c, keys)
	if err != nil{
		println(err.Error())
	}
	*/
}

func (controller *PermissionsController) View(){
	//http.ServeFile(controller.W,controller.R,"views/layout.html")
	c := appengine.NewContext(controller.R)
	if token,err := datastore.DecodeKey(controller.Params[2]); err !=nil{
		controller.RenderError(err.Error())
	}else{
		video := models.Permission{}
		if err := datastore.Get(c,token,&video); err !=nil{
			controller.RenderError(err.Error())
		}else{
			
		}
	}
	//
}

func (controller *PermissionsController) Home(){
	//http.ServeFile(controller.W,controller.R,"views/layout.html")
	//Select a random video and get the first image to beam down to the apps!
	c := appengine.NewContext(controller.R)
	m := make(map[string]interface{})
	m["Result"] = "Failure"
	videos := make([]models.Permission,0,1)
	q := datastore.NewQuery("Permission").Filter("Active =",true).Filter("Category =","Pro").Order("-Views").Limit(1)
	if keys, err := q.GetAll(c,&videos); err !=nil{
		m["Error"] = err.Error()
	}else if len(keys)>0{
		
		
	}
	data, _ := json.Marshal(m)
	fmt.Fprintf(controller.W,"%s",data)
}

func (controller *PermissionsController) Pros(){
	c := appengine.NewContext(controller.R)
	m := make(map[string]interface{})
	m["Result"] = "Failure"
	m["Error"] = "No authorization"
	if !controller.CheckSession(){
		//Checks to see if this is from the web
		controller.UserToken = controller.R.FormValue("token")
	}
		if _,userKey,err := controller.GetUser(c); err !=nil{
			m["Error"] = "Authorization Error"
		}else if userKey != nil{
			var offset int
			if offset,err = strconv.Atoi(controller.R.FormValue("offset")); err !=nil{
				offset = 0
			}
			q := datastore.NewQuery("Permission").Filter("Active =",true).Filter("Category =","Pro").Order("-Created").Offset(offset).Limit(10)
			videos := make([]models.Permission,0,10)
			if keys,err := q.GetAll(c,&videos); err !=nil{
				m["Error"] = err.Error()
			}else{

				m["Result"] = "Success"
				m["Permissions"] = videos
				m["SafeKeys"] = keys
			}
		}
	data, _ := json.Marshal(m)
	fmt.Fprintf(controller.W,"%s",data)
}

//AJAX only 
func (controller *PermissionsController) MobileIndex(){
	//Adding a video from a phone / Get the video + images
	c := appengine.NewContext(controller.R)
	m := make(map[string]interface{})
	m["Result"] = "Failure"
	if !controller.CheckSession(){
		//Checks to see if this is from the web
		controller.UserToken = controller.R.FormValue("token")
	}
		if _,userKey,err := controller.GetUser(c); err !=nil{
			m["Error"] = "Authorization Error"
		}else{
			offset,_ := strconv.Atoi(controller.R.FormValue("offset"))
			videos := make([]models.Permission,0,10)

			q := datastore.NewQuery("Permission").Ancestor(userKey).Filter("Active =",true).Order("-Created").Offset(offset).Limit(10)
			if keys, err := q.GetAll(c,&videos); err != nil{
				m["Error"] = err.Error()
			}else{
				m["SafeKeys"] = keys
				m["Result"] = "Success"
				m["Permissions"] = videos
			}	
		}
	data, _ := json.Marshal(m)
	fmt.Fprintf(controller.W,"%s",data)
}

//Mobile add - Upload from phones
func (controller *PermissionsController) MobileAdd(){
	//Adding a video from a phone / Get the video + images
	c := appengine.NewContext(controller.R)
	m := make(map[string]interface{})
	m["Result"] = "Failure"
	m["Error"] = "No authorization"
	if !controller.CheckSession(){
		//Checks to see if this is from the web
		controller.UserToken = controller.R.FormValue("token")
	}
		if _,userKey,err := controller.GetUser(c); err !=nil{
			m["Error"] = "Authorization Error"
		}else if userKey != nil{
			
		}
	data, _ := json.Marshal(m)
	fmt.Fprintf(controller.W,"%s",data)
}

//Mobile add - Upload from phones
func (controller *PermissionsController) MobileEdit(){
	//Adding a video from a phone / Get the video + images
	c := appengine.NewContext(controller.R)
	m := make(map[string]interface{})
	m["Result"] = "Failure"
	m["Error"] = "No authorization"
	if !controller.CheckSession(){
		//Checks to see if this is from the web
		controller.UserToken = controller.R.FormValue("token")
	}
		if _,userKey,err := controller.GetUser(c); err !=nil{
			m["Error"] = "Authorization Error"
		}else if userKey != nil{
			
		}
	data, _ := json.Marshal(m)
	fmt.Fprintf(controller.W,"%s",data)
}

func (controller *PermissionsController) Count(){
	c := appengine.NewContext(controller.R)
	m := make(map[string]interface{})
	m["Result"] = "Failure"
	if !controller.CheckSession(){
		//Checks to see if this is from the web
		controller.UserToken = controller.R.FormValue("token")
	}
		if _,_,err := controller.GetUser(c); err !=nil{
			m["Error"] = "Authorization Error"
		}else{
			//Decode video key
			if videoKey, err := datastore.DecodeKey(controller.R.FormValue("video")); err != nil{
				m["Error"] = err.Error()
			}else{
				println(videoKey.Encode())
			}	
		}	
}


func (controller *PermissionsController) UpdateAll(){
	c := appengine.NewContext(controller.R)
	m := make(map[string]interface{})
	m["Result"] = "Failure"
	m["Error"] = "No authorization"
	q := datastore.NewQuery("Permission").Limit(100)
	assignments := make([]models.Permission,0,100)
	if keys, err := q.GetAll(c,&assignments); err != nil{
		println(err.Error())
	}else{
		
		datastore.PutMulti(c,keys,assignments)
	}
	data, _ := json.Marshal(m)
	fmt.Fprintf(controller.W,"%s",data)
}

func (controller *PermissionsController) Render(show bool) error{
	//controller.Controller.Render(show)
	controller.Controller.Render(show)
	return nil
}