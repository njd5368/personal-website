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

// GetProjects returns up to 10 projects that meet the required parameters.
//
// page (int): the page number for pagination. If less than, 1, defaults to 1. If over the max pages, an empty Project struct is returned.
// 
func (d *SQLiteDatabase) GetProjects(page int, search string, types []string, languages []string, technologies []string, first string, last string) ([]Project, error) {
	if page < 1 {
		page = 1
	}
	if len(first) != 0 && len(last) != 0 {
		first = ""
		last = ""
	}
	
	selectClause := `
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
	whereClause := `
	WHERE Project.Name LIKE '%` + search + `%'
	`

	if len(types) != 0 {
		whereClause += `
		AND Project.Type IN (`
		for i, t := range(types) {
			if i != 0 {
				whereClause += `, `
			}
			whereClause += `'` + t + `'`
		}
		whereClause += `)`
	}

	for _, l := range(languages) {
		whereClause += `
		AND Languages LIKE '%` + l + `%'
		`
	}

	for _, t := range(technologies) {
		whereClause += `
		AND Technologies LIKE '%` + t + `%'
		`
	}

	orderByClause := `
	ORDER BY Project.Date DESC, Project.ID DESC
	`

	if len(first) != 0 {
		firstList := strings.Split(first, " ")
		whereClause += `
		AND (Project.Date > '` + firstList[0] + `' OR (Project.Date = '` + firstList[0] + `' AND Project.ID > ` + firstList[1] + `))
		`
		orderByClause = `
		ORDER BY Project.Date ASC, Project.ID ASC
		`
	}

	if len(last) != 0 {
		lastList := strings.Split(last, " ")
		whereClause += `
		AND (Project.Date < '` + lastList[0] + `' OR (Project.Date = '` + lastList[0] + `' AND Project.ID < ` + lastList[1] + `))
		`
	}

	limitClause := `
	LIMIT 10
	`
	if len(first) == 0 && len(last) == 0 {
		limitClause += `
		OFFSET ` + strconv.Itoa((page - 1) * 10)
	}

	query := selectClause + whereClause + orderByClause + limitClause
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

	if len(first) != 0 {
		for i, j := 0, len(projects)-1; i < j; i, j = i+1, j-1 {
			projects[i], projects[j] = projects[j], projects[i]
		}
	}

	return projects, nil
}

func (d *SQLiteDatabase) GetProjectCount(search string, types []string, languages []string, technologies []string) int {
	query := `
	SELECT COUNT(*) 
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
	query += `
	WHERE Project.Name LIKE '%` + search + `%'
	`

	if len(types) != 0 {
		query += `
		AND Project.Type IN (`
		for i, t := range(types) {
			if i != 0 {
				query += `, `
			}
			query += `'` + t + `'`
		}
		query += `)`
	}

	for _, l := range(languages) {
		query += `
		AND Languages LIKE '%` + l + `%'
		`
	}

	for _, t := range(technologies) {
		query += `
		AND Technologies LIKE '%` + t + `%'
		`
	}

	query += `;`

	r := d.db.QueryRow(query)
	if r.Err() == sql.ErrNoRows {
		return 0
	}
	var result int
	if err := r.Scan(&result); err != nil {
		return 0
	}
	return result
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
