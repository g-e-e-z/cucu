package presentation

import "github.com/g-e-e-z/cucu/commands"

func GetRequestStrings(request *commands.Request) []string {
    return []string{request.Method, request.Name}
}

