package main
import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"io"
	"time"
	"reflect"
	"strings"
	"strconv"
	"os"
)
type MyForm struct{
	UserName 	string `required:"true" field:"name" name:"Имя пользователя" type:"text" default:"true"`
	Login 		int `required:"true" field:"login" name:"Логин"`
	Age 		int64 `field:"age" name:"Возраст" type:"text" default:"true"`
	Koe 		float64 `field:"koef" name:"Коэфицент" type:"text" default:"true"`
	UserPassword 	string `required:"true" field:"password" name:"Пароль пользователя" type:"password"`
	Gender	 	string `required:"true" field:"gerder" name:"Пол" type:"select" select:"Не	известный=3,Мужской=1;selected,Женский=2"`
	Resident 	bool `field:"resident" type:"radio" radio:"1;checked" name:"Резидент РФ"`
	NoResident 	bool `field:"resident" type:"radio" radio:"2" name:"Не резидент РФ"`		
	Foot 		bool `field:"body" type:"radio" radio:"1;checked" name:"Нога"`
	Hand 		bool `field:"body" type:"radio" radio:"2" name:"Рука"`
	Duration 	time.Duration `required:"true" field:"duration" name:"Длительность" type:"text" default:"true"`
	Comment 	string `field:"tarea" name:"комментарий" type:"textarea"`
}
	

var (
	myform *MyForm
	err		error
	)
	

func main() {
	myform = &MyForm{UserName:"user",Age:21,Koe:2.34,Duration:12345678}
	fmt.Println(myform)
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
					fmt.Println("*** error : ",err)
				}else{						
					StructPrt(myform,os.Stdout)
					StructPrt(myform,w)}
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
	val := reflect.ValueOf(formData).Elem()
	for i := 0; i < val.NumField(); i++ {			
			tag := val.Type().Field(i).Tag		
			field_type := tag.Get("type")
			field := tag.Get("field")
			name := tag.Get("name")
			def := tag.Get("default")
			def_value:=""
			if field_type == "" {field_type="text"}
			
			if (field != "")||(field != "-") {				
					// label
					if name != "" {
						str = str+"<label>"+name+": </label><br>"}
					
					// type select					
					if field_type == "select" {
						SelectTagCreate(&str,tag)
				
					// type radio
					}else if field_type == "radio"{
						RadioTagCreate(&str,tag)
						
					// textarea	
					}else if field_type == "textarea"{
						str = str+"<textarea name="+field+"></textarea><br>"
						
					}else{
						// inputs with default value from variable
						if def == "true"{
							tt := val.Type().Field(i).Type.String()
							if tt == "int" || tt == "int64" {
								value_field := val.Field(i).Int()
								def_value="value="+strconv.FormatInt(value_field,10)
							}else if tt == "string" {
								value_field := val.Field(i).String()
								def_value="value='"+value_field+"'"
							}else if tt == "float64" {
								value_field := val.Field(i).Float()
								def_value="value="+strconv.FormatFloat(value_field,'f',5,64)
							}else if tt == "uint" || tt == "uint64"{
								value_field := val.Field(i).Uint()
								def_value="value="+strconv.FormatUint(value_field,10)
							}else if tt == "bool"{
								value_field := val.Field(i).Bool()
								def_value="value="+strconv.FormatBool(value_field)
							}else if tt == "time.Duration" {
								value_field := val.Field(i).Int()
								def_value="value="+strconv.FormatInt(value_field,10)+"ns"
							}else{
								def_value=""
							}
							fmt.Println("default for ",name," : ",def_value)
						}
													
						str = str+"<input type='"+field_type+"' name='"+field+"' "+def_value+"><br>"						
					}
			}else{
				err = errors.New("*** error : no field in struct metadata")}			
	}
	str=str+"<button type='submit'>send</button> </form>"
	form=str

return
}


func FormRead(formData *MyForm, r *http.Request) (err error){
	fmt.Println("* POST request ",time.Now())
	err = nil	
	val := reflect.ValueOf(formData).Elem()
	r.ParseForm()
	for i := 0; i < val.NumField(); i++ {
		f := val.Field(i)
		tag := val.Type().Field(i).Tag	
	
		field := tag.Get("field")
		form_v := r.PostFormValue(field)		
		struct_field_type := val.Type().Field(i).Type.Name()
		var v reflect.Value
		if tag.Get("required")== "true"&& form_v == "" {
			err=errors.New("*** error :  form field "+field+" requared")
			return
		}
		
		if (tag.Get("type") == "radio"){
			if (form_v == strings.Split(tag.Get("radio"),";")[0]){	
				v = reflect.ValueOf(true)
			}else{
				v = reflect.ValueOf(false)
			}
		}else{
		switch struct_field_type {
			case "string":{
				v = reflect.ValueOf(form_v)}
			case "int":{
				form_vi,e:=strconv.Atoi(form_v)
				v = reflect.ValueOf(form_vi)
				if e != nil {err = e;return}}
			case "int64":{
				form_vi,e:=strconv.ParseInt(form_v,10,64)
				v = reflect.ValueOf(form_vi)
				if e != nil {err = e;return}}
			case "float64":{
				form_vf,e:=strconv.ParseFloat(form_v,64)
				v = reflect.ValueOf(form_vf)
				if e != nil {err = e;return}}
			case "uint":{
				form_vui,e:=strconv.ParseUint(form_v,10,0)
				v = reflect.ValueOf(form_vui)
				if e != nil {err = e;return}}
			case "uint64":{
				form_vui,e:=strconv.ParseUint(form_v,10,64)
				v = reflect.ValueOf(form_vui)
				if e != nil {err = e;return}}
			case "Duration":{
				form_vt,e:=time.ParseDuration(form_v)		
				v = reflect.ValueOf(form_vt)
				if e != nil {err = e;return}}			
			case "bool":{				 
				form_vb,e:=strconv.ParseBool(form_v)
				v = reflect.ValueOf(form_vb)
				if e != nil {err = e;return}}
			default:
					err=errors.New("*** error : type of"+field+"("+struct_field_type+") not supported!")
			
			}
		}
		f.Set(v)	
	}
return
}


func SelectTagCreate(s *string,t reflect.StructTag){
	f:= t.Get("field")
	tag:=t.Get("select")
	*s +="<select name="+f+">"
	options := strings.Split(tag,",")
	var temps1 string
	for _,v := range options{
		one_option := strings.Split(v,";")
		temps1 = ""
		if len(one_option)>1 {temps1="selected"}
		option_values := strings.Split(one_option[0],"=")
		*s += "<option value='"+option_values[1]+"' "+temps1+">"+option_values[0]+"</option>"
		}
	*s += "</select><br>"
return
}

func RadioTagCreate(s *string,t reflect.StructTag){
	f:= t.Get("field")
	tag:=t.Get("radio")
	*s += "<input type='radio' name='"+f+"' "
	radio_value := strings.Split(tag,";")
	*s += "value='"+radio_value[0]+"' "
	if len(radio_value)>1 {
		*s += "checked"}
	*s += "><br>"
return
}

func StructPrt(st *MyForm,w io.Writer){
	var fv string
	
	val := reflect.ValueOf(st).Elem()
	for i := 0; i < val.NumField(); i++ {
		tt := val.Type().Field(i).Type.String()
		fname := val.Type().Field(i).Name
		if tt == "int" || tt == "int64" {
			value_field := val.Field(i).Int()
			fv= fname+" : "+strconv.FormatInt(value_field,10)
		}else if tt == "string" {
			value_field := val.Field(i).String()
			fv=fname+" : "+value_field
		}else if tt == "float64" {
			value_field := val.Field(i).Float()
			fv=fname+" : "+strconv.FormatFloat(value_field,'f',5,64)
		}else if tt == "uint" || tt == "uint64"{
			value_field := val.Field(i).Uint()
			fv=fname+" : "+strconv.FormatUint(value_field,10)
		}else if tt == "bool"{
			value_field := val.Field(i).Bool()
			fv=fname+" : "+strconv.FormatBool(value_field)
		}else if tt == "time.Duration" {
			value_field := val.Field(i).Int()
			fv=fname+" : "+strconv.FormatInt(value_field,10)
		}else{
			fv=fname+" : "+" unknown type"
		}
	fmt.Fprintln(w,fv)
	}
}
