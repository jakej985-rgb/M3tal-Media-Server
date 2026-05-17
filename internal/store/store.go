package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

// Store wraps a SQLite database for M3TAL state persistence.
type Store struct {
	db *sql.DB
}

// RouteRecord is a stored route with its database ID.
type RouteRecord struct {
	ID          int64  `json:"id"`
	Service     string `json:"service"`
	Domain      string `json:"domain"`
	Port        int    `json:"port"`
	Entrypoints string `json:"entrypoints"`
	Stack       string `json:"stack,omitempty"`
	CreatedAt   string `json:"created_at"`
}

// StackRecord is a stored stack reference.
type StackRecord struct {
	Name         string `json:"name"`
	ComposePath  string `json:"compose_path"`
	LastDeployed string `json:"last_deployed,omitempty"`
	Status       string `json:"status"`
}

// MiddlewareRecord is a stored middleware definition.
type MiddlewareRecord struct {
	ID        int64             `json:"id"`
	Name      string            `json:"name"`
	Type      string            `json:"type"`
	Config    map[string]string `json:"config"`
	CreatedAt string            `json:"created_at"`
}

// Open creates or opens a SQLite database at the given path.
// It creates parent directories if they don't exist.
func Open(path string) (*Store, error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("cannot create database directory %s: %w", dir, err)
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("cannot open database %s: %w", path, err)
	}

	// Enable WAL mode for concurrent reads
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		db.Close()
		return nil, fmt.Errorf("cannot set WAL mode: %w", err)
	}

	// Enable foreign keys
	if _, err := db.Exec("PRAGMA foreign_keys=ON"); err != nil {
		db.Close()
		return nil, fmt.Errorf("cannot enable foreign keys: %w", err)
	}

	s := &Store{db: db}
	if err := s.Migrate(); err != nil {
		db.Close()
		return nil, err
	}

	return s, nil
}

// Close closes the database connection.
func (s *Store) Close() error {
	return s.db.Close()
}

// Migrate creates the required tables if they don't exist.
func (s *Store) Migrate() error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS routes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			service TEXT NOT NULL,
			domain TEXT NOT NULL UNIQUE,
			port INTEGER NOT NULL,
			entrypoints TEXT DEFAULT 'web',
			stack TEXT DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS stacks (
			name TEXT PRIMARY KEY,
			compose_path TEXT NOT NULL,
			last_deployed DATETIME,
			status TEXT DEFAULT 'unknown'
		)`,
		`CREATE TABLE IF NOT EXISTS middleware (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			type TEXT NOT NULL,
			config TEXT DEFAULT '{}',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	for _, m := range migrations {
		if _, err := s.db.Exec(m); err != nil {
			return fmt.Errorf("migration failed: %w\nSQL: %s", err, m)
		}
	}

	return nil
}

// --- Routes ---

// CreateRoute inserts a new route and returns its ID.
func (s *Store) CreateRoute(service, domain string, port int, entrypoints, stack string) (int64, error) {
	if entrypoints == "" {
		entrypoints = "web"
	}
	result, err := s.db.Exec(
		`INSERT INTO routes (service, domain, port, entrypoints, stack) VALUES (?, ?, ?, ?, ?)`,
		service, domain, port, entrypoints, stack,
	)
	if err != nil {
		return 0, fmt.Errorf("cannot create route: %w", err)
	}
	return result.LastInsertId()
}

// ListRoutes returns all stored routes.
func (s *Store) ListRoutes() ([]RouteRecord, error) {
	rows, err := s.db.Query(`SELECT id, service, domain, port, entrypoints, stack, created_at FROM routes ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var routes []RouteRecord
	for rows.Next() {
		var r RouteRecord
		if err := rows.Scan(&r.ID, &r.Service, &r.Domain, &r.Port, &r.Entrypoints, &r.Stack, &r.CreatedAt); err != nil {
			return nil, err
		}
		routes = append(routes, r)
	}
	return routes, rows.Err()
}

// GetRouteByDomain looks up a route by its domain.
func (s *Store) GetRouteByDomain(domain string) (*RouteRecord, error) {
	var r RouteRecord
	err := s.db.QueryRow(
		`SELECT id, service, domain, port, entrypoints, stack, created_at FROM routes WHERE domain = ?`,
		domain,
	).Scan(&r.ID, &r.Service, &r.Domain, &r.Port, &r.Entrypoints, &r.Stack, &r.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &r, nil
}

// DeleteRoute removes a route by ID.
func (s *Store) DeleteRoute(id int64) error {
	result, err := s.db.Exec(`DELETE FROM routes WHERE id = ?`, id)
	if err != nil {
		return err
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return fmt.Errorf("route %d not found", id)
	}
	return nil
}

// --- Stacks ---

// UpsertStack creates or updates a stack record.
func (s *Store) UpsertStack(name, composePath string) error {
	_, err := s.db.Exec(
		`INSERT INTO stacks (name, compose_path, status) VALUES (?, ?, 'discovered')
		 ON CONFLICT(name) DO UPDATE SET compose_path = excluded.compose_path`,
		name, composePath,
	)
	return err
}

// UpdateStackStatus updates the deployment status and timestamp.
func (s *Store) UpdateStackStatus(name, status string) error {
	_, err := s.db.Exec(
		`UPDATE stacks SET status = ?, last_deployed = ? WHERE name = ?`,
		status, time.Now().UTC().Format(time.RFC3339), name,
	)
	return err
}

// ListStacks returns all stored stacks.
func (s *Store) ListStacks() ([]StackRecord, error) {
	rows, err := s.db.Query(`SELECT name, compose_path, COALESCE(last_deployed, ''), status FROM stacks ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stacks []StackRecord
	for rows.Next() {
		var st StackRecord
		if err := rows.Scan(&st.Name, &st.ComposePath, &st.LastDeployed, &st.Status); err != nil {
			return nil, err
		}
		stacks = append(stacks, st)
	}
	return stacks, rows.Err()
}

// --- Middleware ---

// CreateMiddleware inserts a new middleware definition.
func (s *Store) CreateMiddleware(name, mwType string, config map[string]string) (int64, error) {
	configJSON, err := json.Marshal(config)
	if err != nil {
		return 0, err
	}
	result, err := s.db.Exec(
		`INSERT INTO middleware (name, type, config) VALUES (?, ?, ?)`,
		name, mwType, string(configJSON),
	)
	if err != nil {
		return 0, fmt.Errorf("cannot create middleware: %w", err)
	}
	return result.LastInsertId()
}

// ListMiddleware returns all stored middleware definitions.
func (s *Store) ListMiddleware() ([]MiddlewareRecord, error) {
	rows, err := s.db.Query(`SELECT id, name, type, config, created_at FROM middleware ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mws []MiddlewareRecord
	for rows.Next() {
		var m MiddlewareRecord
		var configJSON string
		if err := rows.Scan(&m.ID, &m.Name, &m.Type, &configJSON, &m.CreatedAt); err != nil {
			return nil, err
		}
		m.Config = make(map[string]string)
		json.Unmarshal([]byte(configJSON), &m.Config)
		mws = append(mws, m)
	}
	return mws, rows.Err()
}

// DeleteMiddleware removes a middleware by ID.
func (s *Store) DeleteMiddleware(id int64) error {
	result, err := s.db.Exec(`DELETE FROM middleware WHERE id = ?`, id)
	if err != nil {
		return err
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return fmt.Errorf("middleware %d not found", id)
	}
	return nil
}

// GetStatePath returns the default database path, respecting overrides.
func GetStatePath() string {
	if p := os.Getenv("M3TAL_STATE_DB"); p != "" {
		return p
	}
	// System install path
	sysPath := "/var/lib/m3tal/state.db"
	if dir := filepath.Dir(sysPath); dirExists(dir) {
		return sysPath
	}
	// Fallback to local
	return "./data/state.db"
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
