package usecase

import (
	"time"
	"github.com/asumsi/livestream/utility"
	"github.com/asumsi/livestream/stream/models"
	"os"
	"encoding/base64"
	"encoding/json"
	"strconv"
	"fmt"
)


type StreamUsecase interface {
	InitSession(offerSession string, roomId string, userId string) (map[string]interface{}, *models.ErrorObject)
	StartStreaming(roomId string) (map[string]interface{}, *models.ErrorObject)
	EndStreaming(pid string) (map[string]interface{}, *models.ErrorObject)
}

type streamRepo struct{
	contextTimeout   time.Duration
}

func NewStreamUsecase(timeout time.Duration) StreamUsecase {
	return &streamRepo{
		contextTimeout : timeout,
	}
}

// InitSession generate answer session base on offer from browser,
// also spawn rtpforwarder instance with randomized udp's video port and audio port
func (a *streamRepo) InitSession(offerSession string, roomId string, userId string) (map[string]interface{}, *models.ErrorObject) {
	errorObject := models.NewErrorObject()
	
	//channel to receive async answer from rtpforwarder instance
	answerSession := make(chan string, 1)
	pid := make(chan int, 1)

	//init session's related variable
	var sessionRTPByte []byte
	var sessionRTPJson string
	var objectSessionRTP map[string]interface{}
	var portUsedByte []byte
	var portUsedJson string
	var objectPortUsed []string

	//retrieve current sessions exist
	sessionRTPBase64 := os.Getenv("RTP_SESSION")
	sessionRTPByte, _ = base64.StdEncoding.DecodeString(sessionRTPBase64)
	sessionRTPJson = string(sessionRTPByte)
	err := json.Unmarshal([]byte(sessionRTPJson), &objectSessionRTP)
	if err != nil {
		errorObject.ErrorActual = err
		errorObject.ErrorMessage = "fail unmarshal sessionRTPJson to objectSessionRTP"
		return nil, errorObject
	}
	//check wether rtpsession from the userid is still exist
	for i,v := range objectSessionRTP {
		//kill the session if exist userId in objectSessionRTP
		if i == userId {
			utility.KillProcess(v.(string))
		}
	}

	//check wether portAudio and portVideo available 
	isPortExist := false
	portAudio := ""
	portVideo := ""
	portUsedBase64 := os.Getenv("PORT_USED")
	portUsedByte, _ = base64.StdEncoding.DecodeString(portUsedBase64)
	portUsedJson = string(portUsedByte)
	err = json.Unmarshal([]byte(portUsedJson), &objectPortUsed)
	if err != nil {
		errorObject.ErrorActual = err
		errorObject.ErrorMessage = "fail unmarshal portUsedJson to objectPortUsed"
		return nil, errorObject
	}
	for {
		portAudio = utility.RandIntToString(49152, 65535)
		portVideo = utility.RandIntToString(49152, 65535)

		isPortExist, _ = utility.InArray(portAudio, objectPortUsed)
		isPortExist, _ = utility.InArray(portVideo, objectPortUsed)

		if !isPortExist {
			break
		}
	}

	// start rtpforwarder instance
	go utility.ExecuteRTPForwarder(portAudio, portVideo, offerSession, answerSession, pid)

	//retrieve rtpforwarder's pid and answersession through channel
	valuePid := <-pid
	valueAnswerSession := <-answerSession

	if valuePid == 0 || valueAnswerSession == "" {
		errorObject.ErrorActual = nil
		errorObject.ErrorCode = "FR"
		errorObject.ErrorMessage = "fail on executing RTP Forwarder"
		return nil, errorObject
	}
	//save the user's rtp pid, for the next checking
	objectSessionRTP[userId] = strconv.Itoa(valuePid)
	sessionRTPByte, _ = json.Marshal(objectSessionRTP)
	sessionRTPBase64 = base64.StdEncoding.EncodeToString(sessionRTPByte)
	os.Setenv("RTP_SESSION", sessionRTPBase64)

	//save audioport and videoport, on env=PORT_USED, to make sure the port not reused by other user
	objectPortUsed = append(objectPortUsed, portAudio)
	objectPortUsed = append(objectPortUsed, portVideo)
	portUsedByte, _ = json.Marshal(objectPortUsed)
	portUsedBase64 = base64.StdEncoding.EncodeToString(portUsedByte)
	os.Setenv("PORT_USED", portUsedBase64)
	
	fmt.Println(objectPortUsed)

	//generate SDP file containing IP,audioport, and videoport, for next to be used by FFMPEG 
	_ = utility.GenerateSDP(portAudio, portVideo, roomId)

	//return answersession, audioport, videoport, and pid of newly created rtpforwarder
	resultUsecase := map[string]interface{}{
		"answer_session" : valueAnswerSession,
		"port_audio" : portAudio,
		"port_video" : portVideo,
		"pid" : valuePid,
	}
	return resultUsecase, nil
}

// StartStreaming trigger FFMPEG to start running, consuming audiport and videoport produced by rtpforwarder 
// also send stream as chunk to dash server
// as output, it will return pid of FFMPEG and manifest url (that can be played on bifrost video player)
func (a *streamRepo) StartStreaming(roomId string) (map[string]interface{}, *models.ErrorObject) {
	errorObject := models.NewErrorObject()

	// channel to receive pid answer from FFMPEG instance
	pid := make(chan int, 1)

	//init ffmpeg's related variable
	var ffmpegInstancesByte []byte
	var objectFfmpeg []string
	var objectFfmpegByte []byte

	//retrieve current ffmpeg exist
	ffmpegInstancesBase64 := os.Getenv("FFMPEG_INSTANCES")
	ffmpegInstancesByte, _ = base64.StdEncoding.DecodeString(ffmpegInstancesBase64)
	err := json.Unmarshal(ffmpegInstancesByte, &objectFfmpeg)
	if err != nil {
		errorObject.ErrorActual = err
		errorObject.ErrorMessage = "fail unmarshal ffmpegInstancesByte to objectFfmpegInstances"
		return nil, errorObject
	}

	// get path of SDP file to be consumed by FFMPEG instance 
	// also construct url stream which can be used as parameter in FFMPEG (-f) to send chunked stream to dash server
	currentDir := utility.GetCurrentDir()
	fileSDP := currentDir + "/" + os.Getenv("SDP_DIRECTORY") + "/" + roomId + "-" + os.Getenv("SDP_FILENAME")
	generatedStreamId := utility.GenerateUUID()
	urlStream := os.Getenv("DASH_SERVER") + "/ldash/" + roomId + "-" + generatedStreamId + "/" + os.Getenv("MANIFEST_FILENAME")

	// start ffmpeg instance as a goroutine, which give PID as channel result 
	go utility.ExecuteFFMPEG(fileSDP, urlStream, pid)
	valuePid := <-pid

	//save the user's ffmpeg pid, for later checking
	objectFfmpeg = append(objectFfmpeg, strconv.Itoa(valuePid))
	objectFfmpegByte, _ = json.Marshal(objectFfmpeg)
	ffmpegInstancesBase64 = base64.StdEncoding.EncodeToString(objectFfmpegByte)
	os.Setenv("FFMPEG_INSTANCES", ffmpegInstancesBase64)

	fmt.Println("ffmpeg instance reserved: ")
	fmt.Println(objectFfmpeg)

	//return pid and url_stream of newly created ffmpeg
	resultUsecase := map[string]interface{}{
		"pid" : valuePid,
		"url_stream" : urlStream,
	}
	return resultUsecase, nil
}


// EndStreaming will kill one instance of FFMPEG base on user or PID, 
// as output, it will return pid of FFMPEG and manifest url (that can be played on bifrost video player)
func (a *streamRepo) EndStreaming(pid string) (map[string]interface{}, *models.ErrorObject) {
	errorObject := models.NewErrorObject()

	//init ffmpeg's related variable
	var ffmpegInstancesByte []byte
	var objectFfmpeg []string
	var isPidExist bool

	//retrieve current ffmpeg exist
	ffmpegInstancesBase64 := os.Getenv("FFMPEG_INSTANCES")
	ffmpegInstancesByte, _ = base64.StdEncoding.DecodeString(ffmpegInstancesBase64)
	err := json.Unmarshal(ffmpegInstancesByte, &objectFfmpeg)
	if err != nil {
		errorObject.ErrorActual = err
		errorObject.ErrorMessage = "fail unmarshal ffmpegInstancesByte to objectFfmpeg"
		return nil, errorObject
	}

	isPidExist, _ = utility.InArray(pid, objectFfmpeg)

	if !isPidExist {
		errorObject.ErrorCode = "PU"
		errorObject.ErrorMessage = "PID not found in list active FFMPEG"
		return nil, errorObject
	}

	// start ffmpeg instance as a goroutine, which give PID as channel result 
	endResult, _ := utility.EndFFMPEG(pid)

	fmt.Println("ffmpeg instance deleted: " + endResult)

	//return pid and url_stream of newly created ffmpeg
	resultUsecase := map[string]interface{}{
		"pid" : pid,
	}
	return resultUsecase, nil
}