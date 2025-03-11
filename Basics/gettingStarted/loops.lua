fruits = {"apple", "carrot", "pineapple"}

for i, fruit in ipairs(fruits) do 
    print(i, fruit)
end
 --- while loop 

i = 1
while i <= 5 do
    print("count :" ,i)
    i = i + 1
end    

-- repat till forloop 
i = 1
 repeat
    print("Repeating:", i)
    i = i + 1
 until i > 5
-- normal forloop
for i = 1,5 do 
    print("local")
end    

---Loop Type	Key	Value	When to Use
---for i=1,10	index	(just i)	Numeric ranges
---ipairs	index	element	Array-style tables
---pairs	key	value	Key-value/dictionary tables