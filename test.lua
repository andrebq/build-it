local http = require('http')
local json = require('json')
local cmd = require('cmd')

body, status, err = http.get("https://reqres.in/api/users/2")
print(body, status, err)

bodyObj, err = json.decode(body)
print(bodyObj.data.first_name, err)

output, err = cmd.run("catito", "test.lua")
print(output, err)

cmd.fatal(1, "bye bye")
