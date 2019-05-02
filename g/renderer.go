package g

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/radovskyb/watcher"
)

var (
	path      string
	templates map[string]*template.Template
	renderer  *Renderer
	lock      = new(sync.Mutex)
)

// Renderer TODO
type Renderer struct{}

func init() {
	var err error
	path, err = GetCurrentPath()
	if err != nil {
		fmt.Printf("cannot startup project, error: %s\r\n", err.Error())
		os.Exit(1)
	}
}

// GetCurrentPath TODO
func GetCurrentPath() (string, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "", err
	}
	return strings.Replace(dir, "\\", "/", -1), nil
}

// NewRenderer TODO
func NewRenderer() *Renderer {
	lock.Lock()
	if renderer != nil {
		return renderer
	}
	renderer = &Renderer{}
	lock.Unlock()
	flush()
	go watch()
	return renderer
}

func flush() {
	lock.Lock()
	defer lock.Unlock()
	templates = make(map[string]*template.Template)

	templatesDir := path + "/templates/"
	pages := []string{}
	if err := filepath.Walk(templatesDir+"views/", func(page string, f os.FileInfo, err error) error {
		if !f.IsDir() && strings.HasSuffix(page, ".html") {
			if err != nil {
				return err
			}
			if _, err := ioutil.ReadFile(page); err != nil {
				return err
			}
			pages = append(pages, page)
		}
		return nil
	}); err != nil {
		log.Fatal(err)
	}

	layouts, err := filepath.Glob(templatesDir + "layouts/*.html")
	if err != nil {
		log.Fatal(err)
	}

	partials, err := filepath.Glob(templatesDir + "partials/*.html")
	if err != nil {
		log.Fatal(err)
	}
	includes := append(layouts, partials...)
	// Generate our templates map from our layouts/ and partials/ directories
	for _, page := range pages {
		files := append(includes, page)
		name := page[len(templatesDir+"views/"):]
		templates[name] = template.Must(parse(path, name, files...))
	}
}

func watch() {
	w := watcher.New()

	// SetMaxEvents to 1 to allow at most 1 event's to be received
	// on the Event channel per watching cycle.
	//
	// If SetMaxEvents is not set, the default is to send all events.
	w.SetMaxEvents(1)
	w.IgnoreHiddenFiles(true)

	// Only notify rename and move events.
	w.FilterOps(watcher.Rename, watcher.Move, watcher.Create, watcher.Write, watcher.Remove)

	go func() {
		for {
			select {
			case event := <-w.Event:
				fmt.Println(event) // Print the event's info.
				flush()
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	// Watch test_folder recursively for changes.
	if err := w.AddRecursive(path + "/templates"); err != nil {
		log.Fatalln(err)
	}

	// Start the watching process - it'll check for changes every 100ms.
	if err := w.Start(time.Millisecond * 100); err != nil {
		log.Fatalln(err)
	}
}

// Render TODO
func (t *Renderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	// Ensure the template exists in the map.
	lock.Lock()
	tmpl, ok := templates[name]
	lock.Unlock()
	if !ok {
		return fmt.Errorf("The template %s does not exist", name)
	}

	if err := tmpl.ExecuteTemplate(w, name, data); err != nil {
		log.Fatal(err)
	}
	return nil
}

// GetTemplate TODO
func (t *Renderer) GetTemplate(name string) (*template.Template, error) {
	if tmpl, ok := templates[name]; ok {
		return tmpl, nil
	}
	return nil, fmt.Errorf("specified template: %s does not exist", name)
}

func parse(path string, name string, files ...string) (*template.Template, error) {
	if len(files) == 0 {
		// Not really a problem, but be consistent.
		return nil, fmt.Errorf("html/template: no files named in call to parse")
	}
	tmpl := template.New(name)
	for _, file := range files {
		b, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, err
		}
		if _, err := tmpl.Parse(string(b)); err != nil {
			return nil, err
		}
	}
	return tmpl, nil
}
