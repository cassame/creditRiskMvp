package handlers

import (
	"net/http"
	"time"
)

// HandleMockCreditHistory is external "service" emulator that is checking credit history score (0 to 100)
func HandleMockCreditHistory(w http.ResponseWriter, r *http.Request) {
	passport := r.URL.Query().Get("passport")
	if passport == "" {
		writeError(w, http.StatusBadRequest, "passport is required")
		return
	}
	// demo: passport ends on 8 -> bad credit history \ ends on 9 -> mock timeout:
	last := passport[len(passport)-1]
	if last == '9' {
		//imitation sleep
		time.Sleep(3 * time.Second)
	}
	isGood := last != '8'
	score := 75
	if !isGood {
		score = 20
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"is_good": isGood,
		"score":   score,
	})
}

// HandleMockBankruptcy is external "service" emulator that is checking for bankruptcy
func HandleMockBankruptcy(w http.ResponseWriter, r *http.Request) {
	passport := r.URL.Query().Get("passport")
	if passport == "" {
		writeError(w, http.StatusBadRequest, "passport is required")
		return
	}
	//demo:
	//if passport last digit is 0 -> bankruptcy is true
	//if passport last digit is 9 -> timeout
	last := passport[len(passport)-1]
	if last == '9' {
		//imitation sleep
		time.Sleep(3 * time.Second)
	}
	isBankrupt := last == '0'
	if isBankrupt {
		writeJSON(w, http.StatusOK, map[string]any{
			"isBankrupt": isBankrupt,
		})
	}
}

// HandleMockTerroristList is external "service" emulator that is checking terrorist list
func HandleMockTerroristList(w http.ResponseWriter, r *http.Request) {
	//For demo: “terrorists” - passports ending on 7
	list := []string{
		"1234 567897",
		"1111 222227",
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"updated_at": time.Now().UTC().Format(time.RFC3339),
		"passports":  list,
	})
}
