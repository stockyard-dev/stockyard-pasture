package store

import (
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"
	"os"
	"path/filepath"
	"time"
)

type DB struct{ db *sql.DB }
type Post struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Body        string `json:"body,omitempty"`
	Platform    string `json:"platform,omitempty"`
	Status      string `json:"status"`
	ScheduledAt string `json:"scheduled_at,omitempty"`
	PublishedAt string `json:"published_at,omitempty"`
	Tags        string `json:"tags,omitempty"`
	CreatedAt   string `json:"created_at"`
}

func Open(d string) (*DB, error) {
	if err := os.MkdirAll(d, 0755); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", filepath.Join(d, "pasture.db")+"?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		return nil, err
	}
	db.Exec(`CREATE TABLE IF NOT EXISTS posts(id TEXT PRIMARY KEY,title TEXT NOT NULL,body TEXT DEFAULT '',platform TEXT DEFAULT '',status TEXT DEFAULT 'draft',scheduled_at TEXT DEFAULT '',published_at TEXT DEFAULT '',tags TEXT DEFAULT '',created_at TEXT DEFAULT(datetime('now')))`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_posts_status ON posts(status)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_posts_sched ON posts(scheduled_at)`)
	db.Exec(`CREATE TABLE IF NOT EXISTS extras(resource TEXT NOT NULL,record_id TEXT NOT NULL,data TEXT NOT NULL DEFAULT '{}',PRIMARY KEY(resource, record_id))`)
	return &DB{db: db}, nil
}
func (d *DB) Close() error { return d.db.Close() }
func genID() string        { return fmt.Sprintf("%d", time.Now().UnixNano()) }
func now() string          { return time.Now().UTC().Format(time.RFC3339) }
func (d *DB) Create(p *Post) error {
	p.ID = genID()
	p.CreatedAt = now()
	if p.Status == "" {
		p.Status = "draft"
	}
	_, err := d.db.Exec(`INSERT INTO posts VALUES(?,?,?,?,?,?,?,?,?)`, p.ID, p.Title, p.Body, p.Platform, p.Status, p.ScheduledAt, p.PublishedAt, p.Tags, p.CreatedAt)
	return err
}
func (d *DB) Get(id string) *Post {
	var p Post
	if d.db.QueryRow(`SELECT * FROM posts WHERE id=?`, id).Scan(&p.ID, &p.Title, &p.Body, &p.Platform, &p.Status, &p.ScheduledAt, &p.PublishedAt, &p.Tags, &p.CreatedAt) != nil {
		return nil
	}
	return &p
}
func (d *DB) List(status string) []Post {
	q := `SELECT * FROM posts`
	args := []any{}
	if status != "" && status != "all" {
		q += ` WHERE status=?`
		args = append(args, status)
	}
	q += ` ORDER BY CASE status WHEN 'scheduled' THEN 0 WHEN 'draft' THEN 1 WHEN 'published' THEN 2 END, scheduled_at ASC, created_at DESC`
	rows, _ := d.db.Query(q, args...)
	if rows == nil {
		return nil
	}
	defer rows.Close()
	var o []Post
	for rows.Next() {
		var p Post
		rows.Scan(&p.ID, &p.Title, &p.Body, &p.Platform, &p.Status, &p.ScheduledAt, &p.PublishedAt, &p.Tags, &p.CreatedAt)
		o = append(o, p)
	}
	return o
}
func (d *DB) Update(id string, p *Post) error {
	_, err := d.db.Exec(`UPDATE posts SET title=?,body=?,platform=?,status=?,scheduled_at=?,tags=? WHERE id=?`, p.Title, p.Body, p.Platform, p.Status, p.ScheduledAt, p.Tags, id)
	return err
}
func (d *DB) Publish(id string) error {
	_, err := d.db.Exec(`UPDATE posts SET status='published',published_at=? WHERE id=?`, now(), id)
	return err
}
func (d *DB) Delete(id string) error {
	_, err := d.db.Exec(`DELETE FROM posts WHERE id=?`, id)
	return err
}

type Stats struct {
	Total     int `json:"total"`
	Draft     int `json:"draft"`
	Scheduled int `json:"scheduled"`
	Published int `json:"published"`
}

func (d *DB) Stats() Stats {
	var s Stats
	d.db.QueryRow(`SELECT COUNT(*) FROM posts`).Scan(&s.Total)
	d.db.QueryRow(`SELECT COUNT(*) FROM posts WHERE status='draft'`).Scan(&s.Draft)
	d.db.QueryRow(`SELECT COUNT(*) FROM posts WHERE status='scheduled'`).Scan(&s.Scheduled)
	d.db.QueryRow(`SELECT COUNT(*) FROM posts WHERE status='published'`).Scan(&s.Published)
	return s
}

// ─── Extras: generic key-value storage for personalization custom fields ───

func (d *DB) GetExtras(resource, recordID string) string {
	var data string
	err := d.db.QueryRow(
		`SELECT data FROM extras WHERE resource=? AND record_id=?`,
		resource, recordID,
	).Scan(&data)
	if err != nil || data == "" {
		return "{}"
	}
	return data
}

func (d *DB) SetExtras(resource, recordID, data string) error {
	if data == "" {
		data = "{}"
	}
	_, err := d.db.Exec(
		`INSERT INTO extras(resource, record_id, data) VALUES(?, ?, ?)
		 ON CONFLICT(resource, record_id) DO UPDATE SET data=excluded.data`,
		resource, recordID, data,
	)
	return err
}

func (d *DB) DeleteExtras(resource, recordID string) error {
	_, err := d.db.Exec(
		`DELETE FROM extras WHERE resource=? AND record_id=?`,
		resource, recordID,
	)
	return err
}

func (d *DB) AllExtras(resource string) map[string]string {
	out := make(map[string]string)
	rows, _ := d.db.Query(
		`SELECT record_id, data FROM extras WHERE resource=?`,
		resource,
	)
	if rows == nil {
		return out
	}
	defer rows.Close()
	for rows.Next() {
		var id, data string
		rows.Scan(&id, &data)
		out[id] = data
	}
	return out
}
