package tests

import (
	"Templatest/internal/config"
	"Templatest/internal/repo"
	"context"
	"github.com/stretchr/testify/assert"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestRepository(t *testing.T) {
	// Создаем тестовый репозиторий
	cfg := struct {
		Capacity    int
		MaxItemSize int
	}{
		Capacity:    3,
		MaxItemSize: 1024,
	}

	r, err := repo.NewRepository(config.InMemory(cfg))
	assert.NoError(t, err)

	t.Run("Post - успешное создание", func(t *testing.T) {
		obj := repo.DataObject{
			Title: "Test Task",
			Data:  "Test Data",
		}

		id, err := r.Post(context.Background(), obj)
		assert.NoError(t, err)
		assert.Equal(t, 1, id)
	})

	t.Run("Post - переполнение хранилища", func(t *testing.T) {
		// Заполняем хранилище
		for i := 0; i < 2; i++ {
			_, err := r.Post(context.Background(), repo.DataObject{
				Title: "Task " + strconv.Itoa(i+2),
			})
			assert.NoError(t, err)
		}

		// Попытка добавить четвертый элемент (лимит = 3)
		_, err := r.Post(context.Background(), repo.DataObject{Title: "Overflow"})
		assert.Error(t, err)
		assert.EqualError(t, err, "unavailable for recording, memory is full")
	})

	t.Run("Get - существующая задача", func(t *testing.T) {
		obj, err := r.Get(context.Background(), "1")
		assert.NoError(t, err)
		assert.Equal(t, "Test Task", obj.Title)
	})

	t.Run("Get - неверный ID (не число)", func(t *testing.T) {
		_, err := r.Get(context.Background(), "abc")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse id")
	})

	t.Run("Get - несуществующий ID", func(t *testing.T) {
		_, err := r.Get(context.Background(), "999")
		assert.Error(t, err)
		assert.EqualError(t, err, "invalid id")
	})

	t.Run("GetAll - получение всех задач", func(t *testing.T) {
		objs, err := r.GetAll(context.Background())
		assert.NoError(t, err)
		assert.Len(t, *objs, 3)
	})

	t.Run("Put - обновление задачи", func(t *testing.T) {
		beforeUpdate := time.Now()

		obj, err := r.Put(context.Background(), "1")
		assert.NoError(t, err)

		assert.GreaterOrEqual(t, obj.UpdatedAt, beforeUpdate)
		assert.Equal(t, "Test Task", obj.Title)
	})

	t.Run("Put - ошибка обновления", func(t *testing.T) {
		_, err := r.Put(context.Background(), "invalid")
		assert.Error(t, err)
	})

	t.Run("Delete - успешное удаление", func(t *testing.T) {
		obj, err := r.Delete(context.Background(), "2")
		assert.NoError(t, err)
		assert.Equal(t, "Task 2", obj.Title)

		// Проверяем обновление ID после удаления
		remaining, _ := r.GetAll(context.Background())
		assert.Len(t, *remaining, 2)
		assert.Equal(t, "1", (*remaining)[0].ID)
		assert.Equal(t, "2", (*remaining)[1].ID) // Бывший 3-й элемент
	})

	t.Run("Delete - ошибка удаления", func(t *testing.T) {
		_, err := r.Delete(context.Background(), "999")
		assert.Error(t, err)
	})
}

func TestConcurrency(t *testing.T) {
	// Тест на конкурентный доступ
	cfg := struct {
		Capacity    int
		MaxItemSize int
	}{
		Capacity:    100,
		MaxItemSize: 1024,
	}

	r, _ := repo.NewRepository(config.InMemory(cfg))
	var wg sync.WaitGroup

	// Запускаем 50 горутин для записи
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			_, err := r.Post(context.Background(), repo.DataObject{
				Title: "Concurrent " + strconv.Itoa(idx),
			})
			assert.NoError(t, err)
		}(i)
	}

	// Запускаем 20 горутин для чтения
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := r.GetAll(context.Background())
			assert.NoError(t, err)
		}()
	}

	wg.Wait()
	objs, _ := r.GetAll(context.Background())
	assert.Len(t, *objs, 50)
}
