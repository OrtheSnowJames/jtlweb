function handle() 
    local x = document.get(".buttonclass")
    print(x)
    x.remove()
end

document.onEvent(".buttonclass", "click", [[handle()]])