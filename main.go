package main

import (
	"bytes"
	"fmt"
	"log"
	"mime/multipart"
	"net/smtp"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

type ResponseData struct {
	To string `json:"to"`
}

func main() {
	// .env 읽기
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
		return
	}

	app := fiber.New()

	app.Post("/send_email", send_email)

	log.Fatal(app.Listen(":8080"))
}

func send_email(c *fiber.Ctx) (err error) {
	// 환경변수
	smtpUsername := os.Getenv("SMTP_USERNAME")
	smtpPassword := os.Getenv("SMTP_PASSWORD")
	smtpFROM := os.Getenv("SMTP_FROM")

	// POST 요청 body 파싱
	var responseData ResponseData

	err = c.BodyParser(&responseData)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// SMTP 인증 정보 설정
	smtpHost := "smtp.gmail.com"
	smtpPort := 587
	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)

	// 메시지 작성
	headerSubject := "메일 테스트\r\n"
	headerBlank := "\r\n"
	body := "메일 테스트 입니다.\r\n"
	msg := []byte(headerSubject + headerBlank + body)

	// MIME 멀티파트 바디 생성
	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)

	// 텍스트 파트 추가
	textPart, err := writer.CreatePart(map[string][]string{
		"Content-Type": {"text/plain; charset=UTF-8"},
	})
	if err != nil {
		fmt.Println(err)
		return err
	}
	textPart.Write(msg)

	// 메일 헤더 설정
	header := make(map[string]string)
	header["From"] = smtpFROM
	header["To"] = responseData.To
	header["Subject"] = "메일 테스트"
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = writer.FormDataContentType()

	// 메일 보내기
	err = smtp.SendMail(fmt.Sprintf("%s:%d", smtpHost, smtpPort), auth, smtpFROM, []string{responseData.To}, composeEmail(header, buf.Bytes()))
	if err != nil {
		fmt.Println(err)
		return err
	}

	return c.SendStatus(200)
}

func composeEmail(header map[string]string, body []byte) []byte {
	var email bytes.Buffer

	for key, value := range header {
		email.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}

	email.WriteString("\r\n")
	email.Write(body)

	return email.Bytes()
}
