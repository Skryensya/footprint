package tracking

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/Skryensya/footprint/internal/git"
	repodomain "github.com/Skryensya/footprint/internal/repo"
	"github.com/Skryensya/footprint/internal/store"
	"github.com/Skryensya/footprint/internal/ui"
)

type Deps struct {
	// git
	GitIsAvailable func() bool
	RepoRoot       func(string) (string, error)
	OriginURL      func(string) (string, error)
	HeadCommit     func() (string, error)
	CurrentBranch  func() (string, error)
	CommitMessage  func() (string, error)
	CommitAuthor   func() (string, error)

	// repo
	DeriveID    func(string, string) (repodomain.RepoID, error)
	Track       func(repodomain.RepoID) (bool, error)
	Untrack     func(repodomain.RepoID) (bool, error)
	IsTracked   func(repodomain.RepoID) (bool, error)
	ListTracked func() ([]repodomain.RepoID, error)

	// store
	DBPath      func() string
	OpenDB      func(string) (*sql.DB, error)
	InitDB      func(*sql.DB) error
	InsertEvent func(*sql.DB, store.RepoEvent) error
	ListEvents  func(*sql.DB, store.EventFilter) ([]store.RepoEvent, error)

	// io
	Printf  func(string, ...any) (int, error)
	Println func(...any) (int, error)
	Pager   func(string)

	// misc
	Now    func() time.Time
	Getenv func(string) string
}

func DefaultDeps() Deps {
	return Deps{
		GitIsAvailable: git.IsAvailable,
		RepoRoot:       git.RepoRoot,
		OriginURL:      git.OriginURL,
		HeadCommit:     git.HeadCommit,
		CurrentBranch:  git.CurrentBranch,
		CommitMessage:  git.CommitMessage,
		CommitAuthor:   git.CommitAuthor,

		DeriveID:    repodomain.DeriveID,
		Track:       repodomain.Track,
		Untrack:     repodomain.Untrack,
		IsTracked:   repodomain.IsTracked,
		ListTracked: repodomain.ListTracked,

		DBPath:      store.DBPath,
		OpenDB:      store.Open,
		InitDB:      store.Init,
		InsertEvent: store.InsertEvent,
		ListEvents:  store.ListEvents,

		Printf:  fmt.Printf,
		Println: fmt.Println,
		Pager:   ui.Pager,

		Now:    time.Now,
		Getenv: os.Getenv,
	}
}
