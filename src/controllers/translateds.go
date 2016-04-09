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

type TranslatedsController struct{
	Controller
}

func NewTranslatedsController() *TranslatedsController{
	return &TranslatedsController{
	}
}

func (controller *TranslatedsController) Init(r *http.Request, w http.ResponseWriter){
	controller.R = r
	controller.W = w
	controller.Data = make(map[interface{}]interface{})
	controller.C = appengine.NewContext(r)
	controller.Layout = make([]string,0)
	controller.TemplateNames = "body"
	controller.Name = "videos"
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

func (controller *TranslatedsController) Serve(action string){
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
			}else if action =="findline"{
				controller.FindLine()
			}else if action =="add"{
				controller.Add()
			}
}

func (controller *TranslatedsController) MobileServe(action string){
	
}


func (controller *TranslatedsController) Index(){

	controller.Layout = append(controller.Layout,"translateds/index.tmpl")

	controller.Render(true)
}

func (controller *TranslatedsController) Edit(){
	controller.Render(true)
}

func (controller *TranslatedsController) Delete(){
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
		if translatedKey,err := datastore.DecodeKey(controller.R.FormValue("safekey")); err != nil{
			m["Error"] = err.Error()
		}else if translatedKey!=nil{
			translated := models.Translated{}
			if err = datastore.Get(c,translatedKey,&translated); err != nil{
				m["Error"] = err.Error()
			}else if translated.User == userKey.Encode(){
				translated.Active = false
				translated.Modified = time.Now()
				if key, err := datastore.Put(c,translatedKey,&translated); err != nil{
					m["Error"] = err.Error()
				}else{
					m["Result"] = "Success"
					m["SafeKey"] = key.Encode()
					m["Translateds"] = "0"
					line := models.Line{}
					lineKey := translatedKey.Parent()
					if err = datastore.Get(c,lineKey,&line); err != nil{
						m["Error"] = err.Error()
					}else{
						
					}
				}
			}
		}
	}
	data, _ := json.Marshal(m)
	fmt.Fprintf(controller.W,"%s",data)
}

func (controller *TranslatedsController) DeleteAll(){
	c := appengine.NewContext(controller.R)
	q := datastore.NewQuery("Translated").KeysOnly()
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

func (controller *TranslatedsController) View(){
	//http.ServeFile(controller.W,controller.R,"views/layout.html")
	c := appengine.NewContext(controller.R)
	if token,err := datastore.DecodeKey(controller.Params[2]); err !=nil{
		controller.RenderError(err.Error())
	}else{
		video := models.Translated{}
		if err := datastore.Get(c,token,&video); err !=nil{
			controller.RenderError(err.Error())
		}else{
			
		}
	}
	//
}

func (controller *TranslatedsController) Home(){
	//http.ServeFile(controller.W,controller.R,"views/layout.html")
	
	controller.Layout = append(controller.Layout,"general/ontranslated.tmpl")
	controller.Render(true)
	//
}

func (controller *TranslatedsController) Search(){
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
			translateds := make([]models.Translated,0,20)
			q := datastore.NewQuery("Translated").Filter("Name >=",search).Filter("Active =",true).Project("Name").Limit(20)
			if keys, err := q.GetAll(c,&translateds); err !=nil{
				m["Error"] = err.Error()
			}else{
				for i,key := range keys{
					//safekeys[i] = key.Encode()
					//Get the state as well
					stateKey := key.Parent()
					state := models.State{}
					if err = datastore.Get(c,stateKey,&state); err == nil{
						translateds[i].Name = strings.Title(translateds[i].Name)
					}
				}
				m["Result"] = "Success"
				m["Translateds"] = translateds
			}
		}
	data, _ := json.Marshal(m)
	fmt.Fprintf(controller.W,"%s",data)
}

//AJAX only 
func (controller *TranslatedsController) MobileIndex(){
	c := appengine.NewContext(controller.R)
	m := make(map[string]interface{})
	m["Result"] = "Failure"
	if lineKey,err := datastore.DecodeKey(controller.R.FormValue("safekey")); err !=nil{
		m["Error"] = err.Error()
	}else{
			offset,_ := strconv.Atoi(controller.R.FormValue("offset"))
			translateds := make([]models.Translated,0,10)
			q := datastore.NewQuery("Translated").Ancestor(lineKey).Filter("Locale =",controller.R.FormValue("locale")).Filter("Active =",true).Order("-Ranking").Offset(offset).Limit(10)

			if keys, err := q.GetAll(c,&translateds); err != nil{
				m["Error"] = err.Error()
			}else{
				m["Result"] = "Success"
				for i,key := range keys{
					translateds[i].User = key.Encode()
				}
				m["Translateds"] = translateds
			}
	}	
	data, _ := json.Marshal(m)
	fmt.Fprintf(controller.W,"%s",data)
}

func (controller *TranslatedsController) FindLine(){
	c := appengine.NewContext(controller.R)
	m := make(map[string]interface{})
	m["Result"] = "Failure"
	//if controller.R.FormValue("key") == "39234hwersKLnsdk2!sdkdfjd"{
		locale := controller.R.FormValue("locale")
		firsttwo := locale[:2]
		if firsttwo == "en"{
				locale = "en_US"
		}
		if lineKey, err := datastore.DecodeKey(controller.R.FormValue("safekey")); err != nil{
			m["Error"] = err.Error()
		}else{
			translateds := make([]models.Translated,0,1)
			q := datastore.NewQuery("Translated").Ancestor(lineKey).Filter("Locale =",locale).Filter("Active =",true).Limit(1)
			if keys, err := q.GetAll(c,&translateds); err != nil{
				m["Error"] = err.Error()
			}else if len(keys)>0{
				m["Result"] = "Success"
				translateds[0].Nickname = keys[0].Encode()
				m["Translateds"] = translateds[0]
			}	
		}
	//}
	data, _ := json.Marshal(m)
	fmt.Fprintf(controller.W,"%s",data)
}

//Mobile add - Upload from phones
func (controller *TranslatedsController) MobileAdd(){
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
			if lineKey,err := datastore.DecodeKey(controller.R.FormValue("safekey")); err != nil{
				m["Error"] = err.Error()
			}else{
				name := controller.R.FormValue("name")
					//Get the country and put it in here
					translated := models.Translated{
						Name: name,
						Controller: controller.R.FormValue("controller"),
						Action: controller.R.FormValue("safekey"),
						Created:time.Now(),
						Modified:time.Now(),
						Locale: controller.R.FormValue("locale"),
						User: userKey.Encode(),
						Likes:0,
						Flags:0,
						Ranking:0.0,
						Active:true,
					}
					vKey := datastore.NewIncompleteKey(c,"Translated",lineKey)
					if countryKey, err := datastore.Put(c,vKey,&translated); err != nil{
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

func (controller *TranslatedsController) Add(){
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
				
					//Get the country and put it in here
					translated := models.Translated{
						Name: controller.R.FormValue("body"),
						Nickname: user.Name,
						Controller: controller.R.FormValue("controller"),
						Action: controller.R.FormValue("safekey"),
						Created:time.Now(),
						Modified:time.Now(),
						User: userKey.Encode(),
						Locale:controller.R.FormValue("locale"),
						Likes:0,
						Ranking:0,
						Flags:0,
						Active:true,
					}
					vKey := datastore.NewIncompleteKey(c,"Translated",lineKey)
					if countryKey, err := datastore.Put(c,vKey,&translated); err != nil{
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

func (controller *TranslatedsController) FindBest(){
	c := appengine.NewContext(controller.R)
	m := make(map[string]interface{})
	m["Result"] = "Failure"
		cont := controller.R.FormValue("controller")
		action := controller.R.FormValue("action")
		locale := controller.R.FormValue("locale")
	if actionKey,err := datastore.DecodeKey(action); err !=nil{
		m["Error"] = err.Error()
	}else{
		q := datastore.NewQuery("Translated").Ancestor(actionKey).Filter("Controller =",cont).Filter("Locale =",locale).Order("-Likes").Limit(1)
		trans := make([]models.Translated,0,1)
		if keys,err := q.GetAll(c,&trans); err != nil{
			m["Error"] = err.Error()
		}else if len(keys)>0{
			m["Result"] = "Success"
			m["Translated"] = trans[0]
			m["SafeKey"] = keys[0].Encode()
		}
	}
	data, _ := json.Marshal(m)	
	fmt.Fprintf(controller.W,"%s",data)
}

func (controller *TranslatedsController) FindTranslated(c appengine.Context,userstring string, controllerKey *datastore.Key) (*datastore.Key,error){

	translateds := make([]models.Translated,0,1)
	q := datastore.NewQuery("Translated").Ancestor(controllerKey).Filter("User =",userstring).Limit(1)
	if keys, err := q.GetAll(c,&translateds); err != nil{
		return nil,err
	}else if len(keys)>0{
		if translateds[0].Active == false {
			translateds[0].Active = true
			translateds[0].Modified = time.Now()
			if _,err := datastore.Put(c,keys[0],&translateds[0]); err != nil{
				return nil,err
			}else{
			

				return keys[0],nil
			}
		}
		return keys[0],nil
	}
	return nil,nil
}




func (controller *TranslatedsController) MobileView(){
	c := appengine.NewContext(controller.R)
	m := make(map[string]interface{})
	m["Result"] = "Failure"
	m["Error"] = "No authorization"
	m["Mine"] = "false"
	if !controller.CheckSession(){
		//Checks to see if this is from the web
		controller.UserToken = controller.R.FormValue("token")
	}
			if translatedKey,err := datastore.DecodeKey(controller.R.FormValue("safekey")); err != nil{
				m["Error"] = err.Error()
			}else{
				translated := models.Translated{}
				if err = datastore.Get(c,translatedKey,&translated); err != nil{
					m["Error"] = err.Error()
				}else{
					m["Result"] = "Success"
					
					if _,userKey,err := controller.GetUser(c); err !=nil{
						m["Error"] = "Authorization Error"
					}else if userKey != nil && userKey.Encode() == translated.User{
						m["Mine"] = "true"
					}
					m["Translated"] = translated
				}
			}
	data, _ := json.Marshal(m)
	fmt.Fprintf(controller.W,"%s",data)
}

//Mobile add - Upload from phones
func (controller *TranslatedsController) MobileEdit(){
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
			translateds := make([]models.Translated,0,1)
			q := datastore.NewQuery("Translated").Ancestor(userKey).Limit(1)
			if keys,err := q.GetAll(c,&translateds); err != nil{
				m["Error"] = err.Error()+" Step 3"
			}else if len(keys)>0{
					video := translateds[0]
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



func (controller *TranslatedsController) UpdateAll(){

	c := appengine.NewContext(controller.R)
	m := make(map[string]interface{})
	m["Result"] = "Failure"
	m["Error"] = "No authorization"
	q := datastore.NewQuery("Translated").Limit(100)
	translateds := make([]models.Translated,0,100)
	if keys, err := q.GetAll(c,&translateds); err != nil{
		println(err.Error())
	}else{
		for i,_ := range keys{
			translateds[i].Ranking = 0.0
		}
		datastore.PutMulti(c,keys,translateds)
		m["Result"] = "Success"
	}
	data, _ := json.Marshal(m)
	fmt.Fprintf(controller.W,"%s",data)
	/*
	c := appengine.NewContext(controller.R)
	q := datastore.NewQuery("Translated").Limit(100)
	assignments := make([]models.Translated,0,100)
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



func (controller *TranslatedsController) Render(show bool) error{
	//controller.Controller.Render(show)
	controller.Controller.Render(show)
	return nil
}