package repository

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

type Vajbs struct {
	db *sqlx.DB
}

type Vajb struct {
	ID          int            `db:"id"`
	CreatorID   int            `db:"creator_id"`
	Name        string         `db:"name"`
	Description string         `db:"description"`
	Address     string         `db:"address"`
	Region      string         `db:"region"`
	Date        time.Time      `db:"date"`
	HeaderImage sql.NullString `db:"header_image"`
}

func NewVajbs(db *sqlx.DB) *Vajbs {
	return &Vajbs{db: db}
}

func (vajbs *Vajbs) CreateVajb(creatorID int, name, description, address, region, headerImage string, date time.Time) (*Vajb, error) {
	vajb := &Vajb{
		CreatorID:   creatorID,
		Name:        name,
		Description: description,
		Address:     address,
		Region:      region,
		Date:        date,
	}
	if headerImage != "" {
		vajb.HeaderImage = sql.NullString{String: headerImage, Valid: true}
	}

	var id int
	err := vajbs.db.QueryRow(`INSERT INTO vajbs (creator_id, name, description, address, region, date, header_image) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`, creatorID, name, description, address, region, date, headerImage).Scan(&id)
	if err != nil {
		return nil, err
	}
	vajb.ID = id

	return vajb, nil
}

func (vajbs *Vajbs) GetVajbByID(id int) (*Vajb, error) {
	vajb := &Vajb{}
	err := vajbs.db.Get(vajb, "SELECT * FROM vajbs WHERE id = $1", id)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return vajb, nil
}

func (vajbs *Vajbs) GetFullRegionName(region string) string {
	switch region {
	case "praha":
		return "Praha"
	case "plzensky":
		return "Plzeňský kraj"
	case "karlovarsky":
		return "Karlovarský kraj"
	case "ustecky":
		return "Ústecký kraj"
	case "liberecky":
		return "Liberecký kraj"
	case "kralovehradecky":
		return "Královéhradecký kraj"
	case "pardubicky":
		return "Pardubický kraj"
	case "vysocina":
		return "Vysočina"
	case "jihomoravsky":
		return "Jihomoravský kraj"
	case "olomoucky":
		return "Olomoucký kraj"
	case "zlinsky":
		return "Zlínský kraj"
	case "moravskoslezsky":
		return "Moravskoslezský kraj"
	case "stredocesky":
		return "Středočeský kraj"
	case "jihocesky":
		return "Jihočeský kraj"
	}
	return ""
}
