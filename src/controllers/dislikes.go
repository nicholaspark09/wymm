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

type DislikesController struct{
	Controller
}

func NewDislikesController() *DislikesController{
	return &DislikesController{
	}
}

func (controller *DislikesController) Init(r *http.Request, w http.ResponseWriter){
	controller.R = r
	controller.W = w
	controller.Data = make(map[interface{}]interface{})
	controller.C = appengine.NewContext(r)
	controller.Layout = make([]string,0)
	controller.TemplateNames = "body"
	controller.Name = "dislikes"
	controller.DetermineRoutes()

	cookie, err := controller.R.Cookie("WymmLanguage")
	if err!= nil{
		controller.Locale = "en_US"
	}else if cookie.Value != ""{
		controller.Locale = cookie.Value
	}
	targetCookie,err := controller.R.Cookie("WymmTargetLanguage")
	if err != nil{
		controller.Target = "en_US"
	}else if targetCookie.Value != ""{
		controller.Target = targetCookie.Value
	}
	controller.Data["Locale"] = controller.Locale
	controller.Data["Target"] = controller.Target
	controller.Layout = append(controller.Layout,"layouts/header.tmpl")
	controller.Layout = append(controller.Layout,"layouts/footer.tmpl")
	controller.Layout = append(controller.Layout,"layouts/layout.tmpl")
	if controller.CheckSession(){
		controller.Layout = append(controller.Layout,"layouts/logged.tmpl")
	}else{
				controller.Layout = append(controller.Layout,"layouts/nav.tmpl")
	}
}

func (controller *DislikesController) Serve(action string){
			controller.Action = action
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
			}else if action =="createmini"{
				controller.CreateMini()
			}else if action =="topten"{
				controller.TopTen()
			}else if action =="add"{
				controller.Add()
			}
}

func (controller *DislikesController) MobileServe(action string){
	
}


func (controller *DislikesController) Index(){

	controller.Layout = append(controller.Layout,"likes/index.tmpl")

	controller.Render(true)
}

func (controller *DislikesController) Edit(){
	controller.Render(true)
}

func (controller *DislikesController) Delete(){
	//Adding a video from a phone / Get the video + images
	c := appengine.NewContext(controller.R)
	m := make(map[string]interface{})
	m["Result"] = "Failure"
	if !controller.CheckSession(){
		//Checks to see if this is from the web
		controller.UserToken = controller.R.FormValue("token")
	}
	if _,userKey, err := controller.GetUser(c); err != nil{
		m["Error"] = err.Error()
	}else{
		if likeKey,err := datastore.DecodeKey(controller.R.FormValue("safekey")); err != nil{
			m["Error"] = err.Error()
		}else if likeKey!=nil{
			like := models.Dislike{}
			if err = datastore.Get(c,likeKey,&like); err != nil{
				m["Error"] = err.Error()
			}else if like.User == userKey.Encode(){
				like.Active = false
				like.Modified = time.Now()
				if key, err := datastore.Put(c,likeKey,&like); err != nil{
					m["Error"] = err.Error()
				}else{
					m["Result"] = "Success"
					m["SafeKey"] = key.Encode()
					m["Dislikes"] = "0"
					line := models.Line{}
					lineKey := likeKey.Parent()
					if err = datastore.Get(c,lineKey,&line); err != nil{
						m["Error"] = err.Error()
					}else{
						if line.Dislikes > 0{
							line.Dislikes-=1
							datastore.Put(c,lineKey,&line)
							m["Dislikes"] = line.Dislikes
						}
					}
				}
			}
		}
	}
	data, _ := json.Marshal(m)
	fmt.Fprintf(controller.W,"%s",data)
}

func (controller *DislikesController) DeleteAll(){
	c := appengine.NewContext(controller.R)
	q := datastore.NewQuery("Dislike").KeysOnly()
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

func (controller *DislikesController) View(){
	//http.ServeFile(controller.W,controller.R,"views/layout.html")
	c := appengine.NewContext(controller.R)
	if token,err := datastore.DecodeKey(controller.Params[2]); err !=nil{
		controller.RenderError(err.Error())
	}else{
		video := models.Dislike{}
		if err := datastore.Get(c,token,&video); err !=nil{
			controller.RenderError(err.Error())
		}else{
			
		}
	}
	//
}

func (controller *DislikesController) Home(){
	//http.ServeFile(controller.W,controller.R,"views/layout.html")
	
	controller.Layout = append(controller.Layout,"general/onlike.tmpl")
	controller.Render(true)
	//
}

func (controller *DislikesController) Search(){
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
			likes := make([]models.Dislike,0,20)
			q := datastore.NewQuery("Dislike").Filter("Name >=",search).Filter("Active =",true).Project("Name").Limit(20)
			if keys, err := q.GetAll(c,&likes); err !=nil{
				m["Error"] = err.Error()
			}else{
				for i,key := range keys{
					//safekeys[i] = key.Encode()
					//Get the state as well
					stateKey := key.Parent()
					state := models.State{}
					if err = datastore.Get(c,stateKey,&state); err == nil{
						likes[i].Name = strings.Title(likes[i].Name)
					}
				}
				m["Result"] = "Success"
				m["Dislikes"] = likes
			}
		}
	data, _ := json.Marshal(m)
	fmt.Fprintf(controller.W,"%s",data)
}

//AJAX only 
func (controller *DislikesController) MobileIndex(){
	c := appengine.NewContext(controller.R)
	m := make(map[string]interface{})
	m["Result"] = "Failure"

			offset,_ := strconv.Atoi(controller.R.FormValue("offset"))
			likes := make([]models.Dislike,0,50)
			var q *datastore.Query
				q = datastore.NewQuery("Dislike").Filter("Active =",true).Order("-Dislikes").Offset(offset).Limit(10)

			if keys, err := q.GetAll(c,&likes); err != nil{
				m["Error"] = err.Error()
			}else{
				m["Result"] = "Success"
				safekeys := make([]string,len(likes))
				for i,key := range keys{
					safekeys[i] = key.Encode()
				}
				m["SafeKeys"] = safekeys
				m["Dislikes"] = likes
			}	
	data, _ := json.Marshal(m)
	fmt.Fprintf(controller.W,"%s",data)
}

func (controller *DislikesController) TopTen(){
	c := appengine.NewContext(controller.R)
	m := make(map[string]interface{})
	m["Result"] = "Failure"
	key := controller.R.FormValue("key")
	if key == "39234hwersKLnsdk2!sdkdfjd"{

		if langKey, err := datastore.DecodeKey(controller.R.FormValue("locale")); err != nil{
			m["Error"] = err.Error()
		}else if langKey != nil{
			likes := make([]models.Dislike,0,10)
			q := datastore.NewQuery("Dislike").Ancestor(langKey).Filter("Active =",true).Order("-Dislikes").Limit(10)
			if keys, err := q.GetAll(c,&likes); err != nil{
				m["Error"] = err.Error()
			}else{
				m["Result"] = "Success"
				safekeys := make([]string,len(likes))
				for i,key := range keys{
					safekeys[i] = key.Encode()
				}
				m["SafeKeys"] = safekeys
				m["Dislikes"] = likes
			}	
		}else{

		}
	}
	data, _ := json.Marshal(m)
	fmt.Fprintf(controller.W,"%s",data)
}

//Mobile add - Upload from phones
func (controller *DislikesController) MobileAdd(){
	c := appengine.NewContext(controller.R)
	m := make(map[string]interface{})
	m["Result"] = "Failure"
	m["Error"] = "You aren't logged in"
	if !controller.CheckSession(){
		//Checks to see if this is from the web
		controller.UserToken = controller.R.FormValue("token")
	}
		if _,userKey,err := controller.GetUser(c); err !=nil{
			m["Error"] = "You aren't logged in"
		}else if userKey != nil{

				name := controller.R.FormValue("name")
					//Get the country and put it in here
					like := models.Dislike{
						Name: name,
						Created:time.Now(),
						Modified:time.Now(),
						Active:true,
					}
					vKey := datastore.NewIncompleteKey(c,"Dislike",userKey)
					if countryKey, err := datastore.Put(c,vKey,&like); err != nil{
						m["Error"] = err.Error()
					}else{
						m["Result"] = "Success"
						m["SafeKey"] = countryKey.Encode()
					}	
			
		}
	data, _ := json.Marshal(m)
	fmt.Fprintf(controller.W,"%s",data)
}

func (controller *DislikesController) Add(){
	c := appengine.NewContext(controller.R)
	m := make(map[string]interface{})
	m["Result"] = "Failure"
	m["Error"] = "You aren't logged in <a href=\"/users/login\">Login</a>"
	if !controller.CheckSession(){
		//Checks to see if this is from the web
		controller.UserToken = controller.R.FormValue("token")
	}
		if user,userKey,err := controller.GetUser(c); err ==nil && userKey != nil{


			if lineKey,err := datastore.DecodeKey(controller.R.FormValue("safekey")); err !=nil{
				m["Error"] = "Error Decoding Line: "+err.Error()
			}else if lineKey != nil{
				if likeKey,err := controller.FindDislike(c,userKey.Encode(),lineKey); err !=nil{
					m["Error"] = err.Error()
				}else if likeKey != nil{
					m["Result"] = "Success"
					m["SafeKey"] = likeKey.Encode()
				}else{
					//Get the country and put it in here
					like := models.Dislike{
						Name: user.Name,
						Controller: "lines",
						Action: controller.R.FormValue("safekey"),
						Created:time.Now(),
						Modified:time.Now(),
						User: userKey.Encode(),
						Active:true,
					}
					vKey := datastore.NewIncompleteKey(c,"Dislike",lineKey)
					if countryKey, err := datastore.Put(c,vKey,&like); err != nil{
						m["Error"] = err.Error()
					}else{
						m["Result"] = "Success"
						m["SafeKey"] = countryKey.Encode()
						line := models.Line{}
						if err = datastore.Get(c,lineKey,&line); err != nil{
							m["Error"] = err.Error()
						}else{
							line.Dislikes += 1
							datastore.Put(c,lineKey,&line)
							m["Dislikes"] = line.Dislikes
						}
					}
				}
			}
		}
	data, _ := json.Marshal(m)
	fmt.Fprintf(controller.W,"%s",data)
}

func (controller *DislikesController) FindDislike(c appengine.Context,userstring string, controllerKey *datastore.Key) (*datastore.Key,error){

	likes := make([]models.Dislike,0,1)
	q := datastore.NewQuery("Dislike").Ancestor(controllerKey).Filter("User =",userstring).Limit(1)
	if keys, err := q.GetAll(c,&likes); err != nil{
		return nil,err
	}else if len(keys)>0{
		if likes[0].Active == false {
			likes[0].Active = true
			likes[0].Modified = time.Now()
			if _,err := datastore.Put(c,keys[0],&likes[0]); err != nil{
				return nil,err
			}else{
				line := models.Line{}
				if err = datastore.Get(c,controllerKey,&line);err == nil{
					line.Dislikes+=1
					line.Modified = time.Now()
					datastore.Put(c,controllerKey,&line)
				}

				return keys[0],nil
			}
		}
		return keys[0],nil
	}
	return nil,nil
}

func (controller *DislikesController) CreateMini(){
	c := appengine.NewContext(controller.R)
	m := make(map[string]interface{})
	m["Result"] = "Failure"
	q := datastore.NewQuery("Language").Filter("Name =","English").Limit(1)
	langs := make([]models.Language,0,1)
	if keys,err := q.GetAll(c,&langs); err != nil{
		m["Error"] = err.Error()
	}else if len(keys)>0{
		key := keys[0]
		like := models.Dislike{
			Name:"I'm not a photographer, but I can picture me and you together.",
			Created:time.Now(),
			Modified:time.Now(),
			Active:true,
		}
		lKey := datastore.NewIncompleteKey(c,"Dislike",key)
		datastore.Put(c,lKey,&like)
	}
	data, _ := json.Marshal(m)
	fmt.Fprintf(controller.W,"%s",data)
}


func (controller *DislikesController) MobileView(){
	c := appengine.NewContext(controller.R)
	m := make(map[string]interface{})
	m["Result"] = "Failure"
	m["Error"] = "No authorization"

			if likeKey,err := datastore.DecodeKey(controller.R.FormValue("safekey")); err != nil{
				m["Error"] = err.Error()
			}else{
				like := models.Dislike{}
				if err = datastore.Get(c,likeKey,&like); err != nil{
					m["Error"] = err.Error()
				}else{
					m["Result"] = "Success"
					m["Dislike"] = like
					m["Translated"] = ""
				}
			}
	data, _ := json.Marshal(m)
	fmt.Fprintf(controller.W,"%s",data)
}

//Mobile add - Upload from phones
func (controller *DislikesController) MobileEdit(){
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
			likes := make([]models.Dislike,0,1)
			q := datastore.NewQuery("Dislike").Ancestor(userKey).Limit(1)
			if keys,err := q.GetAll(c,&likes); err != nil{
				m["Error"] = err.Error()+" Step 3"
			}else if len(keys)>0{
					video := likes[0]
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



func (controller *DislikesController) UpdateAll(){

	c := appengine.NewContext(controller.R)
	m := make(map[string]interface{})
	m["Result"] = "Failure"
	m["Error"] = "No authorization"
	q := datastore.NewQuery("Dislike").Filter("Active =",false)
	likes := make([]models.Dislike,0,100)
	if keys, err := q.GetAll(c,&likes); err != nil{
		println(err.Error())
	}else{
		for i,_ := range keys{
			likes[i].Active = false
		}
		datastore.PutMulti(c,keys,likes)
		m["Result"] = "Success"
	}
	data, _ := json.Marshal(m)
	fmt.Fprintf(controller.W,"%s",data)
	/*
	c := appengine.NewContext(controller.R)
	q := datastore.NewQuery("Dislike").Limit(100)
	assignments := make([]models.Dislike,0,100)
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



func (controller *DislikesController) Render(show bool) error{
	//controller.Controller.Render(show)
	controller.Controller.Render(show)
	return nil
}