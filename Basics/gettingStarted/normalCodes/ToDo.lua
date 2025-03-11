local filename = "todo.txt"
local undolist = "undo.txt"

local function split(str,sep)
        local fields = {}
        for field in string.gmatch(str,"([^"..sep.."]+)") do
            table.insert(fields, field)
        end
        return fields
    end
local function save_undo()
    local src = io.open(filename, "r")
    if src then
        local dest = io.open(undo_file, "w")
        dest:write(src:read("*a"))
        src:close()
        dest:close()
    end
end    

local function load_tasks()
    local tasks = {}
    local file = io.open(filename, "r")
    if file then
        for line in file:lines() do
            local parts = split(line, "|")
            local done = tonumber(parts[1])
            local text = parts[2]
            table.insert(tasks, {done = done, text = text})
        end
        file:close()
    end
    return tasks
end

local function save_tasks(tasks)
    local file = io.open(filename, "w")
    for _, task in ipairs(tasks) do
        file:write(task.done .. "|" .. task.text .. "\n")
    end
    file:close()
end

local function add(task)
    save_undo()
    local tasks = load_tasks()
    table.insert(tasks, {done = 0, text = task})
    save_tasks(tasks)
    print("Added task: " .. task)
end

local function list()
    local tasks = load_tasks()
    if #tasks == 0 then
        print("📭 No tasks.")
    else
        print("📋 To-Do List:")
        for i, task in ipairs(tasks) do
            local status = task.done == 1 and "[x]" or "[ ]"
            print(i .. ". " .. status .. " " .. task.text)
        end
    end
end

local function remove(index)
    save_undo()
    local tasks = load_tasks()
    index = tonumber(index)
    if index and tasks[index] then
        local removed = table.remove(tasks, index)
        save_tasks(tasks)
        print("🗑️ Removed: " .. removed.text)
    else
        print("Invalid task number.")
    end
end

local function done(index)
    save_undo()
    local tasks = load_tasks()
    index = tonumber(index)
    if index and tasks[index] then
        tasks[index].done = 1
        save_tasks(tasks)
        print("🎉 Marked as done: " .. tasks[index].text)
    else
        print("Invalid task number.")
    end
end

local function undo()
    local file = io.open(undo_file, "r")
    if file then
        local backup = file:read("*a")
        local out = io.open(filename, "w")
        out:write(backup)
        out:close()
        file:close()
        print("↩️  Undo complete.")
    else
        print("❌ Nothing to undo.")
    end
end

local function help()
    print([[
🧾 To-Do CLI - Commands:
  lua todo.lua add "Task Name"       → Add task
  lua todo.lua list                  → List tasks
  lua todo.lua done <task number>    → Mark task as done
  lua todo.lua remove <task number>  → Remove task
  lua todo.lua undo                  → Undo last change
  lua todo.lua help                  → Show help
]])
end

local command = arg[1]

if command == "add" then
    local task = arg[2]
    if task then add(task) else print("❗ Provide a task.") end

elseif command == "list" then
    list()

elseif command == "remove" then
    local index = arg[2]
    if index then remove(index) else print("❗ Provide task number.") end

elseif command == "done" then
    local index = arg[2]
    if index then done(index) else print("❗ Provide task number to mark done.") end

elseif command == "undo" then
    undo()

else
    help()
end