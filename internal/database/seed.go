package database

import (
	"database/sql"
	"math/rand"

	randomDataTime "github.com/duktig-solutions/go-random-date-generator"
	"github.com/go-faker/faker/v4"
	"github.com/jmoiron/sqlx"
)

type User struct {
	ID          int            `db:"id"`
	DisplayName string         `db:"display_name" faker:"name"`
	Username    string         `db:"username" faker:"username"`
	Password    string         `db:"password" faker:"password"`
	ProfilePic  sql.NullString `db:"profile_image"`
}

type Vajb struct {
	ID          int            `db:"id"`
	CreatorID   int            `db:"creator_id"`
	Name        string         `db:"name" faker:"sentence"`
	Description string         `db:"description" faker:"paragraph"`
	Date        string         `db:"date"`
	Region      string         `db:"region" faker:"oneof: praha, pardubicky, jihomoravsky, karlovarsky, ustecky, plzensky, stredocesky, jihocesky, vysocina, liberecky, olomoucky, moravskoslezsky, zlinsky, kralovehradecky"`
	Address     string         `db:"address" faker:"word"`
	HeaderImage sql.NullString `db:"header_image"`
}

func createUsers(db *sqlx.DB) error {
	users := []User{}
	err := faker.FakeData(&users)
	if err != nil {
		return err
	}
	for _, user := range users {
		_, err := db.NamedExec("INSERT INTO users (display_name, username, password) VALUES (:display_name, :username, :password)", user)
		if err != nil {
			return err
		}
	}
	return nil
}

func dropUsers(db *sqlx.DB) error {
	_, err := db.Exec("DELETE FROM users")
	if err != nil {
		return err
	}
	return nil
}

func getUsers(db *sqlx.DB) ([]User, error) {
	users := []User{}
	err := db.Select(&users, "SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	return users, nil
}

func createVajbs(db *sqlx.DB, users []User) error {
	vajbs := []Vajb{}
	err := faker.FakeData(&vajbs)
	if err != nil {
		return err
	}
	for _, vajb := range vajbs {
		vajb.CreatorID = users[rand.Intn(len(users))].ID
		date, err := randomDataTime.GenerateDateTime("2025-01-01 00:00:00", "2025-12-31 00:00:00")
		if err != nil {
			return err
		}
		vajb.Date = date
		_, err = db.NamedExec("INSERT INTO vajbs (creator_id, name, description, address, region, date) VALUES (:creator_id, :name, :description, :address, :region, :date)", vajb)
		if err != nil {
			return err
		}
	}
	return nil
}

func dropVajbs(db *sqlx.DB) error {
	_, err := db.Exec("DELETE FROM vajbs")
	if err != nil {
		return err
	}
	return nil
}

func getVajbs(db *sqlx.DB) ([]Vajb, error) {
	vajbs := []Vajb{}
	err := db.Select(&vajbs, "SELECT * FROM vajbs")
	if err != nil {
		return nil, err
	}
	return vajbs, nil
}

func dropJoinedVajbs(db *sqlx.DB) error {
	_, err := db.Exec("DELETE FROM joined_vajbs")
	if err != nil {
		return err
	}
	return nil
}

func createJoinedVajbs(db *sqlx.DB, users []User, vajbs []Vajb) error {
	for _, user := range users {
		for _, vajb := range vajbs {
			if rand.Intn(2) == 0 || user.ID == vajb.CreatorID {
				continue
			}
			_, err := db.Exec("INSERT INTO joined_vajbs (user_id, vajb_id) VALUES ($1, $2)", user.ID, vajb.ID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func seed(db *sqlx.DB) error {
	err := dropUsers(db)
	if err != nil {
		return err
	}

	err = createUsers(db)
	if err != nil {
		return err
	}
	users, err := getUsers(db)
	if err != nil {
		return err
	}

	err = dropVajbs(db)
	if err != nil {
		return err
	}

	err = createVajbs(db, users)
	if err != nil {
		return err
	}

	vajbs, err := getVajbs(db)
	if err != nil {
		return err
	}

	err = dropJoinedVajbs(db)
	if err != nil {
		return err
	}

	err = createJoinedVajbs(db, users, vajbs)
	if err != nil {
		return err
	}

	return nil
}
