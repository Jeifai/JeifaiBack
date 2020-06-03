package data

import (
	"database/sql"
	"fmt"
	"time"
)

type User struct {
	Id                int            `db:"id"`
	Uuid              string         `db:"uuid"`
	UserName          string         `db:"username" validate:"min=1"`
	Email             string         `db:"email" validate:"email"`
	Password          string         `db:"password"`
	CreatedAt         time.Time      `db:"createdat"`
	DeletedAt         time.Time      `db:"deletedat"`
	FirstName         sql.NullString `db:"firstname"`
	LastName          sql.NullString `db:"lastname"`
	DateOfBirth       sql.NullString `db:"dateofbirth"`
	Country           sql.NullString `db:"country"`
	City              sql.NullString `db:"city"`
	Gender            sql.NullString `db:"gender"`
	CurrentPassword   string         `validate:"required,eqfield=Password"`
	NewPassword       string         `db:"newpassword" validate:"eqfield=RepeatNewPassword"`
	RepeatNewPassword string         `validate:"eqfield=NewPassword"`
}

type Session struct {
	Id        int
	Uuid      string
	Email     string
	UserId    int
	CreatedAt time.Time
}

// Create a new session for an existing user
func (user *User) CreateSession() (session Session, err error) {
	statement := `INSERT INTO sessions (uuid, email, user_id, created_at)
                  VALUES ($1, $2, $3, $4)
                  RETURNING id, uuid, email, user_id, created_at`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	// use QueryRow to return a row and scan the returned id into the Session struct
	err = stmt.QueryRow(
		createUUID(),
		user.Email,
		user.Id,
		time.Now(),
	).Scan(
		&session.Id,
		&session.Uuid,
		&session.Email,
		&session.UserId,
		&session.CreatedAt,
	)
	return
}

// Get the session for an existing user
func (user *User) Session() (session Session, err error) {
	session = Session{}
	err = Db.QueryRow(`SELECT
                        id, 
                        uuid, 
                        email, 
                        user_id, 
                        created_at
                      FROM sessions
                      WHERE user_id = $1`,
		user.Id,
	).
		Scan(
			&session.Id,
			&session.Uuid,
			&session.Email,
			&session.UserId,
			&session.CreatedAt,
		)
	return
}

// Check if session is valid in the database
func (session *Session) Check() (valid bool, err error) {
	err = Db.QueryRow(`SELECT
                        id,
                        uuid,
                        email,
                        user_id,
                        created_at
                      FROM sessions
                      WHERE uuid = $1`,
		session.Uuid,
	).
		Scan(
			&session.Id,
			&session.Uuid,
			&session.Email,
			&session.UserId,
			&session.CreatedAt,
		)
	if err != nil {
		valid = false
		return
	}
	if session.Id != 0 {
		valid = true
	}
	return
}

// Delete session from database
func (session *Session) DeleteByUUID() (err error) {
	statement := `DELETE FROM sessions
                  WHERE uuid = $1`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(session.Uuid)
	return
}

// Create a new user, save user info into the database
func (user *User) Create() (err error) {
	// Postgres does not automatically return the last insert id, because it would be wrong to assume
	// you're always using a sequence.You need to use the RETURNING keyword in your insert to get this
	// information from postgres.
	statement := `INSERT INTO users
                  (uuid, username, email, password, createdat)
                  VALUES ($1, $2, $3, $4, $5)
                  RETURNING id, uuid, createdat`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	// use QueryRow to return a row and scan the returned id into the User struct
	err = stmt.QueryRow(
		createUUID(),
		user.UserName,
		user.Email,
		Encrypt(user.Password),
		time.Now(),
	).Scan(
		&user.Id,
		&user.Uuid,
		&user.CreatedAt,
	)
	return
}

// Get a single user given the email
func UserByEmail(email string) (user User, err error) {
	user = User{}
	err = Db.QueryRow(`SELECT
                        id,
                        uuid,
                        username,
                        email,
                        password,
                        createdat,
                        firstname,
                        lastname,
                        TO_CHAR(dateofbirth, 'YYYY-MM-DD'),
                        country,
                        city,
                        gender
                      FROM users
                      WHERE email = $1`,
		email,
	).
		Scan(
			&user.Id,
			&user.Uuid,
			&user.UserName,
			&user.Email,
			&user.Password,
			&user.CreatedAt,
			&user.FirstName,
			&user.LastName,
			&user.DateOfBirth,
			&user.Country,
			&user.City,
			&user.Gender,
		)
	return
}

// Get a single user given the email
func UserById(id int) (user User, err error) {
	user = User{}
	err = Db.QueryRow(`SELECT
                        id,
                        uuid,
                        username,
                        email,
                        password,
                        createdat,
                        firstname,
                        lastname,
                        TO_CHAR(dateofbirth, 'YYYY-MM-DD'),
                        country,
                        city,
                        gender
                      FROM users
                      WHERE id = $1`,
		id,
	).
		Scan(
			&user.Id,
			&user.Uuid,
			&user.UserName,
			&user.Email,
			&user.Password,
			&user.CreatedAt,
			&user.FirstName,
			&user.LastName,
			&user.DateOfBirth,
			&user.Country,
			&user.City,
			&user.Gender,
		)
	return
}

func (user User) UpdateUser() {
	fmt.Println("Starting UpdateUser...")
	statement := `UPDATE users SET 
                    username=$2,
                    email=$3,
                    password=$4,
                    firstname=$5,
                    lastname=$6,
                    country=$7,
                    city=$8,
                    gender=$9,
                    dateofbirth=$10
                  WHERE id=$1`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		user.Id,
		user.UserName,
        user.Email,
		user.NewPassword,
		user.FirstName.String,
		user.LastName.String,
		user.Country.String,
		user.City.String,
		user.Gender.String,
		user.DateOfBirth.String)
}
