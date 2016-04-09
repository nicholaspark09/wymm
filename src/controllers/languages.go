package controllers

import (
	"net/http"
	"appengine"
	"appengine/datastore"
	"time"
	"github.com/parkn09/wymm/src/models"
	"fmt"
	"encoding/json"
	"strings"
)

type LanguagesController struct{
	Controller
}

func NewLanguagesController() *LanguagesController{
	return &LanguagesController{
	}
}

func (controller *LanguagesController) Init(r *http.Request, w http.ResponseWriter){
	controller.R = r
	controller.W = w
	controller.Data = make(map[interface{}]interface{})
	controller.C = appengine.NewContext(r)
	controller.Layout = make([]string,0)
	controller.TemplateNames = "body"
	controller.Name = "videos"
	controller.DetermineRoutes()

	controller.Layout = append(controller.Layout,"layouts/header.tmpl")
	controller.Layout = append(controller.Layout,"layouts/footer.tmpl")
	controller.Layout = append(controller.Layout,"layouts/layout.tmpl")
	if controller.CheckSession(){
		controller.Layout = append(controller.Layout,"layouts/logged.tmpl")
	}else{
				controller.Layout = append(controller.Layout,"layouts/nav.tmpl")
	}
	cookie, err := controller.R.Cookie("WymmLanguage")
	if err!= nil{
		controller.Locale = "en_US"
	}else if cookie.Value != ""{
		controller.Locale = cookie.Value
	}
	controller.Data["Language"] = controller.Locale
}

func (controller *LanguagesController) Serve(action string){
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
			}else if action =="settarget"{
				controller.SetTarget()
			}
}

func (controller *LanguagesController) MobileServe(action string){
	
}


func (controller *LanguagesController) Index(){

	c := appengine.NewContext(controller.R)
	if user,userKey, err := controller.GetUser(c); err != nil{
		controller.RenderError(err.Error())
	}else if userKey != nil && user.Group ==1{

		languages := make([]models.Language,0,100)
		q := datastore.NewQuery("Language").Filter("Active =",true).Filter("Native =",controller.Locale).Order("Name").Limit(100)
		if keys,err := q.GetAll(c,&languages); err != nil{
			controller.Data["Error"] = err.Error()
		}else{
			safekeys := make([]string,len(keys))
			for i,key := range keys{
				safekeys[i] = key.Encode()
			}
			controller.Data["Languages"] = languages
			controller.Data["SafeKeys"] = safekeys
		}
		controller.Layout = append(controller.Layout,"languages/index.tmpl")
	}else{
		controller.RenderError("You don't have authorization")
	}
	controller.Render(true)
}

func (controller *LanguagesController) Edit(){
	controller.Render(true)
}

func (controller *LanguagesController) Delete(){
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

func (controller *LanguagesController) DeleteAll(){
	/*
	c := appengine.NewContext(controller.R)
	q := datastore.NewQuery("Language").KeysOnly()
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

func (controller *LanguagesController) View(){
	//http.ServeFile(controller.W,controller.R,"views/layout.html")
	c := appengine.NewContext(controller.R)
	if token,err := datastore.DecodeKey(controller.Params[2]); err !=nil{
		controller.RenderError(err.Error())
	}else{
		video := models.Language{}
		if err := datastore.Get(c,token,&video); err !=nil{
			controller.RenderError(err.Error())
		}else{
			
		}
	}
	//
}

func (controller *LanguagesController) Home(){
	//http.ServeFile(controller.W,controller.R,"views/layout.html")
	
	controller.Layout = append(controller.Layout,"general/online.tmpl")
	controller.Render(true)
	//
}

func (controller *LanguagesController) Search(){
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
			languages := make([]models.Language,0,20)
			q := datastore.NewQuery("Language").Filter("Name >=",search).Filter("Active =",true).Project("Name").Limit(20)
			if keys, err := q.GetAll(c,&languages); err !=nil{
				m["Error"] = err.Error()
			}else{
				for i,key := range keys{
					//safekeys[i] = key.Encode()
					//Get the state as well
					stateKey := key.Parent()
					state := models.State{}
					if err = datastore.Get(c,stateKey,&state); err == nil{
						languages[i].Name = strings.Title(languages[i].Name)
					}
				}
				m["Result"] = "Success"
				m["Languages"] = languages
			}
		}
	data, _ := json.Marshal(m)
	fmt.Fprintf(controller.W,"%s",data)
}

//AJAX only 
func (controller *LanguagesController) MobileIndex(){
	c := appengine.NewContext(controller.R)
	m := make(map[string]interface{})
	m["Result"] = "Failure"
		locale := controller.R.FormValue("locale")
			languages := make([]models.Language,0,50)
		firsttwo := locale[:2]
		if firsttwo == "en"{
			locale = "en_US"
		}
			q := datastore.NewQuery("Language").Filter("Active =",true).Filter("Native =",locale).Order("Name").Limit(50)
			if keys, err := q.GetAll(c,&languages); err != nil{
				m["Error"] = err.Error()
			}else{
				m["Result"] = "Success"
				safekeys := make([]string,len(languages))
				for i,key := range keys{
					safekeys[i] = key.Encode()
					languages[i].Name = strings.Title(languages[i].Name)
				}
				m["SafeKeys"] = safekeys
				m["Languages"] = languages
			}	
	data, _ := json.Marshal(m)
	fmt.Fprintf(controller.W,"%s",data)
}

func (controller *LanguagesController) SetTarget(){
	cookie := http.Cookie{Name:"WymmTargetLanguage",Value:controller.R.FormValue("target"),Expires:time.Now().Add(365 * 24 *time.Hour),HttpOnly:true, MaxAge:900000,Path:"/"}
	http.SetCookie(controller.W,&cookie)
}

//Mobile add - Upload from phones
func (controller *LanguagesController) MobileAdd(){
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

				name := strings.ToLower(controller.R.FormValue("name"))
				languages := make([]models.Language,0,1)
				q := datastore.NewQuery("Language").Filter("Name =",name).Filter("Native =",controller.R.FormValue("native")).Filter("Locale =",controller.R.FormValue("locale")).Filter("Active =",true).Limit(1)
				if keys,err := q.GetAll(c,&languages); err != nil{
					m["Error"] = err.Error()
				}else if len(keys) == 0 {
					//Get the country and put it in here
					language := models.Language{
						Name: name,
						Locale: controller.R.FormValue("locale"),
						Native: controller.R.FormValue("native"),
						Created:time.Now(),
						Active:true,
					}
					vKey := datastore.NewIncompleteKey(c,"Language",nil)
					if countryKey, err := datastore.Put(c,vKey,&language); err != nil{
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

func (controller *LanguagesController) CreateMini(){
	c := appengine.NewContext(controller.R)
	m := make(map[string]interface{})
	m["Result"] = "Failure"
	langs := [10]string{"English","Korean","Spanish","Italian","French","German","Portuguese","Japanese","Arabic","Chinese"}
	countries := [10]string{"en_US","ko_KR","es_ES","it_IT","fr_FR","de_LI","pt_BR","ja_JP","ar_EG","zh_TW"}
	for i := 0; i<10; i++ {
		l := models.Language{
			Name: langs[i],
			Locale: countries[i],
			Created:time.Now(),
			Modified:time.Now(),
			Native: "en_US",
			Active:true,
		}
		key := datastore.NewIncompleteKey(c,"Language",nil)
		if k,err := datastore.Put(c,key,&l); err != nil{
			m["Error"] = err.Error()
		}else if k != nil{
			m["Result"] = "Success"
			m["key"] = k.Encode()
			m["Name"] = l.Name
		}
	}
	data, _ := json.Marshal(m)
	fmt.Fprintf(controller.W,"%s",data)
}


func (controller *LanguagesController) MobileView(){
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
			if languageKey,err := datastore.DecodeKey(controller.Params[2]); err != nil{
				m["Error"] = err.Error()
			}else{
				language := models.Language{}
				if err = datastore.Get(c,languageKey,&language); err != nil{
					m["Error"] = err.Error()
				}else{
					m["Result"] = "Success"
					language.Name = strings.Title(language.Name);
					m["Language"] = language
				}
			}
		}
	data, _ := json.Marshal(m)
	fmt.Fprintf(controller.W,"%s",data)
}

//Mobile add - Upload from phones
func (controller *LanguagesController) MobileEdit(){
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
			languages := make([]models.Language,0,1)
			q := datastore.NewQuery("Language").Ancestor(userKey).Limit(1)
			if keys,err := q.GetAll(c,&languages); err != nil{
				m["Error"] = err.Error()+" Step 3"
			}else if len(keys)>0{
					video := languages[0]
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



func (controller *LanguagesController) UpdateAll(){
	/*
	c := appengine.NewContext(controller.R)
	m := make(map[string]interface{})
	m["Result"] = "Failure"
	m["Error"] = "No authorization"
	q := datastore.NewQuery("Language").Filter("Active =",true)
	lines := make([]models.Language,0,100)
	if keys, err := q.GetAll(c,&lines); err != nil{
		println(err.Error())
	}else{
		for i,_ := range keys{
			lines[i].Native = "en_US"
		}
		datastore.PutMulti(c,keys,lines)
		m["Result"] = "Success"
	}
	data, _ := json.Marshal(m)
	fmt.Fprintf(controller.W,"%s",data)
	/*
	c := appengine.NewContext(controller.R)
	q := datastore.NewQuery("Language").Limit(100)
	assignments := make([]models.Language,0,100)
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



func (controller *LanguagesController) Render(show bool) error{
	//controller.Controller.Render(show)
	controller.Controller.Render(show)
	return nil
}