-- lets explore some data type 
-- good thing its dynamically typed 
x = nil                         -- nil
isLuaFun = true                 -- boolean
age = 25                        -- number
pi = 3.14                       -- number (float)
name = "Rashid"                 -- string
person = {name = "Rashid", age = 25}  -- table/dict
greet = function() print("Hello!") end  -- function

print(type(x))
print(type(isLuaFun))
print(type(age))
print(type(name))
print(type(person))
print(type(greet))
