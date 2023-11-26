package main

import (
        "io"
        "fmt"
        "errors"
        "net/http"
        "database/sql"
        _ "github.com/lib/pq"
)

func checkPG() (bool, error) {

    db, err := sql.Open("postgres", "host=/var/run/postgresql")
    if err != nil {
        return false, err
    }

    defer db.Close()

    err = db.Ping()
    if err != nil {
        return false, err
    }

    rows, err := db.Query("SELECT pg_is_in_recovery()")
    if err != nil {
        return false, err
    }

    defer rows.Close()

    if rows.Next() {

        var res bool
        err = rows.Scan(&res)

        if err != nil {
            return false, err
        }

        return res, nil

    } else {
        return false, errors.New("pg_is_in_recovery() returns empty set")
    }
}

func getRoot(w http.ResponseWriter, r *http.Request) {

    res, err := checkPG()

    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        io.WriteString(w, err.Error())

    } else if res {
        w.WriteHeader(http.StatusPartialContent)
        io.WriteString(w, "Database alive and in recovery")
    } else {
        io.WriteString(w, "Database alive and is master")
    }
}

func main() {

    http.HandleFunc("/", getRoot)
    err := http.ListenAndServe(":8009", nil)
    fmt.Printf("%v\n", err)
}
