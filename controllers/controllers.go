package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"github.com/m0rk0vka/avito-tech-backend-trainee-assigment-2023/models"
)

type response struct {
	Time    time.Time `json:"time,omitempty"`
	Code    int       `json:"code,omitempty"`
	Message string    `json:"message,omitempty"`
}

func createConnection() *sql.DB {
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))
	if err != nil {
		log.Fatalf("Unable to open connection to db.\n %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Unable to ping db.\n %v", err)
	}

	res := response{
		Time:    time.Now(),
		Code:    http.StatusOK,
		Message: "Successifully connected to postgres",
	}

	log.Println(res)

	return db
}

func CreateSegment(w http.ResponseWriter, r *http.Request) {
	var segment models.Segments

	err := json.NewDecoder(r.Body).Decode(&segment)
	if err != nil {
		log.Fatalf("Unable to decode request body. %v", err)
	}

	var msg string

	insertID, statusCode := insertSegment(segment)

	switch statusCode {
	case http.StatusCreated:
		msg = fmt.Sprintf("Segment with id %d created successifully", insertID)
	case http.StatusBadRequest:
		msg = fmt.Sprintf("Segment with name %s already exists.", segment.Name)
	case http.StatusInternalServerError:
		msg = fmt.Sprintf("Unable to scan the row")
	}

	res := response{
		Time:    time.Now(),
		Code:    statusCode,
		Message: msg,
	}

	log.Println(res)

	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(res)
}

func insertSegment(segment models.Segments) (int64, int) {
	db := createConnection()
	defer db.Close()

	sqlStatement := `SELECT EXISTS (SELECT * FROM segments WHERE name=$1)`

	var exists bool

	if err := db.QueryRow(sqlStatement, segment.Name).Scan(&exists); err != nil {
		log.Fatalf("Unable to scan the row. %v", err)
		return 0, http.StatusInternalServerError
	}

	if exists {
		return 0, http.StatusBadRequest
	}

	sqlStatement = `INSERT INTO segments(name) VALUES($1) RETURNING id`

	var id int64

	err := db.QueryRow(sqlStatement, segment.Name).Scan(&id)
	if err != nil {
		log.Fatalf("Unable to scan the row. %v", err)
		return 0, http.StatusInternalServerError
	}

	return id, http.StatusCreated
}

func DeleteSegment(w http.ResponseWriter, r *http.Request) {
	var segment models.Segments

	err := json.NewDecoder(r.Body).Decode(&segment)
	if err != nil {
		log.Fatalf("Unable to decode request body")
	}

	var (
		statusCode int
		msg        string
	)

	segment.ID, statusCode = getSegmentIDByName(segment.Name)
	switch statusCode {
	case http.StatusOK:

		sqlStatement := `DELETE FROM relations WHERE segment_id=$1`
		deletedRows := deleteSegmentFromTableByID(int64(segment.ID), sqlStatement)

		sqlStatement = `DELETE FROM segments WHERE id=$1`
		deletedRows += deleteSegmentFromTableByID(int64(segment.ID), sqlStatement)

		msg = fmt.Sprintf("Successifully deleted segment. Total rows/records affected %d.", deletedRows)
	case http.StatusInternalServerError:
		msg = fmt.Sprintf("Unable to scan the row")
	case http.StatusNotFound:
		msg = fmt.Sprintf("No sigment with name %s.", segment.Name)
	}

	res := response{
		Time:    time.Now(),
		Code:    statusCode,
		Message: msg,
	}

	log.Println(res)

	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(res)
}

func getSegmentIDByName(name string) (int64, int) {
	db := createConnection()
	defer db.Close()

	sqlStatement := `SELECT id FROM segments WHERE name=$1`

	var id int64

	if err := db.QueryRow(sqlStatement, name).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return 0, http.StatusNotFound
		}
		return 0, http.StatusInternalServerError
	}

	return id, http.StatusOK
}

func deleteSegmentFromTableByID(id int64, sqlStatement string) int64 {
	db := createConnection()
	defer db.Close()

	res, err := db.Exec(sqlStatement, id)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatalf("Error while checking affected rows. %v", err)
	}

	return rowsAffected
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	var update models.Update

	err := json.NewDecoder(r.Body).Decode(&update)
	if err != nil {
		log.Fatalf("Unable to decode request body. %v", err)
	}

	isUser, err := isUser(update.UserID)
	if err != nil {
		log.Fatalf("Unable to check existing user. %v", err)
	}

	statusCode := http.StatusBadRequest

	msg := fmt.Sprintf("Bad request. Check all segments are exists or all segments you want to remove are related to user")

	if checkInputSegments(update) {
		// if input is valide then create user if it is not exists
		if !isUser {
			insertUser(update.UserID)
		}

		updatedRows := updateRelations(update)

		statusCode = http.StatusOK

		msg = fmt.Sprintf("Successifully updated user segments. Total rows/records affected %v", updatedRows)
	}

	res := response{
		Time:    time.Now(),
		Code:    statusCode,
		Message: msg,
	}

	log.Println(res)

	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(res)
}

func checkInputSegments(update models.Update) bool {
	db := createConnection()
	defer db.Close()

	isSegment := `SELECT EXISTS (SELECT * FROM segments WHERE name=$1)`

	isUser := `SELECT EXISTS (SELECT * FROM relations WHERE user_id=$1 AND segment_id=$2)`

	for _, upd := range update.SegmentsToAdd {
		var exists bool

		if err := db.QueryRow(isSegment, upd).Scan(&exists); err != nil {
			log.Fatalf("Unable to scan the row.\n %v", err)
		}

		if !exists {
			return false
		}
	}

	for _, dlt := range update.SegmentsToDelete {
		var exists bool

		segment_id, code := getSegmentIDByName(dlt)

		switch code {
		case http.StatusInternalServerError:
			log.Fatalln("Unable to scan the row")
		case http.StatusBadRequest:
			return false
		case http.StatusOK:

		}

		if err := db.QueryRow(isUser, update.UserID, segment_id).Scan(&exists); err != nil {
			log.Fatalf("Unable to scan the row.\n %v", err)
		}

		if !exists {
			return false
		}
	}

	return true
}

func updateRelations(update models.Update) int64 {
	db := createConnection()
	defer db.Close()

	sqlUpdate := `INSERT INTO relations(user_id, segment_id) SELECT $1, segments.id FROM segments WHERE segments.name=$2`

	sqlDelete := `DELETE FROM relations WHERE segment_id=(SELECT segments.id FROM segments WHERE segments.name=$1)`

	allRowsAffected := int64(1)

	for _, upd := range update.SegmentsToAdd {
		res, err := db.Exec(sqlUpdate, update.UserID, upd)
		if err != nil {
			log.Fatalf("Unable to execute the query. %v", err)
		}

		rowsAffected, err := res.RowsAffected()
		if err != nil {
			log.Fatalf("Error while checking the affected rows. %v", err)
		}

		allRowsAffected += rowsAffected
	}

	for _, dlt := range update.SegmentsToDelete {
		res, err := db.Exec(sqlDelete, dlt)
		if err != nil {
			log.Fatalf("Unable to execute the query. %v", err)
		}

		rowsAffected, err := res.RowsAffected()
		if err != nil {
			log.Fatalf("Error while checking the affected rows. %v", err)
		}

		allRowsAffected += rowsAffected
	}

	return allRowsAffected
}

func isUser(id int64) (bool, error) {
	db := createConnection()
	defer db.Close()

	sqlStatement := `SELECT EXISTS (SELECT * FROM users WHERE id=$1)`

	var exists bool

	if err := db.QueryRow(sqlStatement, id).Scan(&exists); err != nil {
		return false, fmt.Errorf("Unable to scan the row.\n %v", err)
	}

	return exists, nil
}

func insertUser(id int64) {
	db := createConnection()
	defer db.Close()

	sqlStatement := `INSERT INTO users(id) VALUES($1)`

	_, err := db.Exec(sqlStatement, id)
	if err != nil {
		log.Fatalf("Unable to execute the query (insert into users). %v", err)
	}

	log.Printf("Inserted single record %v.", id)
}

func GetUserSegments(w http.ResponseWriter, r *http.Request) {
	var user models.Users

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Fatalf("Unable to decoe request body. %v", err)
	}

	segments, err := getSegmentsByUser(int64(user.ID))
	if err != nil {
		log.Fatalf("Unable to get user segments. %v", err)
	}

	res := response{
		Time:    time.Now(),
		Code:    http.StatusOK,
		Message: strings.Join(segments, ","),
	}

	log.Println(res)

	json.NewEncoder(w).Encode(res)
}

func getSegmentsByUser(id int64) ([]string, error) {
	db := createConnection()
	defer db.Close()

	sqlStatement := `SELECT T1.segment_id, T2.name FROM relations AS T1 LEFT JOIN segments AS T2 ON T1.segment_id=T2.id AND T1.user_id=$1`

	rows, err := db.Query(sqlStatement, id)
	if err != nil {
		return nil, fmt.Errorf("Unable to exec the query.\n %v", err)
	}
	defer rows.Close()

	var segments []string

	for rows.Next() {
		var seg models.Segments
		if err := rows.Scan(&seg.ID, &seg.Name); err != nil {
			return segments, fmt.Errorf("Unable to scan the row.\n %v", err)
		}

		segments = append(segments, seg.Name)
	}

	if err = rows.Err(); err != nil {
		return segments, err
	}

	return segments, nil
}
