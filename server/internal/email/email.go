package email

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

func getGmailService() (*gmail.Service, error) {
	clientID := os.Getenv("GMAIL_CLIENT_ID")
	clientSecret := os.Getenv("GMAIL_CLIENT_SECRET")
	refreshToken := os.Getenv("GMAIL_REFRESH_TOKEN")

	if clientID == "" || clientSecret == "" || refreshToken == "" {
		return nil, fmt.Errorf("Gmail OAuth2 credentials missing (CLIENT_ID, CLIENT_SECRET, REFRESH_TOKEN)")
	}

	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  "https://developers.google.com/oauthplayground",
		Scopes:       []string{"https://www.googleapis.com/auth/gmail.send"},
	}

	token := &oauth2.Token{
		RefreshToken: refreshToken,
		TokenType:    "Bearer", // Default assumption
	}

	client := config.Client(context.Background(), token)
	srv, err := gmail.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}
	return srv, nil
}

func sendEmailAPI(toEmail, toName, subject, htmlContent string) error {
	srv, err := getGmailService()
	if err != nil {
		return err
	}

	senderEmail := os.Getenv("SMTP_EMAIL")
	if senderEmail == "" {
		return fmt.Errorf("SMTP_EMAIL is not set")
	}

	// Create MIME Email Check
	// Headers
	messageStr := fmt.Sprintf("From: %s < %s >\r\n", "JobHMS Admin", senderEmail)
	messageStr += fmt.Sprintf("To: %s < %s >\r\n", toName, toEmail)
	messageStr += fmt.Sprintf("Subject: %s\r\n", subject)
	messageStr += "MIME-Version: 1.0\r\n"
	messageStr += "Content-Type: text/html; charset=\"utf-8\"\r\n\r\n"
	messageStr += htmlContent // Body

	// Encode to Base64 Url Safe
	rawMessage := base64.URLEncoding.EncodeToString([]byte(messageStr))

	msg := &gmail.Message{
		Raw: rawMessage,
	}

	_, err = srv.Users.Messages.Send("me", msg).Do()
	if err != nil {
		log.Printf("Error sending email via Gmail API: %v", err)
		return err
	}

	log.Printf("Email sent via Gmail API to %s", toEmail)
	return nil
}

func SendWelcomeEmail(toEmail, name string) error {
	subject := "Registration Successful - JobHMS Voting"
	htmlContent := fmt.Sprintf(`
        <h3>Welcome, %s!</h3>
        <p>Your registration for the upcoming election was successful.</p>
        <p>Please wait while the admin verifies your credentials (KTM & Profile).</p>
        <p>You will receive another email with your Voting Token once approved.</p>
    `, name)
	return sendEmailAPI(toEmail, name, subject, htmlContent)
}

func SendReminderEmail(toEmail, name, token string) error {
	subject := "Election Reminder: Your Token Inside"
	htmlContent := fmt.Sprintf(`
        <h3>Hello, %s</h3>
        <p>The election is starting soon!</p>
        <p><strong>Your Voting Token is: %s</strong></p>
        <p>Do not share this token. Use it to log in when the election starts.</p>
        <p>See you at the polls!</p>
    `, name, token)
	return sendEmailAPI(toEmail, name, subject, htmlContent)
}

func SendVoteConfirmation(toEmail, name, candidateName string) error {
	subject := "Vote Confirmed"
	htmlContent := fmt.Sprintf(`
        <h3>Thank you, %s</h3>
        <p>Your vote for <strong>%s</strong> has been successfully recorded.</p>
        <p>Each vote counts!</p>
    `, name, candidateName)
	return sendEmailAPI(toEmail, name, subject, htmlContent)
}

func SendApprovalEmail(toEmail, name, token string) error {
	subject := "Account Approved - Your Voting Token"
	htmlContent := fmt.Sprintf(`
        <h3>Congratulations, %s!</h3>
        <p>Your account has been verified.</p>
        <p><strong>Your Voting Token is: %s</strong></p>
        <p>Please keep this token safe. You will need it to log in when the election starts.</p>
    `, name, token)
	return sendEmailAPI(toEmail, name, subject, htmlContent)
}
