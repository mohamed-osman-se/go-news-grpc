package memstore

import (
	"net/url"
	"sync"
	"time"

	"github.com/google/uuid"
)

type News struct {
	ID uuid.UUID

	Author string

	Title string

	Summary string

	Content string

	Source *url.URL

	Tags []string

	CreatedAt time.Time

	UpdatedAt time.Time

	DeletedAt time.Time
}


type Store struct {
	lock sync.RWMutex
	news []*News
}


func New() *Store {
	return &Store{
		lock: sync.RWMutex{},
		news: make([]*News, 0),
	}
}


func (s *Store) Create(news *News) *News {
	createdNews := &News{
		ID:        uuid.New(),
		Author:    news.Author,
		Title:     news.Title,
		Summary:   news.Summary,
		Content:   news.Content,
		Source:    news.Source,
		Tags:      news.Tags,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	s.lock.Lock()
	defer s.lock.Unlock()
	s.news = append(s.news, createdNews)
	return createdNews
}

// Get news by it's id.
func (s *Store) Get(id uuid.UUID) *News {
	s.lock.RLock()
	defer s.lock.RUnlock()
	for _, news := range s.news {
		if news.ID == id && news.DeletedAt.IsZero() {
			return news
		}
	}
	return nil
}

// GetAll news.
func (s *Store) GetAll() []*News {
	s.lock.RLock()
	defer s.lock.RUnlock()

	result := make([]*News, 0)

	for _, news := range s.news {
		if news.DeletedAt.IsZero() {
			result = append(result, news)
		}
	}

	return result
}

// Update news.
func (s *Store) Update(updatedNews *News) {
	s.lock.Lock()
	defer s.lock.Unlock()
	for idx, news := range s.news {
		if updatedNews.ID == news.ID && news.DeletedAt.IsZero() {
			s.news[idx] = updatedNews
			return
		}
	}
}

// Delete news from store.
func (s *Store) Delete(id uuid.UUID) {
	s.lock.Lock()
	defer s.lock.Unlock()
	for idx, news := range s.news {
		if id == news.ID {
			s.news[idx].DeletedAt = time.Now().UTC()
			return
		}
	}
}