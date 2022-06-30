package main

import (
	"Crawler/models"
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
	"github.com/jmoiron/sqlx"
)

func dbConnection() (*sqlx.DB, error) {
	db, err := sqlx.Connect("mysql", "root:@tcp(127.0.0.1:3306)/crawl")
	if err != nil {
		log.Printf("Error %s when opening DB\n", err)
		return nil, err
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(time.Minute * 5)

	return db, nil
}

func newChromedp(db *sqlx.DB) (context.Context, context.CancelFunc) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:]) // chromedp.Flag("headless", false),

	allocCtx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))

	saveBook(ctx, db)

	return ctx, cancel
}

var Catagories []models.CatagoryModel

var Chapter []models.ChapterModel

func getAllCatagories(db *sqlx.DB) {
	query := "select * from categories"
	rows, err := db.Queryx(query)
	if err != nil {
		panic(err)
	}

	Catagories = make([]models.CatagoryModel, 0)
	for rows.Next() {
		do := models.CatagoryModel{}
		err := rows.StructScan(&do)
		if err != nil {
			panic(err)
		}

		Catagories = append(Catagories, do)
	}
}

func getAllChapter(db *sqlx.DB) {
	query := "select * from chapters"
	rows, err := db.Queryx(query)
	if err != nil {
		panic(err)
	}

	Chapter = make([]models.ChapterModel, 0)
	for rows.Next() {
		do := models.ChapterModel{}
		err := rows.StructScan(&do)
		if err != nil {
			panic(err)
		}

		Chapter = append(Chapter, do)
	}
}

func saveBook(ctx context.Context, db *sqlx.DB) error {
	var dataCategories []string
	var titleBook []string
	var des string
	var title string
	var author string
	var source string
	var status string
	var categories []string
	var titleData string
	// var example string
	// var load []string
	var listData []string
	// var words []string
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

				// _, err := db.Exec(fmt.Sprintf(`INSERT INTO categories (name, link) VALUES (%q, %q)`, text, titleT))

				if err != nil {
					log.Printf("saveBook - Error: %v", err)
				}
				if title {
					dataCategories = append(dataCategories, titleT)
					// db.Exec(fmt.Sprintf(`INSERT INTO list_categories (link) VALUES (%q)`, titleT))
					// fmt.Printf(titleT)
					// fmt.Println(titleT)
				}

			})

			// doc.Find("ul.navbar-nav > li.dropdown > ul.dropdown-menu > li > a").Each(func(index int, info *goquery.Selection) {
			// 	// text := info.Text()
			// 	// fmt.Println(text)
			// 	// _, err := db.Exec(fmt.Sprintf(`INSERT INTO list_book (title) VALUES (%q)`, text))
			// 	// if err != nil {
			// 	// 	log.Printf("saveBook - Error: %v", err)
			// 	// }
			// })

			return nil
		}),
	}

	if err := chromedp.Run(ctx, task); err != nil {
		fmt.Println(err)

	}

	fmt.Println("dataCategories", len(dataCategories))
	for _, num := range dataCategories {
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
					// 	log.Printf("saveBook - Error: %v", err)
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

	fmt.Println(len(titleBook))

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
					text := info.Text()
					// _, err := db.Exec(fmt.Sprintf(`INSERT INTO list_book (title) VALUES (%q)`, text))

					// if err != nil {
					// 	log.Printf("saveBook - Error: %v", err)
					// }
					title = text
					// fmt.Println(text)

				})

				doc.Find("#truyen .col-truyen-main .col-info-desc .desc .desc-text").Each(func(index int, info *goquery.Selection) {
					desc := info.Text()
					// _, err := db.Exec(fmt.Sprintf(`INSERT INTO list_book (des) VALUES (%q)`, desc))
					// if err != nil {
					// 	log.Printf("saveBook - Error: %v", err)
					// }
					des = desc
					// fmt.Println(desc)

				})

				doc.Find("#truyen .col-truyen-main .col-info-desc .info-holder .info div:nth-child(1) a").Each(func(index int, info *goquery.Selection) {
					authors := info.Text()
					// _, err := db.Exec(fmt.Sprintf(`INSERT INTO list_book (author) VALUES (%q)`, authors))
					// if err != nil {
					// 	log.Printf("saveBook - Error: %v", err)
					// }
					author = authors
					// fmt.Println(desc)

				})

				doc.Find("#truyen .col-truyen-main .col-info-desc .info-holder .info div:nth-child(2) a").Each(func(index int, info *goquery.Selection) {
					categorie := info.Text()
					// _, err := db.Exec(fmt.Sprintf(`INSERT INTO list_book (author) VALUES (%q)`, authors))
					// if err != nil {
					// 	log.Printf("saveBook - Error: %v", err)
					// }
					// categories = authors
					categories = append(categories, categorie)
					// fmt.Println(desc)

				})

				doc.Find("#truyen .col-truyen-main .col-info-desc .info-holder .info div:nth-child(3) span").Each(func(index int, info *goquery.Selection) {
					sources := info.Text()
					// _, err := db.Exec(fmt.Sprintf(`INSERT INTO list_book (source) VALUES (%q)`, sources))
					// if err != nil {
					// 	log.Printf("saveBook - Error: %v", err)
					// }
					source = sources
					// fmt.Println(sources)

				})

				doc.Find("#truyen .col-truyen-main .col-info-desc .info-holder .info div:nth-child(4) span").Each(func(index int, info *goquery.Selection) {
					statuss := info.Text()

					status = statuss
					// fmt.Println(category)

				})
				doc.Find(".list-chapter li a").Each(func(index int, info *goquery.Selection) {
					chapter := info.Text()
					titleT, titleS := info.Attr("href")
					if titleS {
						_, err := db.Exec(fmt.Sprintf(`INSERT INTO chapters (name, name_book, link_chapters) VALUES (%q, %q ,%q )`, chapter, title, titleT))
						if err != nil {
							log.Printf("saveBook - Error: %v", err)
						}
					}

					if err != nil {
						log.Printf("saveBook - Error: %v", err)
					}
					source = chapter
					// fmt.Println(sources)

				})
				// if author != nil {
				// 	_, err := db.Exec(fmt.Sprintf(`INSERT INTO chapters (name, name_book, id_name_book, link_chapters) VALUES (%q, %q , %d, %q )`, title, des, author, source, status))
				// }

				rows, err := db.Exec(fmt.Sprintf(`INSERT INTO books (title, des, author, source, status) VALUES (%q, %q , %q, %q , %q)`, title, des, author, source, status))

				if err != nil {
					log.Printf("saveBook - Error: %v", err)
				}

				bookID, err := rows.LastInsertId()

				if err != nil {
					panic(err)
				}

				// Tim catagory cho book
				listCategoriesID := make([]int, 0)
				for _, v := range Catagories {

					if strings.Contains(strings.Join(categories, ","), v.Name) {
						listCategoriesID = append(listCategoriesID, v.ID)

					}
				}

				// string.Contains func trong golang được sử dụng để kiểm tra các chữ cái đã cho có trong chuỗi hay không. nếu có kí tự có trong chuỗi đã cho
				// thì nó trả về true ngược lại trả về false / func Contains(str, substr string) bool / fmt.Println(strings.Contains("GeeksforGeeks", "for")) / for tồn tại trong chuỗi

				// string.Join ()  func nối tất cả các phần từ có trong chuỗi thành một chuỗi duy nhất, hàm này có sẵn trong string package
				// func Join(s []string, sep string) string / nối chúng với nhau và cách nhau bằng dấu ","
				// Save vao bang book_categories
				for _, id := range listCategoriesID {
					query := fmt.Sprintf("INSERT INTO book_categories(book_id, category_id) VALUES (%d, %d)", bookID, id)
					db.Exec(query)
				}

				categories = make([]string, 0)
				return nil
			}),
		}
		if err := chromedp.Run(ctx, titleP); err != nil {
			fmt.Println(err)

		}
	}

	for _, num := range titleBook {
		// fmt.Println("num", num)
		chapterData := chromedp.Tasks{
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
					text := info.Text()
					titleData = text

				})

				// f := func(i int, sel *goquery.Selection) bool {
				// 	return strings.HasPrefix(sel.Text(), "")
				// }

				// doc.Filter(".list-chapter li").Each(func(index int, info *goquery.Selection) {
				// 	fmt.Println("1", index)
				// 	fmt.Println("2", info)
				// })
				// for i := 2; i < 7 && i != 3; {

				doc.Find("#list-chapter .pagination li:nth-child(2) a").Each(func(index int, info *goquery.Selection) {
					// words = append(words, sel.Text())
					link, _ := info.Attr("href")
					// fmt.Println("sel", sel.Text())
					listData = append(listData, link)
					// i = i + 1
					// for i := 2; i++ {

					// }
				})

				// doc.Find("#list-chapter .pagination li:nth-child(4) a").Each(func(index int, info *goquery.Selection) {
				// 	// words = append(words, sel.Text())
				// 	link, _ := info.Attr("href")
				// 	// fmt.Println("sel", sel.Text())
				// 	listData = append(listData, link)
				// 	// i = i + 1
				// 	// for i := 2; i++ {

				// 	// }
				// })

				// doc.Find("#list-chapter .pagination li:nth-child(5) a").Each(func(index int, info *goquery.Selection) {
				// 	// words = append(words, sel.Text())
				// 	link, _ := info.Attr("href")
				// 	// fmt.Println("sel", sel.Text())
				// 	listData = append(listData, link)
				// 	// i = i + 1
				// 	// for i := 2; i++ {

				// 	// }
				// })

				// }

				doc.Find(".list-chapter li a").Each(func(index int, info *goquery.Selection) {
					chapter := info.Text()
					titleT, titleS := info.Attr("href")
					if titleS {
						_, err := db.Exec(fmt.Sprintf(`INSERT INTO chapters (name, name_book, link_chapters) VALUES (%q, %q ,%q )`, chapter, titleData, titleT))
						if err != nil {
							log.Printf("saveBook - Error: %v", err)
						}
					}

					if err != nil {
						log.Printf("saveBook - Error: %v", err)
					}
					// source = chapter

				})

				return nil
			}),
		}

		// fmt.Println("words", words)
		if err := chromedp.Run(ctx, chapterData); err != nil {
			fmt.Println(err)

		}

	}

	// listChapterName := make([]string, 0)
	// for _, v := range Chapter {

	// 	if !!strings.Contains(strings.Join(listData, ","), v.Name) {
	// 		listChapterName = append(listChapterName, v.Link)

	// 	}
	// }

	for _, num := range listData {
		chapterList := chromedp.Tasks{
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
					text := info.Text()
					titleData = text

				})

				// smt := fmt.Fprintf("#list-chapter .pagination li:nth-child(%d) a", i)
				doc.Find(".list-chapter li a").Each(func(index int, info *goquery.Selection) {
					// words = append(words, sel.Text())
					chapter := info.Text()
					titleT, titleS := info.Attr("href")
					// fmt.Println("sel", sel.Text())
					if titleS {

						_, err := db.Exec(fmt.Sprintf(`INSERT INTO chapters (name, name_book, link_chapters) VALUES (%q, %q ,%q )`, chapter, titleData, titleT))
						if err != nil {
							log.Printf("saveBook - Error: %v", err)
						}
					}

					// for i := 2; i++ {

					// }
				})

				return nil
			}),
		}

		fmt.Println("listData", listData)
		if err := chromedp.Run(ctx, chapterList); err != nil {
			fmt.Println(err)

		}
	}

	return nil
}

func main() {
	db, err := dbConnection()

	if err != nil { //Catch error trong quá trình thực thi
		log.Printf("Error %s when getting db connection", err)
		return
	}
	defer db.Close()

	getAllCatagories(db)
	getAllChapter(db)
	_, _ = newChromedp(db) //Khởi tạo biến conection

	log.Printf("Successfully connected to database")

}
