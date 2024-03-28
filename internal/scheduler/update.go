package scheduler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"sync"
	"training/proj/internal/api/models"
	"training/proj/internal/db/repositories"

	"go.uber.org/zap"
)

type Scheduler struct {
	ItemRepository         *repositories.ItemRepository
	CategoryRepository     *repositories.CategoryRepository
	CategoryItemRepository *repositories.CategoryItemRepository
	Logger                 *zap.SugaredLogger
	Wg                     *sync.WaitGroup
}

func NewScheduler(r *repositories.Repositories, l *zap.SugaredLogger, wg *sync.WaitGroup) *Scheduler {
	return &Scheduler{
		ItemRepository:         r.ItemRepository,
		CategoryRepository:     r.CategoryRepository,
		CategoryItemRepository: r.CategoryItemRepository,
		Logger:                 l,
		Wg:                     wg,
	}
}

type ExternalItem struct {
	Name     string `json:"name"`
	Category string `json:"category"`
}

func (s *Scheduler) ExternalDbFill() {
	s.Wg.Add(1)
	defer s.Wg.Done()
	s.Logger.Info("Start filling db")
	externalItems := s.parse()
	s.fillTables(externalItems)
	s.Logger.Info("Finish filling db")
}

func (s *Scheduler) parse() []ExternalItem {
	requestURL := "https://emojihub.yurace.pro/api/all"
	res, httpErr := http.Get(requestURL)
	if httpErr != nil {
		s.Logger.Fatalf(httpErr.Error())
	}

	resBody, ioErr := io.ReadAll(res.Body)
	if ioErr != nil {
		s.Logger.Fatalf(ioErr.Error())
	}
	items := []ExternalItem{}
	unmarshallErr := json.Unmarshal(resBody, &items)
	if unmarshallErr != nil {
		s.Logger.Fatalf(unmarshallErr.Error())
	}
	return items
}

func (s *Scheduler) fillTables(items []ExternalItem) {

	for _, v := range items {
		dbItem, itemErr := s.createItemIfAbsent(v)
		dbCategory, catErr := s.createCategoryIfAbsent(v)

		if itemErr != nil || catErr != nil {
			s.CategoryItemRepository.Create(dbCategory.CategoryID, dbItem.ItemID)
		}
	}

}

func (s *Scheduler) createCategoryIfAbsent(item ExternalItem) (models.Category, error) {
	dbCategory, getCatErr := s.CategoryRepository.GetByName(item.Category)

	if getCatErr == sql.ErrNoRows {
		newCategory := models.Category{
			Category: item.Category,
		}
		newDbCategory, createErr := s.CategoryRepository.Create(&newCategory)

		if createErr != nil {
			s.Logger.Fatalf(createErr.Error())
		}

		return newDbCategory, nil
	} else if getCatErr != nil {
		s.Logger.Fatalf(getCatErr.Error())
	}

	return dbCategory, fmt.Errorf("category already exists")
}

func (s *Scheduler) createItemIfAbsent(item ExternalItem) (models.Item, error) {
	dbItem, getItemErr := s.ItemRepository.GetByName(item.Name)

	if getItemErr == sql.ErrNoRows {
		newItem := models.Item{
			Item:  item.Name,
			Price: rand.Int64N(99900) + 1000,
		}
		newDbItem, createErr := s.ItemRepository.Create(&newItem)

		if createErr != nil {
			s.Logger.Fatalf(createErr.Error())
		}

		return newDbItem, nil
	} else if getItemErr != nil {
		s.Logger.Fatalf(getItemErr.Error())
	}

	return dbItem, fmt.Errorf("item already exists")
}
