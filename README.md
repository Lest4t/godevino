## GoDevino library

GoDevino library for api site http://devinitele.com

###Installation:

```bash
go get github.com/BurntSushi/toml
```

### Examples

```go
package main

import (
    "fmt"
    devino "github.com/maximal/godevino"
)

func main() {
    var sms *devino.Client

    devino.Username = "user"
    devino.Password = "password"

    _ = devino.Initialize()

    balance, err := sms.GetBalance()
    if err != nil {
        fmt.Printf("[ERR] Get balance error: %s\n", err)
    }

    fmt.Printf("[INF] Balance: %s\n", balance)

    message_ids, err := sms.SendMessage("TESTSMS", "79320001112", "Happy New Year", 0, "")
    if err != nil {
       fmt.Printf("[ERR] Send message error: %s\n", err)
    }

    fmt.Printf("[INF] Get message status...\n")
    msg_sts, err := sms.GetMessageState(message_ids[2:len(message_ids)-2])
    if err != nil {
        fmt.Printf("[ERR] Get message state: %s\n", err)
    }
    fmt.Printf("[INF] Message state: %s\n", msg_sts)
}
```
