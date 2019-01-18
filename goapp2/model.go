// model.go

package main
    
import (
	"database/sql"
	"fmt"
)

type models struct {
   Models2id int `json:"models2_id"`
   Name string `json:"name"`
   Height string `json:"height"`
   Weight string `json:"weight"`
   Isdisabled string `json:"is_disabled"`
   }
func (u *models) getModels(db *sql.DB) error {
	statement := fmt.Sprintf(`SELECT models2_id,name, 
height, 
weight, 
is_disabled FROM models2 WHERE models2_id=%d`, u.Models2id)
return db.QueryRow(statement).Scan(&u.Models2id,
&u.Name,
&u.Height,
&u.Weight,
&u.Isdisabled)
}

func (u *models) updateModels(db *sql.DB) error {
	statement := fmt.Sprintf(`UPDATE models2 SET name = '%s', 
height = '%s', 
weight = '%s' 
WHERE models2_id=%d`, u.Name,
u.Height,
u.Weight,
u.Models2id)
_, err := db.Exec(statement)
	return err
}

func (u *models) deleteModels(db *sql.DB) error {
	statement := fmt.Sprintf(`DELETE FROM models2 WHERE models2_id=%d`, u.Models2id)
	_, err := db.Exec(statement)
	return err
}

func (u *models) createModels(db *sql.DB) error {
	statement := fmt.Sprintf(`INSERT INTO models2(
name, 
height, 
weight) VALUES(
'%s', 
'%s', 
'%s' 
)`,u.Name, 
u.Height, 
u.Weight) 

	_, err := db.Exec(statement)

	if err != nil {
		return err
	}

	err = db.QueryRow(`SELECT LAST_INSERT_ID()`).Scan(&u.Models2id)

	if err != nil {
		return err
	}

	return nil
}

func getModels2(db *sql.DB, start, count int) ([]models, error) {
	statement := fmt.Sprintf(`SELECT models2_id,name, 
height, 
weight 
FROM models2 LIMIT %d OFFSET %d`, count, start)
	rows, err := db.Query(statement)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	models2 := []models{}

	for rows.Next() {
		var u models
		if err := rows.Scan(&u.Models2id,
&u.Name, 
&u.Height, 
&u.Weight); 
err != nil {
			return nil, err
		}
		models2 = append(models2, u)
	}

	return models2, nil
}