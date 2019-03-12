package service

import (
	"fmt"
	"os"
	"time"
)

const (
	log_path = "logs/"
)

func Log(tag string, content interface{}) {
	this_week := get_week_begin_date().Format("2006-01-02")[:10]
	file := log_path + this_week + ".out"
	fileObj, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND,0644)
	if err != nil {
		fmt.Println("open file " + file + " fail!")
		return
	}
	defer fileObj.Close()
	_, err = fileObj.Write(gen_content(tag, content))
	if err != err {
		fmt.Println("write content fail!")
		return
	}
}


func get_week_begin_date() time.Time {
	today := time.Now()
	today_week := today.Weekday().String()
	var plus int
	switch today_week {
	case "Sunday":
		plus = 0
	case "Monday":
		plus = 1
	case "Tuesday":
		plus = 2
	case "Wednesday":
		plus = 3
	case "Thursday":
		plus = 4
	case "Friday":
		plus = 5
	case "Saturday":
		plus = 6
	}
	week_begin := today.AddDate(0, 0, -plus)
	return week_begin
}

func gen_content(tag string, con interface{}) []byte {
	log := "====================="+ fmt.Sprint(time.Now()) +"=================\rtag:" + tag + " content:"+fmt.Sprint(con)+"\r\r"
	return []byte(log)
}