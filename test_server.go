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
	"strings"
)
type MyForm struct{
		UserName 	string `required:"true" field:"name" name:"Имя пользователя" type:"text"`
		UserPassword 	string `required:"true" field:"password" name:"Пароль пользователя" type:"password"`
		Gender	 	string `required:"true" field:"gerder" name:"Пол" type:"select" select:"Не	известный=3,Мужской=1;selected,Женский=2"`
		Resident bool `field:"resident" type:"radio" radio:"1;checked" name:"Резидент РФ"`
		NoResident bool `field:"resident" type:"radio" radio:"2" name:"Не резидент РФ"`
		MbResident bool `type:"radio" radio:"3;checked" name:"mb Резидент РФ"`
		Comment 	string `required:"true" field:"tarea" name:"комментарий" type:"textarea"`

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
					fmt.Println(err)
					//log.Fatal(err)
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
	err = nil
	str:="<h1>Hello!</h1> <form method='POST' action='http:\\localhost:1234'>"

//parse struct definition	
	val := reflect.ValueOf(myform).Elem()
	for i := 0; i < val.NumField(); i++ {
			//valueField := val.Field(i)
			//typeField := val.Type().Field(i)
			//tag := typeField.Tag
			struct_elem := val.Type().Field(i).Name
			tag := val.Type().Field(i).Tag
			str=str+TagToHtml(tag)
		//	fmt.Println("type: ",tag.Get("type"),", name: ",tag.Get("name"))			

			field_type := tag.Get("type")
			field := tag.Get("field")
			name := tag.Get("name")
			def := tag.Get("default")
			//requared := tag.Get("requared")
			def_value:=""
			
			if (field_type != "")&&(field != "") {				
					// label
					if name != "" {
						str = str+"<label>"+name+": </label><br>"}
					// default value from variable
					if def != "true" {
						def_value="value='"+myform(struct_elem)+"'"}
					
					// type select
					
					if field_type == "select" {
						str = str+"<select name="+field+">"
						options := strings.Split(tag.Get("select"),",")
						var temps1 string
						for _,v := range options{
							one_option := strings.Split(v,";")
							temps1 = ""
							if len(one_option)>1 {temps1="selected"}
							option_values := strings.Split(one_option[0],"=")
							str = str + "<option value='"+option_values[1]+"' "+temps1+">"+option_values[0]+"</option>"
							}
						str = str+"</select><br>"
					// type radio
					}else if field_type == "radio"{
						str = str+"<input type='radio' name='"+field+"' "
						radio_value := strings.Split(tag.Get("radio"),";")
						str = str + "value='"+radio_value[0]+"' "
						if len(radio_value)>1 {
							str = str + "checked"
							}
						str = str+"><br>"
					// other input
					}else if field_type == "textarea"{
						str = str+"<textarea name="+field+"></textarea><br>"
					}else{
						str = str+"<input type="+field_type+" name="+field+"><br>"						
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

func TagToHtml(tag reflect.StructTag)(htmlstr string){
return
}
