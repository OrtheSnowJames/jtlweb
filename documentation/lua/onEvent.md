# onEvent

document.onEvent loads a string of lua to execute when an event triggers.

EX:
```lua
-- Assuming we have an element of class 'buttonclass'
function handler()
    print("clicked")
end

function handlerrepeat()
    print("spam")
end

document.onEvent(".buttonclass", "click", [[handler()]]) -- will trigger every time the buton is clicked
document.onEvent(".buttonclass", "clickrepeat", [[handlerrepeat()]]) -- will trigger every frame the button is held
```