package main

import (
	"fmt"
	"regexp"

	_ "github.com/lib/pq"
	phonedb "github.com/masahiroyoshida/gophercises/normalizer/db"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "ebuser"
	password = "password"
	dbname   = "gophercise_phone"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable", host, port, user, password)
	must(phonedb.Reset("postgres", psqlInfo, dbname))

	psqlInfo = fmt.Sprintf("%s dbname=%s", psqlInfo, dbname)
	must(phonedb.Migrate("postgres", psqlInfo))

	db, err := phonedb.Open("postgres", psqlInfo)
	must(err)
	defer db.Close()

	if err := db.Seed(); err != nil {
		panic(err)
	}
	phones, err := db.AllPhones()
	must(err)
	for _, p := range phones {
		number := normalize(p.Number)
		if number != p.Number {
			existing, err := db.FindPhone(number)
			must(err)
			if existing != nil {
				// delete
				must(db.DeletePhone(p.ID))

			} else {
				// update
				p.Number = number
				must(db.UpdatePhone(&p))

			}
		} else {
			fmt.Println("No changes required")
		}
	}

}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func normalize(phone string) string {
	re := regexp.MustCompile("\\D")
	return re.ReplaceAllString(phone, "")
}

// func normalize(phone string) string {
// 	var buf bytes.Buffer
// 	for _, ch := range phone {
// 		if ch >= '0' && ch <= '9' {
// 			buf.WriteRune(ch)
// 		}
// 	}
// 	return buf.String()
// }
