package library

// import (
// 	"os"

// 	"github.com/aws/aws-sdk-go/aws"
// 	"github.com/aws/aws-sdk-go/aws/session"
// 	"github.com/aws/aws-sdk-go/service/ses"
// )

// // SendMailAmazon ...
// func SendMailAmazon(title string, body string, email string) error {

// 	// Upload
// 	os.Setenv("AWS_ACCESS_KEY_ID", "AKIAIRXDX4EC2YGQX2RA")
// 	os.Setenv("AWS_SECRET_ACCESS_KEY", "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")

// 	conf := aws.Config{Region: aws.String("eu-west-1")}
// 	awsSession := session.New(&conf)

// 	sesSession := ses.New(awsSession)

// 	sesEmailInput := &ses.SendEmailInput{
// 		Destination: &ses.Destination{
// 			ToAddresses: []*string{aws.String(email)},
// 		},
// 		Message: &ses.Message{
// 			Body: &ses.Body{
// 				Html: &ses.Content{
// 					Data: aws.String(body)},
// 			},
// 			Subject: &ses.Content{
// 				Data: aws.String(title),
// 			},
// 		},
// 		Source: aws.String("hanyny.com@gmail.com"),
// 	}

// 	_, err := sesSession.SendEmail(sesEmailInput)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
