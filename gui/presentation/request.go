package presentation

import "github.com/g-e-e-z/cucu/commands"

func GetRequestStrings(request *commands.Request) []string {
    name := request.Name
    request.CheckModifed()
    if request.Modified {
        name += " *"
    }
    return []string{request.Method, name}
}

