function handle()
    local x = document.create("p")
    x.text = "Hello World!"
    print(x.text)
    print(x)
    document.replace(".buttonclass", x)
end

document.onEvent(".buttonclass", "click", [[handle()]])