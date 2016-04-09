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

type FavoritesController struct{
	Controller
}

func NewFavoritesController() *FavoritesController{
	return &FavoritesController{
	}
}

func (controller *FavoritesController) Init(r *http.Request, w http.ResponseWriter){
	controller.R = r
	controller.W = w
	controller.Data = make(map[interface{}]interface{})
	controller.C = appengine.NewContext(r)
	controller.Layout = make([]string,0)
	controller.TemplateNames = "body"
	controller.Name = "favorites"
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

func (controller *FavoritesController) Serve(action string){
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
			}else if action =="add"{
				controller.Add()
			}
}

func (controller *FavoritesController) MobileServe(action string){
	
}


func (controller *FavoritesController) Index(){

	controller.Layout = append(controller.Layout,"favorites/index.tmpl")

	controller.Render(true)
}

func (controller *FavoritesController) Edit(){
	controller.Render(true)
}

func (controller *FavoritesController) Delete(){
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
		favorites := make([]models.Favorite,0,1)
		q := datastore.NewQuery("Favorite").Ancestor(userKey).Filter("Line =",controller.R.FormValue("safekey")).Filter("Active =",true).Limit(1)
		if keys,err := q.GetAll(c,&favorites); err != nil{
			m["Error"] = "Could not complete. Error: "+err.Error()
		}else if len(keys)>0{
			favorites[0].Active = false
			favorites[0].Modified = time.Now()
			if _,err := datastore.Put(c,keys[0],&favorites[0]); err != nil{
				m["Error"] = err.Error()
			}else{
				m["Result"] = "Success"
			}
		}
	}
	data, _ := json.Marshal(m)
	fmt.Fprintf(controller.W,"%s",data)
}

func (controller *FavoritesController) DeleteAll(){
	c := appengine.NewContext(controller.R)
	q := datastore.NewQuery("Favorite").KeysOnly()
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

func (controller *FavoritesController) View(){
	//http.ServeFile(controller.W,controller.R,"views/layout.html")
	c := appengine.NewContext(controller.R)
	if token,err := datastore.DecodeKey(controller.Params[2]); err !=nil{
		controller.RenderError(err.Error())
	}else{
		video := models.Favorite{}
		if err := datastore.Get(c,token,&video); err !=nil{
			controller.RenderError(err.Error())
		}else{
			
		}
	}
	//
}

func (controller *FavoritesController) Home(){
	//http.ServeFile(controller.W,controller.R,"views/layout.html")
	
	controller.Layout = append(controller.Layout,"general/onfavorite.tmpl")
	controller.Render(true)
	//
}

func (controller *FavoritesController) Search(){
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
			favorites := make([]models.Favorite,0,20)
			q := datastore.NewQuery("Favorite").Filter("Name >=",search).Filter("Active =",true).Project("Name").Limit(20)
			if keys, err := q.GetAll(c,&favorites); err !=nil{
				m["Error"] = err.Error()
			}else{
				for i,key := range keys{
					//safekeys[i] = key.Encode()
					//Get the state as well
					stateKey := key.Parent()
					state := models.State{}
					if err = datastore.Get(c,stateKey,&state); err == nil{
						favorites[i].Name = strings.Title(favorites[i].Name)
					}
				}
				m["Result"] = "Success"
				m["Favorites"] = favorites
			}
		}
	data, _ := json.Marshal(m)
	fmt.Fprintf(controller.W,"%s",data)
}

//AJAX only 
func (controller *FavoritesController) MobileIndex(){
	c := appengine.NewContext(controller.R)
	m := make(map[string]interface{})
	m["Result"] = "Failure"

			offset,_ := strconv.Atoi(controller.R.FormValue("offset"))
			favorites := make([]models.Favorite,0,50)
			var q *datastore.Query
				q = datastore.NewQuery("Favorite").Filter("Active =",true).Order("-Favorites").Offset(offset).Limit(10)

			if keys, err := q.GetAll(c,&favorites); err != nil{
				m["Error"] = err.Error()
			}else{
				m["Result"] = "Success"
				safekeys := make([]string,len(favorites))
				for i,key := range keys{
					safekeys[i] = key.Encode()
				}
				m["SafeKeys"] = safekeys
				m["Favorites"] = favorites
			}	
	data, _ := json.Marshal(m)
	fmt.Fprintf(controller.W,"%s",data)
}


//Mobile add - Upload from phones
func (controller *FavoritesController) MobileAdd(){
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
					favorite := models.Favorite{
						Name: name,
						Created:time.Now(),
						Modified:time.Now(),
						Active:true,
					}
					vKey := datastore.NewIncompleteKey(c,"Favorite",userKey)
					if countryKey, err := datastore.Put(c,vKey,&favorite); err != nil{
						m["Error"] = err.Error()
					}else{
						m["Result"] = "Success"
						m["SafeKey"] = countryKey.Encode()
					}	
			
		}
	data, _ := json.Marshal(m)
	fmt.Fprintf(controller.W,"%s",data)
}

func (controller *FavoritesController) Add(){
	c := appengine.NewContext(controller.R)
	m := make(map[string]interface{})
	m["Result"] = "Failure"
	m["Error"] = "You aren't logged in <a href=\"/users/login\">Login</a>"
	if !controller.CheckSession(){
		//Checks to see if this is from the web
		controller.UserToken = controller.R.FormValue("token")
	}
		if _,userKey,err := controller.GetUser(c); err ==nil && userKey != nil{

			if lineKey,err := datastore.DecodeKey(controller.R.FormValue("safekey")); err !=nil{
				m["Error"] = "Error Decoding Line: "+err.Error()
			}else if lineKey != nil{
				
				//Check to see if a previous favorite exists
				favorites := make([]models.Favorite,0,1)
				q := datastore.NewQuery("Favorite").Ancestor(userKey).Filter("Line =",controller.R.FormValue("safekey")).Limit(1)
				if keys,err := q.GetAll(c,&favorites); err == nil && len(keys)>0{
					//Found a line
					favoriteKey := ""
					favorite := favorites[0]
					if favorite.Active == false{
						favorite.Active = true
						favorite.Modified = time.Now()
						if fKey,err := datastore.Put(c,keys[0],&favorite); err != nil{
							m["Error"] = err.Error()
						}else{
							favoriteKey = fKey.Encode()
						}
					}else{
						favoriteKey = keys[0].Encode()
					}
					m["Result"] = "Success"
					m["SafeKey"] = favoriteKey
				}else if err != nil{
					m["Error"] = err.Error()
				}else{
					//Get the country and put it in here
					favorite := models.Favorite{
						Line:controller.R.FormValue("safekey"),
						Created:time.Now(),
						Modified:time.Now(),
						Active:true,
					}
					vKey := datastore.NewIncompleteKey(c,"Favorite",userKey)
					if countryKey, err := datastore.Put(c,vKey,&favorite); err != nil{
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



func (controller *FavoritesController) MobileView(){
	c := appengine.NewContext(controller.R)
	m := make(map[string]interface{})
	m["Result"] = "Failure"
	m["Error"] = "No authorization"

			if favoriteKey,err := datastore.DecodeKey(controller.R.FormValue("safekey")); err != nil{
				m["Error"] = err.Error()
			}else{
				favorite := models.Favorite{}
				if err = datastore.Get(c,favoriteKey,&favorite); err != nil{
					m["Error"] = err.Error()
				}else{
					m["Result"] = "Success"
					m["Favorite"] = favorite
					m["Translated"] = ""
				}
			}
	data, _ := json.Marshal(m)
	fmt.Fprintf(controller.W,"%s",data)
}

//Mobile add - Upload from phones
func (controller *FavoritesController) MobileEdit(){
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
			favorites := make([]models.Favorite,0,1)
			q := datastore.NewQuery("Favorite").Ancestor(userKey).Limit(1)
			if keys,err := q.GetAll(c,&favorites); err != nil{
				m["Error"] = err.Error()+" Step 3"
			}else if len(keys)>0{
					video := favorites[0]
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



func (controller *FavoritesController) UpdateAll(){

	c := appengine.NewContext(controller.R)
	m := make(map[string]interface{})
	m["Result"] = "Failure"
	m["Error"] = "No authorization"
	q := datastore.NewQuery("Favorite").Filter("Active =",false)
	favorites := make([]models.Favorite,0,100)
	if keys, err := q.GetAll(c,&favorites); err != nil{
		println(err.Error())
	}else{
		for i,_ := range keys{
			favorites[i].Active = false
		}
		datastore.PutMulti(c,keys,favorites)
		m["Result"] = "Success"
	}
	data, _ := json.Marshal(m)
	fmt.Fprintf(controller.W,"%s",data)
	/*
	c := appengine.NewContext(controller.R)
	q := datastore.NewQuery("Favorite").Limit(100)
	assignments := make([]models.Favorite,0,100)
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



func (controller *FavoritesController) Render(show bool) error{
	//controller.Controller.Render(show)
	controller.Controller.Render(show)
	return nil
}