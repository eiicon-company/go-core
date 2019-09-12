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
	fmt.Printf("**************************************************")
	fmt.Printf("TO:%s", strings.Join(data.To, ","))
	fmt.Printf("CC:%s", strings.Join(data.Cc, ","))
	fmt.Printf("BCC:%s", strings.Join(data.Bcc, ","))
	fmt.Printf("From:%s", data.From)
	fmt.Printf("Subject:%s", data.Subject)
	fmt.Println("**************************************************")
	if data.Text != nil {
		fmt.Println(string(data.Text[:]))
	}
	if data.HTML != nil {
		fmt.Println(string(data.HTML[:]))
	}
	fmt.Println("**************************************************")
	return nil
}
