package cmd

import (
	"bytes"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
	ishell "gopkg.in/abiosoft/ishell.v2"

	"github.com/launchdarkly/ldc/api"
)

func addAuditLogCommands(shell *ishell.Shell) {
	root := &ishell.Cmd{
		Name: "log",
		Help: "search audit log entries",
		Func: func(c *ishell.Context) {
			options := make(map[string]interface{})
			if len(c.Args) > 0 {
				options["q"] = strings.Join(c.Args, " ")
				options["limit"] = 5
				//options["spec"] = ""
				//options["after"] = 1518163200000
				//options["before"] = time.Now().UnixNano() / int64(time.Millisecond)
			}
			auth := api.GetAuthCtx(getToken(nil))
			client, err := api.GetClient(getServer(currentConfig))
			if err != nil {
				c.Err(err)
				return
			}
			entries, _, err := client.AuditLogApi.GetAuditLogEntries(auth, options)
			if err != nil {
				c.Err(err)
				return
			}
			buf := bytes.Buffer{}
			table := tablewriter.NewWriter(&buf)
			//table.SetHeader([]string{"Key", "Name"})
			for _, entry := range entries.Items {

				table.Append([]string{time.Unix(entry.Date/1000, 0).Format("2006/01/02 15:04:05"), entry.Title})
			}
			table.Render()
			if buf.Len() > 1000 {
				c.Err(c.ShowPaged(buf.String()))
			} else {
				c.Println(buf.String())
			}
		},
	}

	shell.AddCmd(root)

}
