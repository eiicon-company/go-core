package mail

import (
	"fmt"
	"strings"

	"github.com/eiicon-company/go-core/util/dsn"
)

type (
	stdoutMail struct {
		dsn *dsn.MailDSN

		To      []string
		Bcc     []string
		Cc      []string
		From    string
		Subject string
		Text    []byte
		HTML    []byte
	}
)

func (m *stdoutMail) Send() error {
	fmt.Printf("**************************************************")
	fmt.Printf("TO:%s", strings.Join(m.To, ","))
	fmt.Printf("CC:%s", strings.Join(m.Cc, ","))
	fmt.Printf("BCC:%s", strings.Join(m.Bcc, ","))
	fmt.Printf("From:%s", m.From)
	fmt.Printf("Subject:%s", m.Subject)
	fmt.Println("**************************************************")
	if m.Text != nil {
		fmt.Println(string(m.Text[:]))
	}
	if m.HTML != nil {
		fmt.Println(string(m.HTML[:]))
	}
	fmt.Println("**************************************************")
	return nil
}
