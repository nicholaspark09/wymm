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
	"strings"
)

type CitiesController struct{
	Controller
}

func NewCitiesController() *CitiesController{
	return &CitiesController{
	}
}

func (controller *CitiesController) Init(r *http.Request, w http.ResponseWriter){
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

func (controller *CitiesController) Serve(action string){
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
			}else if action =="search"{
				controller.Search()
			}
}

func (controller *CitiesController) MobileServe(action string){
	
}


func (controller *CitiesController) Index(){

	c := appengine.NewContext(controller.R)
	if user,userKey, err := controller.GetUser(c); err != nil{
		controller.RenderError(err.Error())
	}else if userKey != nil && user.Group ==1{

		

		
		controller.Layout = append(controller.Layout,"cities/index.tmpl")
	}else{
		controller.RenderError("You don't have authorization")
	}
	controller.Render(true)
}

func (controller *CitiesController) Edit(){
	controller.Render(true)
}

func (controller *CitiesController) Delete(){
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
			if err = datastore.Delete(c,videoKey); err != nil{
				m["Error"] = err.Error()
			}else{
				m["Result"] = "Success"
			}
		}
	}
	data, _ := json.Marshal(m)
	fmt.Fprintf(controller.W,"%s",data)
}

func (controller *CitiesController) DeleteAll(){
	c := appengine.NewContext(controller.R)
	q := datastore.NewQuery("City").KeysOnly()
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

func (controller *CitiesController) View(){
	//http.ServeFile(controller.W,controller.R,"views/layout.html")
	c := appengine.NewContext(controller.R)
	if token,err := datastore.DecodeKey(controller.Params[2]); err !=nil{
		controller.RenderError(err.Error())
	}else{
		video := models.City{}
		if err := datastore.Get(c,token,&video); err !=nil{
			controller.RenderError(err.Error())
		}else{
			
		}
	}
	//
}

func (controller *CitiesController) Home(){
	//http.ServeFile(controller.W,controller.R,"views/layout.html")
	
	controller.Layout = append(controller.Layout,"general/online.tmpl")
	controller.Render(true)
	//
}

func (controller *CitiesController) Search(){
	c := appengine.NewContext(controller.R)
	m := make(map[string]interface{})
	m["Result"] = "Failure"
	if !controller.CheckSession(){
		//Checks to see if this is from the web
		controller.UserToken = controller.R.FormValue("token")
	}
		if _,userKey,err := controller.GetUser(c); err !=nil{
			m["Error"] = "Authorization Error"
		}else if userKey != nil{
			search := strings.ToLower(controller.R.FormValue("query"))
			cities := make([]models.City,0,20)
			q := datastore.NewQuery("City").Filter("Name >=",search).Filter("Active =",true).Project("Name","Country").Limit(20)
			if keys, err := q.GetAll(c,&cities); err !=nil{
				m["Error"] = err.Error()
			}else{
				for i,key := range keys{
					//safekeys[i] = key.Encode()
					//Get the state as well
					stateKey := key.Parent()
					state := models.State{}
					if err = datastore.Get(c,stateKey,&state); err == nil{
						cities[i].Name = strings.Title(cities[i].Name)+", "+strings.Title(state.Name)+", "+cities[i].Country
						cities[i].Country = key.Encode()
					}
				}
				m["Result"] = "Success"
				m["Cities"] = cities
			}
		}
	data, _ := json.Marshal(m)
	fmt.Fprintf(controller.W,"%s",data)
}

//AJAX only 
func (controller *CitiesController) MobileIndex(){
	c := appengine.NewContext(controller.R)
	m := make(map[string]interface{})
	m["Result"] = "Failure"
	if countryKey, err := datastore.DecodeKey(controller.R.FormValue("safekey")); err !=nil{
		m["Error"] = err.Error()
	}else{
			offset,_ := strconv.Atoi(controller.R.FormValue("offset"))
			cities := make([]models.City,0,50)
			q := datastore.NewQuery("City").Ancestor(countryKey).Filter("Active =",true).Order("Name").Offset(offset).Limit(10)
			if keys, err := q.GetAll(c,&cities); err != nil{
				m["Error"] = err.Error()
			}else{
				m["Result"] = "Success"
				safekeys := make([]string,len(cities))
				for i,key := range keys{
					safekeys[i] = key.Encode()
					cities[i].Name = strings.Title(cities[i].Name)
				}
				m["SafeKeys"] = safekeys
								m["Cities"] = cities
			}	
	}
	data, _ := json.Marshal(m)
	fmt.Fprintf(controller.W,"%s",data)
}

//Mobile add - Upload from phones
func (controller *CitiesController) MobileAdd(){
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
			if safekey, err := datastore.DecodeKey(controller.R.FormValue("safekey")); err != nil{
				m["Error"] = err.Error()
			}else{
				name := strings.ToLower(controller.R.FormValue("name"))
				cities := make([]models.City,0,1)
				q := datastore.NewQuery("City").Ancestor(safekey).Filter("Name =",name).Filter("Active =",true).Limit(1)
				if keys,err := q.GetAll(c,&cities); err != nil{
					m["Error"] = err.Error()
				}else if len(keys)>0{
						m["Result"] = "Success"
						m["City"] = cities[0]
						m["SafeKey"] = keys[0].Encode()
				}else{
					//Get the country and put it in here
					country := models.Country{}
					countryKey := safekey.Parent()
					datastore.Get(c,countryKey,&country)
					city := models.City{
						Name: name,
						Nickname:controller.R.FormValue("nickname"),
						Created:time.Now(),
						Country: country.Nickname,
						Active:true,
					}
					vKey := datastore.NewIncompleteKey(c,"City",safekey)
					if countryKey, err := datastore.Put(c,vKey,&city); err != nil{
						m["Error"] = err.Error()
					}else{
						m["Result"] = "Success"
						m["SafeKey"] = countryKey.Encode()
					}
				}	
			}
		}
	data, _ := json.Marshal(m)
	fmt.Fprintf(controller.W,"%s",data)
}



func (controller *CitiesController) MobileView(){
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
			if cityKey,err := datastore.DecodeKey(controller.Params[2]); err != nil{
				m["Error"] = err.Error()
			}else{
				city := models.City{}
				if err = datastore.Get(c,cityKey,&city); err != nil{
					m["Error"] = err.Error()
				}else{
					m["Result"] = "Success"
					city.Name = strings.Title(city.Name);
					m["City"] = city
				}
			}
		}
	data, _ := json.Marshal(m)
	fmt.Fprintf(controller.W,"%s",data)
}

//Mobile add - Upload from phones
func (controller *CitiesController) MobileEdit(){
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
			cities := make([]models.City,0,1)
			q := datastore.NewQuery("City").Ancestor(userKey).Limit(1)
			if keys,err := q.GetAll(c,&cities); err != nil{
				m["Error"] = err.Error()+" Step 3"
			}else if len(keys)>0{
					video := cities[0]
						video.Name = controller.R.FormValue("name")
						video.Modified = time.Now()
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



func (controller *CitiesController) UpdateAll(){

	c := appengine.NewContext(controller.R)
	m := make(map[string]interface{})
	m["Result"] = "Failure"
	m["Error"] = "No authorization"
	q := datastore.NewQuery("City").Filter("Active =",false)
	assignments := make([]models.City,0,100)
	if keys, err := q.GetAll(c,&assignments); err != nil{
		println(err.Error())
	}else{
		for _,key := range keys{
			datastore.Delete(c,key)
		}
		datastore.PutMulti(c,keys,assignments)
		m["Result"] = "Success"
	}
	data, _ := json.Marshal(m)
	fmt.Fprintf(controller.W,"%s",data)
	/*
	c := appengine.NewContext(controller.R)
	q := datastore.NewQuery("City").Limit(100)
	assignments := make([]models.City,0,100)
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



func (controller *CitiesController) Render(show bool) error{
	//controller.Controller.Render(show)
	controller.Controller.Render(show)
	return nil
}