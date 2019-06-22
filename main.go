//API
//https://github.com/gin-gonic/gin/blob/master/README.md
package main

import (
	"fmt"
	"database/sql"
	"net/http"
	"os"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Todo struct {
	ID     int
	Title  string
	Status string
}

func pingHandler(c *gin.Context) {
	response := gin.H{
		"message": "pong2",
	}

	c.JSON(http.StatusOK, response)
}

func pingPostHandler(c *gin.Context) {
	response := gin.H{
		"message": "this is ping POST",
	}

	c.JSON(http.StatusOK, response)
}

type Student struct {
	Name string `json:"name"`       //field name สำหรับ json ฝั่ง client   //back quote : Alt+0096
	ID   int    `json:"student_id"` //ส่วนใหญ่ใช้ตัวเล็กมี "_"
}

var students = map[int]Student{
	1: Student{Name: "AnuchitO", ID: 1},
}

func getStudentHandler(c *gin.Context) {
	//c.JSON(http.StatusOK, students)

	ss := []Student{}
	for _, s := range students { //วนลูปจาก map เป็น slice
		ss = append(ss, s)
	}

	c.JSON(http.StatusOK, ss) //return เป็น array
}

func postStudentHandler(c *gin.Context) {
	s := Student{}
	fmt.Printf("before bind % #v\n", s)
	if err := c.ShouldBindJSON(&s); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Printf("after bind % #v\n", s)

	//Gen id
	id := len(students)
	id++
	s.ID = id
	students[id] = s

	c.JSON(http.StatusOK, s)
}

func getTodosHandler(c *gin.Context) {
	db, _ := sql.Open("postgres", os.Getenv("DATABASE_URL")) //ควรดัก err ด้วย
	stmt, _ := db.Prepare("SELECT id, title, status FROM todos") //ป้องกัน SQL injection
	rows, _ := stmt.Query()

	todos := []Todo{}
	for rows.Next() {
		t := Todo{}
		err := rows.Scan(&t.ID, &t.Title, &t.Status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		todos = append(todos, t)
	}

	c.JSON(http.StatusOK, todos)
}

func main() {
	r := gin.Default()
	r.GET("/ping", pingHandler)
	r.POST("/ping", pingPostHandler)
	r.GET("/students", getStudentHandler) //ตั้งชื่อเป็นพหูพจน์
	r.POST("/students", postStudentHandler)

	r.GET("/api/todos", getTodosHandler)
	//r.GET("/api/todos/:id", getTodosByIdHandler)
	//r.POST("/api/todos", postTodosHandler)
	//r.DELETE("/api/todos/:id", deleteTodosByIdHandler)

	r.Run(":1234") //listen and serve on 0.0.0.0.8080
}
