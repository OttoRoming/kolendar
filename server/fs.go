package server

import (
	"os"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
)

const (
	dataDir                 = "data"
	librarysDir             = "data/libraries"
	perms       os.FileMode = 0755
)

type Book struct {
	Authors []string
	Title   string
}

func (s *Server) setupFS() error {
	_, err := os.Stat(dataDir)
	if os.IsNotExist(err) {
		err = os.Mkdir(dataDir, perms)
		if err != nil {
			return err
		}
		s.slog.Info("Created data directory")
	}

	_, err = os.Stat(librarysDir)
	if os.IsNotExist(err) {
		err = os.Mkdir(librarysDir, perms)
		if err != nil {
			return err
		}
		s.slog.Info("Created libraries directory")
	}

	return nil
}

func (s *Server) createLibraryFS(libraryID pgtype.UUID) error {
	return os.Mkdir(librarysDir+libraryID.String(), perms)
}

func (s *Server) deleteLibraryFS(libraryID pgtype.UUID) error {
	return os.RemoveAll(librarysDir + libraryID.String())
}

func (s *Server) listLibraryBooks(libraryID pgtype.UUID) ([]Book, error) {
	books := []Book{}

	entries, err := os.ReadDir(librarysDir + libraryID.String())
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		name := entry.Name()

		dashSplitResult := strings.Split(name, " - ")
		if len(dashSplitResult) != 2 {
			s.slog.Warn("Invalid book directory name, skipping", "name", name)
			continue
		}

		authors := strings.Split(dashSplitResult[0], ", ")
		title := dashSplitResult[1]

		books = append(books, Book{
			Authors: authors,
			Title:   title,
		})
	}

	return books, nil
}
