-- some manipulation wiht strings

a  = "one string"
b = string.gsub(a,"one","another")
print(a,"\n",b)
print(string.len(b))

-- for conversion we can use "tostring"
io.write("enter your number and it will be of type string")
line = io.read()
numberString = tostring(line)
print(type(numberString)," ",numberString)
-- slicing
print(string.sub(a,2,2))
