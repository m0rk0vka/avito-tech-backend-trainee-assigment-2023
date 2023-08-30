package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/m0rk0vka/avito-tech-backend-trainee-assigment-2023/models"
)

type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

func createConnection() *sql.DB {
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successifully connected postgres")

	return db
}

func CreateSegment(w http.ResponseWriter, r *http.Request) {
	var segment models.Segments

	err := json.NewDecoder(r.Body).Decode(&segment)
	if err != nil {
		log.Fatalf("Unable to decode request body. %v", err)
	}

	insertID := insertSegment(segment)
	res := response{
		ID:      insertID,
		Message: "Segment created successifully",
	}

	json.NewEncoder(w).Encode(res)
}

func insertSegment(segment models.Segments) int64 {
	db := createConnection()
	defer db.Close()

	sqlStatement := `INSERT INTO segments(name) VALUES($1) RETURNING id`

	var id int64

	err := db.QueryRow(sqlStatement, segment.Name).Scan(&id)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	fmt.Println("Inserted single record %v.", id)

	return id
}

func DeleteSegment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("Unable to convert id. %v", err)
	}

	deletedRows := deleteSegment(int64(id))

	msg := fmt.Sprintf("Successifully deleted segment. Total rows/records %v.", deletedRows)

	res := response{
		ID:      int64(id),
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)
}

func deleteSegment(id int64) int64 {
	db := createConnection()
	defer db.Close()

	sqlStatement := `DELETE FROM segments WHERE id=$1`

	res, err := db.Exec(sqlStatement, id)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatalf("Error while checking affected rows. %v", err)
	}

	fmt.Printf("Total rows/records affected %v.", rowsAffected)

	return rowsAffected
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {

}

func GetUserByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("Unable to convert id. %v", err)
	}

	user, err := getUserByID(int64(id))
	if err != nil {
		log.Fatalf("Unable to get user. %v", err)
	}

	json.NewEncoder(w).Encode(user)
}

func getUserByID(id int64) (models.Users, error) {
	db := createConnection()
	defer db.Close()

	var user models.Users

	sqlStatement := `SELECT * FROM users WHERE id = $1`

	row := db.QueryRow(sqlStatement, id)

	err := row.Scan(&user.ID, &user.Username)
	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return user, nil
	case nil:
		return user, nil
	default:
		log.Fatalf("Unable to scan the row. %v", err)
	}

	return user, nil
}
