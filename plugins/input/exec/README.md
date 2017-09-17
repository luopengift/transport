### [inputs] plugin exec
```
{
    "exec": {
        "commands": [
            "python test/scripts/hello.py",
            "python test/scripts/world.py",
            "go run test/scripts/go.go"
        ],
        "cron": "* * * * * *"
    }
}
```
