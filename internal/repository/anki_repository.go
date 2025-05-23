package repositories

import (
	"anki-extractor/internal/models"
	"anki-extractor/internal/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

type AnkiRepository interface {
	OpenDB(path string) (*sql.DB, error)
	ExtractNotes(db *sql.DB) ([]models.Note, error)
	ExtractCards(db *sql.DB) ([]models.Card, error)
	ExtractDecks(db *sql.DB) ([]models.Deck, error)
}

type ankiRepository struct{}

func NewAnkiRepository() AnkiRepository {
	return &ankiRepository{}
}

func (r *ankiRepository) OpenDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("erro abrindo DB Anki: %w", err)
	}
	return db, nil
}

func (r *ankiRepository) ExtractNotes(db *sql.DB) ([]models.Note, error) {
	rows, err := db.Query("SELECT id, guid, mid, flds, tags FROM notes")
	if err != nil {
		return nil, fmt.Errorf("erro ao consultar notes: %w", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Printf("erro ao fechar rows: %v\n", err)
		}
	}(rows)

	var notes []models.Note
	for rows.Next() {
		var n models.Note
		var flds, tags string
		err = rows.Scan(&n.ID, &n.GUID, &n.Model, &flds, &tags)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear note: %w", err)
		}
		n.Fields = utils.SplitByRune(flds, '\x1f')
		n.Tags = utils.ParseTags(tags)
		notes = append(notes, n)
	}
	return notes, nil
}

func (r *ankiRepository) ExtractCards(db *sql.DB) ([]models.Card, error) {
	rows, err := db.Query("SELECT id, nid, did, ivl, due, ord FROM cards")
	if err != nil {
		return nil, fmt.Errorf("erro ao consultar cards: %w", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Printf("erro ao fechar rows: %v\n", err)
		}
	}(rows)

	var cards []models.Card
	for rows.Next() {
		var c models.Card
		err = rows.Scan(&c.ID, &c.NoteID, &c.DeckID, &c.Interval, &c.Due, &c.TemplateOrd)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear card: %w", err)
		}
		cards = append(cards, c)
	}
	return cards, nil
}

func (r *ankiRepository) ExtractDecks(db *sql.DB) ([]models.Deck, error) {
	var decksJSON string
	err := db.QueryRow("SELECT decks FROM col").Scan(&decksJSON)
	if err != nil {
		return nil, fmt.Errorf("erro ao consultar decks na tabela 'col': %w", err)
	}

	var decksMap map[string]struct {
		Name string `json:"name"`
	}
	err = json.Unmarshal([]byte(decksJSON), &decksMap)
	if err != nil {
		return nil, fmt.Errorf("erro ao fazer unmarshal dos decks JSON: %w", err)
	}

	var decks []models.Deck
	for idStr, deckData := range decksMap {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {

			continue
		}
		decks = append(decks, models.Deck{
			ID:   id,
			Name: deckData.Name,
		})
	}
	return decks, nil
}
