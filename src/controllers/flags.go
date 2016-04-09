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

type FlagsController struct{
	Controller
}

func NewFlagsController() *FlagsController{
	return &FlagsController{
	}
}

func (controller *FlagsController) Init(r *http.Request, w http.ResponseWriter){
	controller.R = r
	controller.W = w
	controller.Data = make(map[interface{}]interface{})
	controller.C = appengine.NewContext(r)
	controller.Layout = make([]string,0)
	controller.TemplateNames = "body"
	controller.Name = "flags"
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

func (controller *FlagsController) Serve(action string){
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

func (controller *FlagsController) MobileServe(action string){
	
}


func (controller *FlagsController) Index(){

	controller.Layout = append(controller.Layout,"flags/index.tmpl")

	controller.Render(true)
}

func (controller *FlagsController) Edit(){
	controller.Render(true)
}

func (controller *FlagsController) Delete(){
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
		if flagKey,err := datastore.DecodeKey(controller.R.FormValue("safekey")); err != nil{
			m["Error"] = err.Error()
		}else if flagKey!=nil{
			flag := models.Flag{}
			if err = datastore.Get(c,flagKey,&flag); err != nil{
				m["Error"] = err.Error()
			}else if flagKey.Parent().Encode() == userKey.Encode(){
				flag.Active = false
				flag.Modified = time.Now()
				if key, err := datastore.Put(c,flagKey,&flag); err != nil{
					m["Error"] = err.Error()
				}else{
					m["Result"] = "Success"
					m["SafeKey"] = key.Encode()
					if flag.Controller == "translateds"{
						translated := models.Translated{}
						if translatedKey,err := datastore.DecodeKey(flag.Action); err != nil{
							m["Error"] = err.Error()
						}else if err = datastore.Get(c,translatedKey,&translated); err == nil{
							translated.Flags -=1
							translated.Modified = time.Now()
							if translated.Likes < 1 {
								translated.Ranking = float32(translated.Flags)
							}else if translated.Likes > 0 && translated.Flags == 0{
								translated.Ranking = float32(translated.Likes)
							}else if translated.Likes > 0 && translated.Flags > 0{
								translated.Ranking = float32(translated.Likes) / float32(translated.Flags)
							}
							datastore.Put(c,translatedKey,&translated)
						}
						
					}	
				}
			}
		}
	}
	data, _ := json.Marshal(m)
	fmt.Fprintf(controller.W,"%s",data)
}

func (controller *FlagsController) DeleteAll(){
	c := appengine.NewContext(controller.R)
	q := datastore.NewQuery("Flag").KeysOnly()
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

func (controller *FlagsController) View(){
	//http.ServeFile(controller.W,controller.R,"views/layout.html")
	c := appengine.NewContext(controller.R)
	if token,err := datastore.DecodeKey(controller.Params[2]); err !=nil{
		controller.RenderError(err.Error())
	}else{
		video := models.Flag{}
		if err := datastore.Get(c,token,&video); err !=nil{
			controller.RenderError(err.Error())
		}else{
			
		}
	}
	//
}

func (controller *FlagsController) Home(){
	//http.ServeFile(controller.W,controller.R,"views/layout.html")
	
	controller.Layout = append(controller.Layout,"general/onflag.tmpl")
	controller.Render(true)
	//
}

func (controller *FlagsController) Search(){
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
			flags := make([]models.Flag,0,20)
			q := datastore.NewQuery("Flag").Filter("Name >=",search).Filter("Active =",true).Project("Name").Limit(20)
			if keys, err := q.GetAll(c,&flags); err !=nil{
				m["Error"] = err.Error()
			}else{
				for i,key := range keys{
					//safekeys[i] = key.Encode()
					//Get the state as well
					stateKey := key.Parent()
					state := models.State{}
					if err = datastore.Get(c,stateKey,&state); err == nil{
						flags[i].Name = strings.Title(flags[i].Name)
					}
				}
				m["Result"] = "Success"
				m["Flags"] = flags
			}
		}
	data, _ := json.Marshal(m)
	fmt.Fprintf(controller.W,"%s",data)
}

//AJAX only 
func (controller *FlagsController) MobileIndex(){
	c := appengine.NewContext(controller.R)
	m := make(map[string]interface{})
	m["Result"] = "Failure"

			offset,_ := strconv.Atoi(controller.R.FormValue("offset"))
			flags := make([]models.Flag,0,50)
			var q *datastore.Query
				q = datastore.NewQuery("Flag").Filter("Active =",true).Order("-Flags").Offset(offset).Limit(10)

			if keys, err := q.GetAll(c,&flags); err != nil{
				m["Error"] = err.Error()
			}else{
				m["Result"] = "Success"
				safekeys := make([]string,len(flags))
				for i,key := range keys{
					safekeys[i] = key.Encode()
				}
				m["SafeKeys"] = safekeys
				m["Flags"] = flags
			}	
	data, _ := json.Marshal(m)
	fmt.Fprintf(controller.W,"%s",data)
}

func (controller *FlagsController) TopTen(){
	c := appengine.NewContext(controller.R)
	m := make(map[string]interface{})
	m["Result"] = "Failure"
	key := controller.R.FormValue("key")
	if key == "39234hwersKLnsdk2!sdkdfjd"{

		if langKey, err := datastore.DecodeKey(controller.R.FormValue("locale")); err != nil{
			m["Error"] = err.Error()
		}else if langKey != nil{
			flags := make([]models.Flag,0,10)
			q := datastore.NewQuery("Flag").Ancestor(langKey).Filter("Active =",true).Order("-Flags").Limit(10)
			if keys, err := q.GetAll(c,&flags); err != nil{
				m["Error"] = err.Error()
			}else{
				m["Result"] = "Success"
				safekeys := make([]string,len(flags))
				for i,key := range keys{
					safekeys[i] = key.Encode()
				}
				m["SafeKeys"] = safekeys
				m["Flags"] = flags
			}	
		}else{

		}
	}
	data, _ := json.Marshal(m)
	fmt.Fprintf(controller.W,"%s",data)
}

//Mobile add - Upload from phones
func (controller *FlagsController) MobileAdd(){
	c := appengine.NewContext(controller.R)
	m := make(map[string]interface{})
	m["Result"] = "Failure"
	m["Error"] = "You aren't logged in"
	if !controller.CheckSession(){
		//Checks to see if this is from the web
		controller.UserToken = controller.R.FormValue("token")
	}
		if user,userKey,err := controller.GetUser(c); err !=nil{
			m["Error"] = "You aren't logged in"
		}else if userKey != nil{
			//Check to make sure you haven't already flagd it!
			flags := make([]models.Flag,0,1)
			q := datastore.NewQuery("Flag").Ancestor(userKey).Filter("Controller =",controller.R.FormValue("controller")).Filter("Action =",controller.R.FormValue("safekey")).Limit(1)
			if flagKeys,err := q.GetAll(c,&flags); err == nil && len(flagKeys)>0{
				if flags[0].Active == false{
					flags[0].Active = true
					flags[0].Modified = time.Now()
					if newFlagKey,err := datastore.Put(c,flagKeys[0],&flags[0]); err != nil{
						m["Error"] = err.Error()
					}else{
						m["Result"] = "Success"
						m["SafeKey"] = newFlagKey.Encode()
					}
				}else{
					m["Result"] = "Partial"
					m["SafeKey"] = flagKeys[0].Encode()
				}
			}
					//Get the country and put it in here
					flag := models.Flag{
						Name: user.Name,
						Controller: controller.R.FormValue("controller"),
						Action:controller.R.FormValue("safekey"),
						Created:time.Now(),
						Modified:time.Now(),
						Active:true,
					}
					vKey := datastore.NewIncompleteKey(c,"Flag",userKey)
					if countryKey, err := datastore.Put(c,vKey,&flag); err != nil{
						m["Error"] = err.Error()
					}else{
						m["Result"] = "Success"
						m["SafeKey"] = countryKey.Encode()
					}	
			if controller.R.FormValue("controller") == "translateds" && m["Result"] == "Success"{
				if translatedKey,err := datastore.DecodeKey(controller.R.FormValue("safekey")); err == nil{
					translated := models.Translated{}
					if err = datastore.Get(c,translatedKey,&translated); err == nil{
						translated.Flags = translated.Flags+1
						m["Flags"] = translated.Flags
						if translated.Flags == 0{
							translated.Ranking = float32(translated.Likes)
						}else{
							translated.Ranking = float32(translated.Likes) / float32(translated.Flags)
						}
						if _,err = datastore.Put(c,translatedKey,&translated); err != nil{
							m["Error"] = err.Error()
						}else{
							m["Flags"] = translated.Flags
						}
					}
				}
			}
		}
	data, _ := json.Marshal(m)
	fmt.Fprintf(controller.W,"%s",data)
}

func (controller *FlagsController) Add(){
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
				if flagKey,err := controller.FindFlag(c,userKey.Encode(),lineKey); err !=nil{
					m["Error"] = err.Error()
				}else if flagKey != nil{
					m["Result"] = "Success"
					m["SafeKey"] = flagKey.Encode()
				}else{
					//Get the country and put it in here
					flag := models.Flag{
						Name: user.Name,
						Controller: "lines",
						Action: controller.R.FormValue("safekey"),
						Created:time.Now(),
						Modified:time.Now(),
						User: userKey.Encode(),
						Active:true,
					}
					vKey := datastore.NewIncompleteKey(c,"Flag",lineKey)
					if countryKey, err := datastore.Put(c,vKey,&flag); err != nil{
						m["Error"] = err.Error()
					}else{
						m["Result"] = "Success"

						m["SafeKey"] = countryKey.Encode()
						line := models.Line{}
						if err = datastore.Get(c,lineKey,&line); err != nil{
							m["Error"] = err.Error()
						}else{
							line.Flags += 1
							datastore.Put(c,lineKey,&line)
							m["Flags"] = line.Flags
						}
					}
				}
			}
		}
	data, _ := json.Marshal(m)
	fmt.Fprintf(controller.W,"%s",data)
}

func (controller *FlagsController) FindFlag(c appengine.Context,userstring string, controllerKey *datastore.Key) (*datastore.Key,error){

	flags := make([]models.Flag,0,1)
	q := datastore.NewQuery("Flag").Ancestor(controllerKey).Filter("User =",userstring).Limit(1)
	if keys, err := q.GetAll(c,&flags); err != nil{
		return nil,err
	}else if len(keys)>0{
		if flags[0].Active == false {
			flags[0].Active = true
			flags[0].Modified = time.Now()
			if _,err := datastore.Put(c,keys[0],&flags[0]); err != nil{
				return nil,err
			}else{
				line := models.Line{}
				if err = datastore.Get(c,controllerKey,&line);err == nil{
					line.Flags+=1
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

func (controller *FlagsController) CreateMini(){
	c := appengine.NewContext(controller.R)
	m := make(map[string]interface{})
	m["Result"] = "Failure"
	q := datastore.NewQuery("Language").Filter("Name =","English").Limit(1)
	langs := make([]models.Language,0,1)
	if keys,err := q.GetAll(c,&langs); err != nil{
		m["Error"] = err.Error()
	}else if len(keys)>0{
		key := keys[0]
		flag := models.Flag{
			Name:"I'm not a photographer, but I can picture me and you together.",
			Created:time.Now(),
			Modified:time.Now(),
			Active:true,
		}
		lKey := datastore.NewIncompleteKey(c,"Flag",key)
		datastore.Put(c,lKey,&flag)
	}
	data, _ := json.Marshal(m)
	fmt.Fprintf(controller.W,"%s",data)
}



func (controller *FlagsController) MobileView(){
	c := appengine.NewContext(controller.R)
	m := make(map[string]interface{})
	m["Result"] = "Failure"
	m["Error"] = "No authorization"

			if flagKey,err := datastore.DecodeKey(controller.R.FormValue("safekey")); err != nil{
				m["Error"] = err.Error()
			}else{
				flag := models.Flag{}
				if err = datastore.Get(c,flagKey,&flag); err != nil{
					m["Error"] = err.Error()
				}else{
					m["Result"] = "Success"
					m["Flag"] = flag
					m["Translated"] = ""
				}
			}
	data, _ := json.Marshal(m)
	fmt.Fprintf(controller.W,"%s",data)
}

//Mobile add - Upload from phones
func (controller *FlagsController) MobileEdit(){
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
			flags := make([]models.Flag,0,1)
			q := datastore.NewQuery("Flag").Ancestor(userKey).Limit(1)
			if keys,err := q.GetAll(c,&flags); err != nil{
				m["Error"] = err.Error()+" Step 3"
			}else if len(keys)>0{
					video := flags[0]
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



func (controller *FlagsController) UpdateAll(){

	c := appengine.NewContext(controller.R)
	m := make(map[string]interface{})
	m["Result"] = "Failure"
	m["Error"] = "No authorization"
	q := datastore.NewQuery("Flag").Filter("Active =",false)
	flags := make([]models.Flag,0,100)
	if keys, err := q.GetAll(c,&flags); err != nil{
		println(err.Error())
	}else{
		for i,_ := range keys{
			flags[i].Active = false
		}
		datastore.PutMulti(c,keys,flags)
		m["Result"] = "Success"
	}
	data, _ := json.Marshal(m)
	fmt.Fprintf(controller.W,"%s",data)
	/*
	c := appengine.NewContext(controller.R)
	q := datastore.NewQuery("Flag").Limit(100)
	assignments := make([]models.Flag,0,100)
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



func (controller *FlagsController) Render(show bool) error{
	//controller.Controller.Render(show)
	controller.Controller.Render(show)
	return nil
}