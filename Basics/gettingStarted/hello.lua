-- As every new language we write the good old hello!
print("Hello from lua")

-- A simple way use a function in lua!
function greet(name)
    print("Hello", name , "!")
end

-- I used io write to keep things in the same line

io.write("What should I call you, learner? ")
name = io.read(...)
greet(name)