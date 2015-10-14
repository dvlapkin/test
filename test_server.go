package main
import (
	"errors"
	"fmt"
	//"html"
	//"html/template"
	"log"
	//"io/ioutil"
	"net/http"
	//"net"
	"time"
	"reflect"
)
type MyForm struct{
		UserName 		string `required:"true" field:"name" name:"Имя пользователя" type:"text"`
		UserPassword 	string `required:"true" field:"password" name:"Пароль пользователя" type:"password"`
	}
	
var (myform *MyForm
	err		error
	)
	
//form = &MyForm{}

func main() {
	myform = &MyForm{}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			
			if r.Method == "GET" {	
			//fmt.Fprintf(w, "<b>World</b><button>ok</button>")
				//var t = template.Must(template.New("name").Parse("<b>World</b><button>ok</button>"))
				//err = t.Execute(w,t)	
				frm_crt,err := FormCreate(myform)
				if err != nil {
					log.Fatal(err)
					}						
				fmt.Fprint(w,frm_crt)			
				}
				
			if r.Method == "POST" {
			
			//frm_crt,err := FormCreate(myform)
			//	if err != nil {
			//		log.Fatal(err)
			//		}	
			//	fmt.Fprint(w,frm_crt)
				
				err = FormRead(myform,r)
					if err != nil {
						log.Fatal(err)
						}
						

				fmt.Fprint(w,"post")
				}
			})
				
	err := http.ListenAndServe("localhost:1234", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func FormCreate(formData *MyForm) (form string, err error){

	fmt.Println("*GET request, creating form ",time.Now())
	form = "hello!!!"
	err = nil

	str:="<h1>Hello!</h1> <form method='POST' action='http:\\localhost:1234'>"
//parse struct definition	
	val := reflect.ValueOf(myform).Elem()
	for i := 0; i < val.NumField(); i++ {
			//valueField := val.Field(i)
			//typeField := val.Type().Field(i)
			//tag := typeField.Tag
			tag := val.Type().Field(i).Tag
			str=str+TagToHtml(tag)
		//	fmt.Println("type: ",tag.Get("type"),", name: ",tag.Get("name"))			
		//var t = template.Must(template.New("name").Parse("html"))
			field_type := tag.Get("type")
			field := tag.Get("field")
			name := tag.Get("name")
			
			if (field_type != "")&&(field != "") {				
				if field_type == "select" {
					str = str+"<select>"
				}else{
					if name != "" {
						str = str+"<label>"+name+": </label>"
						}
					str = str+"<input type="+field_type+" name="+field							
					str = str+"><br>"
					}
			}else{
				err = errors.New("no type or field in struct metadata")
				}
			}	
	str=str+"<button type='submit'>send</button> </form>"
	//form=html.EscapeString(str)
	form=str
return
}

func FormRead(formData *MyForm, r *http.Request) (err error){
	fmt.Println("*POST request ",time.Now())
	err = nil
return
}

func TagTOHtml(tag reflect.StructTag)(htmlstr string){
return
}
