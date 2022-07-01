package worker

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
)

func dbConnection() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/worker")
	if err != nil {
		log.Printf("Error %s when opening DB\n", err)
		return nil, err
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(time.Minute * 5)

	return db, nil
}

func newChromedp(db *sql.DB) (context.Context, context.CancelFunc) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:]) // chromedp.Flag("headless", false),

	allocCtx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))

	extractItviecTask(ctx, db)

	return ctx, cancel
}

func worker(id int, jobs <-chan string, results chan<- int, ctx context.Context, db *sql.DB) {
	var titleBook []string

	log.Printf("worker %d Waiting", id)

	for j := range jobs {
		fmt.Println("-- > worker", id, "processing job", j)

		time.Sleep(time.Second)
		titleP := chromedp.Tasks{
			chromedp.Navigate(j),
			chromedp.ActionFunc(func(ctx context.Context) error {
				node, err := dom.GetDocument().Do(ctx)
				if err != nil {
					return err
				}
				res, err := dom.GetOuterHTML().WithNodeID(node.NodeID).Do(ctx)
				if err != nil {
					return err
				}
				doc, err := goquery.NewDocumentFromReader(strings.NewReader(res))
				if err != nil {
					return err
				}

				doc.Find(".truyen-title > a").Each(func(index int, info *goquery.Selection) {
					// text := info.Text()

					// fmt.Println(text)
					// fmt.Println("text", text, w)
					titleT, title := info.Attr("href")
					if title {

						titleBook = append(titleBook, titleT)

						// _, err := db.Exec(fmt.Sprintf(`INSERT INTO book (name, name_link) VALUES (%q, %q)`, text, titleT))
						if err != nil {
							log.Printf("saveBook - Error: %v", err)
						}

					}

				})

				return nil
			}),
		}
		if err := chromedp.Run(ctx, titleP); err != nil {
			fmt.Println(err)

		}
		results <- 1

	}

}

func extractItviecTask(ctx context.Context, db *sql.DB) error {
	var dataTitle []string

	// var listData []string{"title", }
	task := chromedp.Tasks{
		chromedp.Navigate("https://truyenfull.vn"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			node, err := dom.GetDocument().Do(ctx)
			if err != nil {
				return err
			}
			res, err := dom.GetOuterHTML().WithNodeID(node.NodeID).Do(ctx)
			if err != nil {
				return err
			}
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(res))
			if err != nil {
				return err
			}

			doc.Find(".list-truyen > .row > .col-xs-6 > a").Each(func(index int, titlefo *goquery.Selection) {
				titleT, title := titlefo.Attr("href")
				// text := titlefo.Text()
				// _, err := db.Exec(fmt.Sprintf(`INSERT INTO list_book (title) VALUES (%q)`, text))
				if err != nil {
					log.Printf("extractItviecTask - Error: %v", err)
				}
				if title {

					dataTitle = append(dataTitle, titleT)
				}

			})

			return nil
		}),
	}

	if err := chromedp.Run(ctx, task); err != nil {
		fmt.Println(err)

	}

	jobs := make(chan string, 100)
	results := make(chan int, 100)

	for w := 1; w <= 3; w++ {
		go worker(w, jobs, results, ctx, db)
	}

	time.Sleep(10 * time.Second)

	for _, link := range dataTitle {
		jobs <- link
	}
	close(jobs)

	for a := 1; a <= len(dataTitle); a++ {
		<-results

	}

	return nil
}

func main() {

	db, err := dbConnection()
	_, _ = newChromedp(db) //Khởi tạo biến conection
	if err != nil {        //Catch error trong quá trình thực thi
		log.Printf("Error %s when getting db connection", err)
		return
	}

	defer db.Close()

	log.Printf("Successfully connected to database")

}
