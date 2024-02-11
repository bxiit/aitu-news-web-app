package main

import (
	"alexedwards.net/snippetbox/pkg/forms"
	"alexedwards.net/snippetbox/pkg/models"
	"html/template"
	"path/filepath"
	"time"
)

// Add a Flash field to the templateData struct.
type templateData struct {
	CurrentYear     int
	Snippet         *models.News
	Snippets        []*models.News
	Form            *forms.Form
	Flash           string
	IsAuthenticated bool
	Role            string
	User            *models.User
	Users           []*models.User
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	// Инициализируем пустое отображение для хранения шаблонов
	cache := map[string]*template.Template{}

	// Находим все файлы, соответствующие шаблону "*.page.tmpl" в указанном каталоге
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	// Перебираем каждый найденный файл шаблона страницы
	for _, page := range pages {
		// Извлекаем базовое имя файла
		name := filepath.Base(page)

		// Создаем новый шаблон с указанным именем и разбираем файл шаблона страницы
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Разбираем шаблоны макетов (*.layout.tmpl) и включаем их в набор шаблонов
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		// Разбираем частичные шаблоны (*.partial.tmpl) и включаем их в набор шаблонов
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		// Добавляем шаблон в кэш
		cache[name] = ts
	}
	return cache, nil
}
