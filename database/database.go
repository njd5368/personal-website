package database

import (
	"database/sql"
	"encoding/hex"
	"errors"
	"strconv"
	"strings"

	"golang.org/x/crypto/sha3"
)

type Type int64

const (
	Project = iota
	Blog
)

func (t Type) String() string {
	return [...]string{"Project", "Blog"}[t]
}

func StringToType(s string) Type {
	s = strings.ToLower(s)
	return map[string]Type{
		"project": Project,
		"blog": Blog,
	}[s]
}

type Database interface {
	StartDatabase() error
	CreatePost(p Post) (*Post, error)
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
		CREATE TABLE IF NOT EXISTS User(
			Username TEXT PRIMARY KEY UNIQUE,
			Password TEXT NOT NULL
		);
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
			Category TEXT NOT NULL,
			Image INTEGER NOT NULL,
			File BLOB NOT NULL,
			FOREIGN KEY(Image) REFERENCES Image(ID)
		);
		`,
		`
		CREATE TABLE IF NOT EXISTS Blog(
			ID INTEGER PRIMARY KEY AUTOINCREMENT,
			Name TEXT NOT NULL UNIQUE,
			Description TEXT NOT NULL,
			Date TEXT NOT NULL,
			Category TEXT NOT NULL,
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
		CREATE TABLE IF NOT EXISTS BlogLanguage(
			Blog INTEGER NOT NULL,
			Language INTEGER NOT NULL,
			PRIMARY KEY(Blog, Language),
			FOREIGN KEY(Blog) REFERENCES Blog(ID),
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
		`
		CREATE TABLE IF NOT EXISTS BlogTechnology(
			Blog INTEGER NOT NULL,
			Technology INTEGER NOT NULL,
			PRIMARY KEY(Blog, Technology),
			FOREIGN KEY(Blog) REFERENCES Blog(ID),
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

func (d *SQLiteDatabase) CreateUser(username string, password string) error {
	hash := sha3.New256()
	hash.Write([]byte(password))
	password = hex.EncodeToString(hash.Sum(nil))

	_, err := d.db.Exec("INSERT INTO User(Username, Password) VALUES(?, ?)", username, password)

 	return err
}

func (d *SQLiteDatabase) UserExists() bool {
	r := d.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM User);`)
	if r.Err() == sql.ErrNoRows {
		return false
	}
	var exists int
	err := r.Scan(&exists)
	if err != nil {
		return false
	}

	return exists == 1
}

func (d *SQLiteDatabase) IsUser(username string, password string) bool {
	r := d.db.QueryRow(`SELECT Password FROM User WHERE Username = '` + username + `'`)
	if r.Err() == sql.ErrNoRows {
		return false
	}

	var p string
	if err := r.Scan(&p); err != nil {
		return false
	}

	hash := sha3.New256()
	hash.Write([]byte(password))
	password = hex.EncodeToString(hash.Sum(nil))
	return password == p
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

func (d *SQLiteDatabase) CreatePost(p Post, i []byte) (*Post, error) {

	r, err := d.db.Exec("INSERT INTO Image(Image) VALUES(?)", i)
	if err != nil {
		return nil, err
	}

	id, err := r.LastInsertId()
	if err != nil {
		return nil, err
	}

	p.Image = id

	r, err = d.db.Exec(`INSERT INTO ` + p.Type.String() + `(Name, Description, Date, Category, Image, File) VALUES(?,?,?,?,?,?)`,
		p.Name, p.Description, p.Date, p.Category, p.Image, p.File)
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
		statement := `INSERT INTO ` + p.Type.String() + `Language(` + p.Type.String() + `, Language) VALUES`
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
		statement := `INSERT INTO ` + p.Type.String() + `Technology(` + p.Type.String() + `, Technology) VALUES`
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

func (d *SQLiteDatabase) GetProjectPosts(page int, search string, categories []string, languages []string, technologies []string, first string, last string) ([]Post, error) {
	return d.GetPosts(Project, page, search, categories, languages, technologies, first, last)
}

func (d *SQLiteDatabase) GetBlogPosts(page int, search string, categories []string, languages []string, technologies []string, first string, last string) ([]Post, error) {
	return d.GetPosts(Blog, page, search, categories, languages, technologies, first, last)
}

// GetPosts returns up to 10 posts that meet the required parameters.
//
// page (int): the page number for pagination. If less than, 1, defaults to 1. If over the max pages, an empty Post struct is returned.
// 
func (d *SQLiteDatabase) GetPosts(t Type, page int, search string, categories []string, languages []string, technologies []string, first string, last string) ([]Post, error) {
	if page < 1 {
		page = 1
	}
	if len(first) != 0 && len(last) != 0 {
		first = ""
		last = ""
	}
	
	table := t.String()
	selectClause := `
	SELECT ` + table + `.ID, ` + table + `.Name, ` + table + `.Description, ` + table + `.Date, ` + table + `.Category, ` + table + `.Image, Languages, Technologies
	FROM ` + table + `
	LEFT JOIN (
		SELECT ` + table + `Language.` + table + ` AS ` + table + `LanguageID, GROUP_CONCAT(Language.Name) AS Languages
		FROM ` + table + `Language
		LEFT JOIN Language
		ON ` + table + `Language.Language=Language.ID
		GROUP BY ` + table + `LanguageID
	)
	ON ` + table + `.ID=` + table + `LanguageID
	LEFT JOIN (
		SELECT ` + table + `Technology.` + table + ` AS ` + table + `TechnologyID, GROUP_CONCAT(Technology.Name) AS Technologies
		FROM ` + table + `Technology
		LEFT JOIN Technology
		ON ` + table + `Technology.Technology=Technology.ID
		GROUP BY ` + table + `TechnologyID
	)
	ON ` + table + `.ID=` + table + `TechnologyID
	`
	whereClause := `
	WHERE ` + table + `.Name LIKE '%` + search + `%'
	`

	if len(categories) != 0 {
		whereClause += `
		AND ` + table + `.Category IN (`
		for i, t := range(categories) {
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
	ORDER BY ` + table + `.Date DESC, ` + table + `.ID DESC
	`

	if len(first) != 0 {
		firstList := strings.Split(first, " ")
		whereClause += `
		AND (` + table + `.Date > '` + firstList[0] + `' OR (` + table + `.Date = '` + firstList[0] + `' AND ` + table + `.ID > ` + firstList[1] + `))
		`
		orderByClause = `
		ORDER BY ` + table + `.Date ASC, ` + table + `.ID ASC
		`
	}

	if len(last) != 0 {
		lastList := strings.Split(last, " ")
		whereClause += `
		AND (` + table + `.Date < '` + lastList[0] + `' OR (` + table + `.Date = '` + lastList[0] + `' AND ` + table + `.ID < ` + lastList[1] + `))
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

	var posts []Post
	for r.Next() {
		var p Post
		var languages sql.NullString
		var techs sql.NullString
		if err := r.Scan(&p.ID, &p.Name, &p.Description, &p.Date, &p.Category, &p.Image, &languages, &techs); err != nil {
			return nil, err
		}
		p.Languages = strings.Split(languages.String, ",")
		p.Technologies = strings.Split(techs.String, ",")
		posts = append(posts, p)
	}

	if len(first) != 0 {
		for i, j := 0, len(posts)-1; i < j; i, j = i+1, j-1 {
			posts[i], posts[j] = posts[j], posts[i]
		}
	}

	return posts, nil
}

func (d *SQLiteDatabase) GetProjectPostCount(search string, categories []string, languages []string, technologies []string) int {
	return d.GetPostCount(Project, search, categories, languages, technologies)
}

func (d *SQLiteDatabase) GetBlogPostCount(search string, categories []string, languages []string, technologies []string) int {
	return d.GetPostCount(Blog, search, categories, languages, technologies)
}


func (d *SQLiteDatabase) GetPostCount(t Type, search string, categories []string, languages []string, technologies []string) int {
	table := t.String()

	query := `
	SELECT COUNT(*) 
	FROM ` + table + `
	LEFT JOIN (
		SELECT ` + table + `Language.` + table + ` AS ` + table + `LanguageID, GROUP_CONCAT(Language.Name) AS Languages
		FROM ` + table + `Language
		LEFT JOIN Language
		ON ` + table + `Language.Language=Language.ID
		GROUP BY ` + table + `LanguageID
	)
	ON ` + table + `.ID=` + table + `LanguageID
	LEFT JOIN (
		SELECT ` + table + `Technology.` + table + ` AS ` + table + `TechnologyID, GROUP_CONCAT(Technology.Name) AS Technologies
		FROM ` + table + `Technology
		LEFT JOIN Technology
		ON ` + table + `Technology.Technology=Technology.ID
		GROUP BY ` + table + `TechnologyID
	)
	ON ` + table + `.ID=` + table + `TechnologyID
	`
	query += `
	WHERE ` + table + `.Name LIKE '%` + search + `%'
	`

	if len(categories) != 0 {
		query += `
		AND ` + table + `.Category IN (`
		for i, t := range(categories) {
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

func (d *SQLiteDatabase) GetProjectPostByName(n string) (*Post, error) {
	return d.GetPostByName(Project, n)
}

func (d *SQLiteDatabase) GetBlogPostByName(n string) (*Post, error) {
	return d.GetPostByName(Blog, n)
}

func (d *SQLiteDatabase) GetPostByName(t Type, n string) (*Post, error) {
	table := t.String()

	query := `
	SELECT ` + table + `.ID, ` + table + `.Name, ` + table + `.Description, ` + table + `.Date, ` + table + `.Category, ` + table + `.Image, ` + table + `.File, Languages, Technologies
	FROM ` + table + `
	LEFT JOIN (
		SELECT ` + table + `Language.` + table + ` AS ` + table + `LanguageID, GROUP_CONCAT(Language.Name) AS Languages
		FROM ` + table + `Language
		LEFT JOIN Language
		ON ` + table + `Language.Language=Language.ID
		GROUP BY ` + table + `LanguageID
	)
	ON ` + table + `.ID=` + table + `LanguageID
	LEFT JOIN (
		SELECT ` + table + `Technology.` + table + ` AS ` + table + `TechnologyID, GROUP_CONCAT(Technology.Name) AS Technologies
		FROM ` + table + `Technology
		LEFT JOIN Technology
		ON ` + table + `Technology.Technology=Technology.ID
		GROUP BY ` + table + `TechnologyID
	)
	ON ` + table + `.ID=` + table + `TechnologyID
	WHERE ` + table + `.Name='` + n + `'
	`

	r := d.db.QueryRow(query)
	if r.Err() == sql.ErrNoRows {
		return nil, errors.New("no post with name '" + n + "'")
	}
	
	var p Post
	var languages sql.NullString
	var techs sql.NullString
	if err := r.Scan(&p.ID, &p.Name, &p.Description, &p.Date, &p.Category, &p.Image, &p.File, &languages, &techs); err != nil {
		return nil, err
	}
	p.Languages = strings.Split(languages.String, ",")
	p.Technologies = strings.Split(techs.String, ",")
	return &p, nil
}

func (d *SQLiteDatabase) GetLatestBlogPostName() (string, error) {
	query := `SELECT Name FROM Blog ORDER BY Blog.Date DESC, Blog.ID DESC LIMIT 1;`

	r := d.db.QueryRow(query)
	if r.Err() == sql.ErrNoRows {
		return "", errors.New("no posts")
	}

	var name string
	if err := r.Scan(&name); err != nil {
		return "", err
	}

	return name, nil
}

func (d *SQLiteDatabase) GetAllLanguages() ([]string, error) {
	query := `
	SELECT Name
	FROM Language
	ORDER BY Name ASC
	`
	r, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var languages []string
	for r.Next() {
		var language string
		if err := r.Scan(&language); err != nil {
			return nil, err
		}
		languages = append(languages, language)
	}
	return languages, nil
}

func (d *SQLiteDatabase) GetAllTechnologies() ([]string, error) {
	query := `
	SELECT Name
	FROM Technology
	ORDER BY Name ASC
	`
	r, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var technologies []string
	for r.Next() {
		var technology string
		if err := r.Scan(&technology); err != nil {
			return nil, err
		}
		technologies = append(technologies, technology)
	}
	return technologies, nil
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
