package mail

import (
	"fmt"
	"strings"

	"github.com/eiicon-company/go-core/util/dsn"
)

type (
	stdoutMail struct {
		dsn *dsn.MailDSN
	}
)

func (m *stdoutMail) Send(data *Data) error {
	fmt.Printf("**************************************************\n")
	fmt.Printf("TO:%s\n", strings.Join(data.To, ","))
	fmt.Printf("CC:%s\n", strings.Join(data.Cc, ","))
	fmt.Printf("BCC:%s\n", strings.Join(data.Bcc, ","))
	fmt.Printf("From:%s\n", data.From)
	fmt.Printf("Subject:%s\n", data.Subject)
	fmt.Printf("**************************************************\n")
	if data.Text != nil {
		fmt.Println(string(data.Text))
	}
	if data.HTML != nil {
		fmt.Println(string(data.HTML))
	}
	fmt.Println("**************************************************")
	return nil
}
