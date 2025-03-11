math.randomseed(os.time())

io.write("Put the max range of number you want to guess: ")
local numberInput = tonumber(io.read())
local secret = math.random(0, numberInput)
io.write("How many lives do you want ♥️: ")
local lives = tonumber(io.read())
print("🎯 Welcome to the Guessing Game!")
print("I'm thinking of a number between 0 and " .. numberInput .. ". You have " .. lives .. " tries.")

for i = 1, lives do
    io.write("Attempt " .. i .. ": Your guess? ")
    local guess = tonumber(io.read())
    if guess == nil then 
        print("That is not a valid number you have entered.")
    elseif guess == secret then 
        print("Congratulations! You guessed it right!!!")
        return
    elseif guess > secret then 
        print("The number you have entered is too high!")
    else
        print("Ooo, you are close! Just aim for a higher number.")   
    end
end

print("😢 Out of tries! The number was: " .. secret)
