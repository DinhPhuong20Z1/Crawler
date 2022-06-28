package main

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
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/crawl")
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

func extractItviecTask(ctx context.Context, db *sql.DB) error {
	var dataTitle []string
	var titleBook []string
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

			doc.Find("ul.navbar-nav > li.dropdown > ul.dropdown-menu > li > a").Each(func(index int, titlefo *goquery.Selection) {
				titleT, title := titlefo.Attr("href")
				text := titlefo.Text()
				_, err := db.Exec(fmt.Sprintf(`INSERT INTO list_book (title) VALUES (%q)`, text))
				if err != nil {
					log.Printf("extractItviecTask - Error: %v", err)
				}
				if title {
					dataTitle = append(dataTitle, titleT)
					// fmt.Printf(titleT)
					// fmt.Println(titleT)

				}

			})

			// doc.Find("ul.navbar-nav > li.dropdown > ul.dropdown-menu > li > a").Each(func(index int, info *goquery.Selection) {
			// 	// text := info.Text()
			// 	// fmt.Println(text)
			// 	// _, err := db.Exec(fmt.Sprintf(`INSERT INTO list_book (title) VALUES (%q)`, text))
			// 	// if err != nil {
			// 	// 	log.Printf("extractItviecTask - Error: %v", err)
			// 	// }
			// })

			return nil
		}),
	}

	if err := chromedp.Run(ctx, task); err != nil {
		fmt.Println(err)

	}

	fmt.Println("dataTitle", len(dataTitle))
	for _, num := range dataTitle {
		// fmt.Println("num", num)
		titleP := chromedp.Tasks{
			chromedp.Navigate(num),
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
					// _, err := db.Exec(fmt.Sprintf(`INSERT INTO list_book (title) VALUES (%q)`, text))
					// if err != nil {
					// 	log.Printf("extractItviecTask - Error: %v", err)
					// }
					// fmt.Println(text)
					titleT, title := info.Attr("href")
					if title {
						titleBook = append(titleBook, titleT)
						// fmt.Printf(titleT)
						// fmt.Println(titleT)

					}

				})

				return nil
			}),
		}
		if err := chromedp.Run(ctx, titleP); err != nil {
			fmt.Println(err)

		}
	}

	for _, num := range titleBook {
		// fmt.Println("num", num)
		titleP := chromedp.Tasks{
			chromedp.Navigate(num),
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

				doc.Find("#truyen .col-truyen-main .col-info-desc h3.title").Each(func(index int, info *goquery.Selection) {
					// text := info.Text()
					// _, err := db.Exec(fmt.Sprintf(`INSERT INTO list_book (title) VALUES (%q)`, text))
					// if err != nil {
					// 	log.Printf("extractItviecTask - Error: %v", err)
					// }
					// fmt.Println(text)

				})

				doc.Find("#truyen .col-truyen-main .col-info-desc .desc .desc-text-full").Each(func(index int, info *goquery.Selection) {
					// desc := info.Text()
					// _, err := db.Exec(fmt.Sprintf(`INSERT INTO list_book (title) VALUES (%q)`, text))
					// if err != nil {
					// 	log.Printf("extractItviecTask - Error: %v", err)
					// }
					// fmt.Println(desc)

				})

				doc.Find("#truyen .col-truyen-main .col-info-desc .info-holder .desc .info > div:nth-child(1) a").Each(func(index int, info *goquery.Selection) {
					// desc := info.Text()
					// _, err := db.Exec(fmt.Sprintf(`INSERT INTO list_book (title) VALUES (%q)`, text))
					// if err != nil {
					// 	log.Printf("extractItviecTask - Error: %v", err)
					// }
					// fmt.Println(desc)

				})

				doc.Find("#truyen .col-truyen-main .col-info-desc .info-holder .desc .info > div:nth-child(2) > a").Each(func(index int, info *goquery.Selection) {
					category := info.Text()
					// _, err := db.Exec(fmt.Sprintf(`INSERT INTO list_book (title) VALUES (%q)`, text))
					// if err != nil {
					// 	log.Printf("extractItviecTask - Error: %v", err)
					// }
					fmt.Println(category)

				})

				return nil
			}),
		}
		if err := chromedp.Run(ctx, titleP); err != nil {
			fmt.Println(err)

		}
	}

	fmt.Println(len(titleBook))

	return nil
}

func main() {

	// close chrome
	// _, cancel := newChromedp()
	// defer cancel()
	db, err := dbConnection()
	_, _ = newChromedp(db) //Khởi tạo biến conection
	if err != nil {        //Catch error trong quá trình thực thi
		log.Printf("Error %s when getting db connection", err)
		return
	}
	defer db.Close()

	log.Printf("Successfully connected to database")

}
