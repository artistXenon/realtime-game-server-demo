package HTTP

import (
	// "fmt"
	"net/http"
	// "strconv"
	// "time"

	globalContext "inagame/global-context"
)

func HTTPHandler(w http.ResponseWriter, req *http.Request) {

	
    keys, ok := req.URL.Query()["key"]
    
    if !ok || len(keys[0]) < 1 {
		w.Write([]byte("no key sry"))
        return
    }



	idEmpty := -1;
	l1:
	for index, element := range globalContext.GameIds {
		if len(element) == 0 {
			idEmpty = index
			break l1
		}
	}

	if idEmpty == -1 {
		w.Write([]byte("no chance sry"))
	} else {
		w.Write([]byte("Hell yeah come on!"))
		globalContext.GameIds[idEmpty] = keys[0]
		// strconv.FormatInt(time.Now().UnixMilli(), 10)
	}
}