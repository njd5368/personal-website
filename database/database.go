package database

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"
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
		PRAGMA foreign_keys = ON;
		`,
		`
		CREATE TABLE IF NOT EXISTS Image(
			ID INTEGER PRIMARY KEY AUTOINCREMENT,
			Image BLOB
		);
		`,
		`
		CREATE TABLE IF NOT EXISTS Project(
			ID INTEGER PRIMARY KEY AUTOINCREMENT,
			Name TEXT NOT NULL UNIQUE,
			Description TEXT NOT NULL,
			Date TEXT NOT NULL,
			Type TEXT NOT NULL,
			Image INTEGER NOT NULL,
			File BLOB NOT NULL,
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

func (d *SQLiteDatabase) CreateImage(i []byte) (int64, error) {
	r, err := d.db.Exec("INSERT INTO Image(Image) VALUES(?)", i)
	if err != nil {
		return -1, err
	}

	id, err := r.LastInsertId()
	if err != nil {
		return -1, err
	}

	return id, nil
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
		existanceChecks := `INSERT OR IGNORE INTO Language(Name) VALUES`
		statement := `INSERT INTO ProjectLanguage(Project, Language) VALUES`
		for i, l := range(p.Languages) {
			if i != 0 {
				existanceChecks += `, `
				statement += `, `
			}
			existanceChecks += `('` + l + `')`
			statement += `(` + strconv.FormatInt(p.ID, 10) + `, (SELECT ID FROM Language WHERE Name = '` + l + `'))`
		}
		existanceChecks += `;`
		statement += `;`

		_, err := d.db.Exec(existanceChecks + statement)
		if err != nil {
			return nil, err
		}
	}

	if len(p.Technologies) != 0 {
		existanceChecks := `INSERT OR IGNORE INTO Technology(Name) VALUES`
		statement := `INSERT INTO ProjectTechnology(Project, Technology) VALUES`
		for i, t := range(p.Technologies) {
			if i != 0 {
				existanceChecks += `, `
				statement += `, `
			}
			existanceChecks += `('` + t + `')`
			statement += `(` + strconv.FormatInt(p.ID, 10) + `, (SELECT ID FROM Technology WHERE Name = '` + t + `'))`
		}
		existanceChecks += `;`
		statement += `;`

		_, err := d.db.Exec(existanceChecks + statement)
		if err != nil {
			return nil, err
		}
	}

	return &p, nil
}

func (d *SQLiteDatabase) GetProjects() ([]Project, error) {
	query := `
	SELECT Project.ID, Project.Name, Project.Description, Project.Date, Project.Type, Project.Image, Languages, Technologies
	FROM Project
	LEFT JOIN (
		SELECT ProjectLanguage.Project AS ProjectLanguageID, GROUP_CONCAT(Language.Name) AS Languages
		FROM ProjectLanguage
		LEFT JOIN Language
		ON ProjectLanguage.Language=Language.ID
		GROUP BY ProjectLanguageID
	)
	ON Project.ID=ProjectLanguageID
	LEFT JOIN (
		SELECT ProjectTechnology.Project AS ProjectTechnologyID, GROUP_CONCAT(Technology.Name) AS Technologies
		FROM ProjectTechnology
		LEFT JOIN Technology
		ON ProjectTechnology.Technology=Technology.ID
		GROUP BY ProjectTechnologyID
	)
	ON Project.ID=ProjectTechnologyID
	`
	r, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var projects []Project
	for r.Next() {
		var p Project
		var languages string
		var techs string
		if err := r.Scan(&p.ID, &p.Name, &p.Description, &p.Date, &p.Type, &p.Image, &languages, &techs); err != nil {
			return nil, err
		}
		p.Languages = strings.Split(languages, ",")
		p.Technologies = strings.Split(techs, ",")
		projects = append(projects, p)
	}
	return projects, nil
}

func (d *SQLiteDatabase) GetProjectByName(n string) (*Project, error) {
	query := `
	SELECT Project.ID, Project.Name, Project.Description, Project.Date, Project.Type, Project.Image, Project.File, Languages, Technologies
	FROM Project
	LEFT JOIN (
		SELECT ProjectLanguage.Project AS ProjectLanguageID, GROUP_CONCAT(Language.Name) AS Languages
		FROM ProjectLanguage
		LEFT JOIN Language
		ON ProjectLanguage.Language=Language.ID
		GROUP BY ProjectLanguageID
	)
	ON Project.ID=ProjectLanguageID
	LEFT JOIN (
		SELECT ProjectTechnology.Project AS ProjectTechnologyID, GROUP_CONCAT(Technology.Name) AS Technologies
		FROM ProjectTechnology
		LEFT JOIN Technology
		ON ProjectTechnology.Technology=Technology.ID
		GROUP BY ProjectTechnologyID
	)
	ON Project.ID=ProjectTechnologyID
	WHERE Project.Name='` + n + `'
	`

	r := d.db.QueryRow(query)
	if r.Err() == sql.ErrNoRows {
		return nil, errors.New("No project with name '" + n + "'")
	}
	
	var p Project
	var languages string
	var techs string
	if err := r.Scan(&p.ID, &p.Name, &p.Description, &p.Date, &p.Type, &p.Image, &p.File, &languages, &techs); err != nil {
		return nil, err
	}
	p.Languages = strings.Split(languages, ",")
	p.Technologies = strings.Split(techs, ",")
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
