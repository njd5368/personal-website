package database

import (
	"database/sql"
	"errors"
	"strconv"
)

type Database interface {
	StartDatabase() error
	CreateProject(p Project) (*Project, error)
}

type SQLiteDatabase struct {
	db *sql.DB
}

func NewSQLiteDatabase(db *sql.DB) *SQLiteDatabase {
	return &SQLiteDatabase{
		db: db,
	}
}

func (d *SQLiteDatabase) StartDatabase() error {
	queries := []string{
		`
		CREATE TABLE IF NOT EXISTS Image(
			ID INTEGER PRIMARY KEY AUTOINCREMENT,
			Image BLOB
		);
		`,
		`
		CREATE TABLE IF NOT EXISTS Project(
			ID INTEGER PRIMARY KEY AUTOINTCREMENT,
			Name TEXT NOT NULL UNIQUE,
			Description TEXT NOT NULL,
			Date TEXT NOT NULL,
			Type TEXT NOT NULL,
			Image INTEGER,
			File BLOB NOT NULL
			FOREIGN KEY(Image) REFERENCES Image(ID)
		);
		`,
		`
		CREATE TABLE IF NOT EXISTS Language(
			ID INTEGER PRIMARY KEY AUTOINCREMENT,
			Name TEXT NOT NULL UNIQUE
		);
		`,
		`
		CREATE TABLE IF NOT EXISTS ProjectLanguage(
			Project INTEGER NOT NULL,
			Language INTEGER NOT NULL,
			PRIMARY KEY(Project, Language),
			FOREIGN KEY(Project) REFERENCES Project(ID),
			FOREIGN KEY(Language) REFERENCES Language(ID)
		);
		`,
		`
		CREATE TABLE IF NOT EXISTS Technology(
			ID INTEGER PRIMARY KEY AUTOINCREMENT,
			Name TEXT NOT NULL UNIQUE
		);
		`,
		`
		CREATE TABLE IF NOT EXISTS ProjectTechnology(
			Project INTEGER NOT NULL,
			Technology INTEGER NOT NULL,
			PRIMARY KEY(Project, Technology),
			FOREIGN KEY(Project) REFERENCES Project(ID),
			FOREIGN KEY(Technology) REFERENCES Technology(ID)
		);
		`,
	}

	for _, query := range queries {
		_, err := d.db.Exec(query)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *SQLiteDatabase) CreateProject(p Project, i []byte) (*Project, error) {
	r, err := d.db.Exec("INSERT INTO Image(Image) VALUES(?)", i)
	if err != nil {
		return nil, err
	}

	id, err := r.LastInsertId()
	if err != nil {
		return nil, err
	}

	p.Image = id

	r, err = d.db.Exec("INSERT INTO Project(Name, Description, Date, Type, Image, File) VALUES(?,?,?,?,?,?)",
		p.Name, p.Description, p.Date, p.Type, p.Image, p.File)
	if err != nil {
		return nil, err
	}

	id, err = r.LastInsertId()
	if err != nil {
		return nil, err
	}

	p.ID = id

	if len(p.Languages) != 0 {
		statement := `INSERT INTO ProjectLanguage(Project, Language) VALUES`
		for i, l := range(p.Languages) {
			if i != 0 {
				statement += `, `
			}
			statement +=
			`
			(` + strconv.FormatInt(p.ID, 10) + `, (IF EXISTS ( SELECT 1 FROM Language WHERE Name = '` + l + `')
			BEGIN
				SELECT ID FROM Language WHERE Name = '` + l + `'
			END
			ELSE
			BEGIN
				INSERT INTO Language(Name) VALUES('` + l + `') AND SELECT ID FROM Language WHERE Name = '` + l + `'
			END
			END IF))`
		}
		statement += ";"

		_, err := d.db.Exec(statement)
		if err != nil {
			return nil, err
		}
	}

	if len(p.Technologies) != 0 {
		statement := `INSERT INTO ProjectTechnology(Project, Technology) VALUES`
		for i, t := range(p.Technologies) {
			if i != 0 {
				statement += `, `
			}
			statement +=
			`
			(` + strconv.FormatInt(p.ID, 10) + `, (IF EXISTS ( SELECT 1 FROM Technology WHERE Name = '` + t + `')
			BEGIN
				SELECT ID FROM Technology WHERE Name = '` + t + `'
			END
			ELSE
			BEGIN
				INSERT INTO Technology(Name) VALUES('` + t + `') AND SELECT ID FROM Technology WHERE Name = '` + t + `'
			END
			END IF))`
		}
		statement += `;`

		_, err := d.db.Exec(statement)
		if err != nil {
			return nil, err
		}
	}

	return &p, nil
}

func (d *SQLiteDatabase) GetProjects() ([]Project, error) {
	query := `
	SELECT Project.ID, Project.Name, Project.Description, Project.Date, Project.Type, Project.Image, Language.Name, Technology.Name
	FROM Project 
	LEFT JOIN (
		ProjectLanguage 
		INNER JOIN Language 
		ON ProjectLanguage.Language=Language.ID
	)
	ON Project.ID=ProjectLanguage.Project
	LEFT JOIN (
		ProjectTechnology
		INNER JOIN Technology
		ON ProjectTechnology.Technology=Technology.ID
	)
	ON Project.ID=ProjectTechnology.Project
	`
	r, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var projects []Project
	for r.Next() {
		var p Project
		if err := r.Scan(&p.ID, &p.Name, &p.Description, &p.Date, &p.Type, &p.Image, &p.Languages, &p.Technologies); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	return projects, nil
}

func (d *SQLiteDatabase) GetProjectByName(n string) (*Project, error) {
	query := `
	SELECT Project.ID, Project.Name, Project.Description, Project.Date, Project.Type, Project.Image, Project.File, Language.Name, Technology.Name
	FROM Project 
	LEFT JOIN (
		ProjectLanguage 
		INNER JOIN Language 
		ON ProjectLanguage.Language=Language.ID
	)
	ON Project.ID=ProjectLanguage.Project
	LEFT JOIN (
		ProjectTechnology
		INNER JOIN Technology
		ON ProjectTechnology.Technology=Technology.ID
	)
	ON Project.ID=ProjectTechnology.Project
	WHERE Project.Name='` + n + `'
	`

	r := d.db.QueryRow(query)
	if r.Err() == sql.ErrNoRows {
		return nil, errors.New("No project with name '" + n + "'")
	}

	var p Project
	if err := r.Scan(&p.ID, &p.Name, &p.Description, &p.Date, &p.Type, &p.Image, &p.File, &p.Languages, &p.Technologies); err != nil {
		return nil, err
	}

	return &p, nil
}

func (d *SQLiteDatabase) GetImageByID(id int) ([]byte, error) {
	query := `
	SELECT Image
	FROM Image
	WHERE ID=` + strconv.Itoa(id) + `
	`

	r := d.db.QueryRow(query)
	if r.Err() == sql.ErrNoRows {
		return nil, errors.New("No image with id " + strconv.Itoa(id))
	}

	var i []byte
	if err := r.Scan(&i); err != nil {
		return nil, err
	}

	return i, nil
}
