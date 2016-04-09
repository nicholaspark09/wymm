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

type LinesController struct{
	Controller
}

func NewLinesController() *LinesController{
	return &LinesController{
	}
}

func (controller *LinesController) Init(r *http.Request, w http.ResponseWriter){
	controller.R = r
	controller.W = w
	controller.Data = make(map[interface{}]interface{})
	controller.C = appengine.NewContext(r)
	controller.Layout = make([]string,0)
	controller.TemplateNames = "body"
	controller.Name = "lines"
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

func (controller *LinesController) Serve(action string){
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
			}else if action =="mylines"{
				controller.MyLines()
			}else if action =="myindex"{
				controller.MyIndex()
			}
}

func (controller *LinesController) MobileServe(action string){
	
}


func (controller *LinesController) Index(){

	controller.Layout = append(controller.Layout,"lines/index.tmpl")

	controller.Render(true)
}

func (controller *LinesController) Edit(){
	controller.Render(true)
}

func (controller *LinesController) Delete(){
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
		}else if videoKey!=nil{
			verified := false
			if videoKey.Parent().Encode() == userKey.Encode(){
				verified = true
			}else if user.Group ==1{
				verified = true
			}
			if verified{
				line := models.Line{}
				if err = datastore.Get(c,videoKey,&line); err != nil{
					m["Error"] = err.Error()
				}else{
					line.Modified = time.Now()
					line.Active = false
					if _, err := datastore.Put(c,videoKey,&line); err !=nil{
						m["Error"] = err.Error()
					}else{
						m["Result"] = "Success"
						//Set all the likes and dislikes to delte
					}
				}
			}
		}
	}
	data, _ := json.Marshal(m)
	fmt.Fprintf(controller.W,"%s",data)
}

func (controller *LinesController) DeleteAll(){
	/*
	c := appengine.NewContext(controller.R)
	q := datastore.NewQuery("Line").KeysOnly()
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

func (controller *LinesController) View(){
	//http.ServeFile(controller.W,controller.R,"views/layout.html")
	c := appengine.NewContext(controller.R)
	if token,err := datastore.DecodeKey(controller.Params[2]); err !=nil{
		controller.RenderError(err.Error())
	}else{
		video := models.Line{}
		if err := datastore.Get(c,token,&video); err !=nil{
			controller.RenderError(err.Error())
		}else{
			
		}
	}
	//
}

func (controller *LinesController) Home(){
	//http.ServeFile(controller.W,controller.R,"views/layout.html")
	
	controller.Layout = append(controller.Layout,"general/online.tmpl")
	controller.Render(true)
	//
}

func (controller *LinesController) Search(){
	c := appengine.NewContext(controller.R)
	m := make(map[string]interface{})
	m["Result"] = "Failure"
	if !controller.CheckSession(){
		//Checks to see if this is from the web
		controller.UserToken = controller.R.FormValue("token")
	}
	_,userKey,_ := controller.GetUser(c)
	userkey := ""
	if userKey != nil{
		userkey = userKey.Encode()
	}
			offset,_ := strconv.Atoi(controller.R.FormValue("offset"))
			search := strings.ToLower(controller.R.FormValue("query"))
			locale := controller.R.FormValue("locale")
			translated := controller.R.FormValue("translated")
			lines := make([]models.Line,0,10)
			var q *datastore.Query
			order := controller.R.FormValue("order")
			queryOrder := "-Created"
			switch order{
				case "1":
					queryOrder = "-Views"
					break
				case "2":
					queryOrder = "-Likes"
					break
				case "3":
					queryOrder = "-Dislikes"
					break
				case "4":
					queryOrder = "-Created"
					break
				case "5":
					queryOrder = "Created"
					break
				default: 
					break
			}
			if search == ""{
				q = datastore.NewQuery("Line").Filter("Active =",true).Filter("Locale =",locale).Order(queryOrder).Offset(offset).Limit(10)
			}else{
				q = datastore.NewQuery("Line").Filter("Name >=",search).Filter("Locale =",locale).Filter("Active =",true).Offset(offset).Limit(10)
			}
			if keys, err := q.GetAll(c,&lines); err !=nil{
				m["Error"] = err.Error()
			}else{
				safekeys := make([]string,len(keys))
				for i,key := range keys{
					//safekeys[i] = key.Encode()
					//Get the state as well
					safekeys[i] = key.Encode()
					translateds := make([]models.Translated,0,1)
					q = datastore.NewQuery("Translated").Ancestor(key).Filter("Locale =",translated).Filter("Active =",true).Order("-Ranking").Limit(1)
					if translatedKeys,err := q.GetAll(c,&translateds); err != nil{
						m["Error"] = err.Error()
					}else if len(translatedKeys) > 0{
						lines[i].Translated = translatedKeys[0].Encode()
					}
				}
				m["Result"] = "Success"
				m["Lines"] = lines
				m["SafeKeys"] = safekeys
				controller.SaveInfo(c,controller.R.FormValue("device"),search,userkey)
			}
	data, _ := json.Marshal(m)
	fmt.Fprintf(controller.W,"%s",data)
}

//AJAX only 
func (controller *LinesController) MobileIndex(){
	c := appengine.NewContext(controller.R)
	m := make(map[string]interface{})
	m["Result"] = "Failure"
	if !controller.CheckSession(){
		//Checks to see if this is from the web
		controller.UserToken = controller.R.FormValue("token")
	}
	_,userKey,_ := controller.GetUser(c)
	userkey := ""
	if userKey != nil{
		userkey = userKey.Encode()
	}
			offset,_ := strconv.Atoi(controller.R.FormValue("offset"))
			locale := controller.R.FormValue("locale")
			lines := make([]models.Line,0,50)
			var q *datastore.Query
			if locale == ""{
				q = datastore.NewQuery("Line").Filter("Active =",true).Order("-Likes").Offset(offset).Limit(10)
			}else{
				q = datastore.NewQuery("Line").Filter("Active =",true).Filter("Locale =",locale).Order("-Likes").Offset(offset).Limit(10)
			}
			if keys, err := q.GetAll(c,&lines); err != nil{
				m["Error"] = err.Error()
			}else{
				m["Result"] = "Success"
				safekeys := make([]string,len(lines))
				for i,key := range keys{
					safekeys[i] = key.Encode()
				}
				m["SafeKeys"] = safekeys
				m["Lines"] = lines
				controller.SaveInfo(c,controller.R.FormValue("device"),"",userkey)
			}	
	data, _ := json.Marshal(m)
	fmt.Fprintf(controller.W,"%s",data)
}

func (controller *LinesController) TopTen(){
	c := appengine.NewContext(controller.R)
	m := make(map[string]interface{})
	m["Result"] = "Failure"
	key := controller.R.FormValue("key")
	if key == "39234hwersKLnsdk2!sdkdfjd"{

			lines := make([]models.Line,0,25)
			q := datastore.NewQuery("Line").Filter("Locale =",controller.R.FormValue("locale")).Filter("Active =",true).Order("-Likes").Limit(25)
			if keys, err := q.GetAll(c,&lines); err != nil{
				m["Error"] = "Error was here with locale: "+controller.R.FormValue("locale")+" and error: "+err.Error()
			}else{
				m["Result"] = "Success"
				safekeys := make([]string,len(lines))
				for i,key := range keys{
					safekeys[i] = key.Encode()
				}
				m["SafeKeys"] = safekeys
				m["Lines"] = lines
			}	
	}
	data, _ := json.Marshal(m)
	fmt.Fprintf(controller.W,"%s",data)
}

//Mobile add - Upload from phones
func (controller *LinesController) MobileAdd(){
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
					line := models.Line{
						Name: name,
						Locale: controller.R.FormValue("locale"),
						Created:time.Now(),
						Modified:time.Now(),
						Likes:0,
						Dislikes:0,
						Views:0,
						Rank:0,
						User: userKey.Encode(),
						Active:true,
					}
					vKey := datastore.NewIncompleteKey(c,"Line",userKey)
					if countryKey, err := datastore.Put(c,vKey,&line); err != nil{
						m["Error"] = err.Error()
					}else{
						m["Result"] = "Success"
						m["SafeKey"] = countryKey.Encode()
					}	
			
		}
	data, _ := json.Marshal(m)
	fmt.Fprintf(controller.W,"%s",data)
}

func (controller *LinesController) Add(){
	c := appengine.NewContext(controller.R)
	m := make(map[string]interface{})
	m["Result"] = "Failure"
	m["Error"] = "You aren't logged in <a href=\"/users/login\">Login</a>"
	if !controller.CheckSession(){
		//Checks to see if this is from the web
		controller.UserToken = controller.R.FormValue("token")
	}
		if _,userKey,err := controller.GetUser(c); err ==nil && userKey != nil{

				name := controller.R.FormValue("name")
					//Get the country and put it in here
					line := models.Line{
						Name: name,
						Locale: controller.R.FormValue("locale"),
						Created:time.Now(),
						Modified:time.Now(),
						Likes:0,
						Dislikes:0,
						Views:0,
						Rank:0,
						User: userKey.Encode(),
						Active:true,
					}
					vKey := datastore.NewIncompleteKey(c,"Line",userKey)
					if countryKey, err := datastore.Put(c,vKey,&line); err != nil{
						m["Error"] = err.Error()
					}else{
						m["Result"] = "Success"
						m["SafeKey"] = countryKey.Encode()
					}	
			
		}
	data, _ := json.Marshal(m)
	fmt.Fprintf(controller.W,"%s",data)
}

func (controller *LinesController) CreateMini(){
	c := appengine.NewContext(controller.R)
	m := make(map[string]interface{})
	m["Result"] = "Failure"
	q := datastore.NewQuery("Language").Filter("Name =","English").Limit(1)
	langs := make([]models.Language,0,1)
	if keys,err := q.GetAll(c,&langs); err != nil{
		m["Error"] = err.Error()
	}else if len(keys)>0{
		key := keys[0]
		line := models.Line{
			Name:"I'm not a photographer, but I can picture me and you together.",
			Translated:"You need help",
			Likes:0,
			Dislikes:0,
			Created:time.Now(),
			Modified:time.Now(),
			Flags:0,
			Locale: "en_US",
			Active:true,
		}
		lKey := datastore.NewIncompleteKey(c,"Line",key)
		datastore.Put(c,lKey,&line)
	}
	data, _ := json.Marshal(m)
	fmt.Fprintf(controller.W,"%s",data)
}


func (controller *LinesController) MobileView(){
	c := appengine.NewContext(controller.R)
	m := make(map[string]interface{})
	m["Result"] = "Failure"
	m["Error"] = "No authorization"
	
			locale := controller.R.FormValue("locale")
			firsttwo := locale[:2]
			if firsttwo == "en"{
				locale = "en_US"
			}
			if lineKey,err := datastore.DecodeKey(controller.R.FormValue("safekey")); err != nil{
				m["Error"] = err.Error()
			}else{
				line := models.Line{}
				if err = datastore.Get(c,lineKey,&line); err != nil{
					m["Error"] = err.Error()
				}else{
					m["Result"] = "Success"
					m["Line"] = line
					m["Translated"] = ""
					m["TranslatedKey"] = ""
					m["Like"] = ""
					m["Dislike"] = ""
					userkey := ""
					line.Views = line.Views+1
					if _,err = datastore.Put(c,lineKey,&line); err !=nil{
						m["Error"] = err.Error()
					}
					t := make([]models.Translated,0,1)
					q := datastore.NewQuery("Translated").Ancestor(lineKey).Filter("Locale =",locale).Order("-Likes").Limit(1)
					if ks,err := q.GetAll(c,&t); err != nil{
						m["Error"] =err.Error()
					}else if len(ks)>0{
						m["Translated"] = t[0]
						m["TranslatedKey"] = ks[0].Encode()
					}
					if !controller.CheckSession(){
							//Checks to see if this is from the web
						controller.UserToken = controller.R.FormValue("token")
					}
					if _,userKey,err := controller.GetUser(c); err ==nil && userKey != nil{
						userkey = userKey.Encode()
						likes := make([]models.Like,0,1)
						q = datastore.NewQuery("Like").Ancestor(lineKey).Filter("User =",userKey.Encode()).Filter("Active =",true).Limit(1)
						if likeKeys, err := q.GetAll(c,&likes); err != nil{
							m["Error"] = err.Error()
						}else if len(likeKeys)>0{
							m["Like"] = likeKeys[0].Encode()
						}else{
							m["Message"] = "You couldn't find it"
						}
						dislikes := make([]models.Dislike,0,1)
						q = datastore.NewQuery("Dislike").Ancestor(lineKey).Filter("User =",userKey.Encode()).Filter("Active =",true).Limit(1)
						if dislikeKeys, err := q.GetAll(c,&dislikes); err != nil{
							m["Error"] = err.Error()
						}else if len(dislikeKeys)>0{
							m["Dislike"] = dislikeKeys[0].Encode()
						}else{
							m["Message"] = "You couldn't find it"
						}
					}
					controller.SaveInfo(c,controller.R.FormValue("device"),controller.R.FormValue("safekey"),userkey)
				}
			}
	data, _ := json.Marshal(m)
	fmt.Fprintf(controller.W,"%s",data)
}

//Mobile add - Upload from phones
func (controller *LinesController) MobileEdit(){
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
			if lineKey,err := datastore.DecodeKey(controller.R.FormValue("safekey")); err != nil{
				m["Error"] = err.Error()
			}else{
				line := models.Line{}
				if err = datastore.Get(c,lineKey,&line); err != nil{
					m["Error"] = err.Error()
				}else{
					line.Modified = time.Now()
					line.Name = controller.R.FormValue("name")
					line.Locale = controller.R.FormValue("locale")
					if key, err := datastore.Put(c,lineKey,&line); err != nil{
						m["Error"] = err.Error()
					}else{
						m["Result"] = "Success"
						m["SafeKey"] = key.Encode()
					}
				}
			}
		}
	data, _ := json.Marshal(m)
	fmt.Fprintf(controller.W,"%s",data)
}

func (controller *LinesController) MyLines(){

	if !controller.CheckSession(){
		controller.RenderError("You aren't logged in.")
	}else{
		controller.Layout = append(controller.Layout,"lines/mylines.tmpl")
		controller.Render(true)
	}
}


func (controller *LinesController) MyIndex(){
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
			offset,_ := strconv.Atoi(controller.R.FormValue("offset"))
			lines := make([]models.Line,0,10)

			q := datastore.NewQuery("Line").Filter("Active =",true).Filter("User =",userKey.Encode()).Order("-Created").Offset(offset).Limit(10)
			if keys, err := q.GetAll(c,&lines); err != nil{
				m["Error"] = err.Error()
			}else{
				m["Result"] = "Success"
				safekeys := make([]string,len(lines))
				for i,key := range keys{
					safekeys[i] = key.Encode()
				}
				m["SafeKeys"] = safekeys
				m["Lines"] = lines
			}	

		}
	data, _ := json.Marshal(m)
	fmt.Fprintf(controller.W,"%s",data)
}


func (controller *LinesController) UpdateAll(){
	/*
	c := appengine.NewContext(controller.R)
	m := make(map[string]interface{})
	m["Result"] = "Failure"
	m["Error"] = "No authorization"
	q := datastore.NewQuery("Line").Filter("Active =",false)
	lines := make([]models.Line,0,100)
	if keys, err := q.GetAll(c,&lines); err != nil{
		println(err.Error())
	}else{
		for i,_ := range keys{
			lines[i].Views = 0
		}
		datastore.PutMulti(c,keys,lines)
		m["Result"] = "Success"
	}
	data, _ := json.Marshal(m)
	fmt.Fprintf(controller.W,"%s",data)
	/*
	c := appengine.NewContext(controller.R)
	q := datastore.NewQuery("Line").Limit(100)
	assignments := make([]models.Line,0,100)
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



func (controller *LinesController) Render(show bool) error{
	//controller.Controller.Render(show)
	controller.Controller.Render(show)
	return nil
}