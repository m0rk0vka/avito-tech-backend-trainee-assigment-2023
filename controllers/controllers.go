package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

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

	deletedRows := deleteSegmentFromRelations(int64(id))
	msg := fmt.Sprintf("Successifully deleted relations. Total rows/records %v.", deletedRows)

	deletedRows = deleteSegmentFromSegment(int64(id))
	msg = fmt.Sprintf(msg, "Successifully deleted segment. Total rows/records %v.", deletedRows)

	res := response{
		ID:      int64(id),
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)
}

func deleteSegmentFromSegment(id int64) int64 {
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

func deleteSegmentFromRelations(id int64) int64 {
	db := createConnection()
	defer db.Close()

	sqlStatement := `DELETE FROM relations WHERE segment_id=$1`

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
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("Unable to convert id. %v", err)
	}

	var update models.Update

	err = json.NewDecoder(r.Body).Decode(&update)
	if err != nil {
		log.Fatalf("Unable to decode request body. %v", err)
	}

	updatedRows := updateUser(int64(id), update)

	msg := fmt.Sprintf("Successifully updated user segments. Total rows/records affected %v", updatedRows)

	res := response{
		ID:      int64(id),
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)
}

func updateUser(id int64, update models.Update) int64 {
	db := createConnection()
	defer db.Close()

	sqlUpdate := `INSERT INTO relations(user_id, segment_id) SELECT $1, segments.id FROM segments WHERE segments.name=$2`

	sqlDelete := `DELETE FROM relations WHERE segment_id=(SELECT segments.id FROM segments WHERE segments.name=$1)`

	allRowsAffected := int64(0)

	for _, upd := range update.SegmentsToAdd {
		res, err := db.Exec(sqlUpdate, id, upd)
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
		res, err := db.Exec(sqlDelete, id, dlt)
		if err != nil {
			log.Fatalf("Unable to execute the query. %v", err)
		}

		rowsAffected, err := res.RowsAffected()
		if err != nil {
			log.Fatalf("Error while checking the affected rows. %v", err)
		}

		allRowsAffected += rowsAffected
	}

	fmt.Println("Total rows/records affected %v", allRowsAffected)

	return allRowsAffected
}

func GetUserByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("Unable to convert id. %v", err)
	}

	segments, err := getSegmentsByUser(int64(id))
	if err != nil {
		log.Fatalf("Unable to get user segments. %v", err)
	}

	res := response{
		ID:      int64(id),
		Message: strings.Join(segments, ","),
	}

	json.NewEncoder(w).Encode(res)
}

func getSegmentsByUser(id int64) ([]string, error) {
	db := createConnection()
	defer db.Close()

	sqlStatement := `SELECT T1.segment_id, T2.name FROM relations AS T1 LEFT JOIN segments AS T2 ON T1.segment_id=T2.id AND T1.user_id=$1`

	rows, err := db.Query(sqlStatement, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var segments []string

	for rows.Next() {
		var seg models.Segments
		if err := rows.Scan(&seg.ID, &seg.Name); err != nil {
			fmt.Println("Unable to scan the row. %v", err)
			return segments, err
		}

		segments = append(segments, seg.Name)
	}

	if err = rows.Err(); err != nil {
		return segments, err
	}

	fmt.Println("Successifully got user segments.")

	return segments, nil
}
