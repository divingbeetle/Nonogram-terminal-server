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

	return db.Ping()
}

func FetchPuzzles(offset, limit int) ([]types.Puzzle, error) {
	rows, err := db.Query(`
    WITH selected_puzzles AS (
    SELECT id, row_size, col_size, clues
    FROM puzzles
    OFFSET $1
    LIMIT $2
    )

    SELECT 
        p.id, p.row_size, p.col_size, p.clues,
        COALESCE(pt.title, 'Untitled') AS title,
        COALESCE(a.name, 'anonymous') AS author
    FROM selected_puzzles p 
    LEFT JOIN 
        puzzle_titles pt ON p.id = pt.puzzle_id AND pt.language = 'ko' 
    LEFT JOIN
        puzzle_metadata pm ON p.id = pm.puzzle_id
    LEFT JOIN
        authors a ON pm.author_id = a.id
    `, offset, limit)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var puzzles []types.Puzzle
	for rows.Next() {
		var puzzle types.Puzzle
		err := rows.Scan(&puzzle.ID, &puzzle.RowSize, &puzzle.ColSize, &puzzle.Clues,
			&puzzle.Title, &puzzle.Author)
		if err != nil {
			return nil, err
		}
		puzzles = append(puzzles, puzzle)
	}

	return puzzles, nil
}

func FetchRandomPuzzle() (types.Puzzle, error) {
	row := db.QueryRow(`
        WITH selected_puzzle AS (
        SELECT id, row_size, col_size, clues
        FROM puzzles TABLESAMPLE BERNOULLI(2)
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

func FetchPuzzle(id int) (types.Puzzle, error) {

	row := db.QueryRow(`
        SELECT 
            p.id, p.row_size, p.col_size, p.clues,
            COALESCE(pt.title, 'Untitled') AS title,
            COALESCE(a.name, 'anonymous') AS author
        FROM puzzles p 
        LEFT JOIN 
            puzzle_titles pt ON p.id = pt.puzzle_id AND pt.language = 'ko' 
        LEFT JOIN
            puzzle_metadata pm ON p.id = pm.puzzle_id
        LEFT JOIN
            authors a ON pm.author_id = a.id
        WHERE p.id = $1
    `, id)

	var puzzle types.Puzzle
	err := row.Scan(&puzzle.ID, &puzzle.RowSize, &puzzle.ColSize, &puzzle.Clues,
		&puzzle.Title, &puzzle.Author)

	if err != nil {
		log.Println(err)
	}

	return puzzle, err
}
