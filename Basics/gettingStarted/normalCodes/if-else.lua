io.write("hey welcome, what is your name ")
local name  = io.read()
io.write("And what is your age " .. name)
local age = tonumber(io.read())
if age > 18 and age < 30 then
    print("You’re in your 20s!")
elseif age < 18 then
    print("lol you are just a kid ".. name)
else
    print("wow old man ")
end    