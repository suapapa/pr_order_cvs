package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/suapapa/pr_label/order"
)

func main() {
	file, _ := os.Open(os.Args[1])
	rdr := csv.NewReader(bufio.NewReader(file))
	rows, _ := rdr.ReadAll()
	for i, row := range rows {
		if i == 0 {
			continue
		}
		// already sent
		if row[11] == "1" {
			continue
		}

		ord := &order.Order{
			ID: atoi(row[0]),
			From: &order.Addr{
				Line1:      row[1],
				Line2:      row[2],
				Name:       row[3],
				PostNumber: row[4],
			},
			To: &order.Addr{
				Line1:      row[5],
				Line2:      row[6],
				Name:       row[7],
				PostNumber: row[8],
			},
			Items: []*order.Item{},
		}
		dungeon01Cnt := atoi(row[9])
		if dungeon01Cnt > 0 {
			ord.Items = append(ord.Items, &order.Item{Name: "dungeon_01", Cnt: dungeon01Cnt})
		}

		defer01Cnt := atoi(row[10])
		if defer01Cnt > 0 {
			ord.Items = append(ord.Items, &order.Item{Name: "defer_01", Cnt: defer01Cnt})
		}

		// params := url.Values{}
		// params.Add("{ID:1234567890,from:{line1:경기 성남시 분당구 판교역로 235 (에이치 스퀘어 엔동),line2:7층,name:카카오 엔터프라이즈,phone_number:010-1234-5678},to:{line1:경기도 성남시 분당구 판교역로 166,name:판교 아지트,phone_number:010-1234-5678}}", ``)
		// body := strings.NewReader(params.Encode())
		buf := &bytes.Buffer{}
		json.NewEncoder(buf).Encode(ord)
		req, err := http.NewRequest("POST", "http://rpi-airplay.local:8080/v1/order", buf)
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		resp.Body.Close()
	}
}

func atoi(str string) int {
	ret, err := strconv.Atoi(str)
	if err != nil {
		log.Fatal(err)
	}
	return ret
}
