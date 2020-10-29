# line-login

## _Usage_
``` 
$ go get -u -v github.com/5hields/line-login
```

## _Note_
1. Enable a LetterSealing
2. Allow SecondaryLogin (email & password)

## _Example_
```go
impoet (
    "fmt"
    ll "github.com/5hields/line-login"
)

func main() {
    token, err := ll.LoginWithCredential("Your Email", "Your Pass")
    if err != nil {
    // hoge
    }
    fmt.Println("AuthToken:" + token)
}
```

## _support_
LineVersion: 10.17.0 ~<br>
Go Version:  go1.15.3 ~

## _BugReport_
0authn@protonmail.com
