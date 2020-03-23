/*
 * Created by Dac Hai on 20/10/2018
 */

package view

import (
	"encoding/json"
	"net/http"
)

// View bala
type View struct {
	Template  string
	Layout    string
	Extension string
	Vars      map[string]interface{}
}

// Message ...
func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

// Respond json
func Respond(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(data)
}
