package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"strconv"
	"time"
)

var db *sql.DB

func main() {
	//инициализируем соединение
	var err error
	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/Calorie")
	if err != nil {
		panic(err)
	}
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	//Создание нового экземпляра сервера
	server := echo.New()
	//создание роута POST, который будет доступен по пути /post и будет обрабатывать функцию HandlerPOST
	server.POST("/post", HandlerPOST)
	//создание роута GET, который будет доступен по пути /get и будет обрабатывать функцию HandlerGET
	server.GET("/get", HandlerGET)

	//запускаем сервер и делаем проверку на возможную ошибку
	errs := server.Start(":8080")
	if errs != nil {
		fmt.Println("Ошибка")
	}
}

func conv_int(indicator string) int {
	var err error
	var indicator_int int
	if indicator != "" {
		if indicator_int, err = strconv.Atoi(indicator); err != nil {
			log.Fatal(err)
		} else {
			if indicator_int >= 0 {
				return indicator_int
			} else {
				log.Panicln("Число должно быть положительным!")
			}
		}
	} else {
		return indicator_int
	}
	return indicator_int
}

func conv_date(date string) string {
	if date == "" {
		date_c := time.Now()
		date_now_new := date_c.Format("02.01.2006")
		return date_now_new
	}
	date_c, er := time.Parse("02.01.2006", date)
	if er != nil {
		log.Fatal(er)
	}
	date_new := date_c.Format("02.01.2006")
	return date_new
}

// функция-обработчик роута POST
func HandlerPOST(c echo.Context) (err error) {
	prot := c.QueryParam("prot")
	if err != nil {
		log.Fatal(err)
	}
	prot_int := conv_int(prot)
	fmt.Println("Prot =", prot_int)

	fats := c.QueryParam("fats")
	if err != nil {
		log.Fatal(err)
	}
	fats_int := conv_int(fats)
	fmt.Println("Fats =", fats_int)

	carb := c.QueryParam("carb")
	if err != nil {
		log.Fatal(err)
	}
	carb_int := conv_int(carb)
	fmt.Println("Carb =", carb_int)

	date_imp := c.QueryParam("date")
	if err != nil {
		log.Fatal(err)
	}
	date_imp_new := conv_date(date_imp)
	fmt.Println("Date =", date_imp_new)

	result, e := db.Exec("INSERT INTO diet (date, protein, fats, carb) VALUES (?,?,?,?) ON DUPLICATE KEY UPDATE protein=protein+VALUES(protein),fats=fats+VALUES(fats),carb=carb+VALUES(carb)", date_imp_new, prot_int, fats_int, carb_int)
	if e != nil {
		log.Fatal(e)
	}
	id, er := result.LastInsertId()
	if er != nil {
		log.Panicln(err)
	}
	fmt.Println(id)
	return nil

}

type Information struct {
	Date    string
	Protein int
	Fats    int
	Carb    int
}

// функция-обработчик роута GET
func HandlerGET(c echo.Context) error {
	var info Information
	date_imp := c.QueryParam("date")
	date_imp_new := conv_date(date_imp)
	result := db.QueryRow("SELECT * FROM diet WHERE date=?", date_imp_new)
	if err := result.Scan(&info.Date, &info.Protein, &info.Fats, &info.Carb); err != nil {
		if err == sql.ErrNoRows {
			fmt.Errorf("Нет данных по дате: %s", date_imp_new)
		}
		fmt.Errorf("%s, %v", date_imp_new, err)
	}
	fmt.Println(info)
	return nil
}
