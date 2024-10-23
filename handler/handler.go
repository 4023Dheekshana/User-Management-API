package handler

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"userapi/database"
	"userapi/utils"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type User struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
	Age      int64  `json:"age"`
}

type Usercrd struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func UserSignUp(context *gin.Context) {

	var newUser Usercrd

	if context.Request.Header.Get("Content-Type") != "application/json" {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid Content-Type"})
		fmt.Println("Invalid Content-Tye")
		return
	}

	if err := context.ShouldBindJSON(&newUser); err != nil {
		log.Println("Error binding json", err)
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Error binding json"})
	}
	_, err := database.Db.Exec("INSERT INTO userlp (username, password)VALUES($1, $2) ", newUser.Username, newUser.Password)
	if err != nil {
		context.AbortWithStatusJSON(400, "Siging up failed")
		fmt.Println("Error inserting user into database:", err)
	} else {
		context.JSON(http.StatusOK, "User signed up successfully")
	}

}

func Userlogin(context *gin.Context) {
	var newuser Usercrd
	err := context.ShouldBindJSON(&newuser)
	if err != nil {
		log.Println("Error binding json", err)
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Error binding json"})
		return
	}
	query := "SELECT password FROM userlp WHERE username = $1"
	row := database.Db.QueryRow(query, newuser.Username)
	var retrievedpassword string
	err = row.Scan(&retrievedpassword)
	if err != nil {
		log.Fatalf("Error getting retrieved password %v", err)
		return
	}
	if newuser.Password != retrievedpassword {
		log.Fatal("Login password is wrong")
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Login password is wrong"})
		return
	}

	token, err := utils.GenerateToken(newuser.Username, newuser.Password)
	if err != nil {
		log.Fatalf("Error generating a token %v", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not generate a token."})
	}
	context.JSON(http.StatusOK, gin.H{"message": "Login Successful.", "token": token})
}

func AddUser(context *gin.Context) {
	var body User

	if context.Request.Header.Get("Content-Type") != "application/json" {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid Content-Type"})
		fmt.Println("Invalid Content-Type")
		return
	}

	if err := context.BindJSON(&body); err != nil {
		log.Println("Error binding json", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Error binding json"})
		return
	}

	_, err := database.Db.Exec("INSERT INTO userinfo (name, location, age) VALUES($1, $2, $3)", body.Name, body.Location, body.Age)
	if err != nil {
		context.AbortWithStatusJSON(400, "Couldnt create the user.")
		fmt.Println("Error inserting user into database:", err)
	} else {
		context.JSON(http.StatusOK, "User is successfully created.")
	}
}

func GetAllUser(context *gin.Context) {
	context.Header("Content-Type", "application/json")

	rows, err := database.Db.Query("SELECT id, name, location, age FROM userinfo")
	if err != nil {
		log.Println("Error querying the data rows", err)
		return
	}
	defer rows.Close()
	var users []User
	for rows.Next() {
		var u User
		err := rows.Scan(&u.Id, &u.Name, &u.Location, &u.Age)
		if err != nil {
			log.Println("Error scanning the rows", err)
			return
		}
		users = append(users, u)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	context.JSON(http.StatusOK, users)
}

func GetUserById(context *gin.Context) {
	context.Header("Content_Type", "application/json")

	id := context.Param("id")

	var user User
	err := database.Db.QueryRow(
		`SELECT id, name, location, age FROM userinfo
		WHERE id = $1`, id).Scan(&user.Id, &user.Name, &user.Location, &user.Age)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("User does not found in the rows(Empty row)", err)
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting empty row"})
		} else {
			log.Println("Error querying the user by id", err)
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying the user"})
		}
		return
	}
	context.JSON(http.StatusOK, user)
}

func UpdateUser(context *gin.Context) {
	context.Header("Content_Type", "application/json")

	id, _ := strconv.Atoi(context.Param("id"))
	var user User
	if err := context.BindJSON(&user); err != nil {
		log.Println("Error binding json", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Error binding json"})
		return
	}

	log.Printf("Received data: %+v\n", user)

	result, err := database.Db.Exec("UPDATE userinfo SET name = $1, location = $2, age = $3 where id = $4",
		user.Name, user.Location, user.Age, id)
	if err != nil {
		log.Println("Error updating the user", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating the user."})
		return
	}

	resultrows, err := result.RowsAffected()
	if err != nil {
		log.Println("Error getting number of rows affected", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting number of rows"})
		return
	}

	if resultrows == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	context.JSON(http.StatusOK, "User updated successfully")
}

func DeleteUser(context *gin.Context) {
	if context.Request.Header.Get("Content-Type") != "application/json" {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid Content-Type"})
		fmt.Println("Invalid Content-Type")
		return
	}

	id := context.Param("id")

	result, err := database.Db.Exec("DELETE FROM userinfo WHERE id = $1", id)
	if err != nil {
		log.Print("Error deleting the user", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting the row"})
		return
	}

	resultedRows, err := result.RowsAffected()
	if err != nil {
		log.Println("Error getting number of rows affected", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": "Error getting number of rows"})
		return
	}

	if resultedRows == 0 {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
		return
	}

	context.JSON(http.StatusOK, "User deleted successfully.")
}
