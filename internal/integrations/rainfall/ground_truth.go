package rainfall

import (
	"fmt"
	"io"
	"net/http"
)

func getGroundTruth(location string) bool {
	requestUrl := "https://api-open.data.gov.sg/v2/real-time/api/rainfall?date=&paginationToken="

	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Add("x-api-key", "")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))

	panic("getGroundTruth is not implemented")
}
