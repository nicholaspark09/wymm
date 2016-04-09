package index

import (
	"bytes"
	"net/http"
	"strings"
	"github.com/parkn09/wymm/src/controllers"
)


func handler(w http.ResponseWriter, r *http.Request){
	m := make(map[string]controllers.ControllerInterface)
	m["users"] = controllers.NewUsersController()
	m["general"] = controllers.NewGeneralController()
	m["languages"] = controllers.NewLanguagesController()
	m["lines"] = controllers.NewLinesController()
	m["likes"] = controllers.NewLikesController()
	m["dislikes"] = controllers.NewDislikesController()
	m["translateds"] = controllers.NewTranslatedsController()
	m["favorites"] = controllers.NewFavoritesController()
	m["flags"] = controllers.NewFlagsController()
	params := string(r.URL.Path[1:])
	arr := strings.Split(params,"/")
	if(params==""){
		m["general"].Init(r,w)
		m["general"].Serve("home")
	}else{
		count := len(arr)
		if count > 0{
			if arr[0]=="images"{
				buffer := bytes.NewBufferString("views/")
				buffer.WriteString(string(params))
				http.ServeFile(w, r, buffer.String())
			}else if arr[0] == "js" || arr[0] == "css" || arr[0]=="fonts"{
				buffer := bytes.NewBufferString("views/")
				buffer.WriteString(string(params))
				http.ServeFile(w, r, buffer.String())
			}else{
				if _, ok := m[arr[0]]; ok{
					m[arr[0]].Init(r,w)
					action := "index"
					if count >1{
						if arr[1] == "mobile"{
							m[arr[0]].MobileServe(action)
						}else{
							action = arr[1]
							m[arr[0]].Serve(action)
						}
					}else{
						m[arr[0]].Serve(action)
					}
				}
				
			}
		}	
	}
	
}


func init(){
	
	//http.HandleFunc("/images*",imageHandler)
	http.HandleFunc("/",handler)

}