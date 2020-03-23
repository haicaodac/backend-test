package botchat

import (
	"encoding/json"
	library "hanyny/app/library"
	models "hanyny/app/models"
	v "hanyny/app/utils/view"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"

	"github.com/asaskevich/govalidator"
)

// Nyny ...
type Nyny struct {
	listData  []map[string]interface{}
	lengthAll int
}

// New ...
func (n *Nyny) New() {
	// message := "Chào bạn"
	training, err := ioutil.ReadFile("app/botchat/data.json")
	if err != nil {
		panic(err.Error())
	}
	err = json.Unmarshal(training, &n.listData)
	if err != nil {
		panic(err.Error())
	}

	n.lengthAll = 0
	for _, data := range n.listData {
		patterns := data["patterns"].([]interface{})
		for _, item := range patterns {
			listWords := strings.Split(item.(string), " ")
			for range listWords {
				n.lengthAll++
			}
		}
	}
}

// QuestionAnswer ...
func (n *Nyny) QuestionAnswer(w http.ResponseWriter, r *http.Request) {
	type Question struct {
		Message string `json:"message" valid:"required~Nội dung không thể trống."`
	}
	var question Question
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&question)
	if err != nil {
		v.Respond(w, v.Message(false, err.Error()))
		return
	}
	defer r.Body.Close()

	// Validator format
	if _, err := govalidator.ValidateStruct(question); err != nil {
		v.Respond(w, v.Message(false, err.Error()))
		return
	}

	pAll := make(map[string]float64)
	for _, data := range n.listData {
		tag := data["tag"].(string)
		patterns := data["patterns"].([]interface{})
		lenPatterns := 0
		for _, item := range patterns {
			listWords := strings.Split(item.(string), " ")
			for range listWords {
				lenPatterns++
			}
		}
		pAll[tag] = float64(lenPatterns) / float64(n.lengthAll)

		listWords := strings.Split(question.Message, " ")
		for _, word := range listWords {
			if word == "bạn" {
				continue
			}
			pWord := checkWord(patterns, word)
			dPatterns := float64(1) / float64(lenPatterns+1)
			if pWord != dPatterns {
				pAll[tag] *= pWord
			}
		}
	}

	max := 0.0
	tag := ""
	for key, val := range pAll {
		if max < val {
			max = val
			tag = key
		}
	}
	for _, data := range n.listData {
		tagStore := data["tag"].(string)
		if tagStore == tag {
			responses := data["responses"].([]interface{})
			ran := randomInt(0, len(responses)-1)

			db := models.OpenDB()
			var botchat models.Botchat
			botchat.Message = question.Message
			botchat.Tag = tag
			botchat.Response = responses[ran].(string)
			db.Create(&botchat)

			data := v.Message(true, "Trả lời.")
			data["answer"] = responses[ran].(string)
			v.Respond(w, data)
			return
		}
	}

	v.Respond(w, v.Message(false, "Có gì đó gây ra lỗi"))
	logger := library.Logger{Type: "ERROR"}
	log := logger.Open()
	log.Println(question.Message + "----" + tag)
	logger.Close()
	return
}

func randomInt(min, max int) int {
	return min + rand.Intn(max-min)
}

func checkWord(list []interface{}, word string) float64 {
	k := 0
	for _, item := range list {
		str1 := strings.ToLower(item.(string))
		str2 := strings.ToLower(word)
		if strings.Contains(str1, str2) {
			// Mỗi lần word xuất hiện trong câu cũ thì thêm biến k++
			k++
		}
	}
	// P(xi|nhan= nonspam)= (k+1)/(sothuthuong+1);
	// trong do: k la so cac mail nonspam xuat hien xi
	// sothuthuong la so mail nonspam
	p := float64((k + 1)) / float64((len(list) + 1))
	return p

}
