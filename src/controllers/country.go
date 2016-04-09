package controllers

import (
	"net/http"
	"appengine"
	"appengine/datastore"
	"github.com/parkn09/wymm/src/models"
	"fmt"
	"strconv"
	"encoding/json"
	"time"
)

type CountriesController struct{
	Controller
}

func NewCountriesController() *CountriesController{
	return &CountriesController{
	}
}

func (controller *CountriesController) Init(r *http.Request, w http.ResponseWriter){
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

func (controller *CountriesController) Serve(action string){
			if action == "home"{
				controller.Home()
			}else if action =="mobileadd"{
				controller.MobileAdd()
			}else if action =="mobileindex"{
				controller.MobileIndex()
			}else if action =="mobileedit"{
				controller.MobileEdit()
			}else if action =="mobileview"{
				controller.MobileView()
			}else if action =="index"{
				controller.Index()
			}else if action =="deleteall"{
				controller.DeleteAll()
			}else if action =="view"{
				controller.View()
			}else if action =="delete"{
				controller.Delete()
			}else if action =="updateall"{
				controller.UpdateAll()
			}
}

func (controller *CountriesController) MobileServe(action string){
	
}


func (controller *CountriesController) Index(){

	c := appengine.NewContext(controller.R)
	if user,userKey, err := controller.GetUser(c); err != nil{
		controller.RenderError(err.Error())
	}else if userKey != nil && user.Group ==1{
		controller.Layout = append(controller.Layout,"countries/index.tmpl")
	}else{
		controller.RenderError("You don't have authorization")
	}
	controller.Render(true)
}

func (controller *CountriesController) Edit(){
	controller.Render(true)
}

func (controller *CountriesController) Delete(){
	//Adding a video from a phone / Get the video + images
	c := appengine.NewContext(controller.R)
	m := make(map[string]interface{})
	m["Result"] = "Failure"
	if !controller.CheckSession(){
		//Checks to see if this is from the web
		controller.UserToken = controller.R.FormValue("token")
	}
	if user,userKey, err := controller.GetUser(c); err != nil{
		m["Error"] = err.Error()
	}else{
		if videoKey,err := datastore.DecodeKey(controller.R.FormValue("safekey")); err != nil{
			m["Error"] = err.Error()
		}else if videoKey!=nil && user.Group ==1 && userKey != nil{
			video := models.State{}
			if err := datastore.Get(c,videoKey,&video); err != nil{
				m["Error"] = err.Error()
			}else{
				video.Active = false
				if _, err := datastore.Put(c,videoKey,&video); err != nil{
					m["Error"] = err.Error()
				}else{
					m["Result"] = "Success"
					//Go ahead and deactivate the images, slate them for deletion in two days
				}
			}
		}
	}
	data, _ := json.Marshal(m)
	fmt.Fprintf(controller.W,"%s",data)
}

func (controller *CountriesController) DeleteAll(){
	c := appengine.NewContext(controller.R)
	q := datastore.NewQuery("Country").KeysOnly()
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

func (controller *CountriesController) View(){
	//http.ServeFile(controller.W,controller.R,"views/layout.html")
	//Adding a video from a phone / Get the video + images
	c := appengine.NewContext(controller.R)
	m := make(map[string]interface{})
	m["Result"] = "Failure"
	if !controller.CheckSession(){
		//Checks to see if this is from the web
		controller.UserToken = controller.R.FormValue("token")
	}
	if user,userKey, err := controller.GetUser(c); err != nil{
		m["Error"] = err.Error()
	}else if userKey != nil && user.Group == 1{
		if token,err := datastore.DecodeKey(controller.Params[2]); err !=nil{
			controller.RenderError(err.Error())
		}else{
			country := models.Country{}
			if err := datastore.Get(c,token,&country); err !=nil{
				controller.RenderError(err.Error())
			}else{
				controller.Data["Country"] = country
				controller.Data["SafeKey"] = controller.Params[2]
			}
		}
	}else{
		controller.RenderError("You aren't authorized")
	}
	
	controller.Layout = append(controller.Layout,"countries/view.tmpl")
	controller.Render(true)
	//
}

func (controller *CountriesController) Home(){
	//http.ServeFile(controller.W,controller.R,"views/layout.html")
	
	controller.Layout = append(controller.Layout,"general/online.tmpl")
	controller.Render(true)
	//
}

//AJAX only 
func (controller *CountriesController) MobileIndex(){
	//Adding a video from a phone / Get the video + images
	c := appengine.NewContext(controller.R)
	m := make(map[string]interface{})
	m["Result"] = "Failure"

			offset,_ := strconv.Atoi(controller.R.FormValue("offset"))
			videos := make([]models.Country,0,10)

			q := datastore.NewQuery("Country").Filter("Active =",true).Order("Name").Offset(offset).Limit(50)
			if keys, err := q.GetAll(c,&videos); err != nil{
				m["Error"] = err.Error()
			}else{
				m["Result"] = "Success"
				m["Countries"] = videos
				m["SafeKeys"] = keys
			}	
	data, _ := json.Marshal(m)
	fmt.Fprintf(controller.W,"%s",data)
}

//Mobile add - Upload from phones
func (controller *CountriesController) MobileAdd(){
	//Adding a video from a phone / Get the video + images
	c := appengine.NewContext(controller.R)
	m := make(map[string]interface{})
	m["Result"] = "Failure"
	m["Error"] = "No authorization"
	if !controller.CheckSession(){
		//Checks to see if this is from the web
		controller.UserToken = controller.R.FormValue("token")
	}
		if user,userKey,err := controller.GetUser(c); err !=nil{
			m["Error"] = "Authorization Error"
		}else if userKey != nil && user.Group ==1{
			
			countries:= make([]models.Country,0,1)
			q := datastore.NewQuery("Country").Filter("Name =",controller.R.FormValue("name")).Filter("Active =",true).Limit(1)
			if keys,err := q.GetAll(c,&countries); err != nil{
				m["Error"] = err.Error()
			}else if len(keys)>0{
				country := countries[0]
					m["Result"] = "Success"
					m["Country"] = country
					m["SafeKey"] = keys[0].Encode()
			}else{
				country := models.Country{
					Name: controller.R.FormValue("name"),
					Nickname:controller.R.FormValue("nickname"),
					Created:time.Now(),
					Active:true,
				}
				vKey := datastore.NewIncompleteKey(c,"Country",userKey)
				if countryKey, err := datastore.Put(c,vKey,&country); err != nil{
					m["Error"] = err.Error()
				}else{
					m["Result"] = "Success"
					m["SafeKey"] = countryKey.Encode()
				}
			}	
		}
	data, _ := json.Marshal(m)
	fmt.Fprintf(controller.W,"%s",data)
}



func (controller *CountriesController) MobileView(){
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
			countries := make([]models.Country,0,1)
			q := datastore.NewQuery("Country").Ancestor(userKey).Limit(1)
			if keys,err := q.GetAll(c,&countries); err != nil{
				m["Error"] = err.Error()
			}else if len(keys)>0{
				m["Result"] = "Success"
				m["Country"] = countries[0]
			}else{
				m["Result"] = "Partial"
			}
		}
	data, _ := json.Marshal(m)
	fmt.Fprintf(controller.W,"%s",data)
}

//Mobile add - Upload from phones
func (controller *CountriesController) MobileEdit(){
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
			countries := make([]models.Country,0,1)
			q := datastore.NewQuery("Country").Ancestor(userKey).Limit(1)
			if keys,err := q.GetAll(c,&countries); err != nil{
				m["Error"] = err.Error()+" Step 3"
			}else if len(keys)>0{
					video := countries[0]
						video.Name = controller.R.FormValue("name")
						if _,err = datastore.Put(c,keys[0],&video); err !=nil{
							m["Error"] = err.Error()+"You're in this arena"
						}else{
							m["Result"] = "Success"
						}
			}
		}
	data, _ := json.Marshal(m)
	fmt.Fprintf(controller.W,"%s",data)
}



func (controller *CountriesController) UpdateAll(){
	/*
	c := appengine.NewContext(controller.R)
	m := make(map[string]interface{})
	m["Result"] = "Failure"
	m["Error"] = "No authorization"
	q := datastore.NewQuery("Video").Limit(100)
	assignments := make([]models.Country,0,100)
	if keys, err := q.GetAll(c,&assignments); err != nil{
		println(err.Error())
	}else{
		for i,_ := range keys{
			assignments[i].Modified = time.Now()
		}
		datastore.PutMulti(c,keys,assignments)
	}
	data, _ := json.Marshal(m)
	fmt.Fprintf(controller.W,"%s",data)
	/*
	c := appengine.NewContext(controller.R)
	q := datastore.NewQuery("Country").Limit(100)
	assignments := make([]models.Country,0,100)
	if keys, err := q.GetAll(c,&assignments); err != nil{
		println(err.Error())
	}else{
		for i,_ := range keys{
			assignments[i].Id = 0
		}
		datastore.PutMulti(c,keys,assignments)
	}
	*/
	/*
	m := make(map[string]interface{})
	c := appengine.NewContext(controller.R)
	m["Result"] = "Failure"
	videos := make([]models.Video,0,1)
	q := datastore.NewQuery("Video").Order("-Created").Limit(1)
	if videoKeys, err := q.GetAll(c,&videos); err != nil{
		m["Error"] = err.Error()
	}else if len(videoKeys)>0{
		videoKey := videoKeys[0]

		q = datastore.NewQuery("VideoImage").Ancestor(videoKey).Order("Created").Limit(300).KeysOnly()
		if keys, err := q.GetAll(c,nil); err !=nil{
			m["Error"] = err.Error()+" In finding shit"
		}else{
			for i,key := range keys{
				if i % 2 == 0{
					if err = datastore.Delete(c,key); err !=nil{
						m["Error"] = err.Error()+" in delete"
					}
				}
			}
		}
	}
	data, _ := json.Marshal(m)
	fmt.Fprintf(controller.W,"%s",data)
	*/
}



func (controller *CountriesController) Render(show bool) error{
	//controller.Controller.Render(show)
	controller.Controller.Render(show)
	return nil
}