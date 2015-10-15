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
	"strconv"
)
type MyForm struct{
		UserName 	string `required:"true" field:"name" name:"Имя пользователя" type:"text" default:"true"`
		Child int `field:"child" name:"Дети" type:"text"`
		Age int64 `field:"age" name:"Возраст" type:"text" default:"true"`
		UserPassword 	string `required:"true" field:"password" name:"Пароль пользователя" type:"password"`
		Gender	 	string `required:"true" field:"gerder" name:"Пол" type:"select" select:"Не	известный=3,Мужской=1;selected,Женский=2"`
		Resident bool `field:"resident" type:"radio" radio:"1;checked" name:"Резидент РФ"`
		NoResident bool `field:"resident" type:"radio" radio:"2" name:"Не резидент РФ"`
		MbResident bool `field:"resident" type:"radio" radio:"3;checked" name:"mb Резидент РФ"`
		Comment 	string `required:"true" field:"tarea" name:"комментарий" type:"textarea"`

}
	


var (myform *MyForm
	err		error
	)
	

func main() {
	myform = &MyForm{UserName:"user",Age:21}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			
			if r.Method == "GET" {		
				frm_crt,err := FormCreate(myform)
				if err != nil {
					fmt.Println(err)
					}						
				fmt.Fprint(w,frm_crt)			
				}
				
			if r.Method == "POST" {
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
	fmt.Println("* GET request, creating form ",time.Now())
	err = nil
	str:="<h1>Hello!</h1> <form method='POST' action='http:\\localhost:1234'>"

//parse struct definition	
	val := reflect.ValueOf(myform).Elem()
	for i := 0; i < val.NumField(); i++ {			
			tag := val.Type().Field(i).Tag
			str=str+TagToHtml(tag)			

			field_type := tag.Get("type")
			field := tag.Get("field")
			name := tag.Get("name")
			def := tag.Get("default")
			def_value:=""
			
			if (field_type != "")&&(field != "") {				
					// label
					if name != "" {
						str = str+"<label>"+name+": </label><br>"}
					
					// type select					
					if field_type == "select" {
						str = str + SelectTagCreate(field,tag.Get("select"))
				
					// type radio
					}else if field_type == "radio"{
						str = str + RadioTagCreate(field,tag.Get("radio"))
						
					// textarea	
					}else if field_type == "textarea"{
						str = str+"<textarea name="+field+"></textarea><br>"
						
					}else{
						// default value from variable
						if def == "true"{
									tt := val.Type().Field(i).Type.String()
									fmt.Println("======================",tt)
									if tt == "int" || tt == "int64" {
										value_field := val.Field(i).Int()
										def_value="value="+strconv.FormatInt(value_field,10)
									}else if tt == "string" {
										value_field := val.Field(i).String()
										def_value="value='"+value_field+"'"
									}else{
										def_value=""
									}
						}
							//def_value="value='"+value_field+"'"}
						fmt.Println(def_value)
						str = str+"<input type='"+field_type+"' name='"+field+"' "+def_value+"><br>"						
					}
			}else{
				err = errors.New("*** error : no type or field in struct metadata")}			
	}
	str=str+"<button type='submit'>send</button> </form>"
	form=str
return
}


func FormRead(formData *MyForm, r *http.Request) (err error){
	fmt.Println("* POST request ",time.Now())
	err = nil	
	val := reflect.ValueOf(myform).Elem()
	r.ParseForm()
	for i := 0; i < val.NumField(); i++ {
		f := val.Field(i)
		//value_field := f.String()
		tag := val.Type().Field(i).Tag		
		field := tag.Get("field")
		form_v := r.PostFormValue(field)
		fmt.Println("type : ",reflect.TypeOf(form_v))
		v := reflect.ValueOf(form_v)
		f.Set(v)
		
		name := val.Type().Field(i).Name
		value_field := f.String()
		
		fmt.Println("=====")
		fmt.Println(name," : ",value_field)
		//fmt.Println("settability of f:", f.CanSet())
		fmt.Println("=====")
		
	}
return
}

func TagToHtml(tag reflect.StructTag)(htmlstr string){
return
}

func SelectTagCreate(f,tag string)(s string){
	s ="<select name="+f+">"
	options := strings.Split(tag,",")
	var temps1 string
	for _,v := range options{
		one_option := strings.Split(v,";")
		temps1 = ""
		if len(one_option)>1 {temps1="selected"}
		option_values := strings.Split(one_option[0],"=")
		s = s + "<option value='"+option_values[1]+"' "+temps1+">"+option_values[0]+"</option>"
		}
	s = s+"</select><br>"
return
}

func RadioTagCreate(f,tag string)(s string){
	s = "<input type='radio' name='"+f+"' "
	radio_value := strings.Split(tag,";")
	s = s + "value='"+radio_value[0]+"' "
	if len(radio_value)>1 {
		s = s + "checked"}
	s = s+"><br>"
return
}
