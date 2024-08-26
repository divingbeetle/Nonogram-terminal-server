package storage

import (
	"database/sql"
	"log"

	"github.com/divingbeetle/Nonogram-terminal-server/types"
	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB(username, password, host, port, dbname string) error {
	var err error

	db, err = sql.Open("postgres", "postgres://"+username+":"+password+"@"+host+":"+port+"/"+dbname+"?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	conns := 64
	db.SetMaxOpenConns(conns)
	db.SetMaxIdleConns(conns)

	return nil
}

func FetchPuzzle() (types.Puzzle, error) {
	row := db.QueryRow(`
        WITH selected_puzzle AS (
        SELECT id, row_size, col_size, clues
        FROM puzzles TABLESAMPLE BERNOULLI(1)
        ORDER BY RANDOM()
        LIMIT 1
    )

    SELECT 
        p.id, p.row_size, p.col_size, p.clues,
        COALESCE(pt.title, 'Untitled') AS title,
        COALESCE(a.name, 'anonymous') AS author
    FROM selected_puzzle p 
    LEFT JOIN 
        puzzle_titles pt ON p.id = pt.puzzle_id AND pt.language = 'ko' 
    LEFT JOIN
        puzzle_metadata pm ON p.id = pm.puzzle_id
    LEFT JOIN
        authors a ON pm.author_id = a.id
    `)

	var puzzle types.Puzzle
	err := row.Scan(&puzzle.ID, &puzzle.RowSize, &puzzle.ColSize, &puzzle.Clues,
		&puzzle.Title, &puzzle.Author)

	if err != nil {
		log.Println(err)
	}

	return puzzle, err
}
