package controllers

import (
	"net/http"
	"appengine"
	"fmt"
	"encoding/json"
	"appengine/datastore"
	"time"
	"github.com/parkn09/wymm/src/models"
	"net/url"
	"github.com/parkn09/wymm/src/custom"
	"appengine/mail"
	"strings"
	"errors"
	"strconv"
)

type UsersController struct{
	Controller
}

func NewUsersController() *UsersController{
	return &UsersController{
	}
}

func (controller *UsersController) Init(r *http.Request, w http.ResponseWriter){
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
	cookie, err := controller.R.Cookie("Language")
	if err!= nil{
		controller.Locale = "en"
	}else if cookie.Value != ""{
		controller.Locale = cookie.Value
	}
	controller.Data["Language"] = controller.Locale
}

func (controller *UsersController) Serve(action string){
			if action == "view"{
				controller.View()
			}else if action =="add"{
				if controller.CheckSession(){
					//controller.Add()
				}else{
					controller.Render(true)
				}
			}else if action=="index"{
				controller.Index()
			}else if action =="complete"{
				controller.Complete()
			}else if action == "login"{
				controller.Login()
			}else if action == "logout"{
				controller.Logout()
			}else if action =="deleteall"{
				controller.DeleteAll()
			}else if action == "home"{
				controller.Home()
			}else if action == "languages"{
				controller.Languages()
			}else if action =="addlittle"{
				controller.AddLittle()
			}else if action =="phonelogin"{
				controller.PhoneLogin()
			}else if action =="register"{
				controller.Register()
			}else if action =="facebook"{
				controller.Facebook()
			}else if action =="delete"{
				controller.Delete()
			}else if action =="facebooklogin"{
				controller.FacebookLogin()
			}else if action =="mobileregister"{
				controller.MobileRegister()
			}else if action =="updateall"{
				controller.UpdateAll()
			}else if action == "loginverified"{
				controller.LoginVerified()
			}else if action =="mobilelogin"{
				controller.MobileLogin()
			}else if action =="policy"{
				controller.Policy()
			}else if action =="iphonelogin"{
				controller.IphoneLogin()
			}else if action =="openlink"{
				controller.OpenLink()
			}
}

func (controller *UsersController) MobileServe(action string){
	
}


func (controller *UsersController) Index(){
	controller.Render(true)
}

func (controller *UsersController) Edit(){
	controller.Render(true)
}

func (controller *UsersController) Delete(){
	c := appengine.NewContext(controller.R)
	q := datastore.NewQuery("User").Filter("Email =","butterphinger@hotmail.com").KeysOnly()
	if keys,err := q.GetAll(c,nil); err == nil{
		if len(keys)>0{
			q = datastore.NewQuery("Role").Filter("User =",keys[0].Encode()).KeysOnly()
			if roleKeys, err := q.GetAll(c,nil); err == nil{
				if len(roleKeys)>0{
					datastore.DeleteMulti(c,roleKeys)
				}
			}
			if err := datastore.DeleteMulti(c,keys); err != nil{
				println("Error")
				println(err.Error())
			}
		}
	}
}

func (controller *UsersController) View(){
	//http.ServeFile(controller.W,controller.R,"views/layout.html")
	if controller.CheckSession(){
		controller.Render(true)
	}else{
		controller.RedirectToLogin()
	}

	//
}

/*
		!!! This is for apps only !!!
			All web logins must go to the regular Login() method 

*/
func (controller *UsersController) MobileAdd(){
	
}

//This is the final step to login an Android App User
func (controller *UsersController) CompleteApp(){
	
}

//Login the user
/*
		This is the internet only Login from a browser, not an native app

*/
func (controller *UsersController) Login(){

	if controller.R.Method=="POST"{
		m := make(map[string]interface{})
		email := controller.R.FormValue("email")
		
		if email ==""{
			m["Error"] = "No email given"
		}else{
			email = strings.TrimSpace(email)
			c := appengine.NewContext(controller.R)
			if _,_, err := controller.CheckUser(c,email); err != nil{
				m["Error"] = err.Error()
				//check if it's you nick!
				if email == "nicholaspark09@gmail.com" || email == "peteracorina@gmail.com"{
					name := "Nicholas Park"
					if email == "peteracorina@gmail.com"{
						name = "Peter Corina"
					}
					if _, err := controller.CreateUser(c,name,email); err !=nil{
						m["Error"] = err.Error()
					}else{
						token, err := controller.CreateToken(c,email)
						if err != nil{
							m["Error"] = err.Error()
						}else{
							if err = controller.SendWebLink(c,token,email); err != nil{
									m["Error"] = err.Error()
							}else{
									m["Result"] = "Success"
							}
						}
					}
				}else{
					m["Result"] = "Failure"
					m["Error"] = "Sorry, but you aren't registered. Please sign up."
				}
			}else{
				token, err := controller.CreateToken(c,email)
				if err != nil{
					m["Error"] = err.Error()
				}else{
					if err = controller.SendWebLink(c,token,email); err != nil{
						m["Error"] = err.Error()
					}else{
						m["Result"] = "Success"
						m["SafeKey"] = token
					}
				}
			}

		}



		data, _ := json.Marshal(m)
		fmt.Fprintf(controller.W,"%s",data)
	}else{
		controller.Layout = append(controller.Layout,"users/login.tmpl")
		controller.Render(true)
	}
}

//Create a login for iphones
func (controller *UsersController) IphoneLogin(){

	if controller.R.Method=="POST"{
		m := make(map[string]interface{})
		email := controller.R.FormValue("email")
		
		if email ==""{
			m["Error"] = "No email given"
		}else{
			email = strings.TrimSpace(email)
			c := appengine.NewContext(controller.R)
			if _,userKey, err := controller.CheckUser(c,email); err != nil{
				m["Error"] = err.Error()
				//check if it's you nick!
				if email == "nicholaspark09@gmail.com"{
					name := "Nicholas Park"
					if userKey, err := controller.CreateUser(c,name,email); err !=nil{
						m["Error"] = err.Error()
					}else{
						token, err := controller.CreateSession(c,userKey)
						if err != nil{
							m["Error"] = err.Error()
						}else{
							if err = controller.SendIphoneLink(c,token,email); err != nil{
									m["Error"] = err.Error()
							}else{
									m["Result"] = "Success"
							}
						}
					}
				}else{
					m["Result"] = "Failure"
					m["Error"] = "Sorry, but you aren't registered. Please sign up."
				}
			}else{
				token, err := controller.CreateSession(c,userKey)
				if err != nil{
					m["Error"] = err.Error()
				}else{
					if err = controller.SendIphoneLink(c,token,email); err != nil{
						m["Error"] = err.Error()
					}else{
						m["Result"] = "Success"
						m["SafeKey"] = token
					}
				}
			}

		}



		data, _ := json.Marshal(m)
		fmt.Fprintf(controller.W,"%s",data)
	}else{
		controller.Layout = append(controller.Layout,"users/login.tmpl")
		controller.Render(true)
	}
}


//Login the user
/*
		This is the internet only Login from a browser, not an native app

*/
func (controller *UsersController) MobileLogin(){

	if controller.R.Method=="POST"{
		m := make(map[string]interface{})
		email := controller.R.FormValue("email")
		
		if email ==""{
			m["Error"] = "No email given"
		}else{
			email = strings.TrimSpace(email)
			c := appengine.NewContext(controller.R)
			if _,userKey, err := controller.CheckUser(c,email); err != nil{
				m["Error"] = err.Error()
				//check if it's you nick!
				if email == "nicholaspark09@gmail.com"{
					name := "Nicholas Park"
					if userKey, err := controller.CreateUser(c,name,email); err !=nil{
						m["Error"] = err.Error()
					}else{
						token, err := controller.CreateSession(c,userKey)
						if err != nil{
							m["Error"] = err.Error()
						}else{
							if err = controller.SendAppLink(c,token,email); err != nil{
									m["Error"] = err.Error()
							}else{
									m["Result"] = "Success"
							}
						}
					}
				}else{
					m["Result"] = "Failure"
					m["Error"] = "Sorry, but you aren't registered. Please sign up."
				}
			}else{
				token, err := controller.CreateSession(c,userKey)
				if err != nil{
					m["Error"] = err.Error()
				}else{
					if err = controller.SendAppLink(c,token,email); err != nil{
						m["Error"] = err.Error()
					}else{
						m["Result"] = "Success"
						m["SafeKey"] = token
					}
				}
			}

		}



		data, _ := json.Marshal(m)
		fmt.Fprintf(controller.W,"%s",data)
	}else{
		controller.Layout = append(controller.Layout,"users/login.tmpl")
		controller.Render(true)
	}
}


func (controller *UsersController) Register(){
	c := appengine.NewContext(controller.R)
	if controller.R.Method=="POST"{
		m := make(map[string]interface{})
		name := strings.TrimSpace(controller.R.FormValue("name"))
		email := strings.TrimSpace(controller.R.FormValue("email"))

		var userKey *datastore.Key
		var err error
		if _,userKey,err = controller.CheckUser(c,email); err !=nil{
			//Couldn't find the user, you can now register them!
			if userKey, err = controller.CreateUser(c,name,email);err != nil{
				m["Error"] = err.Error()
			}
		}
		//Create the role
		if userKey != nil{
			token, err := controller.CreateToken(c,email)
						if err != nil{
							m["Error"] = err.Error()
						}else{
							if err = controller.SendWebLink(c,token,email); err != nil{
									m["Error"] = err.Error()
							}else{
									m["Result"] = "Success"
							}
						}
			/*
			if campusKey, err := datastore.DecodeKey(controller.R.FormValue("campus")); err != nil{
				m["Error"] = err.Error()
			}else{
				role := models.Role{
					Name: name,
					Level:level,
					Phone:phone,
					User: userKey.Encode(),
					Created:time.Now(),
					Modified:time.Now(),
					ParentKey:controller.R.FormValue("campus"),
					Controller:"campuses",
					Active:true,
				}
				rKey := datastore.NewIncompleteKey(c,"Role",campusKey);
				if _, err := datastore.Put(c,rKey,&role); err != nil{
					m["Error"] = err.Error()
				}else{
					token, err := controller.CreateToken(c,email)
					if err != nil{
						m["Error"] = err.Error()
					}else{
						if err = controller.SendWebLink(c,token,email); err != nil{
							m["Error"] = err.Error()
						}else{
							m["Result"] = "Success"
							m["Message"] = "Thank you. Please check your email for a login link at "+email
						}
					}
				}
			}	
			*/
		}
		data, _ := json.Marshal(m)
		fmt.Fprintf(controller.W,"%s",data)
	}else if controller.R.Method=="GET"{
		controller.Layout = append(controller.Layout,"users/register.tmpl")
		controller.Render(true)
	}
}

//Mobile Register
//This is without a facebook account
func (controller *UsersController) MobileRegister(){
	c := appengine.NewContext(controller.R)
	if controller.R.Method=="POST"{
		m := make(map[string]interface{})
		name := strings.TrimSpace(controller.R.FormValue("name"))
		email := strings.TrimSpace(controller.R.FormValue("email"))

		m["Email"] = email
		m["Name"] = name
		if email != ""{
			if user,userKey,err := controller.CheckUser(c,email); err != nil{
				//Couldn't find the user, you can now register them!
					//No user, so create them!
				if userKey,err = controller.CreateUser(c,name,email); err != nil{
					m["Error"] = err.Error()
				}else if userKey == nil{
					m["Error"] = "Sorry, but we couldn't create the key"
				}else{
					//Send a web link!
					//Create a session and send the token back!
					if token,err := controller.CreateSession(c,userKey); err != nil{
						m["Error"] = err.Error()
					}else{
						m["Result"] = "Success"
						m["SafeKey"] = token
					}
				}
			}else if user.Group != -1{
				m["Error"] = "Already registered"
				//Send a weblink!
				//Send a web link!
					//Create a session and send the token back!
				if token,err := controller.CreateSession(c,userKey); err != nil{
						m["Error"] = err.Error()
				}else if err = controller.SendAppLink(c,token,email); err != nil{
							m["Error"] = err.Error()
				}else{
							m["Result"] = "Partial"
							m["SafeKey"] = token
							m["Error"] = "User was already registered. Sending a login link to email"
				}
			}	
		}else{
			m["Error"] = "Your email was blank"
		}
		data, _ := json.Marshal(m)
		fmt.Fprintf(controller.W,"%s",data)
	}else{
		fmt.Fprintf(controller.W,"%s","You're here")
	}
}


//Activation Link has been clicked, Create the user
	// This is to complete the registration
func (controller *UsersController) Complete(){
	//Decode the token
	c := appengine.NewContext(controller.R)
	u, err := url.Parse(controller.R.URL.String())
	if err != nil{
		panic("Problem parsing")
	}
	token := u.Query().Get("token")
	if len(token)<1 {
		
		fmt.Fprintf(controller.W,"%s","404 Page not found")
	}else{

		q := datastore.NewQuery("Temp").Filter("Temptoken = ",token).Filter("Activated = ",false).Order("-Created").Limit(10)
		temps := make([]models.Temp,0,10)
		if keys, err := q.GetAll(c, &temps); err != nil{
			c.Errorf("Fetching temps: %v", err)
			c.Infof("Received Error: %s", "The error was here")
			                return
		}else{
		
		
			for i, t := range temps{
				k := keys[i]
			
			//Updating the temp datastore entity, it doesn't matter when it's finished
			
					t.Activated = true
					_, err := datastore.Put(c,k,t)
					if err != nil{
						fmt.Println(controller.W,"%s","Couldn't find the temp token")
					}

					//temp token has been found and verified
					//Create the Actual user
					user, key, err := controller.CheckUser(c, t.Email)
					if err != nil {
						fmt.Println(controller.W,"%s",err)
						return
					}

					//User was created, now go ahead and Log them in and reroute them
					token, err := controller.CreateSession(c,key)
					if err != nil{
						fmt.Println(controller.W,"%s",err)
						return 
					}
					//User session has been created, take the token and set it into a cookie
					cookie := http.Cookie{Name:"WymmUser",Value:token,Expires:time.Now().Add(365 * 24 *time.Hour),HttpOnly:true, MaxAge:900000,Path:"/"}
					http.SetCookie(controller.W,&cookie)
					group_cookie := http.Cookie{Name:"WymmGroup",Value:string(user.Group),Expires:time.Now().Add(365 * 24 * time.Hour),HttpOnly:true,MaxAge:900000,Path:"/"}
					http.SetCookie(controller.W,&group_cookie)
					controller.UserToken = cookie.Value
					http.Redirect(controller.W,controller.R,"/lines",301)
					//Cookie has been set and the registration is done!
					//Send the user to the home page
					//Check the user, if it's an admin send them to the dashboard
					//http.Redirect(controller.W,controller.R,"/users/home",301)
					//Cookie has been set and the registration is done!
			//Send the user to the home page
			//Check the user, if it's an admin send them to the dashboard
			
				
				}
		}
	}
}


//The users home base
func (controller *UsersController) Home(){
	if controller.CheckSession(){

	}else{
		controller.RenderError("Sorry, but you aren't loggedin")
	}	
}

func (controller *UsersController) Policy(){
	controller.Layout = append(controller.Layout,"users/policy.tmpl")
	controller.Render(true)
}

//Create a user login session in datastore
//			Returns a token string that should be used for the cookie or app
func(controller *UsersController) CreateSession(c appengine.Context,key *datastore.Key) (string, error){
	//Search for previous session, delete them so that we know only this session is logged in
	Token := custom.Token{}
	token := Token.CreateToken(64)
	session := &models.Session{
		Token: token,
		Created: time.Now(),
		Modified: time.Now(),
		Online: true,
		Ttl: time.Now(),
		Active: true,
	}

	sessionKey := datastore.NewIncompleteKey(c,"Session",key)
	_, err := datastore.Put(c,sessionKey,session)
	if err != nil{
		return "", err
	}
	return token,nil
} 

//Actually Create a User Login Session and set the Cookies for Web users
func (controller *UsersController) WebLogin(){

}

//Check the user with facebook's login
func (controller *UsersController) Facebook(){
	m := make(map[string]interface{})
	m["Result"] = "Failure"
	if controller.CheckSession(){
		m["Result"] = "Success"
		m["Message"] = "Already Signed in"
	}else{
		m["Message"] = "Step here"
		name := controller.R.FormValue("name")
		oauth := controller.R.FormValue("oauth")
		id,_ := strconv.ParseInt(controller.R.FormValue("id"),10,64)
		email := controller.R.FormValue("email")
		c := appengine.NewContext(controller.R)
		users := make([]models.User,0,1)
		q := datastore.NewQuery("User").Filter("Email =",email).Filter("Active =",true).Limit(1)
		if keys,err := q.GetAll(c,&users); err != nil{
			m["Error"] = err.Error()
		}else if len(keys)>0{
			userKey := keys[0]
			if users[0].OauthUid == id{
				m["Error"] = "Already registered"
			}else{
				users[0].OauthProvider = oauth
				users[0].OauthUid = id
				if key, err := datastore.Put(c,keys[0],&users[0]); err != nil{
					m["Error"] = err.Error()
				}else{
					m["Result"] = "Success"
					userKey = key
				}
			}
			token, err := controller.CreateSession(c,userKey)
				if err != nil{
					m["Result"] = "Failure"
					m["Error"] = err.Error()
					return 
				}
								m["Result"] = "Success"
				//User session has been created, take the token and set it into a cookie
				cookie := http.Cookie{Name:"WymmUser",Value:token,HttpOnly:true, MaxAge:100000,Path:"/"}
				http.SetCookie(controller.W,&cookie)
				group_cookie := http.Cookie{Name:"WymmGroup",Value:userKey.Parent().Encode(),HttpOnly:true,MaxAge:50000,Path:"/"}
				http.SetCookie(controller.W,&group_cookie)
				controller.Controller.UserToken = cookie.Value
		}else{
			//Create a user
			if userKey,err := controller.CreateUser(c,name,email); err != nil{
				m["Error"] = err.Error()
			}else{
				user := models.User{}
				if err = datastore.Get(c,userKey,&user); err != nil{
					m["Error"] = err.Error()
				}else{
					user.OauthProvider = oauth
					user.OauthUid = id

					if userKey, err = datastore.Put(c,userKey,&user); err != nil{
						m["Error"] = err.Error()
					}else{
						m["Result"] = "Success"
						token, err := controller.CreateSession(c,userKey)
						if err != nil{
							m["Result"] = "Failure"
							m["Error"] = err.Error()
						}
								m["Result"] = "Success"
								//User session has been created, take the token and set it into a cookie
								cookie := http.Cookie{Name:"WymmUser",Value:token,HttpOnly:true, MaxAge:50000,Path:"/"}
								http.SetCookie(controller.W,&cookie)
								group_cookie := http.Cookie{Name:"WymmGroup",Value:userKey.Parent().Encode(),HttpOnly:true,MaxAge:50000,Path:"/"}
								http.SetCookie(controller.W,&group_cookie)
								controller.Controller.UserToken = cookie.Value
					}
				}
			}
		}
	}
	data,_ := json.Marshal(m)
	fmt.Fprintf(controller.W,"%s",data)
}

func (controller *UsersController) FacebookLogin(){
	m := make(map[string]interface{})
	m["Result"] = "Failure"
	c := appengine.NewContext(controller.R)
	users := make([]models.User,0,1)
	user_id,_ := strconv.ParseInt(controller.R.URL.Query()["user"][0],10,64)
	q := datastore.NewQuery("User").Filter("OauthProvider =","facebook").Filter("OauthUid =",user_id).Limit(1)
	if keys,err := q.GetAll(c,&users); err != nil{
		m["Error"] = err.Error()
	}else if len(keys)>0{
		if token,err := controller.CreateSession(c,keys[0]); err != nil{
			m["Error"] = err.Error()
		}else{
			m["Result"] = "Success"
			cookie := http.Cookie{Name:"WymmUser",Value:token,Expires:time.Now().Add(365 * 24 *time.Hour),HttpOnly:true, MaxAge:100000,Path:"/"}
			http.SetCookie(controller.W,&cookie)
			group_cookie := http.Cookie{Name:"WymmGroup",Value:keys[0].Parent().Encode(),Expires:time.Now().Add(365 * 24 * time.Hour),HttpOnly:true,MaxAge:10e0000,Path:"/"}
			http.SetCookie(controller.W,&group_cookie)
			controller.Controller.UserToken = cookie.Value
		}
	}else{
		m["Error"] = "NO user found"
	}
	data,_ := json.Marshal(m)
	fmt.Fprintf(controller.W,"%s",data)
}

//Check if the user exists
func (controller *UsersController) CheckUser(c appengine.Context,email string) (models.User,*datastore.Key,error){
		//Check for user first
		q := datastore.NewQuery("User").Filter("Email =",email).Limit(10)
		users := make([]models.User,0,1)
		err := errors.New("Couldn't find the user")
		if keys, err := q.GetAll(c, &users); err != nil{
			return models.User{Group:-1},nil,err
		}else{
			count := len(users)
			if count >0{
				return users[0],keys[0],nil
			}
		}
		return models.User{Group:-1},nil,err
}

func (controller *UsersController) CreateUser(c appengine.Context,name string, email string) (*datastore.Key,error){
	user := models.User{}
	err := errors.New("Couldn't create the user")
	//This assumes you've already checked for a user with this email
	user.Name = name
	user.Email = email
	user.Created = time.Now()
	user.Modified = time.Now()
	user.Active = true
	//This is the real web key
				//Admin only
				//groupKey, er := datastore.DecodeKey("ahJzfnRyaXBmcmllbmRzLTEwMThyEgsSBUdyb3VwGICAgID4iYwKDA")

				//Real Key
				//groupKey, er := datastore.DecodeKey("ahJzfnRyaXBmcmllbmRzLTEwMThyEgsSBUdyb3VwGICAgIDvxJ0KDA")
				//This is for offline testing only
				//groupKey, er := datastore.DecodeKey("ahRkZXZ-dHJpcGZyaWVuZHMtMTAxOHISCxIFR3JvdXAYgICAgICAgAoM")
				//Check to see if the user's nick!
	group := 3
	if email == "nicholaspark09@gmail.com" || email == "parkn@email.arizona.edu" || email =="peteracorina@gmail.com"{

			group = 1
	}
	user.Group = group
	uKey := datastore.NewIncompleteKey(c,"User",nil)
	if key, err := datastore.Put(c,uKey,&user); err ==nil{
		return key,nil
	}
	return nil,err
}


//Create a token for registration/login
func (controller *UsersController) CreateToken(c appengine.Context,email string) (string, error){
	Token := custom.Token{}
	token := Token.CreateToken(32)
	temp := models.Temp{
		Email: email,
		Temptoken: token,
		Created: time.Now(),
		Modified: time.Now(),
		Activated: false,
		Active: true,
	}
	key, err := datastore.Put(c, datastore.NewIncompleteKey(c, "Temp", nil), &temp)
    if err != nil {
        return "", err
    }
    c.Infof("The key: %v",key)
    return token, nil
}

/*
		Create an Email with the login link
			//Returns an error or nil
			//Requires an appengine Context variable and token string, email string
*/
func (controller *UsersController) SendWebLink(c appengine.Context,token string, email string) error{
		var Url *url.URL
    	path := "www.wymm-1242.appspot.com"
    	Url, err := url.Parse(path)
    	if err !=nil{
    		return err
    	}
    	Url.Path +="/users/complete/"
    	parameters := url.Values{}
    	parameters.Add("token",token)
    	Url.RawQuery = parameters.Encode()
    	c.Infof("Requested token: %s", Url.String())
    	confirmMessage := `Welcome to Lines. Please click on the link to login: %s`
    	addr := email
        msg := &mail.Message{
                Sender:  "nicholaspark09@gmail.com",
                To:      []string{addr},
                Subject: "Confirm your registration with Lines",
                Body:    fmt.Sprintf(confirmMessage,Url.String()),
        }
        if err := mail.Send(c, msg); err != nil {
                c.Errorf("Couldn't send email: %v", err)
                return err
        }
        return nil
}

/*
		Create a login link for apps
*/
func (controller *UsersController) SendAppLink(c appengine.Context,token string, email string) error{
		var Url *url.URL
    	path := "www.wymm-1242.appspot.com"
    	Url, err := url.Parse(path)
    	if err !=nil{
    		return err
    	}
    	Url.Path +="/users/loginverified/"
    	parameters := url.Values{}
    	parameters.Add("token",token)
    	Url.RawQuery = parameters.Encode()
    	c.Infof("Requested token: %s", Url.String())
    	confirmMessage := `Welcome to Lines. Please click on the link to login: %s`
    	addr := email
        msg := &mail.Message{
                Sender:  "nicholaspark09@gmail.com",
                To:      []string{addr},
                Subject: "Confirm your registration with Lines",
                Body:    fmt.Sprintf(confirmMessage,Url.String()),
        }
        if err := mail.Send(c, msg); err != nil {
                c.Errorf("Couldn't send email: %v", err)
                return err
        }
        return nil
}

/*
		Create a login link for apps
*/
func (controller *UsersController) SendIphoneLink(c appengine.Context,token string, email string) error{
		
		if Url,err := url.Parse("http://www.wymm-1242.appspot.com"); err == nil{

	    	Url.Host="wymm-1242.appspot.com"
	    	Url.Path +="/users/openlink"
	    	parameters := url.Values{}
	    	parameters.Add("token",token)
	    	Url.RawQuery = parameters.Encode()
	    	c.Infof("Requested token: %s", Url.String())
	    	confirmMessage := `Welcome <b>to</b> Lines. Please click the link below to login: %s`
	    	link := Url.String()
	    	addr := email
	        msg := &mail.Message{
	                Sender:  "nicholaspark09@gmail.com",
	                To:      []string{addr},
	                Subject: "Confirm your registration with Lines",
	                HTMLBody:    fmt.Sprintf(confirmMessage,link),
	        }
	        if err := mail.Send(c, msg); err != nil {
	                c.Errorf("Couldn't send email: %v", err)
	                return err
	        }
	    }
        return nil
}

func (controller *UsersController) OpenLink(){
	//controller.R.FormValue("token")
	controller.Data["Token"] = controller.R.FormValue("token")
	controller.Layout = append(controller.Layout,"users/openlink.tmpl")
	controller.Render(true)
}

func (controller *UsersController) DeleteAll(){
	/*
	c := appengine.NewContext(controller.R)
	q := datastore.NewQuery("User").KeysOnly()
	keys, err := q.GetAll(c, nil)
	if err !=nil{
		println(err.Error())
	}
	err = datastore.DeleteMulti(c, keys)
	if err != nil{
		println(err.Error())
	}
	*/
}

func (controller *UsersController) Languages(){

	if controller.R.Method=="GET"{
		controller.TemplateNames = "body"
		controller.Layout = append(controller.Layout,"users/languages.tmpl")
		controller.Render(true)
	}else if controller.R.Method=="POST"{
		m := make(map[string]interface{})
		m["Result"] = "Failure"
		c := appengine.NewContext(controller.R)
		locale := controller.R.FormValue("locale")
		cookie := http.Cookie{Name:"WymmLanguage",Value:locale,Expires:time.Now().Add(365 * 24 *time.Hour),HttpOnly:true, MaxAge:900000,Path:"/"}
		http.SetCookie(controller.W,&cookie)
		if controller.UserToken != ""{
			if user,userKey,_ := controller.GetUser(c); userKey != nil{
				user.Locale = locale
				user.Modified = time.Now()
				
				if _,err := datastore.Put(c,userKey,&user); err != nil{
					println(err.Error())
				}else{
					m["Result"] = "Success"
				}
			}
		}
		m["Result"] = "Success"

		data,_ := json.Marshal(m)
		fmt.Fprintf(controller.W,"%s",data)
	}
}

//To Register a user with no email
//You need a name, phone, and pin
func (controller *UsersController) AddLittle(){
	m := make(map[string]interface{})
	m["Result"] = "Failure"
	c := appengine.NewContext(controller.R)
	//c := appengine.NewContext(controller.R)
	if _,userKey,err := controller.GetUser(c); err != nil{
		m["Error"] = err.Error()
	}else{
		if campuskey,err := datastore.DecodeKey(controller.R.FormValue("campuskey")); err != nil{
			m["Error"] = err.Error()
		}else{
			println(userKey.Encode())
			println(campuskey.Encode())
		}
	}
}

//Login with just phone number and pin
func (controller *UsersController) PhoneLogin(){
	m := make(map[string]interface{})
	m["Result"] = "Failure"
		phone := controller.R.FormValue("phone")
		phone = strings.TrimSpace(strings.Replace(phone,"-","",-1))
		pin := strings.TrimSpace(controller.R.FormValue("pin"))
		c := appengine.NewContext(controller.R)
		q := datastore.NewQuery("User").Filter("Phone = ",phone).Filter("Active = ",true).Limit(1)
		t := q.Run(c)
		for{
			var user models.User
			key, err := t.Next(&user)
			if err == datastore.Done{
				break
			}
			if err != nil{
				m["Error"] = err.Error()
			}
			if user.Token == pin{
				token, err := controller.CreateSession(c,key)
				if err != nil{
					m["Result"] = "Failure"
					m["Error"] = err.Error()
					return 
				}
								m["Result"] = "Success"
				//User session has been created, take the token and set it into a cookie
				cookie := http.Cookie{Name:"WymmUser",Value:token,Expires:time.Now().Add(365 * 24 *time.Hour),HttpOnly:true, MaxAge:50000,Path:"/"}
				http.SetCookie(controller.W,&cookie)
				group_cookie := http.Cookie{Name:"WymmGroup",Value:key.Parent().Encode(),Expires:time.Now().Add(365 * 24 * time.Hour),HttpOnly:true,MaxAge:50000,Path:"/"}
				http.SetCookie(controller.W,&group_cookie)
				controller.Controller.UserToken = cookie.Value
			}else{
					m["Error"] = "The pin didn't match"
			}
		}
}

func (controller *UsersController) LoginVerified(){

	}

func (controller *UsersController) UpdateAll(){
	/*
	c := appengine.NewContext(controller.R)
	q := datastore.NewQuery("Session").Limit(100)
	assignments := make([]models.Session,0,100)
	if keys, err := q.GetAll(c,&assignments); err != nil{
		println(err.Error())
	}else{
		for i,_ := range keys{
			assignments[i].Token = ""
		}
		datastore.PutMulti(c,keys,assignments)
	}
	*/
}


func (controller *UsersController) Render(show bool) error{
	//controller.Controller.Render(show)
	controller.Controller.Render(show)
	return nil
}