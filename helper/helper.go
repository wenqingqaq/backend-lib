package helper

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/tcolgate/mp3"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func ReceiveEmail(str string) string {
	startIndex := strings.LastIndex(str, "email:[")
	endIndex := strings.Index(str, "],")
	return str[startIndex+len("email:[") : endIndex]
}
func ReceiveWebhook(str string) string {
	startIndex := strings.LastIndex(str, "webhook:[")
	endIndex := strings.LastIndex(str, "]")
	return str[startIndex+len("webhook:[") : endIndex]
}

func EventIdStr(eventId string) string {
	eventId = strings.ReplaceAll(eventId, "[", "")
	return strings.ReplaceAll(eventId, "]", "")
}

// HashPassword 使用 bcrypt 对密码进行哈希处理
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// ComparePasswords 验证密码是否匹配哈希值
func ComparePasswords(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err
}

func GetMp3UrlDuration(url string) (int64, error) {
	path, err := DownloadAudioFileToTmp(url)
	if err != nil {
		return 0, err
	}
	return ParseFileDuration(path)
}

func DownloadAudioFileToTmp(urlStr string) (string, error) {
	// 解析一下url
	ul, err := url.Parse(urlStr)
	if err != nil {
		log.Error(err)
	}
	resp, err := http.Get(urlStr)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	splitS := strings.Split(ul.Path, "/")
	s := splitS[len(splitS)-1]
	pathStr := "../../tmp/" + s
	out, err := os.Create(pathStr)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)

	return pathStr, nil
}

// ParseFileDuration 解析音频文件时长 原生go的方式
func ParseFileDuration(filePath string) (int64, error) {
	fd, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer fd.Close()

	d := mp3.NewDecoder(fd)
	var f mp3.Frame
	skipped := 0

	var t float64
	for {
		if err := d.Decode(&f, &skipped); err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
			break
		}
		t = t + f.Duration().Seconds()
	}

	return int64(t), nil
}
