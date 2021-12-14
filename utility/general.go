package utility

import (
	"strconv"
	"math/rand"
	"time"
	"github.com/google/uuid"
	"strings"
	"fmt"
	"os/exec"
    "reflect"
)

func RandIntToString(min int, max int) string{
	rand.Seed(time.Now().UnixNano())
	randValInt := rand.Intn(max - min + 1) + min
	return strconv.Itoa(randValInt)
}

func GenerateUUID() string {
	requestIDObject := uuid.Must(uuid.NewRandom())
	requestID := strings.Replace(fmt.Sprintf("%v", requestIDObject), "-", "", -1)
	return requestID
}

func KillProcess(pid string){
	kill := exec.Command("kill", "-9", pid)
    err := kill.Run()
    if err != nil {
        //log.WithError(err).Error("Error killing chromium process")
		fmt.Println(err)
		return
    }
	fmt.Println("successfully kill process : " + pid)
}

func InArray(val interface{}, array interface{}) (exists bool, index int) {
    exists = false
    index = -1

    switch reflect.TypeOf(array).Kind() {
    case reflect.Slice:
        s := reflect.ValueOf(array)

        for i := 0; i < s.Len(); i++ {
            if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
                index = i
                exists = true
                return
            }
        }
    }

    return
}