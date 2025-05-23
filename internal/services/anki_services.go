package services

import (
	"anki-extractor/internal/models"
	repositories "anki-extractor/internal/repository"
	"database/sql"
	"fmt"
)

type AnkiService interface {
	GetFlashcardsWithDeckNames(ankiDBPath string) ([]models.Flashcard, error)
}

type ankiService struct {
	repo repositories.AnkiRepository
}

func NewAnkiService(repo repositories.AnkiRepository) AnkiService {
	return &ankiService{repo: repo}
}

func (s *ankiService) GetFlashcardsWithDeckNames(ankiDBPath string) ([]models.Flashcard, error) {
	db, err := s.repo.OpenDB(ankiDBPath)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir o DB Anki: %w", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			fmt.Printf("erro ao fechar DB Anki: %v\n", err)
		} else {
			fmt.Println("DB Anki fechado com sucesso.")
		}
	}(db)

	notes, err := s.repo.ExtractNotes(db)
	if err != nil {
		return nil, fmt.Errorf("erro extraindo notes: %w", err)
	}

	cards, err := s.repo.ExtractCards(db)
	if err != nil {
		return nil, fmt.Errorf("erro extraindo cards: %w", err)
	}

	decks, err := s.repo.ExtractDecks(db)
	if err != nil {
		return nil, fmt.Errorf("erro extraindo decks: %w", err)
	}

	noteMap := make(map[int64]models.Note)
	for _, n := range notes {
		noteMap[n.ID] = n
	}

	deckMap := make(map[int64]models.Deck)
	for _, d := range decks {
		deckMap[d.ID] = d
	}

	var flashcards []models.Flashcard
	for _, c := range cards {
		if note, ok := noteMap[c.NoteID]; ok {
			deckName := "Unknown Deck"
			if deck, found := deckMap[c.DeckID]; found {
				deckName = deck.Name
			}

			fc := models.Flashcard{
				Note:     note,
				Card:     c,
				DeckName: deckName,
			}
			flashcards = append(flashcards, fc)
		}
	}

	return flashcards, nil
}
