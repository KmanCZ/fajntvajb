package repository

import (
	"bytes"
	"database/sql"
	"fajntvajb/internal/files"
	"strconv"
	"strings"
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
	err := vajbs.db.QueryRow(`INSERT INTO vajbs (creator_id, name, description, address, region, date, header_image) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`, creatorID, name, description, address, region, date, vajb.HeaderImage).Scan(&id)
	if err != nil {
		return nil, err
	}
	vajb.ID = id

	return vajb, nil
}

func (vajbs *Vajbs) DeleteVajb(id int) error {
	_, err := vajbs.db.Exec("DELETE FROM vajbs WHERE id = $1", id)
	return err
}

func (vajbs *Vajbs) UpdateVajb(id, creatorID int, name, description, address, region, headerImage string, date time.Time) error {
	var headerImageNull sql.NullString
	if headerImage != "" {
		headerImageNull = sql.NullString{String: headerImage, Valid: true}
	}
	_, err := vajbs.db.Exec(`UPDATE vajbs SET creator_id = $1, name = $2, description = $3, address = $4, region = $5, date = $6, header_image = $7 WHERE id = $8`, creatorID, name, description, address, region, date, headerImageNull, id)
	return err
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

func (vajbs *Vajbs) JoinVajb(id, userId int) error {
	_, err := vajbs.db.Exec(`INSERT INTO joined_vajbs (user_id, vajb_id) VALUES ($1, $2)`, userId, id)
	return err
}

func (vajbs *Vajbs) UnjoinVajb(id, userId int) error {
	_, err := vajbs.db.Exec(`DELETE FROM joined_vajbs WHERE user_id = $1 AND vajb_id = $2`, userId, id)
	return err
}

func (vajbs *Vajbs) GetIsJoinedToVajb(id, userId int) (bool, error) {
	var res bool
	err := vajbs.db.Get(&res, "SELECT COUNT(*) > 0 FROM joined_vajbs WHERE user_id = $1 AND vajb_id = $2", userId, id)
	return res, err
}

func (vajbs *Vajbs) GetVajbParticipants(id int) ([]User, error) {
	users := []User{}

	err := vajbs.db.Select(&users, `SELECT users.id, users.display_name, users.profile_image FROM users, joined_vajbs WHERE users.id = joined_vajbs.user_id AND joined_vajbs.vajb_id = $1`, id)
	return users, err
}

func (vajbs *Vajbs) GetVajbs(name, region string, from, to time.Time, number, offset int) ([]Vajb, error) {
	var query bytes.Buffer
	var dateSet bool
	var values []interface{}
	query.WriteString("SELECT * FROM vajbs WHERE ")
	if name != "" {
		name = strings.ToLower(name)
		name = "%" + name + "%"
		query.WriteString("LOWER(name) LIKE $")
		query.WriteString(strconv.Itoa(len(values) + 1))
		values = append(values, name)
	}
	if region != "" {
		if len(values) > 0 {
			query.WriteString(" AND ")
		}
		query.WriteString("region = $")
		query.WriteString(strconv.Itoa(len(values) + 1))
		values = append(values, region)
	}
	if !from.IsZero() {
		if len(values) > 0 {
			query.WriteString(" AND ")
		}
		query.WriteString("date >= $")
		query.WriteString(strconv.Itoa(len(values) + 1))
		dateSet = true
		values = append(values, from)
	}
	if !to.IsZero() {
		if len(values) > 0 {
			query.WriteString(" AND ")
		}
		query.WriteString("date <= $")
		query.WriteString(strconv.Itoa(len(values) + 1))
		dateSet = true
		values = append(values, to)
	}
	if !dateSet {
		if len(values) > 0 {
			query.WriteString(" AND ")
		}
		query.WriteString("date >= CURRENT_DATE")
	}
	query.WriteString(" ORDER BY date")
	if number > 0 {
		query.WriteString(" LIMIT $")
		query.WriteString(strconv.Itoa(len(values) + 1))
		values = append(values, number)
	}
	if offset > 0 {
		query.WriteString(" OFFSET $")
		query.WriteString(strconv.Itoa(len(values) + 1))
		values = append(values, offset)
	}

	res := []Vajb{}
	err := vajbs.db.Select(&res, query.String(), values...)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(res); i++ {
		vajb := &res[i]
		if vajb.HeaderImage.Valid {
			vajb.HeaderImage.String = files.GetVajbPicPath(vajb.HeaderImage)
		}

		vajb.Region = vajbs.GetFullRegionName(vajb.Region)
	}

	return res, nil
}

func (vajbs *Vajbs) GetMyVajbs(userID int) ([]Vajb, error) {
	res := []Vajb{}
	err := vajbs.db.Select(&res, "SELECT * FROM vajbs WHERE creator_id = $1 AND date >= CURRENT_DATE ORDER BY date", userID)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(res); i++ {
		vajb := &res[i]
		if vajb.HeaderImage.Valid {
			vajb.HeaderImage.String = files.GetVajbPicPath(vajb.HeaderImage)
		}

		vajb.Region = vajbs.GetFullRegionName(vajb.Region)
	}

	return res, nil
}

func (vajbs *Vajbs) GetJoinedVajbs(userID int) ([]Vajb, error) {
	res := []Vajb{}
	err := vajbs.db.Select(&res, "SELECT vajbs.* FROM vajbs, joined_vajbs WHERE vajbs.id = joined_vajbs.vajb_id AND joined_vajbs.user_id = $1 AND vajbs.date >= CURRENT_DATE ORDER BY vajbs.date", userID)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(res); i++ {
		vajb := &res[i]
		if vajb.HeaderImage.Valid {
			vajb.HeaderImage.String = files.GetVajbPicPath(vajb.HeaderImage)
		}

		vajb.Region = vajbs.GetFullRegionName(vajb.Region)
	}

	return res, nil
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
