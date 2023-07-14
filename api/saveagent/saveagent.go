package main

import (
	"bytes"
	"edetector_API/pkg/logger"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

type Request struct {
	IP         string `json:"ip"`
	Port       string `json:"port"`
	DetectPort string `json:"detect_port"`
}

func main() {

	request := Request{
		IP:         "192.168.200.161",
		Port:       "5000",
		DetectPort: "5001",
	}

	// Marshal payload into JSON
	payload, err := json.Marshal(request)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	// Create an HTTP request
	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://192.168.200.167:8080/agent", bytes.NewBuffer(payload)) // need port change
	if err != nil {
		logger.Error("Error creating HTTP request: " + err.Error())
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the HTTP request
	response, err := client.Do(req)
	if err != nil {
		logger.Error("Error sending HTTP request: " + err.Error())
		return
	}
	defer response.Body.Close()

	// Check the response status code
	if response.StatusCode != http.StatusOK {
		fmt.Println("Request failed with status code:", response.StatusCode)
		return
	}

	// Create the output file
	agentFile, err := os.Create("agent.exe")
	if err != nil {
		logger.Error("Error creating output file: " + err.Error())
		return
	}
	defer agentFile.Close()

	// Save the response body to the output file
	_, err = io.Copy(agentFile, response.Body)
	if err != nil {
		logger.Error("Error saving response: " + err.Error())
		return
	}

	fmt.Println("Agent saved successfully!")

	// Send the file as an email attachment
	err = sendEmail("sender@gmail.com", "recipient@gmail.com", "Agent File", "Please find the agent file attached.", "agent.exe")
	if err != nil {
		fmt.Println("Error sending email:", err)
		return
	}

	fmt.Println("Email sent successfully!")
}

func sendEmail(sender, recipient, subject, body, attachmentPath string) error {
	// Load the Google API credentials file
	credentials, err := ioutil.ReadFile("credentials.json") // Replace with your credentials file path
	if err != nil {
		return fmt.Errorf("failed to read credentials file: %v", err)
	}

	// Initialize the Google OAuth config
	config, err := google.JWTConfigFromJSON(credentials, gmail.GmailSendScope)
	if err != nil {
		return fmt.Errorf("failed to initialize OAuth config: %v", err)
	}

	// Create a new Gmail client using the OAuth config
	client := config.Client(oauth2.NoContext)

	// Create a new Gmail service
	service, err := gmail.New(client)
	if err != nil {
		return fmt.Errorf("failed to create Gmail service: %v", err)
	}

	// Create the email message
	message := &gmail.Message{
		Raw: encodeEmailMessage(sender, recipient, subject, body, attachmentPath),
	}

	// Send the email
	_, err = service.Users.Messages.Send("me", message).Do()
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}

func encodeEmailMessage(sender, recipient, subject, body, attachmentPath string) string {
	// Read the attachment file
	attachmentFile, err := os.Open(attachmentPath)
	if err != nil {
		fmt.Println("Error opening attachment file:", err)
		return ""
	}
	defer attachmentFile.Close()

	// Read the attachment file contents
	attachmentData, err := ioutil.ReadAll(attachmentFile)
	if err != nil {
		fmt.Println("Error reading attachment file:", err)
		return ""
	}

	// Create the email message
	emailMessage := ""
	emailMessage += "From: " + sender + "\r\n"
	emailMessage += "To: " + recipient + "\r\n"
	emailMessage += "Subject: " + subject + "\r\n"
	emailMessage += "MIME-Version: 1.0\r\n"
	emailMessage += "Content-Type: multipart/mixed; boundary=boundary1234567890\r\n"
	emailMessage += "\r\n"
	emailMessage += "--boundary1234567890\r\n"
	emailMessage += "Content-Type: text/plain; charset=UTF-8\r\n"
	emailMessage += "Content-Transfer-Encoding: quoted-printable\r\n"
	emailMessage += "\r\n"
	emailMessage += body + "\r\n"
	emailMessage += "\r\n"
	emailMessage += "--boundary1234567890\r\n"
	emailMessage += "Content-Type: application/octet-stream\r\n"
	emailMessage += "Content-Transfer-Encoding: base64\r\n"
	emailMessage += "Content-Disposition: attachment; filename=\"" + attachmentPath + "\"\r\n"
	emailMessage += "\r\n"
	emailMessage += base64Encode(attachmentData) + "\r\n"
	emailMessage += "\r\n"
	emailMessage += "--boundary1234567890--\r\n"

	return base64Encode([]byte(emailMessage))
}

func base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}
