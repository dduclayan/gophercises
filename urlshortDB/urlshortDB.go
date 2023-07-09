package urlshortDB

import (
	"encoding/json"
	"fmt"
	"net/http"

	bolt "go.etcd.io/bbolt"
)

var (
	dbname = "urlpath.db"
)

type pathToURL struct {
	Path string `json:"path"`
	Url  string `json:"url"`
}

func AddEntry(path, url string) error {
	db, err := bolt.Open(dbname, 0600, nil)
	if err != nil {
		return fmt.Errorf("bolt.Open(urlpath.db): %v", err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		err = tx.Bucket([]byte("PATHTOURL")).Put([]byte(path), []byte(url))
		if err != nil {
			return fmt.Errorf("tx.Bucket(PATHTOURL).Put(%v,%v): %v", path, url, err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("could not setup buckets: %v", err)
	}
	fmt.Println("entry added")
	return nil
}

func SetupDB() (*bolt.DB, error) {
	db, err := bolt.Open(dbname, 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("could not open db: %v", err)
	}
	defer db.Close()
	// setup db
	err = db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte("PATHTOURL"))
		if err != nil {
			return fmt.Errorf("could not create pathtourl bucket: %v", err)
		}
		err = tx.Bucket([]byte("PATHTOURL")).Put([]byte("/yt"), []byte("https://youtube.com"))
		if err != nil {
			return fmt.Errorf("could not insert path/url: %v", err)
		}
		fmt.Println("Added default /yt -> https://youtube.com path")
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("could not setup buckets: %v", err)
	}
	fmt.Println("DB setup done")
	return db, nil
}

func DBHandler(fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		// see if path exists in DB
		url, err := checkDBForPath(path)
		if err != nil {
			fmt.Printf("checkDBForPath(b,%v): %v\n", path, err)
		}
		// redirect to found path
		if url != "" {
			http.Redirect(w, r, url, http.StatusPermanentRedirect)
		}
		// else
		fallback.ServeHTTP(w, r)
	}
}

func checkDBForPath(path string) (string, error) {
	var url string
	db, err := bolt.Open(dbname, 0600, nil)
	if err != nil {
		return "", fmt.Errorf("bolt.Open(urlpath.db): %v", err)
	}
	defer db.Close()
	err = db.View(func(tx *bolt.Tx) error {
		resp := tx.Bucket([]byte("PATHTOURL")).Get([]byte(path))
		if resp != nil {
			url = string(resp)
		}
		// read the DB
		// entry := tx.Bucket([]byte("PATHTOURL"))
		// entry.ForEach(func(k, v []byte) error {
		// 	fmt.Println(string(k), string(v))
		// 	return nil
		// })
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("could not read db: %v", err)
	}
	return url, nil
}

func ParseJSON(data []byte) (map[string]string, error) {
	var pathURLs []pathToURL
	pathURLsMap := make(map[string]string)
	if err := json.Unmarshal(data, &pathURLs); err != nil {
		return nil, err
	}
	for _, val := range pathURLs {
		pathURLsMap[val.Path] = val.Url
	}
	return pathURLsMap, nil
}
