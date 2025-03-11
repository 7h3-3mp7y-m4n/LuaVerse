math.randomseed(os.time())

function greet()
    local greetTable = {
        morning = {
            "Good morning, champion of as",
            "I hope you are having a nice start " .. name,
            "Hope you slept well. Time to work again, " .. name,
            "Wake the fuck up " .. name .. ", we got code to fix!"
        },
        evening = {
            "How’s your day going, " .. name .. "?",
            "Little do you know we slept too in the morning",
            "I hope you are doing okay, " .. name
        },
        night = {
            "Good evening, " .. name .. "! 🌙",
            "Winding down, " .. name .. "?",
            "Don't forget to rest, " .. name .. "!"
        }
    }

    local timeOfDay = "evening"
    if hour < 12 then
        timeOfDay = "morning"
    elseif hour >= 18 then
        timeOfDay = "night"
    end

    local options = greetTable[timeOfDay]
    if not options or #options == 0 then
        return "No greeting available."
    end

    local index = math.random(1, #options)
    return options[index]
end
