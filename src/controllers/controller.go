package controllers

import(
	"net/http"
	"text/template"
	"path"
	"appengine"
	"appengine/datastore"
	"github.com/parkn09/wymm/src/models"
	"time"
	"strings"
	"errors"
	"net"
)

type Controller struct{
	Name string
	R *http.Request
	W http.ResponseWriter
	Data map[interface{}]interface{}
	Tpl *template.Template
	C appengine.Context
	Layout []string
	TemplateNames string
	UserToken string
	User models.User
	Locale string
	Target string
	Action string
	Params []string
}

type ControllerInterface interface{
	Init(r *http.Request, w http.ResponseWriter)
	Serve(action string)
	MobileServe(action string)
	Render(show bool) error
}

func NewController() *Controller{
	return &Controller{}
}

func (controller *Controller) Init(r *http.Request, w http.ResponseWriter){
	controller.R = r
	controller.W = w
	controller.Data = make(map[interface{}]interface{})
	controller.C = appengine.NewContext(r)
	controller.Layout = make([]string,0)
	controller.TemplateNames = "body"
	
	
	cookie, err := controller.R.Cookie("WymmLanguage")
	if err!= nil{
		controller.Locale = "en_US"
	}else if cookie.Value != ""{
		controller.Locale = cookie.Value
	}
	targetCookie,err := controller.R.Cookie("WymmTargetLanguage")
	if err != nil{
		controller.Locale = "en_US"
	}else if cookie.Value != ""{
		controller.Target = targetCookie.Value
	}
	controller.Data["Language"] = controller.Locale
	controller.Data["Target"] = controller.Target
	controller.Layout = append(controller.Layout,"layouts/header.tmpl")
	controller.Layout = append(controller.Layout,"layouts/footer.tmpl")
	if controller.CheckSession(){
		controller.Layout = append(controller.Layout,"layouts/logged.tmpl")
	}else{
				controller.Layout = append(controller.Layout,"layouts/nav.tmpl")
	}
}

func (controller *Controller) DetermineRoutes(){
	controller.Params = strings.Split(string(controller.R.URL.Path[1:]),"/")
}

func (controller *Controller) Render(show bool) error{
	
		var filenames []string
        for _, file := range controller.Layout {
            filenames = append(filenames, path.Join("views/", file))
        }
		t,err := template.ParseFiles(filenames...)
		if err != nil{
			return err
		}
		err = t.ExecuteTemplate(controller.W, controller.TemplateNames,controller.Data)
		if err != nil{
			println("There was an error")
			println(err.Error())
			return err
		}
	return nil
}

//Checks to see if a Session cookie has been set
func (controller *Controller) CheckSession() bool {
	cookie, err := controller.R.Cookie("WymmUser")
	if err!= nil{
		controller.UserToken = ""
		return false
	}else if cookie.Value != ""{
		controller.UserToken = cookie.Value
		return true
	}
	return false
}

func (controller *Controller) View(){
	
}


//Get the User
func (controller *Controller) GetUser(c appengine.Context) (models.User,*datastore.Key,error){
	user := models.User{
		Active: true,
	}
	//c := appengine.NewContext(controller.R)
	q := datastore.NewQuery("Session").Filter("Token = ",controller.UserToken).Filter("Active =",true).Order("-Created").Limit(1)
	err := errors.New("No User")
	sessions := make([]models.Session,0,1)
	if keys, err := q.GetAll(c, &sessions); err != nil{
		return user,nil,err
	}else{
		//Now check for the user with a session
		count := len(keys)
		if count == 0{
			return user,nil,err
		}else{
			userKey := keys[0].Parent()
			if err := datastore.Get(c,userKey,&user); err != nil{
				return user,nil,err
			}else{
				controller.User = user
				return user,userKey,nil
			}
		}
	}
	return user,nil,err
}

func (controller *Controller) GetUserKey(c appengine.Context) (*datastore.Key, error){
	err := errors.New("No Key")
	sessions := make([]models.Session,0,1)
	q := datastore.NewQuery("Session").Filter("Token = ",controller.UserToken).Filter("Active =",true).Order("-Created").Limit(1)
	if keys, err := q.GetAll(c,&sessions); err != nil{
		return nil,err
	}else if len(keys)>0{
		return keys[0].Parent(),nil
	}
	return nil,err
}


func (controller *Controller) RedirectToLogin(){
	controller.Layout = append(controller.Layout,"users/login.tmpl")
	controller.Render(true)
}

func (controller *Controller) Logout(){
	cookie := http.Cookie{Name:"WymmUser", Value:"",Expires:time.Now().Add(-time.Hour),HttpOnly:true,MaxAge:1,Path:"/"}
	http.SetCookie(controller.W, &cookie)
	othercookie := http.Cookie{Name:"WymmTeacher", Value:"",Expires:time.Now().Add(-time.Hour),HttpOnly:true,MaxAge:1,Path:"/"}
	http.SetCookie(controller.W, &othercookie)
	blahcookie := http.Cookie{Name:"WymmStudent", Value:"",Expires:time.Now().Add(-time.Hour),HttpOnly:true,MaxAge:10,Path:"/"}
	http.SetCookie(controller.W, &blahcookie)
	controller.UserToken = ""
	http.Redirect(controller.W,controller.R,"/",301)

}





func (controller *Controller) RenderError(error string){
	controller.Data["Error"] = error
	controller.Layout = append(controller.Layout,"layouts/error.tmpl")
	controller.Render(true)
}

func (controller *Controller) SaveInfo(c appengine.Context, device string, safekey string, user string) (string,error){
	
	//Get Host
	ip,_,_ := net.SplitHostPort(controller.R.RemoteAddr)

	view := models.View{
		Controller: controller.Name,
		Action:controller.Action,
		Device:device,
		IP:ip,
		Created:time.Now(),
		Modified:time.Now(),
		Safekey:safekey,
		User:user,
		Active:true,
	}
	viewKey := datastore.NewIncompleteKey(c,"View",nil)
	if realKey, err := datastore.Put(c,viewKey,&view); err != nil{
		return "",err
	}else{
		return realKey.Encode(),nil
	}
	return "",nil
}
/*
func (controller *Controller) Render(name string) error{

	//var templates = template.Must(template.New("Base").ParseFiles("views/layout.html"))
	
	//templates := template.Must(template.ParseFiles("views/templates/header.tmpl","views/templates/footer.tmpl","views/templates/main.tmpl"))
	t := template.Must(template.ParseFiles(name,"views/templates/header.tmpl","views/templates/footer.tmpl"))
	err := t.ExecuteTemplate(controller.Response,"body",controller.Data)
	
/*
	var templates = template.Must(template.New("body").ParseFiles("views/templates/header.tmpl",name))
	err := templates.ExecuteTemplate(controller.Response,"body",page)
	*/
//	if err != nil{
//		return err
//	}
//	return nil
//}